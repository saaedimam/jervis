package session

import (
	"testing"

	"github.com/saaedimam/jervis/internal/runtime/session/errors"
	"github.com/saaedimam/jervis/internal/runtime/types"
)

func TestSessionManager(t *testing.T) {
	m := New()

	t.Run("Create and Retrieve Session", func(t *testing.T) {
		id := "test-session"
		s, err := m.CreateSession(id)
		if err != nil {
			t.Fatalf("failed to create session: %v", err)
		}

		if s.ID().String() != id {
			t.Errorf("expected ID %s, got %s", id, s.ID())
		}

		if s.State() != types.StateRunning {
			t.Errorf("expected state Running, got %s", s.State())
		}

		retrieved, ok := m.GetSession(s.ID())
		if !ok {
			t.Fatal("failed to retrieve session")
		}
		if retrieved != s {
			t.Error("retrieved session does not match created one")
		}
	})

	t.Run("Metadata Isolation", func(t *testing.T) {
		s1, _ := m.CreateSession("s1")
		s2, _ := m.CreateSession("s2")

		s1.SetMetadata("key", "val1")
		s2.SetMetadata("key", "val2")

		v1, _ := s1.GetMetadata("key")
		v2, _ := s2.GetMetadata("key")

		if v1 != "val1" || v2 != "val2" {
			t.Errorf("metadata mismatch: s1=%s, s2=%s", v1, v2)
		}

		allMeta := s1.Metadata()
		if len(allMeta) != 1 || allMeta["key"] != "val1" {
			t.Error("Metadata() failed")
		}
	})

	t.Run("Close Session", func(t *testing.T) {
		id, _ := types.NewSessionID("close-me")
		_, _ = m.CreateSession("close-me")

		if err := m.CloseSession(id); err != nil {
			t.Fatalf("failed to close session: %v", err)
		}

		if _, ok := m.GetSession(id); ok {
			t.Error("session still exists after close")
		}

		if err := m.CloseSession(id); err == nil {
			t.Error("expected error closing non-existent session")
		}
	})

	t.Run("Validation and Errors", func(t *testing.T) {
		_, err := m.CreateSession("")
		if err == nil {
			t.Error("expected error for empty session ID")
		}

		_, _ = m.CreateSession("duplicate")
		_, err = m.CreateSession("duplicate")
		if err != errors.ErrSessionAlreadyExists {
			t.Errorf("expected ErrSessionAlreadyExists, got %v", err)
		}
	})
}
