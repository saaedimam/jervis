package registry

import (
	"testing"

	"github.com/ioriimasu/jervis/internal/runtime/session/errors"
	"github.com/ioriimasu/jervis/internal/runtime/types"
)

type mockSession struct {
	id types.SessionID
}

func (m *mockSession) ID() types.SessionID                 { return m.id }
func (m *mockSession) Metadata() map[string]string         { return nil }
func (m *mockSession) SetMetadata(k, v string)             {}
func (m *mockSession) GetMetadata(k string) (string, bool) { return "", false }
func (m *mockSession) State() types.State                  { return types.StateCreated }

func TestRegistry(t *testing.T) {
	r := New()
	id1, _ := types.NewSessionID("s1")
	s1 := &mockSession{id: id1}

	t.Run("Register", func(t *testing.T) {
		if err := r.Register(s1); err != nil {
			t.Fatal(err)
		}
		if err := r.Register(s1); err != errors.ErrSessionAlreadyExists {
			t.Errorf("expected ErrSessionAlreadyExists, got %v", err)
		}
	})

	t.Run("Get and Count", func(t *testing.T) {
		if _, ok := r.Get(id1); !ok {
			t.Error("Get failed")
		}
		if r.Count() != 1 {
			t.Errorf("expected 1, got %d", r.Count())
		}
	})

	t.Run("Unregister", func(t *testing.T) {
		id2, _ := types.NewSessionID("missing")
		if err := r.Unregister(id2); err != errors.ErrSessionNotFound {
			t.Errorf("expected ErrSessionNotFound, got %v", err)
		}
		if err := r.Unregister(id1); err != nil {
			t.Fatal(err)
		}
		if r.Count() != 0 {
			t.Error("expected 0 after unregister")
		}
	})

	t.Run("All and Clear", func(t *testing.T) {
		_ = r.Register(s1)
		id2, _ := types.NewSessionID("s2")
		s2 := &mockSession{id: id2}
		_ = r.Register(s2)

		all := r.All()
		if len(all) != 2 {
			t.Errorf("expected 2, got %d", len(all))
		}

		r.Clear()
		if r.Count() != 0 {
			t.Error("Clear failed")
		}
	})
}

func TestRegistry_UnregisterMultiple(t *testing.T) {
	r := New()
	id1, _ := types.NewSessionID("s1")
	id2, _ := types.NewSessionID("s2")
	_ = r.Register(&mockSession{id: id1})
	_ = r.Register(&mockSession{id: id2})

	_ = r.Unregister(id1)
	if r.Count() != 1 {
		t.Error("unregister failed")
	}
	all := r.All()
	if all[0].ID().String() != "s2" {
		t.Error("ordering or unregister failed")
	}
}
