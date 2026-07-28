# Dispatcher Contracts Specification

## 1. Overview
This document specifies the frozen interface contracts for the Jervis Event Bus Dispatcher component (`internal/runtime/eventbus/dispatcher`).

These contracts comply with `internal/runtime/eventbus/contracts`.

---

## 2. Core Interfaces

```go
package contracts

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

// Dispatcher defines the component responsible for routing events to handlers.
type Dispatcher interface {
	Dispatch(event Event, handlers []Handler) error
}

// Middleware defines pipeline hooks executed before and after event handling.
type Middleware interface {
	Execute(event Event, next func(event Event) error) error
}
```

---

## 3. Contract Rules & Guarantees

1. **Zero `context.Context`**: Method signatures accept `Event` and `[]Handler` directly without `context.Context` parameters.
2. **Immutable Input**: `event` MUST NOT be modified by `Dispatcher` or `Handler` implementations.
3. **Synchronous Return**: `Dispatch(event, handlers)` MUST NOT return until all matching handlers have executed to completion.
4. **Panic Isolation**: If a `Handler` panics during `Handle(event)`, the `Dispatcher` traps the panic and returns a wrapped error.
