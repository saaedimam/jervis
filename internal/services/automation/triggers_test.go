package automation_test

import (
	"testing"

	events "github.com/saaedimam/jervis/internal/runtime/eventbus/contracts"
	"github.com/saaedimam/jervis/internal/runtime/types"
	"github.com/saaedimam/jervis/internal/services/automation"
)

// mockRuntimeEvent implements events.Event for testing
type mockRuntimeEvent struct {
	eventType string
}

func (m *mockRuntimeEvent) ID() types.EventID {
	id, _ := types.NewEventID("test-id")
	return id
}
func (m *mockRuntimeEvent) Type() string                { return m.eventType }
func (m *mockRuntimeEvent) Source() string              { return "test-source" }
func (m *mockRuntimeEvent) Timestamp() types.Timestamp  { return types.Now() }
func (m *mockRuntimeEvent) CorrelationID() string       { return "" }
func (m *mockRuntimeEvent) CausationID() string         { return "" }
func (m *mockRuntimeEvent) Priority() uint8             { return 0 }
func (m *mockRuntimeEvent) Payload() any                { return nil }
func (m *mockRuntimeEvent) Metadata() map[string]string { return nil }
func (m *mockRuntimeEvent) Version() string             { return "1.0" }

func TestEventTrigger_IsTriggered(t *testing.T) {
	trigger := &automation.EventTrigger{
		EventType: "task.created",
	}

	tests := []struct {
		name     string
		event    events.Event
		expected bool
	}{
		{
			name:     "Matching event type",
			event:    &mockRuntimeEvent{eventType: "task.created"},
			expected: true,
		},
		{
			name:     "Non-matching event type",
			event:    &mockRuntimeEvent{eventType: "task.completed"},
			expected: false,
		},
		{
			name:     "Nil event",
			event:    nil,
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := trigger.IsTriggered(tc.event)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestCronTrigger_IsTriggered(t *testing.T) {
	trigger := &automation.CronTrigger{
		Expression: "* * * * *",
	}

	// Currently CronTrigger evaluation is deferred and should always return false
	result := trigger.IsTriggered(&mockRuntimeEvent{})
	if result != false {
		t.Errorf("expected CronTrigger.IsTriggered to return false (deferred), got %v", result)
	}
}
