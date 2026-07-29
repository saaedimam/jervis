package registry_test

import (
	"fmt"
	"testing"

	"github.com/ioriimasu/jervis/internal/runtime/observer/contracts"
	obserrors "github.com/ioriimasu/jervis/internal/runtime/observer/errors"
	obsregistry "github.com/ioriimasu/jervis/internal/runtime/observer/registry"
)

type mockObserver struct {
	id string
}

func (m *mockObserver) ID() string {
	return m.id
}

func (m *mockObserver) Handle(n contracts.Notification) error {
	return nil
}

func TestRegistry(t *testing.T) {
	reg := obsregistry.NewRegistry()

	if reg.Count() != 0 {
		t.Errorf("Expected count 0, got %d", reg.Count())
	}

	obs1 := &mockObserver{id: "obs-1"}
	obs2 := &mockObserver{id: "obs-2"}

	err := reg.Register(obs1)
	if err != nil {
		t.Fatalf("Unexpected register error: %v", err)
	}

	err = reg.Register(obs2)
	if err != nil {
		t.Fatalf("Unexpected register error: %v", err)
	}

	if reg.Count() != 2 {
		t.Errorf("Expected count 2, got %d", reg.Count())
	}

	if !reg.Contains("obs-1") || !reg.Contains("obs-2") {
		t.Error("Registry missing registered observers")
	}

	if reg.Contains("obs-missing") {
		t.Error("Registry reported containing non-existent observer")
	}

	// Duplicate registration error
	err = reg.Register(obs1)
	if err != obserrors.ErrDuplicateObserver {
		t.Errorf("Expected ErrDuplicateObserver, got %v", err)
	}

	// Nil observer error
	err = reg.Register(nil)
	if err != obserrors.ErrObserverNotFound {
		t.Errorf("Expected ErrObserverNotFound for nil observer, got %v", err)
	}

	// Unregister
	err = reg.Unregister("obs-1")
	if err != nil {
		t.Fatalf("Unexpected unregister error: %v", err)
	}

	if reg.Count() != 1 {
		t.Errorf("Expected count 1, got %d", reg.Count())
	}

	if reg.Contains("obs-1") {
		t.Error("Unregistered observer still present")
	}

	// Unregister non-existent
	err = reg.Unregister("obs-99")
	if err != obserrors.ErrObserverNotFound {
		t.Errorf("Expected ErrObserverNotFound, got %v", err)
	}

	// Defensive copy test
	obsList := reg.Observers()
	if len(obsList) != 1 {
		t.Fatalf("Expected 1 observer, got %d", len(obsList))
	}

	obsList[0] = nil
	if reg.Observers()[0] == nil {
		t.Error("Registry Observers() did not return defensive copy")
	}

	// Clear
	reg.Clear()
	if reg.Count() != 0 {
		t.Errorf("Expected count 0 after clear, got %d", reg.Count())
	}
}

func TestRegistryFIFOOrder(t *testing.T) {
	reg := obsregistry.NewRegistry()

	for i := 0; i < 5; i++ {
		_ = reg.Register(&mockObserver{id: fmt.Sprintf("obs-%d", i)})
	}

	obsList := reg.Observers()
	for i := 0; i < 5; i++ {
		expectedID := fmt.Sprintf("obs-%d", i)
		if obsList[i].ID() != expectedID {
			t.Errorf("Expected FIFO index %d to be %s, got %s", i, expectedID, obsList[i].ID())
		}
	}
}
