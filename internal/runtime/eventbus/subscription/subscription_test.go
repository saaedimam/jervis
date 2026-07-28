package subscription_test

import (
	"errors"
	"testing"

	"github.com/ioriimasu/jervis/internal/runtime/eventbus/contracts"
	errs "github.com/ioriimasu/jervis/internal/runtime/eventbus/errors"
	"github.com/ioriimasu/jervis/internal/runtime/eventbus/events"
	"github.com/ioriimasu/jervis/internal/runtime/eventbus/subscription"
)

type mockHandler struct {
	id string
}

func (m *mockHandler) ID() string                         { return m.id }
func (m *mockHandler) Handle(event contracts.Event) error { return nil }

func TestSubscriptionValidCreation(t *testing.T) {
	h := &mockHandler{id: "h-1"}
	subID := subscription.SubscriptionID("sub-1")

	sub, err := subscription.New(subID, "runtime.started", events.PriorityHigh, h)
	if err != nil {
		t.Fatalf("unexpected error creating subscription: %v", err)
	}

	if sub.ID() != subID || sub.ID().String() != "sub-1" || sub.ID().IsZero() {
		t.Errorf("SubscriptionID accessor mismatch")
	}
	if sub.Pattern() != "runtime.started" {
		t.Errorf("Pattern() = %s, want runtime.started", sub.Pattern())
	}
	if sub.Priority() != events.PriorityHigh {
		t.Errorf("Priority() = %v, want PriorityHigh", sub.Priority())
	}
	if sub.Handler() != h {
		t.Errorf("Handler() mismatch")
	}
	if sub.Seq() != 0 {
		t.Errorf("expected default seq 0, got %d", sub.Seq())
	}

	subWithSeq, err := subscription.NewWithSeq("sub-2", "runtime.*", events.PriorityNormal, h, 42)
	if err != nil {
		t.Fatalf("unexpected error creating sub with seq: %v", err)
	}
	if subWithSeq.Seq() != 42 {
		t.Errorf("expected seq 42, got %d", subWithSeq.Seq())
	}
}

func TestSubscriptionValidationFailures(t *testing.T) {
	h := &mockHandler{id: "h-1"}

	// Empty ID
	_, err := subscription.New("", "runtime.started", events.PriorityNormal, h)
	if !errors.Is(err, errs.ErrValidationFailed) {
		t.Fatalf("expected ErrValidationFailed for empty ID")
	}

	// Nil Handler
	_, err = subscription.New("sub-1", "runtime.started", events.PriorityNormal, nil)
	if !errors.Is(err, errs.ErrValidationFailed) {
		t.Fatalf("expected ErrValidationFailed for nil handler")
	}

	// Empty Handler ID
	_, err = subscription.New("sub-1", "runtime.started", events.PriorityNormal, &mockHandler{id: ""})
	if !errors.Is(err, errs.ErrValidationFailed) {
		t.Fatalf("expected ErrValidationFailed for empty handler ID")
	}

	// Invalid Priority
	_, err = subscription.New("sub-1", "runtime.started", events.PriorityCritical+1, h)
	if !errors.Is(err, errs.ErrInvalidPriority) {
		t.Fatalf("expected ErrInvalidPriority for invalid priority")
	}

	// Empty Pattern
	_, err = subscription.New("sub-1", "", events.PriorityNormal, h)
	if !errors.Is(err, errs.ErrValidationFailed) {
		t.Fatalf("expected ErrValidationFailed for empty pattern")
	}

	// Validation on NewWithSeq
	_, err = subscription.NewWithSeq("", "runtime.started", events.PriorityNormal, h, 1)
	if !errors.Is(err, errs.ErrValidationFailed) {
		t.Fatalf("expected ErrValidationFailed on NewWithSeq with empty ID")
	}
}
