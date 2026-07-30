package eventbus_test

import (
	"errors"
	"testing"

	"github.com/saaedimam/jervis/internal/runtime/eventbus"
	"github.com/saaedimam/jervis/internal/runtime/eventbus/contracts"
	"github.com/saaedimam/jervis/internal/runtime/eventbus/dispatcher"
	errs "github.com/saaedimam/jervis/internal/runtime/eventbus/errors"
	"github.com/saaedimam/jervis/internal/runtime/eventbus/events"
	"github.com/saaedimam/jervis/internal/runtime/eventbus/middleware"
	"github.com/saaedimam/jervis/internal/runtime/eventbus/subscription"
	"github.com/saaedimam/jervis/internal/runtime/types"
)

type mockHandler struct {
	id       string
	failErr  error
	panicVal any
	called   *[]string
}

func (m *mockHandler) ID() string { return m.id }
func (m *mockHandler) Handle(event contracts.Event) error {
	if m.called != nil {
		*m.called = append(*m.called, m.id)
	}
	if m.panicVal != nil {
		panic(m.panicVal)
	}
	return m.failErr
}

func createTestEvent(t *testing.T, eventType string) contracts.Event {
	evtID, _ := types.NewEventID("evt-001")
	env, err := events.NewBuilder().
		SetID(evtID).
		SetType(events.EventType(eventType)).
		SetSource("test").
		SetPayload("payload").
		Build()
	if err != nil {
		t.Fatalf("failed to build test event: %v", err)
	}
	return env
}

func TestBusPublishSubscribeUnsubscribeSuccess(t *testing.T) {
	bus := eventbus.New()
	if bus.Count() != 0 {
		t.Fatalf("expected count 0, got %d", bus.Count())
	}

	var called []string
	h1 := &mockHandler{id: "h1", called: &called}
	h2 := &mockHandler{id: "h2", called: &called}

	subID1, err := bus.Subscribe("system.user.*", h1, events.PriorityHigh)
	if err != nil {
		t.Fatalf("unexpected subscribe error: %v", err)
	}
	if subID1.IsZero() {
		t.Fatalf("expected non-zero subID")
	}

	subID2, err := bus.Subscribe("system.user.created", h2, events.PriorityLow)
	if err != nil {
		t.Fatalf("unexpected subscribe error: %v", err)
	}
	_ = subID2

	if bus.Count() != 2 {
		t.Fatalf("expected count 2, got %d", bus.Count())
	}

	evt := createTestEvent(t, "system.user.created")
	if err := bus.Publish(evt); err != nil {
		t.Fatalf("unexpected publish error: %v", err)
	}

	if len(called) != 2 || called[0] != "h1" || called[1] != "h2" {
		t.Fatalf("expected execution order [h1, h2], got %v", called)
	}

	// Test Unsubscribe
	if err := bus.Unsubscribe(subID1); err != nil {
		t.Fatalf("unexpected unsubscribe error: %v", err)
	}
	if bus.Count() != 1 {
		t.Fatalf("expected count 1 after unsubscribe, got %d", bus.Count())
	}

	called = nil
	if err := bus.Publish(evt); err != nil {
		t.Fatalf("unexpected publish error: %v", err)
	}
	if len(called) != 1 || called[0] != "h2" {
		t.Fatalf("expected execution order [h2], got %v", called)
	}
}

func TestBusPublishZeroHandlers(t *testing.T) {
	bus := eventbus.New()
	evt := createTestEvent(t, "system.unsubscribed.event")

	if err := bus.Publish(evt); err != nil {
		t.Fatalf("expected nil error for publish with zero handlers, got %v", err)
	}
}

func TestBusPublishInvalidEvent(t *testing.T) {
	bus := eventbus.New()

	if err := bus.Publish(nil); !errors.Is(err, errs.ErrInvalidEvent) {
		t.Fatalf("expected ErrInvalidEvent, got %v", err)
	}
}

func TestBusMiddlewareOrderingAndShortCircuit(t *testing.T) {
	bus := eventbus.New()

	var sequence []string
	m1 := middleware.Func(func(evt contracts.Event, next func(contracts.Event) error) error {
		sequence = append(sequence, "M1-Pre")
		err := next(evt)
		sequence = append(sequence, "M1-Post")
		return err
	})
	m2 := middleware.Func(func(evt contracts.Event, next func(contracts.Event) error) error {
		sequence = append(sequence, "M2-Pre")
		err := next(evt)
		sequence = append(sequence, "M2-Post")
		return err
	})

	bus.Use(m1, m2)

	h := &mockHandler{id: "h1", called: &sequence}
	_, err := bus.Subscribe("test.event", h, events.PriorityNormal)
	if err != nil {
		t.Fatalf("unexpected subscribe error: %v", err)
	}

	evt := createTestEvent(t, "test.event")
	if err := bus.Publish(evt); err != nil {
		t.Fatalf("unexpected publish error: %v", err)
	}

	expectedSeq := []string{"M1-Pre", "M2-Pre", "h1", "M2-Post", "M1-Post"}
	if len(sequence) != len(expectedSeq) {
		t.Fatalf("expected sequence %v, got %v", expectedSeq, sequence)
	}
	for i, s := range expectedSeq {
		if sequence[i] != s {
			t.Errorf("pos %d: got %s, want %s", i, sequence[i], s)
		}
	}
}

func TestBusPublishHandlerFailureAndPanicRecovery(t *testing.T) {
	bus := eventbus.New()

	h1 := &mockHandler{id: "h1", failErr: errors.New("h1 failed")}
	h2 := &mockHandler{id: "h2", panicVal: "h2 boom"}

	_, _ = bus.Subscribe("test.fail", h1, events.PriorityNormal)
	_, _ = bus.Subscribe("test.fail", h2, events.PriorityLow)

	evt := createTestEvent(t, "test.fail")
	err := bus.Publish(evt)
	if err == nil {
		t.Fatalf("expected error from aggregate publish failure")
	}

	var aggErr *dispatcher.AggregateError
	if !errors.As(err, &aggErr) {
		t.Fatalf("expected *dispatcher.AggregateError, got %T", err)
	}
	if aggErr.Count() != 2 {
		t.Fatalf("expected 2 errors in aggregate, got %d", aggErr.Count())
	}
}

func TestBusSubscribeValidationFailuresAndDuplicates(t *testing.T) {
	bus := eventbus.New()
	h := &mockHandler{id: "h1"}

	// Invalid pattern
	_, err := bus.Subscribe("INVALID PATTERN", h, events.PriorityNormal)
	if !errors.Is(err, errs.ErrValidationFailed) {
		t.Fatalf("expected ErrValidationFailed for invalid pattern, got %v", err)
	}

	// Nil handler
	_, err = bus.Subscribe("valid.pattern", nil, events.PriorityNormal)
	if !errors.Is(err, errs.ErrValidationFailed) {
		t.Fatalf("expected ErrValidationFailed for nil handler, got %v", err)
	}

	// Empty handler ID
	_, err = bus.Subscribe("valid.pattern", &mockHandler{id: ""}, events.PriorityNormal)
	if !errors.Is(err, errs.ErrValidationFailed) {
		t.Fatalf("expected ErrValidationFailed for empty handler ID, got %v", err)
	}

	// Invalid Priority
	_, err = bus.Subscribe("valid.pattern", h, events.PriorityCritical+1)
	if !errors.Is(err, errs.ErrInvalidPriority) {
		t.Fatalf("expected ErrInvalidPriority for invalid priority, got %v", err)
	}

	// Duplicate registration
	subID1, err := bus.Subscribe("valid.pattern", h, events.PriorityNormal)
	if err != nil {
		t.Fatalf("unexpected subscribe error: %v", err)
	}
	_ = subID1

	_, err = bus.Subscribe("valid.pattern", h, events.PriorityNormal)
	if !errors.Is(err, errs.ErrDuplicateSubscriber) {
		t.Fatalf("expected ErrDuplicateSubscriber for duplicate registration, got %v", err)
	}
}

func TestBusUnsubscribeFailures(t *testing.T) {
	bus := eventbus.New()

	// Zero ID
	var zeroID subscription.SubscriptionID
	if err := bus.Unsubscribe(zeroID); !errors.Is(err, errs.ErrValidationFailed) {
		t.Fatalf("expected ErrValidationFailed for zero ID, got %v", err)
	}

	// Non-existent ID
	h := &mockHandler{id: "h1"}
	subID, _ := bus.Subscribe("test.topic", h, events.PriorityNormal)
	_ = bus.Unsubscribe(subID)

	// Second unsubscribe on same ID should fail
	if err := bus.Unsubscribe(subID); err == nil {
		t.Fatalf("expected error when unsubscribing already removed subscription")
	}
}
