package automation

import (
	"context"
	"errors"
	"strings"
	"testing"

	events "github.com/saaedimam/jervis/internal/runtime/eventbus/contracts"
	eventspb "github.com/saaedimam/jervis/internal/runtime/eventbus/events"
	"github.com/saaedimam/jervis/internal/runtime/types"
)

// Error variables for test comparison.
var (
	errActionFailure  = errors.New("action failure")
	errPublishFailure = errors.New("publish failure")
)

// failureTestMockPublisher records the event it receives and can be forced to error.
type failureTestMockPublisher struct {
	event events.Event
	err   error
}

func (m *failureTestMockPublisher) Publish(event events.Event) error {
	m.event = event
	return m.err
}

// alwaysTrigger fires for any incoming event.
type alwaysTrigger struct{}

func (alwaysTrigger) IsTriggered(events.Event) bool { return true }

// failingAction always returns an error.
type failingAction struct{}

func (failingAction) Execute(context.Context, map[string]any) error {
	return errActionFailure
}

// helper to create a dummy source event.
func dummySourceEvent(t *testing.T) events.Event {
	id, err := types.NewEventID("dummy-1")
	if err != nil {
		t.Fatalf("event ID creation: %v", err)
	}
	evt, err := eventspb.NewBuilder().
		SetID(id).
		SetType(eventspb.EventType("dummy.event")).
		SetSource("test").
		SetTimestamp(types.Now()).
		SetPayload("payload").
		Build()
	if err != nil {
		t.Fatalf("build dummy event: %v", err)
	}
	return evt
}

func TestAutomationFailureEventPublished(t *testing.T) {
	ctx := context.Background()
	pub := &failureTestMockPublisher{}
	svc := NewService(pub)

	// Register a workflow whose single action always fails.
	if err := svc.Registry().Register(Workflow{
		ID:      "wf-1",
		Name:    "failure-test",
		Actions: []Action{failingAction{}},
		Trigger: alwaysTrigger{},
	}); err != nil {
		t.Fatalf("register workflow: %v", err)
	}

	// Run the handler - it should return the original execution error.
	src := dummySourceEvent(t)
	err := svc.HandleEvent(ctx, src)
	if err == nil {
		t.Fatalf("expected error from failed workflow")
	}
	if !errors.Is(err, errActionFailure) {
		t.Fatalf("expected action failure error, got %v", err)
	}

	// Verify that a failure event was published.
	if pub.event == nil {
		t.Fatalf("no event was published")
	}
	if pub.event.Type() != "automation.failed" {
		t.Fatalf("expected event type automation.failed, got %s", pub.event.Type())
	}
	payload := pub.event.Payload().(map[string]any)
	if payload["workflow_id"] != "wf-1" {
		t.Fatalf("unexpected workflow_id payload: %v", payload["workflow_id"])
	}
	if !strings.Contains(payload["original_error"].(string), "action failure") {
		t.Fatalf("original_error payload missing expected text")
	}
	if payload["trigger_event_id"] != src.ID() {
		t.Fatalf("trigger_event_id does not match source event")
	}
}

// If publishing fails we must see both errors via errors.Join.
func TestAutomationPublishErrorPreserved(t *testing.T) {
	ctx := context.Background()
	pub := &failureTestMockPublisher{err: errPublishFailure}
	svc := NewService(pub)

	if err := svc.Registry().Register(Workflow{
		ID:      "wf-2",
		Name:    "publish-error",
		Actions: []Action{failingAction{}},
		Trigger: alwaysTrigger{},
	}); err != nil {
		t.Fatalf("register workflow: %v", err)
	}

	src := dummySourceEvent(t)
	err := svc.HandleEvent(ctx, src)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, errActionFailure) {
		t.Fatalf("original execution error not found in returned error")
	}
	if !errors.Is(err, errPublishFailure) {
		t.Fatalf("publish error not preserved")
	}
}
