package lifecycle

import (
	"fmt"

	"github.com/saaedimam/jervis/internal/runtime/contracts"
	"github.com/saaedimam/jervis/internal/runtime/errors"
	"github.com/saaedimam/jervis/internal/runtime/types"
)

var _ contracts.LifecycleManager = (*Manager)(nil)

// Manager enforces a deterministic synchronous lifecycle state machine.
type Manager struct {
	state types.State
}

// NewManager initializes a new Lifecycle Manager in the Created state.
func NewManager() *Manager {
	return &Manager{
		state: types.StateCreated,
	}
}

// State returns the current lifecycle state string.
func (m *Manager) State() string {
	return m.state.String()
}

// CurrentState returns the current canonical types.State.
func (m *Manager) CurrentState() types.State {
	return m.state
}

// CanTransitionTo determines if transitioning from the current state to target is permitted.
func (m *Manager) CanTransitionTo(target types.State) bool {
	if !target.IsValid() {
		return false
	}
	switch m.state {
	case types.StateCreated:
		return target == types.StateInitializing || target == types.StateFailed
	case types.StateInitializing:
		return target == types.StateRunning || target == types.StateFailed
	case types.StateRunning:
		return target == types.StateStopping || target == types.StateFailed
	case types.StateStopping:
		return target == types.StateStopped || target == types.StateFailed
	case types.StateStopped:
		return target == types.StateInitializing || target == types.StateCreated
	case types.StateFailed:
		return false
	default:
		return false
	}
}

// TransitionTo attempts a state transition to targetState.
func (m *Manager) TransitionTo(targetState string) error {
	target := types.State(targetState)
	if !m.CanTransitionTo(target) {
		return fmt.Errorf("%w: cannot transition from %s to %s", errors.ErrInvalidState, m.state, targetState)
	}
	m.state = target
	return nil
}

// Start transitions the state machine to Running (via Initializing if in Created state).
func (m *Manager) Start() error {
	if m.state == types.StateRunning {
		return errors.ErrAlreadyRunning
	}
	if m.state == types.StateStopped {
		return errors.ErrShutdown
	}
	if m.state == types.StateCreated {
		m.state = types.StateInitializing
	}
	if m.state == types.StateInitializing {
		return m.TransitionTo(string(types.StateRunning))
	}
	return fmt.Errorf("%w: cannot start from state %s", errors.ErrInvalidState, m.state)
}

// Stop transitions the state machine from Running to Stopping then Stopped.
func (m *Manager) Stop() error {
	if m.state == types.StateStopped {
		return errors.ErrNotRunning
	}
	if m.state == types.StateRunning {
		m.state = types.StateStopping
	}
	if m.state == types.StateStopping {
		return m.TransitionTo(string(types.StateStopped))
	}
	return fmt.Errorf("%w: cannot stop from state %s", errors.ErrInvalidState, m.state)
}

// Fail transitions the state machine directly to Failed from any non-failed state.
func (m *Manager) Fail(err error) error {
	if m.state == types.StateFailed {
		return fmt.Errorf("%w: lifecycle is already in Failed state", errors.ErrInvalidState)
	}
	m.state = types.StateFailed
	return nil
}
