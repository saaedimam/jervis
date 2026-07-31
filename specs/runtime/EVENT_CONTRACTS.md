# Event Bus Contracts Specification

## 1. Overview
This document defines the formal Go interface contracts for the Jervis Event Bus (`internal/runtime/eventbus/contracts`). All implementations MUST conform strictly to these interface definitions.

---

## 2. Core Interface Definitions

```go
package contracts

import (
	"context"

	"github.com/saaedimam/jervis/internal/runtime/types"
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
```

---

## 3. Contract Rules & Constraints

1. **Dependency Boundaries**: The `contracts` package MUST NOT import any packages outside `internal/runtime/types`, `internal/runtime/errors`, and the standard library (`context`).
2. **Context Propagation**: All execution contracts (`Publish`, `Handle`, `Dispatch`, `Execute`) MUST accept `context.Context` as their first parameter for cancellation and deadline propagation.
3. **No Implementation Leaks**: Interfaces MUST NOT expose private struct fields or vendor-specific payload types.
