package contracts

import (
	"context"
	"testing"
)

// Dummy implementations for testing.

type dummyHandler struct {
	id string
}

func (d *dummyHandler) ID() string                                    { return d.id }
func (d *dummyHandler) Handle(ctx context.Context, event Event) error { return nil }

type dummyPublisher struct{}

func (d *dummyPublisher) Publish(ctx context.Context, event Event) error { return nil }

type dummySubscriber struct {
	subscriptions map[string][]Handler
}

func (d *dummySubscriber) Subscribe(eventType string, handler Handler, priority int) error {
	if d.subscriptions == nil {
		d.subscriptions = make(map[string][]Handler)
	}
	d.subscriptions[eventType] = append(d.subscriptions[eventType], handler)
	return nil
}
func (d *dummySubscriber) Unsubscribe(eventType string, handlerID string) error { return nil }

type dummyDispatcher struct{}

func (d *dummyDispatcher) Dispatch(ctx context.Context, event Event, handlers []Handler) error {
	return nil
}

type dummyValidator struct{}

func (d *dummyValidator) Validate(event Event) error { return nil }

type dummyMiddleware struct{}

func (d *dummyMiddleware) Execute(ctx context.Context, event Event, next func(context.Context, Event) error) error {
	return next(ctx, event)
}

type dummyEventFilter struct{}

func (d *dummyEventFilter) Matches(eventType string, targetPattern string) bool {
	return eventType == targetPattern
}

// Test that our dummy types satisfy the interfaces.
func TestContractsImplementations(t *testing.T) {
	var _ Handler = &dummyHandler{id: "test"}
	var _ Publisher = &dummyPublisher{}
	var _ Subscriber = &dummySubscriber{}
	var _ Dispatcher = &dummyDispatcher{}
	var _ Validator = &dummyValidator{}
	var _ Middleware = &dummyMiddleware{}
	var _ EventFilter = &dummyEventFilter{}
}

// Test the Event interface methods with a dummy event.
type dummyEvent struct {
	id            EventID
	eventType     string
	source        string
	timestamp     Timestamp
	correlationID string
	causationID   string
	priority      int
	payload       []byte
	metadata      map[string]string
	version       string
}

func (d *dummyEvent) ID() EventID           { return d.id }
func (d *dummyEvent) Type() string          { return d.eventType }
func (d *dummyEvent) Source() string        { return d.source }
func (d *dummyEvent) Timestamp() Timestamp  { return d.timestamp }
func (d *dummyEvent) CorrelationID() string { return d.correlationID }
func (d *dummyEvent) CausationID() string   { return d.causationID }
func (d *dummyEvent) Priority() int         { return d.priority }
func (d *dummyEvent) Payload() any          { return d.payload }
func (d *dummyEvent) Metadata() map[string]string {
	copy := make(map[string]string, len(d.metadata))
	for k, v := range d.metadata {
		copy[k] = v
	}
	return copy
}
func (d *dummyEvent) Version() string { return d.version }

func TestEventInterface(t *testing.T) {
	e := &dummyEvent{
		id:            "evt-123",
		eventType:     "test.event",
		source:        "test.source",
		timestamp:     1234567890,
		correlationID: "corr-1",
		causationID:   "cause-1",
		priority:      5,
		payload:       []byte("payload"),
		metadata:      map[string]string{"key": "value"},
		version:       "2.0.0",
	}
	if e.ID() != "evt-123" {
		t.Errorf("ID() = %v, want evt-123", e.ID())
	}
	if e.Type() != "test.event" {
		t.Errorf("Type() = %v, want test.event", e.Type())
	}
	if e.Source() != "test.source" {
		t.Errorf("Source() = %v, want test.source", e.Source())
	}
	if e.Timestamp() != 1234567890 {
		t.Errorf("Timestamp() = %v, want 1234567890", e.Timestamp())
	}
	if e.CorrelationID() != "corr-1" {
		t.Errorf("CorrelationID() = %v, want corr-1", e.CorrelationID())
	}
	if e.CausationID() != "cause-1" {
		t.Errorf("CausationID() = %v, want cause-1", e.CausationID())
	}
	if e.Priority() != 5 {
		t.Errorf("Priority() = %v, want 5", e.Priority())
	}
	if string(e.Payload().([]byte)) != "payload" {
		t.Errorf("Payload() = %v, want payload", e.Payload())
	}
	if e.Metadata()["key"] != "value" {
		t.Errorf("Metadata() = %v, want map[key:value]", e.Metadata())
	}
	if e.Version() != "2.0.0" {
		t.Errorf("Version() = %v, want 2.0.0", e.Version())
	}
}
