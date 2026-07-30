package model

import (
	"fmt"
	"time"

	"github.com/ioriimasu/jervis/internal/runtime/scheduler/contracts"
)

// IntervalSchedule triggers a job at fixed time intervals.
type IntervalSchedule struct {
	Interval time.Duration
}

func NewIntervalSchedule(d time.Duration) (contracts.Schedule, error) {
	if d <= 0 {
		return nil, fmt.Errorf("interval must be positive")
	}
	return &IntervalSchedule{Interval: d}, nil
}

func (s *IntervalSchedule) NextRun(ref time.Time) time.Time {
	return ref.Add(s.Interval)
}

// OnceSchedule triggers a job exactly once at a specific time.
type OnceSchedule struct {
	ExecutionTime time.Time
}

func NewOnceSchedule(t time.Time) contracts.Schedule {
	return &OnceSchedule{ExecutionTime: t}
}

func (s *OnceSchedule) NextRun(now time.Time) time.Time {
	if now.After(s.ExecutionTime) {
		return time.Time{}
	}
	return s.ExecutionTime
}
