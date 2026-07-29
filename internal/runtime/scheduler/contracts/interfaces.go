package contracts

import (
	"context"
	"time"
)

// Job represents a scheduled background task.
type Job interface {
	// ID returns the unique identifier for the job.
	ID() string

	// Name returns a human-readable name for the job.
	Name() string

	// Schedule returns the job's schedule.
	Schedule() Schedule

	// Handle executes the job logic.
	Handle(ctx context.Context) error
}

// Schedule defines the interface for determining when a job should run.
type Schedule interface {
	// NextRun returns the next execution time based on the provided reference time.
	// Returns a zero time if the job should never run again.
	NextRun(now time.Time) time.Time
}

// Registry manages in-memory storage and lookup of scheduled jobs.
type Registry interface {
	// Register adds a job to the registry.
	Register(job Job) error

	// Unregister removes a job by its ID.
	Unregister(id string) error

	// Get retrieves a job by its ID.
	Get(id string) (Job, bool)

	// All returns a defensive copy slice of all registered jobs.
	All() []Job

	// Count returns the number of registered jobs.
	Count() int

	// Clear removes all registered jobs.
	Clear()
}

// Scheduler owns the execution engine for triggered jobs.
type Scheduler interface {
	// Start begins the background execution of the scheduler.
	Start(ctx context.Context) error

	// Stop gracefully shuts down the scheduler.
	Stop() error

	// Tick manually triggers a check of all jobs for execution at the provided time.
	// Useful for deterministic testing and manual engine driving.
	Tick(now time.Time) error
}
