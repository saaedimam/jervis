package engine

import (
	"context"
	"sync"
	"time"

	"github.com/ioriimasu/jervis/internal/runtime/scheduler/contracts"
)

type Engine struct {
	registry contracts.Registry
	lastRuns map[string]time.Time
	mu       sync.Mutex
}

func New(registry contracts.Registry) *Engine {
	return &Engine{
		registry: registry,
		lastRuns: make(map[string]time.Time),
	}
}

// Tick checks all registered jobs and executes those that are due.
func (e *Engine) Tick(now time.Time) error {
	jobs := e.registry.All()

	for _, job := range jobs {
		e.mu.Lock()
		lastRun, ok := e.lastRuns[job.ID()]
		if !ok {
			e.lastRuns[job.ID()] = now
			lastRun = now
		}
		e.mu.Unlock()

		nextRun := job.Schedule().NextRun(lastRun)
		if !nextRun.IsZero() && (now.After(nextRun) || now.Equal(nextRun)) {
			// Job is due
			e.safeHandle(job, now)
		}
	}

	return nil
}

func (e *Engine) safeHandle(job contracts.Job, now time.Time) {
	defer func() {
		if r := recover(); r != nil {
			_ = r
		}

		e.mu.Lock()
		e.lastRuns[job.ID()] = now
		e.mu.Unlock()
	}()

	ctx := context.Background()
	_ = job.Handle(ctx)
}
