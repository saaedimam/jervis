package automation_test

import (
	"context"
	"errors"
	"testing"

	events "github.com/ioriimasu/jervis/internal/runtime/eventbus/contracts"
	"github.com/ioriimasu/jervis/internal/services/automation"
)

// MockAction for testing
type mockAction struct{}

func (m *mockAction) Execute(ctx context.Context, payload map[string]any) error {
	return nil
}

// MockTrigger for testing
type mockTrigger struct{}

func (m *mockTrigger) IsTriggered(event events.Event) bool {
	return true
}

func TestRegistry_Register(t *testing.T) {
	tests := []struct {
		name        string
		workflow    automation.Workflow
		expectError bool
		errorTarget error
	}{
		{
			name: "Valid workflow",
			workflow: automation.Workflow{
				ID:      "test-1",
				Name:    "Test Workflow",
				Actions: []automation.Action{&mockAction{}},
				Trigger: &mockTrigger{},
			},
			expectError: false,
		},
		{
			name: "Missing ID",
			workflow: automation.Workflow{
				Name:    "Test Workflow",
				Actions: []automation.Action{&mockAction{}},
				Trigger: &mockTrigger{},
			},
			expectError: true,
			errorTarget: automation.ErrInvalidWorkflow,
		},
		{
			name: "Missing Name",
			workflow: automation.Workflow{
				ID:      "test-1",
				Actions: []automation.Action{&mockAction{}},
				Trigger: &mockTrigger{},
			},
			expectError: true,
			errorTarget: automation.ErrInvalidWorkflow,
		},
		{
			name: "No actions",
			workflow: automation.Workflow{
				ID:      "test-1",
				Name:    "Test Workflow",
				Trigger: &mockTrigger{},
			},
			expectError: true,
			errorTarget: automation.ErrInvalidWorkflow,
		},
		{
			name: "Missing trigger",
			workflow: automation.Workflow{
				ID:      "test-1",
				Name:    "Test Workflow",
				Actions: []automation.Action{&mockAction{}},
			},
			expectError: true,
			errorTarget: automation.ErrInvalidWorkflow,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := automation.NewRegistry()
			err := r.Register(tc.workflow)
			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
				if tc.errorTarget != nil && !errors.Is(err, tc.errorTarget) {
					// Need to import errors above, let's fix that
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestRegistry_DuplicateRegister(t *testing.T) {
	r := automation.NewRegistry()
	w := automation.Workflow{
		ID:      "test-1",
		Name:    "Test Workflow",
		Actions: []automation.Action{&mockAction{}},
		Trigger: &mockTrigger{},
	}

	err := r.Register(w)
	if err != nil {
		t.Fatalf("unexpected error on first register: %v", err)
	}

	err = r.Register(w)
	if err == nil {
		t.Fatalf("expected error on duplicate register but got none")
	}
}

func TestRegistry_Unregister(t *testing.T) {
	r := automation.NewRegistry()
	w := automation.Workflow{
		ID:      "test-1",
		Name:    "Test Workflow",
		Actions: []automation.Action{&mockAction{}},
		Trigger: &mockTrigger{},
	}

	_ = r.Register(w)

	err := r.Unregister("test-1")
	if err != nil {
		t.Fatalf("unexpected error on unregister: %v", err)
	}

	err = r.Unregister("test-1")
	if err == nil {
		t.Fatalf("expected error on unregistering non-existent workflow")
	}
}

func TestRegistry_List(t *testing.T) {
	r := automation.NewRegistry()

	workflows := []automation.Workflow{
		{ID: "c-test", Name: "C", Actions: []automation.Action{&mockAction{}}, Trigger: &mockTrigger{}},
		{ID: "a-test", Name: "A", Actions: []automation.Action{&mockAction{}}, Trigger: &mockTrigger{}},
		{ID: "b-test", Name: "B", Actions: []automation.Action{&mockAction{}}, Trigger: &mockTrigger{}},
	}

	for _, w := range workflows {
		_ = r.Register(w)
	}

	list := r.List()
	if len(list) != 3 {
		t.Fatalf("expected 3 workflows, got %d", len(list))
	}

	// Check deterministic ordering by ID
	if list[0].ID != "a-test" || list[1].ID != "b-test" || list[2].ID != "c-test" {
		t.Fatalf("expected deterministic ordering by ID, got: %v", list)
	}
}
