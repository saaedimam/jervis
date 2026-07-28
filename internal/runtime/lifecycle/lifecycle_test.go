package lifecycle

import (
	"errors"
	"testing"

	errs "github.com/ioriimasu/jervis/internal/runtime/errors"
	"github.com/ioriimasu/jervis/internal/runtime/types"
)

func TestLifecycleManagerFullCycle(t *testing.T) {
	mgr := NewManager()

	if mgr.State() != string(types.StateCreated) {
		t.Fatalf("expected state Created, got %s", mgr.State())
	}
	if mgr.CurrentState() != types.StateCreated {
		t.Fatalf("expected state Created, got %s", mgr.CurrentState())
	}

	// Start (Created -> Initializing -> Running)
	if err := mgr.Start(); err != nil {
		t.Fatalf("unexpected error starting: %v", err)
	}
	if mgr.CurrentState() != types.StateRunning {
		t.Fatalf("expected state Running, got %s", mgr.State())
	}

	// Double start should return ErrAlreadyRunning
	if err := mgr.Start(); !errors.Is(err, errs.ErrAlreadyRunning) {
		t.Fatalf("expected ErrAlreadyRunning, got %v", err)
	}

	// Stop (Running -> Stopping -> Stopped)
	if err := mgr.Stop(); err != nil {
		t.Fatalf("unexpected error stopping: %v", err)
	}
	if mgr.CurrentState() != types.StateStopped {
		t.Fatalf("expected state Stopped, got %s", mgr.State())
	}

	// Stop when stopped returns ErrNotRunning
	if err := mgr.Stop(); !errors.Is(err, errs.ErrNotRunning) {
		t.Fatalf("expected ErrNotRunning, got %v", err)
	}

	// Start when stopped returns ErrShutdown
	if err := mgr.Start(); !errors.Is(err, errs.ErrShutdown) {
		t.Fatalf("expected ErrShutdown, got %v", err)
	}
}

func TestLifecycleManagerTransitions(t *testing.T) {
	mgr := NewManager()

	// Invalid state target string
	if err := mgr.TransitionTo("InvalidState"); !errors.Is(err, errs.ErrInvalidState) {
		t.Fatalf("expected ErrInvalidState, got %v", err)
	}

	// Invalid direct jump (Created -> Running)
	if err := mgr.TransitionTo(string(types.StateRunning)); !errors.Is(err, errs.ErrInvalidState) {
		t.Fatalf("expected ErrInvalidState, got %v", err)
	}

	// Transition to Initializing
	if err := mgr.TransitionTo(string(types.StateInitializing)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mgr.CurrentState() != types.StateInitializing {
		t.Fatalf("expected state Initializing, got %s", mgr.CurrentState())
	}

	// Cannot stop from Initializing
	if err := mgr.Stop(); !errors.Is(err, errs.ErrInvalidState) {
		t.Fatalf("expected ErrInvalidState when stopping from Initializing, got %v", err)
	}

	// Start from Initializing -> Running
	if err := mgr.Start(); err != nil {
		t.Fatalf("unexpected error starting from Initializing: %v", err)
	}

	// Transition to Stopping directly
	if err := mgr.TransitionTo(string(types.StateStopping)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Stop from Stopping -> Stopped
	if err := mgr.Stop(); err != nil {
		t.Fatalf("unexpected error stopping from Stopping: %v", err)
	}
}

func TestLifecycleManagerFailure(t *testing.T) {
	mgr := NewManager()

	if err := mgr.Fail(errors.New("something went wrong")); err != nil {
		t.Fatalf("unexpected error failing: %v", err)
	}
	if mgr.CurrentState() != types.StateFailed {
		t.Fatalf("expected state Failed, got %s", mgr.State())
	}

	// Cannot fail again if already Failed
	if err := mgr.Fail(errors.New("another error")); !errors.Is(err, errs.ErrInvalidState) {
		t.Fatalf("expected ErrInvalidState when failing an already failed manager, got %v", err)
	}

	// Cannot transition out of Failed
	if err := mgr.TransitionTo(string(types.StateCreated)); !errors.Is(err, errs.ErrInvalidState) {
		t.Fatalf("expected ErrInvalidState transitioning out of Failed, got %v", err)
	}

	// Cannot start from Failed
	if err := mgr.Start(); !errors.Is(err, errs.ErrInvalidState) {
		t.Fatalf("expected ErrInvalidState starting from Failed, got %v", err)
	}

	// Cannot stop from Failed
	if err := mgr.Stop(); !errors.Is(err, errs.ErrInvalidState) {
		t.Fatalf("expected ErrInvalidState stopping from Failed, got %v", err)
	}
}

func TestLifecycleManagerInternalBranches(t *testing.T) {
	mgr := NewManager()

	if mgr.CanTransitionTo(types.State("Unknown")) {
		t.Fatalf("expected CanTransitionTo false for invalid target state")
	}

	// Force invalid state to hit default branch in CanTransitionTo, Start, Stop
	mgr.state = types.State("CorruptedState")

	if mgr.CanTransitionTo(types.StateCreated) {
		t.Fatalf("expected CanTransitionTo false for corrupted state")
	}

	if err := mgr.Start(); !errors.Is(err, errs.ErrInvalidState) {
		t.Fatalf("expected ErrInvalidState starting from corrupted state, got %v", err)
	}

	if err := mgr.Stop(); !errors.Is(err, errs.ErrInvalidState) {
		t.Fatalf("expected ErrInvalidState stopping from corrupted state, got %v", err)
	}

	// Reset via new manager to test Stopped state transitions
	mgr2 := NewManager()
	_ = mgr2.TransitionTo(string(types.StateInitializing))
	_ = mgr2.TransitionTo(string(types.StateRunning))
	_ = mgr2.TransitionTo(string(types.StateStopping))
	_ = mgr2.TransitionTo(string(types.StateStopped))

	if !mgr2.CanTransitionTo(types.StateInitializing) {
		t.Fatalf("expected CanTransitionTo true for Stopped -> Initializing")
	}
	if !mgr2.CanTransitionTo(types.StateCreated) {
		t.Fatalf("expected CanTransitionTo true for Stopped -> Created")
	}
	if mgr2.CanTransitionTo(types.StateRunning) {
		t.Fatalf("expected CanTransitionTo false for Stopped -> Running")
	}
}
