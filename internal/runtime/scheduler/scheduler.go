package scheduler

import (
	"context"
	"time"

	"github.com/ioriimasu/jervis/internal/runtime/scheduler/contracts"
	"github.com/ioriimasu/jervis/internal/runtime/scheduler/engine"
	"github.com/ioriimasu/jervis/internal/runtime/scheduler/errors"
	"github.com/ioriimasu/jervis/internal/runtime/scheduler/registry"
)

type Scheduler struct {
	contracts.Registry
	engine    *engine.Engine
	cancel    context.CancelFunc
	isRunning bool
}

func New() *Scheduler {
	reg := registry.New()
	eng := engine.New(reg)
	return &Scheduler{
		Registry: reg,
		engine:   eng,
	}
}

func (s *Scheduler) Tick(now time.Time) error {
	return s.engine.Tick(now)
}

func (s *Scheduler) Start(ctx context.Context) error {
	if s.isRunning {
		return errors.ErrSchedulerRunning
	}

	runCtx, cancel := context.WithCancel(ctx)
	s.cancel = cancel
	s.isRunning = true

	go s.run(runCtx)
	return nil
}

func (s *Scheduler) Stop() error {
	if !s.isRunning {
		return errors.ErrSchedulerStopped
	}
	if s.cancel != nil {
		s.cancel()
	}
	s.isRunning = false
	return nil
}

func (s *Scheduler) run(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			_ = s.Tick(t)
		}
	}
}
