package contracts_test

import (
	"testing"

	"github.com/ioriimasu/jervis/internal/runtime/eventbus/contracts"
	"github.com/ioriimasu/jervis/internal/runtime/types"
)

type dummyHandler struct {
	id string
}

func (d *dummyHandler) ID() string { return d.id }
func (d *dummyHandler) Handle(event contracts.Event) error {
	return nil
}

type dummyPublisher struct{}

func (d *dummyPublisher) Publish(event contracts.Event) error {
	return nil
}

type dummySubscriber struct{}

func (d *dummySubscriber) Subscribe(eventType string, handler contracts.Handler, priority uint8) error {
	return nil
}

func (d *dummySubscriber) Unsubscribe(eventType, handlerID string) error {
	return nil
}

type dummyDispatcher struct{}

func (d *dummyDispatcher) Dispatch(event contracts.Event, handlers []contracts.Handler) error {
	return nil
}

type dummyValidator struct{}

func (d *dummyValidator) Validate(event contracts.Event) error {
	return nil
}

type dummyMiddleware struct{}

func (d *dummyMiddleware) Execute(event contracts.Event, next func(contracts.Event) error) error {
	return next(event)
}

type dummyEventFilter struct{}

func (d *dummyEventFilter) Matches(eventType, targetPattern string) bool {
	return eventType == targetPattern
}

type dummyEvent struct {
	id            types.EventID
	eventType     string
	source        string
	timestamp     types.Timestamp
	correlationID string
	causationID   string
	priority      uint8
	payload       any
	metadata      map[string]string
	version       string
}

func (d *dummyEvent) ID() types.EventID           { return d.id }
func (d *dummyEvent) Type() string                { return d.eventType }
func (d *dummyEvent) Source() string              { return d.source }
func (d *dummyEvent) Timestamp() types.Timestamp  { return d.timestamp }
func (d *dummyEvent) CorrelationID() string       { return d.correlationID }
func (d *dummyEvent) CausationID() string         { return d.causationID }
func (d *dummyEvent) Priority() uint8             { return d.priority }
func (d *dummyEvent) Payload() any                { return d.payload }
func (d *dummyEvent) Metadata() map[string]string { return d.metadata }
func (d *dummyEvent) Version() string             { return d.version }

func TestContractsInterfaces(t *testing.T) {
	h := &dummyHandler{id: "h1"}
	var _ contracts.Handler = h
	if h.ID() != "h1" {
		t.Errorf("ID() = %s, want h1", h.ID())
	}
	if err := h.Handle(nil); err != nil {
		t.Errorf("Handle() err = %v", err)
	}

	var _ contracts.Publisher = &dummyPublisher{}
	var _ contracts.Subscriber = &dummySubscriber{}
	var _ contracts.Dispatcher = &dummyDispatcher{}
	var _ contracts.Validator = &dummyValidator{}
	var _ contracts.Middleware = &dummyMiddleware{}
	var _ contracts.EventFilter = &dummyEventFilter{}

	filter := &dummyEventFilter{}
	if !filter.Matches("a.b", "a.b") {
		t.Errorf("Matches failed")
	}

	evtID, _ := types.NewEventID("evt-1")
	now := types.Now()
	evt := &dummyEvent{
		id:            evtID,
		eventType:     "test.event",
		source:        "test.source",
		timestamp:     now,
		correlationID: "corr-1",
		causationID:   "cause-1",
		priority:      1,
		payload:       "data",
		metadata:      map[string]string{"k": "v"},
		version:       "1.0.0",
	}

	if evt.ID() != evtID || evt.Type() != "test.event" || evt.Source() != "test.source" ||
		evt.Timestamp() != now || evt.CorrelationID() != "corr-1" || evt.CausationID() != "cause-1" ||
		evt.Priority() != 1 || evt.Payload() != "data" || evt.Metadata()["k"] != "v" || evt.Version() != "1.0.0" {
		t.Errorf("dummyEvent accessor mismatch")
	}
}
