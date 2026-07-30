package engine

import (
	"context"
	"testing"
	"time"

	"github.com/saaedimam/jervis/internal/runtime/scheduler/model"
	"github.com/saaedimam/jervis/internal/runtime/scheduler/registry"
)

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
