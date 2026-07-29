package scheduler

import (
	"context"
	"testing"
	"time"
	"github.com/ioriimasu/jervis/internal/runtime/scheduler/model"
)

func TestScheduler_Tick(t *testing.T) {
	s := New()
	
	callCount := 0
	handler := func(ctx context.Context) error {
		callCount++
		return nil
	}

	// 1s interval job
	job := model.NewJob("job1", "Test Job", &model.IntervalSchedule{Interval: 1 * time.Second}, handler)
	_ = s.Register(job)

	now := time.Now()
	
	// First tick (at 'now') - should not run yet if NextRun is now + interval
	// Wait, my IntervalSchedule.NextRun(zero) returns now + interval.
	// So first run should be at now + 1s.
	
	_ = s.Tick(now)
	if callCount != 0 {
		t.Errorf("expected 0 calls, got %d", callCount)
	}

	// Second tick at now + 1s - should run
	_ = s.Tick(now.Add(1 * time.Second))
	if callCount != 1 {
		t.Errorf("expected 1 call, got %d", callCount)
	}

	// Third tick at now + 1.5s - should not run again
	_ = s.Tick(now.Add(1500 * time.Millisecond))
	if callCount != 1 {
		t.Errorf("expected 1 call still, got %d", callCount)
	}

	// Fourth tick at now + 2s - should run again
	_ = s.Tick(now.Add(2 * time.Second))
	if callCount != 2 {
		t.Errorf("expected 2 calls, got %d", callCount)
	}
}

func TestScheduler_Cron(t *testing.T) {
	s := New()
	
	callCount := 0
	handler := func(ctx context.Context) error {
		callCount++
		return nil
	}

	// Job to run every hour at minute 0
	sched, _ := model.NewCronSchedule("0 * * * *")
	job := model.NewJob("job-cron", "Cron Job", sched, handler)
	_ = s.Register(job)

	// Reference time: 2026-07-29 05:45:00
	ref := time.Date(2026, 7, 29, 5, 45, 0, 0, time.UTC)
	
	// First tick at 05:45 - should not run (NextRun is 06:00)
	_ = s.Tick(ref)
	if callCount != 0 {
		t.Errorf("expected 0 calls, got %d", callCount)
	}

	// Tick at 06:00 - should run
	_ = s.Tick(time.Date(2026, 7, 29, 6, 0, 0, 0, time.UTC))
	if callCount != 1 {
		t.Errorf("expected 1 call, got %d", callCount)
	}

	// Tick at 06:01 - should not run
	_ = s.Tick(time.Date(2026, 7, 29, 6, 1, 0, 0, time.UTC))
	if callCount != 1 {
		t.Errorf("expected 1 call still, got %d", callCount)
	}

	// Tick at 07:00 - should run again
	_ = s.Tick(time.Date(2026, 7, 29, 7, 0, 0, 0, time.UTC))
	if callCount != 2 {
		t.Errorf("expected 2 calls, got %d", callCount)
	}
}

func TestScheduler_Lifecycle(t *testing.T) {
	s := New()
	ctx := context.Background()

	if err := s.Stop(); err == nil {
		t.Error("expected error stopping non-running scheduler")
	}

	if err := s.Start(ctx); err != nil {
		t.Fatal(err)
	}

	if err := s.Start(ctx); err == nil {
		t.Error("expected error starting already running scheduler")
	}

	if err := s.Stop(); err != nil {
		t.Error(err)
	}
}

func TestScheduler_PanicIsolation(t *testing.T) {
	s := New()
	
	handler := func(ctx context.Context) error {
		panic("boom")
	}

	job := model.NewJob("panic-job", "Panic Job", &model.IntervalSchedule{Interval: 1 * time.Second}, handler)
	_ = s.Register(job)

	// Should not crash
	err := s.Tick(time.Now().Add(1 * time.Second))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
