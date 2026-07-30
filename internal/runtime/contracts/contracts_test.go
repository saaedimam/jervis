package contracts_test

import (
	"testing"

	"github.com/saaedimam/jervis/internal/runtime/contracts"
)

type mockComponent struct {
	frozen bool
	state  string
}

func (m *mockComponent) Validate() error {
	return nil
}

func (m *mockComponent) IsFrozen() bool {
	return m.frozen
}

func (m *mockComponent) Version() string {
	return "1.0.0"
}

func (m *mockComponent) GitCommit() string {
	return "abc1234"
}

func (m *mockComponent) BuildDate() string {
	return "2026-01-01"
}

func (m *mockComponent) State() string {
	return m.state
}

func (m *mockComponent) Start() error {
	m.state = "Running"
	return nil
}

func (m *mockComponent) Stop() error {
	m.state = "Stopped"
	return nil
}

func (m *mockComponent) TransitionTo(targetState string) error {
	m.state = targetState
	return nil
}

func (m *mockComponent) Fail(err error) error {
	m.state = "Failed"
	return nil
}

func TestContractsInterfaces(t *testing.T) {
	mock := &mockComponent{state: "Created"}

	var validator contracts.Validator = mock
	if err := validator.Validate(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	var freezable contracts.Freezable = mock
	if freezable.IsFrozen() != false {
		t.Fatalf("expected false, got true")
	}

	var versionProvider contracts.VersionProvider = mock
	if versionProvider.Version() != "1.0.0" || versionProvider.GitCommit() != "abc1234" || versionProvider.BuildDate() != "2026-01-01" {
		t.Fatalf("unexpected version provider values")
	}

	var manager contracts.LifecycleManager = mock
	if manager.State() != "Created" {
		t.Fatalf("expected Created, got %s", manager.State())
	}
	if err := manager.Start(); err != nil || manager.State() != "Running" {
		t.Fatalf("start failed")
	}
	if err := manager.TransitionTo("Stopping"); err != nil || manager.State() != "Stopping" {
		t.Fatalf("transition failed")
	}
	if err := manager.Stop(); err != nil || manager.State() != "Stopped" {
		t.Fatalf("stop failed")
	}
	if err := manager.Fail(nil); err != nil || manager.State() != "Failed" {
		t.Fatalf("fail failed")
	}
}
