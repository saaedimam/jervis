package registry

import (
	"github.com/ioriimasu/jervis/internal/runtime/session/contracts"
	"github.com/ioriimasu/jervis/internal/runtime/session/errors"
	"github.com/ioriimasu/jervis/internal/runtime/types"
	"sync"
)

type Registry struct {
	mu       sync.RWMutex
	sessions map[string]contracts.Session
	order    []types.SessionID
}

func New() contracts.Registry {
	return &Registry{
		sessions: make(map[string]contracts.Session),
		order:    make([]types.SessionID, 0),
	}
}

func (r *Registry) Register(s contracts.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := s.ID()
	if _, exists := r.sessions[id.String()]; exists {
		return errors.ErrSessionAlreadyExists
	}

	r.sessions[id.String()] = s
	r.order = append(r.order, id)
	return nil
}

func (r *Registry) Unregister(id types.SessionID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.sessions[id.String()]; !exists {
		return errors.ErrSessionNotFound
	}

	delete(r.sessions, id.String())

	// Remove from order slice
	newOrder := make([]types.SessionID, 0, len(r.order)-1)
	for _, oid := range r.order {
		if oid.String() != id.String() {
			newOrder = append(newOrder, oid)
		}
	}
	r.order = newOrder

	return nil
}

func (r *Registry) Get(id types.SessionID) (contracts.Session, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.sessions[id.String()]
	return s, ok
}

func (r *Registry) All() []contracts.Session {
	r.mu.RLock()
	defer r.mu.RUnlock()

	sessions := make([]contracts.Session, len(r.order))
	for i, id := range r.order {
		sessions[i] = r.sessions[id.String()]
	}
	return sessions
}

func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.sessions)
}

func (r *Registry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions = make(map[string]contracts.Session)
	r.order = make([]types.SessionID, 0)
}
