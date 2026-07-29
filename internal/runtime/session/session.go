package session

import (
	"fmt"
	"github.com/ioriimasu/jervis/internal/runtime/session/contracts"
	"github.com/ioriimasu/jervis/internal/runtime/session/model"
	"github.com/ioriimasu/jervis/internal/runtime/session/registry"
	"github.com/ioriimasu/jervis/internal/runtime/types"
)

type SessionManager struct {
	contracts.Registry
}

func New() *SessionManager {
	return &SessionManager{
		Registry: registry.New(),
	}
}

func (m *SessionManager) CreateSession(id string) (contracts.Session, error) {
	sID, err := types.NewSessionID(id)
	if err != nil {
		return nil, err
	}

	s := model.NewSession(sID)
	if err := m.Register(s); err != nil {
		return nil, err
	}

	// For Phase 1.6, we just mark it as Running immediately
	if mod, ok := s.(*model.Session); ok {
		mod.TransitionTo(types.StateRunning)
	}

	return s, nil
}

func (m *SessionManager) GetSession(id types.SessionID) (contracts.Session, bool) {
	return m.Get(id)
}

func (m *SessionManager) CloseSession(id types.SessionID) error {
	s, ok := m.Get(id)
	if !ok {
		return fmt.Errorf("session not found: %s", id)
	}

	if mod, ok := s.(*model.Session); ok {
		mod.TransitionTo(types.StateStopped)
	}

	return m.Unregister(id)
}
