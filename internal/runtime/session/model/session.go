package model

import (
	"sync"

	"github.com/ioriimasu/jervis/internal/runtime/session/contracts"
	"github.com/ioriimasu/jervis/internal/runtime/types"
)

type Session struct {
	mu       sync.RWMutex
	id       types.SessionID
	metadata map[string]string
	state    types.State
}

func NewSession(id types.SessionID) contracts.Session {
	return &Session{
		id:       id,
		metadata: make(map[string]string),
		state:    types.StateCreated,
	}
}

func (s *Session) ID() types.SessionID {
	return s.id
}

func (s *Session) Metadata() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	meta := make(map[string]string, len(s.metadata))
	for k, v := range s.metadata {
		meta[k] = v
	}
	return meta
}

func (s *Session) SetMetadata(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.metadata[key] = value
}

func (s *Session) GetMetadata(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.metadata[key]
	return val, ok
}

func (s *Session) State() types.State {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state
}

// TransitionTo internal method to update state
func (s *Session) TransitionTo(state types.State) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state = state
}
