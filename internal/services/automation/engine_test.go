package automation_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/saaedimam/jervis/internal/services/automation"
)

type recordedAction struct {
	executed bool
	delay    time.Duration
	err      error
}

func (r *recordedAction) Execute(ctx context.Context, payload map[string]any) error {
	if r.delay > 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(r.delay):
		}
	}
	r.executed = true
	return r.err
}

func TestEngine_Execute_Sequential(t *testing.T) {
	engine := automation.NewEngine()

	action1 := &recordedAction{}
	action2 := &recordedAction{}

	workflow := automation.Workflow{
		ID:      "test-flow",
		Actions: []automation.Action{action1, action2},
	}

	err := engine.Execute(context.Background(), workflow, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !action1.executed || !action2.executed {
		t.Fatalf("expected both actions to execute")
	}
}

func TestEngine_Execute_ErrorPropagation(t *testing.T) {
	engine := automation.NewEngine()

	expectedErr := errors.New("action failed")
	action1 := &recordedAction{err: expectedErr}
	action2 := &recordedAction{}

	workflow := automation.Workflow{
		ID:      "test-flow",
		Actions: []automation.Action{action1, action2},
	}

	err := engine.Execute(context.Background(), workflow, nil)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}

	if !action1.executed {
		t.Fatalf("expected action1 to execute")
	}
	if action2.executed {
		t.Fatalf("expected action2 NOT to execute after error")
	}
}

func TestEngine_Execute_Cancellation(t *testing.T) {
	engine := automation.NewEngine()

	action1 := &recordedAction{delay: 10 * time.Millisecond}
	action2 := &recordedAction{}

	workflow := automation.Workflow{
		ID:      "test-flow",
		Actions: []automation.Action{action1, action2},
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Cancel immediately
	cancel()

	err := engine.Execute(ctx, workflow, nil)
	if err == nil {
		t.Fatalf("expected cancellation error")
	}

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled error, got %v", err)
	}

	if action2.executed {
		t.Fatalf("expected action2 NOT to execute due to cancellation")
	}
}
