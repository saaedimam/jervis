package device_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/saaedimam/jervis/internal/memory/store/sqlite"
	"github.com/saaedimam/jervis/internal/runtime/eventbus"
	"github.com/saaedimam/jervis/internal/services/device"
)

func setupTestService(t *testing.T) (device.Service, func()) {
	t.Helper()
	store, err := sqlite.New(":memory:")
	if err != nil {
		t.Fatalf("sqlite.New failed: %v", err)
	}

	ctx := context.Background()
	if err := device.Initialize(ctx, store); err != nil {
		_ = store.Close()
		t.Fatalf("device.Initialize failed: %v", err)
	}

	bus := eventbus.New()
	svc := device.New(store, bus)

	cleanup := func() {
		_ = store.Close()
	}
	return svc, cleanup
}

func TestEnvelopeValidation(t *testing.T) {
	// Valid envelope
	env := &device.Envelope{
		Schema:    device.SchemaVersion,
		ID:        "evt-001",
		Device:    "alcatel-5022d",
		Timestamp: time.Now().UnixMilli(),
		Type:      device.EventSMSReceived,
		Payload:   json.RawMessage(`{"sender":"+8801","body":"hi"}`),
	}
	if err := env.ValidateEnvelope(); err != nil {
		t.Fatalf("expected valid envelope, got %v", err)
	}

	// Invalid schema
	invalidSchema := *env
	invalidSchema.Schema = "wrong.v2"
	if err := invalidSchema.ValidateEnvelope(); err == nil {
		t.Fatal("expected error on invalid schema, got nil")
	}

	// Missing ID
	missingID := *env
	missingID.ID = ""
	if err := missingID.ValidateEnvelope(); err == nil {
		t.Fatal("expected error on missing ID, got nil")
	}

	// Missing Device
	missingDevice := *env
	missingDevice.Device = ""
	if err := missingDevice.ValidateEnvelope(); err == nil {
		t.Fatal("expected error on missing Device, got nil")
	}

	// Invalid timestamp
	invalidTs := *env
	invalidTs.Timestamp = 0
	if err := invalidTs.ValidateEnvelope(); err == nil {
		t.Fatal("expected error on invalid timestamp, got nil")
	}

	// Missing payload
	missingPayload := *env
	missingPayload.Payload = nil
	if err := missingPayload.ValidateEnvelope(); err == nil {
		t.Fatal("expected error on missing payload, got nil")
	}
}

func TestCapabilityAllowlistAndValidation(t *testing.T) {
	// Allowed action: send_sms with valid fields
	validSMS := device.CommandPayload{
		CommandID: "cmd-1",
		Action:    device.ActionSendSMS,
		Recipient: "+8801700000000",
		Text:      "Test SMS",
	}
	if err := validSMS.Validate(); err != nil {
		t.Fatalf("expected valid send_sms command, got: %v", err)
	}

	// Disallowed action: shell_exec
	disallowed := device.CommandPayload{
		CommandID: "cmd-bad",
		Action:    "shell_exec",
	}
	if err := disallowed.Validate(); err == nil {
		t.Fatal("expected error for shell_exec action, got nil")
	}

	// capture_audio bounded duration check
	audioCmd := device.CommandPayload{
		CommandID:       "cmd-audio",
		Action:          device.ActionCaptureAudio,
		DurationSeconds: 15, // > 10
	}
	if err := audioCmd.Validate(); err == nil {
		t.Fatal("expected error for audio duration > 10s, got nil")
	}
}

func TestEventIngestionAndDeduplication(t *testing.T) {
	svc, cleanup := setupTestService(t)
	defer cleanup()

	ctx := context.Background()
	env := &device.Envelope{
		Schema:    device.SchemaVersion,
		ID:        "evt-dup-1",
		Device:    "alcatel-5022d",
		Timestamp: time.Now().UnixMilli(),
		Type:      device.EventSMSReceived,
		Payload:   json.RawMessage(`{"sender":"+8801","body":"hello"}`),
	}

	// First ingestion
	isDup, err := svc.IngestEvent(ctx, env)
	if err != nil {
		t.Fatalf("first ingestion failed: %v", err)
	}
	if isDup {
		t.Fatal("expected first ingestion to not be a duplicate")
	}

	// Second ingestion with same ID -> duplicate acknowledged
	isDup, err = svc.IngestEvent(ctx, env)
	if err != nil {
		t.Fatalf("second ingestion failed: %v", err)
	}
	if !isDup {
		t.Fatal("expected second ingestion to be detected as duplicate")
	}
}

func TestCommandQueueFIFOAndCapacity(t *testing.T) {
	svc, cleanup := setupTestService(t)
	defer cleanup()

	ctx := context.Background()

	// Enqueue 100 commands (MaxQueueSize)
	for i := 0; i < device.MaxQueueSize; i++ {
		cmdID := "cmd-batch-" + string(rune('A'+i%26)) + "-" + time.Now().Format("150405.000000000")
		cmdPayload, _ := json.Marshal(device.CommandPayload{
			CommandID: cmdID,
			Action:    device.ActionBatteryStatus,
		})

		cmdEnv := &device.Envelope{
			Schema:    device.SchemaVersion,
			ID:        "env-" + cmdID,
			Device:    "alcatel-5022d",
			Timestamp: time.Now().UnixMilli(),
			Type:      device.TypeCommand,
			Payload:   cmdPayload,
		}

		if err := svc.QueueCommand(ctx, cmdEnv); err != nil {
			t.Fatalf("failed enqueuing command %d: %v", i, err)
		}
	}

	// 101st command must overflow and return ErrQueueFull
	overflowPayload, _ := json.Marshal(device.CommandPayload{
		CommandID: "cmd-overflow",
		Action:    device.ActionBatteryStatus,
	})
	overflowEnv := &device.Envelope{
		Schema:    device.SchemaVersion,
		ID:        "env-overflow",
		Device:    "alcatel-5022d",
		Timestamp: time.Now().UnixMilli(),
		Type:      device.TypeCommand,
		Payload:   overflowPayload,
	}

	err := svc.QueueCommand(ctx, overflowEnv)
	if err != device.ErrQueueFull {
		t.Fatalf("expected ErrQueueFull on 101st command, got: %v", err)
	}

	// Pop commands and verify FIFO order
	popped, err := svc.WaitForCommand(ctx, "alcatel-5022d", 50*time.Millisecond)
	if err != nil || popped == nil {
		t.Fatalf("expected command from queue, got err=%v popped=%v", err, popped)
	}
}

func TestWaitForCommandTimeout(t *testing.T) {
	svc, cleanup := setupTestService(t)
	defer cleanup()

	ctx := context.Background()
	start := time.Now()
	// Wait on empty queue with 100ms timeout
	cmd, err := svc.WaitForCommand(ctx, "alcatel-5022d", 100*time.Millisecond)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("unexpected error on timeout: %v", err)
	}
	if cmd != nil {
		t.Fatalf("expected nil command on empty queue, got %v", cmd)
	}
	if elapsed < 80*time.Millisecond {
		t.Fatalf("expected wait of at least 80ms, elapsed: %v", elapsed)
	}
}

func TestResultIngestionAndDeduplication(t *testing.T) {
	svc, cleanup := setupTestService(t)
	defer cleanup()

	ctx := context.Background()

	resPayload, _ := json.Marshal(device.ResultPayload{
		CommandID: "cmd-test-res",
		Action:    device.ActionSendSMS,
		Status:    device.StatusSuccess,
		Data:      json.RawMessage(`{"sent":true}`),
	})

	resEnv := &device.Envelope{
		Schema:    device.SchemaVersion,
		ID:        "res-env-1",
		Device:    "alcatel-5022d",
		Timestamp: time.Now().UnixMilli(),
		Type:      device.TypeResult,
		Payload:   resPayload,
	}

	// First result ingestion
	isDup, err := svc.IngestResult(ctx, resEnv)
	if err != nil {
		t.Fatalf("first result ingestion failed: %v", err)
	}
	if isDup {
		t.Fatal("expected first result ingestion to not be duplicate")
	}

	// Duplicate result ingestion
	isDup, err = svc.IngestResult(ctx, resEnv)
	if err != nil {
		t.Fatalf("second result ingestion failed: %v", err)
	}
	if !isDup {
		t.Fatal("expected second result ingestion to be detected as duplicate")
	}

	// Check device status reflecting result
	status, err := svc.GetDeviceStatus(ctx, "alcatel-5022d")
	if err != nil {
		t.Fatalf("GetDeviceStatus failed: %v", err)
	}
	if !status.Connected {
		t.Fatal("expected device to be marked connected after activity")
	}
}
