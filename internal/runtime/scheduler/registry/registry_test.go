package registry

import (
	"context"
	"github.com/ioriimasu/jervis/internal/runtime/scheduler/contracts"
	"github.com/ioriimasu/jervis/internal/runtime/scheduler/errors"
	"testing"
)

type mockJob struct {
	id string
}

func (m *mockJob) ID() string                       { return m.id }
func (m *mockJob) Name() string                     { return m.id }
func (m *mockJob) Schedule() contracts.Schedule     { return nil }
func (m *mockJob) Handle(ctx context.Context) error { return nil }

func TestRegistry(t *testing.T) {
	r := New()

	t.Run("Register and Unregister", func(t *testing.T) {
		j1 := &mockJob{id: "j1"}
		j2 := &mockJob{id: "j2"}

		_ = r.Register(j1)
		_ = r.Register(j2)
		if r.Count() != 2 {
			t.Errorf("expected 2 jobs, got %d", r.Count())
		}

		if err := r.Register(j1); err != errors.ErrJobAlreadyExists {
			t.Errorf("expected ErrJobAlreadyExists, got %v", err)
		}

		_ = r.Unregister("j1")
		if r.Count() != 1 {
			t.Errorf("expected 1 job, got %d", r.Count())
		}

		if _, exists := r.Get("j1"); exists {
			t.Error("j1 should not exist")
		}

		if err := r.Unregister("j1"); err != errors.ErrJobNotFound {
			t.Errorf("expected ErrJobNotFound, got %v", err)
		}
	})

	t.Run("All and Clear", func(t *testing.T) {
		r.Clear()
		_ = r.Register(&mockJob{id: "a"})
		_ = r.Register(&mockJob{id: "b"})

		all := r.All()
		if len(all) != 2 {
			t.Errorf("expected 2 jobs in All(), got %d", len(all))
		}

		r.Clear()
		if r.Count() != 0 {
			t.Error("expected 0 jobs after Clear")
		}
	})
}
