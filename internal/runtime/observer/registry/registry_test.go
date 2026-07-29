package registry_test

import (
	"testing"

	"github.com/ioriimasu/jervis/internal/runtime/observer/errors"
	"github.com/ioriimasu/jervis/internal/runtime/observer/registry"
	"github.com/ioriimasu/jervis/internal/runtime/observer/testutils"
)

func TestRegistry(t *testing.T) {
	reg := registry.New()

	if reg.Count() != 0 {
		t.Errorf("expected count 0, got %d", reg.Count())
	}

	obs1 := &testutils.MockObserver{IDVal: "obs-1"}
	obs2 := &testutils.MockObserver{IDVal: "obs-2"}

	if err := reg.Register(obs1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := reg.Register(obs2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if reg.Count() != 2 {
		t.Errorf("expected count 2, got %d", reg.Count())
	}

	if !reg.Contains("obs-1") || !reg.Contains("obs-2") {
		t.Errorf("missing observers")
	}

	t.Run("duplicate registration", func(t *testing.T) {
		if err := reg.Register(obs1); err != errors.ErrDuplicateObserver {
			t.Errorf("expected ErrDuplicateObserver, got %v", err)
		}
	})

	t.Run("nil observer", func(t *testing.T) {
		if err := reg.Register(nil); err != errors.ErrObserverNotFound {
			t.Errorf("expected ErrObserverNotFound, got %v", err)
		}
	})

	t.Run("defensive copy", func(t *testing.T) {
		obsList := reg.Observers()
		obsList[0] = nil
		if reg.Observers()[0] == nil {
			t.Errorf("Observers() did not return defensive copy")
		}
	})

	t.Run("unregister", func(t *testing.T) {
		if err := reg.Unregister("obs-1"); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if reg.Count() != 1 {
			t.Errorf("expected count 1, got %d", reg.Count())
		}
		if reg.Contains("obs-1") {
			t.Errorf("unregistered observer still present")
		}
	})

	t.Run("unregister non-existent", func(t *testing.T) {
		if err := reg.Unregister("obs-missing"); err != errors.ErrObserverNotFound {
			t.Errorf("expected ErrObserverNotFound, got %v", err)
		}
	})

	t.Run("clear", func(t *testing.T) {
		reg.Clear()
		if reg.Count() != 0 {
			t.Errorf("expected count 0 after clear, got %d", reg.Count())
		}
	})
}
