package contracts

import (
	"github.com/ioriimasu/jervis/internal/runtime/eventbus/contracts"
	"github.com/ioriimasu/jervis/internal/runtime/types"
)

// Notification defines the immutable notification interface wrapping an Event.
type Notification interface {
	// Event returns the wrapped canonical Event interface.
	Event() contracts.Event

	// ObservedAt returns the observation timestamp.
	ObservedAt() types.Timestamp
}

// Observer defines the interface for components that passively observe runtime notifications.
type Observer interface {
	// ID returns a unique identifier for the observer instance.
	ID() string

	// Handle receives an immutable Notification for observation.
	// Returning an error records a failure in AggregateError but DOES NOT halt remaining observers.
	Handle(n Notification) error
}

// Observable defines the interface for components that dispatch notifications to observers.
type Observable interface {
	// Notify dispatches a notification synchronously to all registered observers.
	Notify(n Notification) error
}

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

// Dispatcher owns the synchronous, panic-protected execution of registered observers.
type Dispatcher interface {
	// Dispatch executes all registered observers sequentially in FIFO order.
	Dispatch(n Notification) error
}