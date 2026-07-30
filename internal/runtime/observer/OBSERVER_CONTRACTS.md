# Observer Contracts (Phase 1.4.0 - FROZEN)

This document defines the frozen interface contracts and data models for the Runtime Observer subsystem.

---

## 1. Notification Contract

```go
package contracts

import (
	eventcontracts "github.com/saaedimam/jervis/internal/runtime/eventbus/contracts"
	"github.com/saaedimam/jervis/internal/runtime/types"
)

// Notification defines the immutable notification interface wrapping an Event.
type Notification interface {
	// Event returns the wrapped canonical Event interface.
	Event() eventcontracts.Event

	// ObservedAt returns the observation timestamp.
	ObservedAt() types.Timestamp
}
```

---

## 2. Observer Contract

```go
package contracts

// Observer defines the interface for components that passively observe runtime notifications.
type Observer interface {
	// ID returns a unique identifier for the observer instance.
	ID() string

	// Handle receives an immutable Notification for observation.
	// Returning an error records a failure in AggregateError but DOES NOT halt remaining observers.
	Handle(n Notification) error
}
```

---

## 3. Observable Contract

```go
package contracts

// Observable defines the interface for components that dispatch notifications to observers.
type Observable interface {
	// Notify dispatches a notification synchronously to all registered observers.
	Notify(n Notification) error
}
```

---

## 4. Registry Contract

```go
package contracts

// Registry manages in-memory storage, lookup, and unregistration of observers.
type Registry interface {
	// Register adds an observer in deterministic FIFO sequence.
	Register(obs Observer) error

	// Unregister removes an observer by its ID.
	Unregister(id string) error

	// Observers returns a defensive copy slice of registered observers sorted deterministically by registration sequence (FIFO).
	Observers() []Observer

	// Count returns the number of registered observers.
	Count() int

	// Contains reports whether an observer with the specified ID exists in the registry.
	Contains(id string) bool

	// Clear removes all registered observers.
	Clear()
}
```

---

## 5. Dispatcher Contract

```go
package contracts

// Dispatcher owns the synchronous, panic-protected execution of registered observers.
type Dispatcher interface {
	// Dispatch executes all registered observers sequentially in FIFO order.
	Dispatch(n Notification) error
}
```

---

## 6. Design & Concurrency Invariants

1. **No Mutexes, Channels, or Goroutines**: All implementations are 100% synchronous and single-threaded.
2. **Defensive Copies**: `Registry.Observers()` returns a defensive copy slice.
3. **Panic Protection & Isolation**: Dispatcher wraps each `Observer.Handle()` call in `recover()` and compiles errors into `AggregateError`.
4. **Read-Only Invariant**: Observers MUST NOT modify the notification or perform state-altering runtime calls.