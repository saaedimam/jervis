package model

import (
	"testing"
	"github.com/ioriimasu/jervis/internal/runtime/types"
)

func TestSession(t *testing.T) {
	id, _ := types.NewSessionID("test")
	s := NewSession(id)

	if s.ID() != id {
		t.Error("ID mismatch")
	}

	if s.State() != types.StateCreated {
		t.Errorf("expected state Created, got %s", s.State())
	}

	s.SetMetadata("a", "1")
	if val, ok := s.GetMetadata("a"); !ok || val != "1" {
		t.Error("GetMetadata failed")
	}

	if _, ok := s.GetMetadata("missing"); ok {
		t.Error("expected ok=false for missing key")
	}

	// Internal transition
	if mod, ok := s.(*Session); ok {
		mod.TransitionTo(types.StateRunning)
		if s.State() != types.StateRunning {
			t.Error("state transition failed")
		}
	}
}

func TestSession_MetadataCopy(t *testing.T) {
	id, _ := types.NewSessionID("copy-test")
	s := NewSession(id)
	s.SetMetadata("k", "v")
	
	m := s.Metadata()
	m["k"] = "mutated" // Should not affect s
	
	val, _ := s.GetMetadata("k")
	if val != "v" {
		t.Error("defensive copy failed")
	}
}
