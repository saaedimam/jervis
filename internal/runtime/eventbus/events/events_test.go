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
	priorities := []int{
		events.MinPriority,
		events.MaxPriority,
		events.Low,
		events.Normal,
		events.High,
		events.Critical,
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

	if err := events.ValidatePriority(-101); !errors.Is(err, errs.ErrInvalidPriority) {
		t.Fatalf("expected ErrInvalidPriority for -101")
	}
	if err := events.ValidatePriority(101); !errors.Is(err, errs.ErrInvalidPriority) {
		t.Fatalf("expected ErrInvalidPriority for 101")
	}
}

func TestValidateEventType(t *testing.T) {
	validTypes := []string{"runtime.lifecycle.booted", "memory.timeline.appended", "a.b"}
	for _, vt := range validTypes {
		if err := events.ValidateEventType(vt); err != nil {
			t.Errorf("expected event type %q to be valid: %v", vt, err)
		}
	}

	invalidTypes := []struct {
		eventType string
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

func TestEventBuilderSuccessfulBuild(t *testing.T) {
	evtID, err := types.NewEventID("evt-001")
	if err != nil {
		t.Fatalf("unexpected error creating EventID: %v", err)
	}

	now := types.NewTimestamp(time.Now())
	evt, err := events.NewBuilder().
		SetID(evtID).
		SetType("runtime.lifecycle.booted").
		SetSource("internal/runtime/lifecycle").
		SetTimestamp(now).
		SetCorrelationID("corr-001").
		SetCausationID("cause-001").
		SetPriority(events.High).
		SetPayload("payload-data").
		SetMetadata("key1", "val1").
		SetVersion("1.0.0").
		Build()

	if err != nil {
		t.Fatalf("unexpected error building event: %v", err)
	}

	if evt.ID() != evtID {
		t.Errorf("expected ID %v, got %v", evtID, evt.ID())
	}
	if evt.Type() != "runtime.lifecycle.booted" {
		t.Errorf("expected Type runtime.lifecycle.booted, got %s", evt.Type())
	}
	if evt.Source() != "internal/runtime/lifecycle" {
		t.Errorf("expected Source internal/runtime/lifecycle, got %s", evt.Source())
	}
	if evt.Timestamp() != now {
		t.Errorf("expected Timestamp %v, got %v", now, evt.Timestamp())
	}
	if evt.CorrelationID() != "corr-001" {
		t.Errorf("expected CorrelationID corr-001, got %s", evt.CorrelationID())
	}
	if evt.CausationID() != "cause-001" {
		t.Errorf("expected CausationID cause-001, got %s", evt.CausationID())
	}
	if evt.Priority() != events.High {
		t.Errorf("expected Priority High (%d), got %d", events.High, evt.Priority())
	}
	if evt.Payload() != "payload-data" {
		t.Errorf("expected Payload payload-data, got %v", evt.Payload())
	}
	if evt.Version() != "1.0.0" {
		t.Errorf("expected Version 1.0.0, got %s", evt.Version())
	}

	meta := evt.Metadata()
	if meta["key1"] != "val1" {
		t.Errorf("expected metadata key1=val1, got %s", meta["key1"])
	}

	// Test metadata immutability copy
	meta["key1"] = "mutated"
	if evt.Metadata()["key1"] != "val1" {
		t.Errorf("metadata was mutated directly, defensive copy failed")
	}
}

func TestEventBuilderDefaultsAndAutoPopulate(t *testing.T) {
	evtID, _ := types.NewEventID("evt-auto")

	evt, err := events.NewBuilder().
		SetID(evtID).
		SetType("runtime.test.event").
		SetSource("test").
		SetPayload("payload").
		Build()

	if err != nil {
		t.Fatalf("unexpected build error: %v", err)
	}

	if evt.Timestamp().IsZero() {
		t.Errorf("expected auto-populated timestamp")
	}
	if evt.CorrelationID() != "evt-auto" {
		t.Errorf("expected correlation ID auto-populated from EventID, got %s", evt.CorrelationID())
	}
	if evt.CausationID() != "evt-auto" {
		t.Errorf("expected causation ID auto-populated from EventID, got %s", evt.CausationID())
	}
	if evt.Priority() != events.DefaultPriority {
		t.Errorf("expected default priority 0, got %d", evt.Priority())
	}
	if evt.Version() != events.DefaultVersion {
		t.Errorf("expected default version 1.0.0, got %s", evt.Version())
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
	builderPriority := events.NewBuilder().SetID(evtID).SetType("runtime.test.event").SetSource("test").SetPayload("p").SetPriority(200)
	if _, err := builderPriority.Build(); !errors.Is(err, errs.ErrInvalidPriority) {
		t.Fatalf("expected ErrInvalidPriority for priority 200, got %v", err)
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
	evt, err := uninitBuilder.Build()
	if err != nil {
		t.Fatalf("unexpected error building uninit builder: %v", err)
	}
	if evt.Metadata()["k"] != "v" {
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
func (m *mockEventInvalidTS) Priority() int               { return 0 }
func (m *mockEventInvalidTS) Payload() any                { return "payload" }
func (m *mockEventInvalidTS) Metadata() map[string]string { return nil }
func (m *mockEventInvalidTS) Version() string             { return "1.0.0" }
