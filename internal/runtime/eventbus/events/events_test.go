package events_test

import (
	"errors"
	"testing"
	"time"

	errs "github.com/ioriimasu/jervis/internal/runtime/eventbus/errors"
	"github.com/ioriimasu/jervis/internal/runtime/eventbus/events"
	"github.com/ioriimasu/jervis/internal/runtime/types"
)

func TestValidatePriority(t *testing.T) {
	priorities := []events.Priority{
		events.PriorityLow,
		events.PriorityNormal,
		events.PriorityHigh,
		events.PriorityCritical,
	}

	for _, p := range priorities {
		if err := events.ValidatePriority(p); err != nil {
			t.Errorf("expected priority %d to be valid: %v", p, err)
		}
	}

	if err := events.ValidatePriority(events.PriorityCritical + 1); !errors.Is(err, errs.ErrInvalidPriority) {
		t.Fatalf("expected ErrInvalidPriority for PriorityCritical + 1")
	}
}

func TestValidateEventType(t *testing.T) {
	validTypes := []events.EventType{"runtime.lifecycle.booted", "memory.timeline.appended", "a.b"}
	for _, vt := range validTypes {
		if err := events.ValidateEventType(vt); err != nil {
			t.Errorf("expected event type %q to be valid: %v", vt, err)
		}
		if vt.String() != string(vt) {
			t.Errorf("String() mismatch for %q", vt)
		}
	}

	invalidTypes := []struct {
		eventType events.EventType
		desc      string
	}{
		{"", "empty"},
		{"Runtime.Lifecycle", "uppercase"},
		{"runtime lifecycle", "spaces"},
		{"singleword", "no namespace dot"},
		{"runtime..booted", "empty namespace segment"},
		{".runtime.booted", "leading dot"},
		{"runtime.booted.", "trailing dot"},
	}

	for _, tt := range invalidTypes {
		t.Run(tt.desc, func(t *testing.T) {
			if err := events.ValidateEventType(tt.eventType); !errors.Is(err, errs.ErrValidationFailed) {
				t.Errorf("expected ErrValidationFailed for %q (%s), got %v", tt.eventType, tt.desc, err)
			}
		})
	}
}

func TestEventBuilderAndClone(t *testing.T) {
	evtID, err := types.NewEventID("evt-001")
	if err != nil {
		t.Fatalf("unexpected error creating EventID: %v", err)
	}

	now := types.NewTimestamp(time.Now())
	env, err := events.NewBuilder().
		SetID(evtID).
		SetType("runtime.lifecycle.booted").
		SetSource("internal/runtime/lifecycle").
		SetTimestamp(now).
		SetCorrelationID("corr-001").
		SetCausationID("cause-001").
		SetPriority(events.PriorityHigh).
		SetPayload("payload-data").
		SetMetadata("key1", "val1").
		SetVersion("1.0.0").
		Build()
	if err != nil {
		t.Fatalf("unexpected error building event: %v", err)
	}

	if env.ID() != evtID {
		t.Errorf("expected ID %v, got %v", evtID, env.ID())
	}
	if env.Type() != "runtime.lifecycle.booted" {
		t.Errorf("expected Type runtime.lifecycle.booted, got %s", env.Type())
	}
	if env.Source() != "internal/runtime/lifecycle" {
		t.Errorf("expected Source internal/runtime/lifecycle, got %s", env.Source())
	}
	if env.Timestamp() != now {
		t.Errorf("expected Timestamp %v, got %v", now, env.Timestamp())
	}
	if env.CorrelationID() != "corr-001" {
		t.Errorf("expected CorrelationID corr-001, got %s", env.CorrelationID())
	}
	if env.CausationID() != "cause-001" {
		t.Errorf("expected CausationID cause-001, got %s", env.CausationID())
	}
	if env.Priority() != uint8(events.PriorityHigh) {
		t.Errorf("expected Priority High (%d), got %d", uint8(events.PriorityHigh), env.Priority())
	}
	if env.Payload() != "payload-data" {
		t.Errorf("expected Payload payload-data, got %v", env.Payload())
	}
	if env.Version() != "1.0.0" {
		t.Errorf("expected Version 1.0.0, got %s", env.Version())
	}

	// Test Clone
	cloned := env.Clone()
	if cloned == nil {
		t.Fatalf("expected non-nil cloned envelope")
	}
	if cloned.ID() != env.ID() || cloned.Type() != env.Type() || cloned.Payload() != env.Payload() {
		t.Fatalf("cloned envelope mismatch")
	}

	// Test nil Clone
	var nilEnv *events.Envelope
	if nilEnv.Clone() != nil {
		t.Fatalf("expected nil for nil envelope Clone()")
	}

	// Test metadata immutability copy
	meta := env.Metadata()
	meta["key1"] = "mutated"
	if env.Metadata()["key1"] != "val1" {
		t.Errorf("metadata was mutated directly, defensive copy failed")
	}
}

func TestEventBuilderDefaultsAndAutoPopulate(t *testing.T) {
	evtID, _ := types.NewEventID("evt-auto")

	env, err := events.NewBuilder().
		SetID(evtID).
		SetType("runtime.test.event").
		SetSource("test").
		SetPayload("payload").
		Build()
	if err != nil {
		t.Fatalf("unexpected build error: %v", err)
	}

	if env.Timestamp().IsZero() {
		t.Errorf("expected auto-populated timestamp")
	}
	if env.CorrelationID() != "evt-auto" {
		t.Errorf("expected correlation ID auto-populated from EventID, got %s", env.CorrelationID())
	}
	if env.CausationID() != "evt-auto" {
		t.Errorf("expected causation ID auto-populated from EventID, got %s", env.CausationID())
	}
	if env.Priority() != uint8(events.DefaultPriority) {
		t.Errorf("expected default priority %d, got %d", events.DefaultPriority, env.Priority())
	}
	if env.Version() != events.DefaultVersion {
		t.Errorf("expected default version 1.0.0, got %s", env.Version())
	}
}

func TestValidateEventValidationFailures(t *testing.T) {
	if err := events.ValidateEvent(nil); !errors.Is(err, errs.ErrInvalidEvent) {
		t.Fatalf("expected ErrInvalidEvent for nil event")
	}

	// Missing ID
	builder := events.NewBuilder().SetType("runtime.test.event").SetSource("test").SetPayload("p")
	if _, err := builder.Build(); !errors.Is(err, errs.ErrValidationFailed) {
		t.Fatalf("expected ErrValidationFailed for zero EventID")
	}

	evtID, _ := types.NewEventID("evt-123")

	// Nil Payload
	builderNilPayload := events.NewBuilder().SetID(evtID).SetType("runtime.test.event").SetSource("test").SetPayload(nil)
	if _, err := builderNilPayload.Build(); !errors.Is(err, errs.ErrValidationFailed) {
		t.Fatalf("expected ErrValidationFailed for nil payload")
	}

	// Invalid Type
	builderInvalidType := events.NewBuilder().SetID(evtID).SetType("INVALID_TYPE").SetSource("test").SetPayload("p")
	if _, err := builderInvalidType.Build(); !errors.Is(err, errs.ErrValidationFailed) {
		t.Fatalf("expected ErrValidationFailed for invalid event type")
	}

	// Empty Source
	builderSource := events.NewBuilder().SetID(evtID).SetType("runtime.test.event").SetSource("").SetPayload("p")
	if _, err := builderSource.Build(); !errors.Is(err, errs.ErrValidationFailed) {
		t.Fatalf("expected ErrValidationFailed for empty Source")
	}

	// Zero Timestamp
	evtZeroTS := &mockEventInvalidTS{id: evtID}
	if err := events.ValidateEvent(evtZeroTS); !errors.Is(err, errs.ErrValidationFailed) {
		t.Fatalf("expected ErrValidationFailed for zero timestamp, got %v", err)
	}

	// Out of range Priority
	builderPriority := events.NewBuilder().SetID(evtID).SetType("runtime.test.event").SetSource("test").SetPayload("p").SetPriority(events.PriorityCritical + 1)
	if _, err := builderPriority.Build(); !errors.Is(err, errs.ErrInvalidPriority) {
		t.Fatalf("expected ErrInvalidPriority for priority out of bounds, got %v", err)
	}

	// Empty Version
	builderVersion := events.NewBuilder().SetID(evtID).SetType("runtime.test.event").SetSource("test").SetPayload("p").SetVersion("")
	if _, err := builderVersion.Build(); !errors.Is(err, errs.ErrValidationFailed) {
		t.Fatalf("expected ErrValidationFailed for empty version")
	}
}

func TestUninitializedBuilderAndNilMetadata(t *testing.T) {
	uninitBuilder := &events.Builder{}
	uninitBuilder.SetMetadata("k", "v").SetVersion("1.0.0")

	evtID, _ := types.NewEventID("evt-uninit")
	uninitBuilder.SetID(evtID).SetType("runtime.test.event").SetSource("test").SetPayload("data")
	env, err := uninitBuilder.Build()
	if err != nil {
		t.Fatalf("unexpected error building uninit builder: %v", err)
	}
	if env.Metadata()["k"] != "v" {
		t.Fatalf("expected metadata k=v")
	}
}

type mockEventInvalidTS struct {
	id types.EventID
}

func (m *mockEventInvalidTS) ID() types.EventID           { return m.id }
func (m *mockEventInvalidTS) Type() string                { return "runtime.test.event" }
func (m *mockEventInvalidTS) Source() string              { return "test" }
func (m *mockEventInvalidTS) Timestamp() types.Timestamp  { return types.Timestamp{} }
func (m *mockEventInvalidTS) CorrelationID() string       { return "corr" }
func (m *mockEventInvalidTS) CausationID() string         { return "cause" }
func (m *mockEventInvalidTS) Priority() uint8             { return uint8(events.PriorityNormal) }
func (m *mockEventInvalidTS) Payload() any                { return "payload" }
func (m *mockEventInvalidTS) Metadata() map[string]string { return nil }
func (m *mockEventInvalidTS) Version() string             { return "1.0.0" }
