package contracts

import (
	"github.com/ioriimasu/jervis/internal/runtime/types"
)

// Event represents the canonical immutable event envelope interface.
type Event interface {
	ID() types.EventID
	Type() string
	Source() string
	Timestamp() types.Timestamp
	CorrelationID() string
	CausationID() string
	Priority() uint8
	Payload() any
	Metadata() map[string]string
	Version() string
}

// Handler defines the subscriber callback contract for processing events.
type Handler interface {
	ID() string
	Handle(event Event) error
}

// Publisher defines the contract for publishing events to the Event Bus.
type Publisher interface {
	Publish(event Event) error
}

// Subscriber defines the contract for registering and unregistering event handlers.
type Subscriber interface {
	Subscribe(eventType string, handler Handler, priority uint8) error
	Unsubscribe(eventType, handlerID string) error
}

// Dispatcher defines the component responsible for routing events to handlers.
type Dispatcher interface {
	Dispatch(event Event, handlers []Handler) error
}

// Validator defines the event envelope validation contract.
type Validator interface {
	Validate(event Event) error
}

// Middleware defines pipeline hooks executed before and after event handling.
type Middleware interface {
	Execute(event Event, next func(event Event) error) error
}

// EventFilter defines matching logic for topic subscriptions.
type EventFilter interface {
	Matches(eventType, targetPattern string) bool
}
