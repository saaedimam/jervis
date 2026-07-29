package registry

import (
	"strings"

	"github.com/ioriimasu/jervis/internal/runtime/observer/contracts"
	obserrors "github.com/ioriimasu/jervis/internal/runtime/observer/errors"
)

type registry struct {
	observers   []contracts.Observer
	observerMap map[string]contracts.Observer
}

// NewRegistry creates a new, un-synchronized FIFO observer registry.
func NewRegistry() contracts.Registry {
	return &registry{
		observers:   make([]contracts.Observer, 0),
		observerMap: make(map[string]contracts.Observer),
	}
}

func (r *registry) Register(obs contracts.Observer) error {
	if obs == nil || strings.TrimSpace(obs.ID()) == "" {
		return obserrors.ErrObserverNotFound
	}
	id := obs.ID()
	if _, exists := r.observerMap[id]; exists {
		return obserrors.ErrDuplicateObserver
	}

	r.observers = append(r.observers, obs)
	r.observerMap[id] = obs
	return nil
}

func (r *registry) Unregister(id string) error {
	if strings.TrimSpace(id) == "" {
		return obserrors.ErrObserverNotFound
	}
	if _, exists := r.observerMap[id]; !exists {
		return obserrors.ErrObserverNotFound
	}

	delete(r.observerMap, id)

	newObservers := make([]contracts.Observer, 0, len(r.observers)-1)
	for _, obs := range r.observers {
		if obs.ID() != id {
			newObservers = append(newObservers, obs)
		}
	}
	r.observers = newObservers
	return nil
}

func (r *registry) Observers() []contracts.Observer {
	if len(r.observers) == 0 {
		return []contracts.Observer{}
	}
	cp := make([]contracts.Observer, len(r.observers))
	copy(cp, r.observers)
	return cp
}

func (r *registry) Count() int {
	return len(r.observers)
}

func (r *registry) Contains(id string) bool {
	if strings.TrimSpace(id) == "" {
		return false
	}
	_, exists := r.observerMap[id]
	return exists
}

func (r *registry) Clear() {
	r.observers = make([]contracts.Observer, 0)
	r.observerMap = make(map[string]contracts.Observer)
}
