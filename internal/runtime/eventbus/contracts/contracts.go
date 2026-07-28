package contracts

import (
	"context"

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
	Priority() int
	Payload() any
	Metadata() map[string]string
	Version() string
}

// Handler defines the subscriber callback contract for processing events.
type Handler interface {
	ID() string
	Handle(ctx context.Context, event Event) error
}

// Publisher defines the contract for publishing events to the Event Bus.
type Publisher interface {
	Publish(ctx context.Context, event Event) error
}

// Subscriber defines the contract for registering and unregistering event handlers.
type Subscriber interface {
	Subscribe(eventType string, handler Handler, priority int) error
	Unsubscribe(eventType string, handlerID string) error
}

// Dispatcher defines the component responsible for routing events to handlers.
type Dispatcher interface {
	Dispatch(ctx context.Context, event Event, handlers []Handler) error
}

// Validator defines the event envelope validation contract.
type Validator interface {
	Validate(event Event) error
}

// Middleware defines pipeline hooks executed before and after event handling.
type Middleware interface {
	Execute(ctx context.Context, event Event, next func(ctx context.Context, event Event) error) error
}

// EventFilter defines matching logic for topic subscriptions.
type EventFilter interface {
	Matches(eventType string, targetPattern string) bool
}
