package automation_test

import (
	"context"
	"testing"

	events "github.com/ioriimasu/jervis/internal/runtime/eventbus/contracts"
	"github.com/ioriimasu/jervis/internal/services/automation"
)

type mockPublisher struct{}

func (m *mockPublisher) Publish(event events.Event) error {
	return nil
}

func TestService_StartStop(t *testing.T) {
	svc := automation.NewService(&mockPublisher{})

	if err := svc.Start(context.Background()); err != nil {
		t.Fatalf("unexpected start error: %v", err)
	}

	if err := svc.Stop(); err != nil {
		t.Fatalf("unexpected stop error: %v", err)
	}

	if svc.Registry() == nil {
		t.Fatalf("expected registry to be initialized")
	}
}

type mockHandleAction struct {
	executed bool
}

func (m *mockHandleAction) Execute(ctx context.Context, payload map[string]any) error {
	m.executed = true
	return nil
}

func TestService_HandleEvent(t *testing.T) {
	svc := automation.NewService(&mockPublisher{})

	action := &mockHandleAction{}
	w := automation.Workflow{
		ID:      "test-handle",
		Name:    "Test Handle",
		Actions: []automation.Action{action},
		Trigger: &automation.EventTrigger{EventType: "trigger.event"},
	}

	_ = svc.Registry().Register(w)

	// Non-matching event
	_ = svc.HandleEvent(context.Background(), &mockRuntimeEvent{eventType: "other.event"})
	if action.executed {
		t.Fatalf("expected action NOT to execute on mismatch")
	}

	// Matching event
	_ = svc.HandleEvent(context.Background(), &mockRuntimeEvent{eventType: "trigger.event"})
	if !action.executed {
		t.Fatalf("expected action to execute on match")
	}
}
