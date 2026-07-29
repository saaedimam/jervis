package engine

import (
	"context"
	"testing"
	"time"
	"github.com/ioriimasu/jervis/internal/runtime/scheduler/contracts"
	"github.com/ioriimasu/jervis/internal/runtime/scheduler/model"
	"github.com/ioriimasu/jervis/internal/runtime/scheduler/registry"
)

type mockJob struct {
	id      string
	handleFn func(ctx context.Context) error
}
func (m *mockJob) ID() string { return m.id }
func (m *mockJob) Name() string { return m.id }
func (m *mockJob) Schedule() contracts.Schedule { return nil }
func (m *mockJob) Handle(ctx context.Context) error { return m.handleFn(ctx) }

func TestEngine_Tick(t *testing.T) {
	reg := registry.New()
	eng := New(reg)

	t.Run("Job handle returning error should continue", func(t *testing.T) {
		job := model.NewJob("err-job", "Error Job", &model.IntervalSchedule{Interval: 1 * time.Second}, func(ctx context.Context) error {
			return context.DeadlineExceeded
		})
		_ = reg.Register(job)

		err := eng.Tick(time.Now().Add(1 * time.Second))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
