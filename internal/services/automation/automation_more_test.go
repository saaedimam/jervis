package automation

import (
	"context"
	"testing"

	events "github.com/ioriimasu/jervis/internal/runtime/eventbus/contracts"
	"github.com/ioriimasu/jervis/internal/runtime/types"
)

type mockAction struct {
	ExecuteFunc func(ctx context.Context, payload map[string]any) error
}

func (m *mockAction) Execute(ctx context.Context, payload map[string]any) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, payload)
	}
	return nil
}

type mockPublisher struct{}

func (p *mockPublisher) Publish(event events.Event) error { return nil }

type mockEvent struct {
	t  string
	id string
}

func (m *mockEvent) ID() types.EventID           { id, _ := types.NewEventID(m.id); return id }
func (m *mockEvent) Type() string                { return m.t }
func (m *mockEvent) Payload() any                { return nil }
func (m *mockEvent) Source() string              { return "" }
func (m *mockEvent) Timestamp() types.Timestamp  { return types.Now() }
func (m *mockEvent) CorrelationID() string       { return "" }
func (m *mockEvent) CausationID() string         { return "" }
func (m *mockEvent) Priority() uint8             { return 0 }
func (m *mockEvent) Metadata() map[string]string { return nil }
func (m *mockEvent) Version() string             { return "1.0" }
func (m *mockEvent) String() string              { return "" }

func TestService_HandleEvent(t *testing.T) {
	service := NewService(&mockPublisher{})
	reg := service.Registry()

	actionExecuted := false
	action := &mockAction{
		ExecuteFunc: func(ctx context.Context, payload map[string]any) error {
			actionExecuted = true
			return nil
		},
	}

	wf := Workflow{
		ID:   "test-wf",
		Name: "Test WF",
		Trigger: &EventTrigger{
			EventType: "test.event",
		},
		Actions: []Action{action},
	}
	_ = reg.Register(wf)

	ctx := context.Background()

	err := service.Start(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Handle matching event
	ev := &mockEvent{t: "test.event", id: "1"}
	if err := service.HandleEvent(ctx, ev); err != nil {
		t.Fatal(err)
	}

	if !actionExecuted {
		t.Error("expected action to be executed")
	}

	_ = service.Stop()
}

func TestRegistry_Get(t *testing.T) {
	reg := NewRegistry()

	wf := Workflow{
		ID:      "test-wf",
		Name:    "Test WF",
		Trigger: &EventTrigger{EventType: "x"},
		Actions: []Action{&mockAction{}},
	}
	_ = reg.Register(wf)

	got, ok := reg.Get("test-wf")
	if !ok {
		t.Error("expected to find test-wf")
	}
	if got.ID != "test-wf" {
		t.Errorf("expected test-wf, got %s", got.ID)
	}

	_, ok = reg.Get("missing")
	if ok {
		t.Error("did not expect to find missing")
	}
}
