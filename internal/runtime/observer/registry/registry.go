package registry

import (
	"github.com/ioriimasu/jervis/internal/runtime/observer/contracts"
	"github.com/ioriimasu/jervis/internal/runtime/observer/errors"
)

// registry implements contracts.Registry interface.
type registry struct {
	observers []contracts.Observer
}

// New creates a new in-memory observer registry.
func New() contracts.Registry {
	return &registry{
		observers: make([]contracts.Observer, 0),
	}
}

// Register adds an observer in deterministic FIFO sequence.
func (r *registry) Register(obs contracts.Observer) error {
	if obs == nil {
		return errors.ErrObserverNotFound
	}

	if r.Contains(obs.ID()) {
		return errors.ErrDuplicateObserver
	}

	r.observers = append(r.observers, obs)
	return nil
}

// Unregister removes an observer by its ID.
func (r *registry) Unregister(id string) error {
	for i, obs := range r.observers {
		if obs.ID() == id {
			r.observers = append(r.observers[:i], r.observers[i+1:]...)
			return nil
		}
	}
	return errors.ErrObserverNotFound
}

// Observers returns a defensive copy slice of registered observers sorted by registration sequence.
func (r *registry) Observers() []contracts.Observer {
	cp := make([]contracts.Observer, len(r.observers))
	copy(cp, r.observers)
	return cp
}

// Count returns the number of registered observers.
func (r *registry) Count() int {
	return len(r.observers)
}

// Contains reports whether an observer with the specified ID exists in the registry.
func (r *registry) Contains(id string) bool {
	for _, obs := range r.observers {
		if obs.ID() == id {
			return true
		}
	}
	return false
}

// Clear removes all registered observers.
func (r *registry) Clear() {
	r.observers = make([]contracts.Observer, 0)
}
