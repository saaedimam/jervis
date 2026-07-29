package errors

import "fmt"

var (
	// ErrJobAlreadyExists indicates a job with the same ID is already registered.
	ErrJobAlreadyExists = fmt.Errorf("job already exists: ID must be unique")

	// ErrJobNotFound indicates a job ID was not found in the registry.
	ErrJobNotFound = fmt.Errorf("job not found: ID does not exist")

	// ErrInvalidSchedule indicates a schedule definition is invalid (e.g. negative duration).
	ErrInvalidSchedule = fmt.Errorf("invalid schedule: configuration is mathematically impossible or out of bounds")

	// ErrSchedulerRunning indicates an operation cannot be performed because the scheduler is already running.
	ErrSchedulerRunning = fmt.Errorf("scheduler already running")

	// ErrSchedulerStopped indicates an operation cannot be performed because the scheduler is stopped.
	ErrSchedulerStopped = fmt.Errorf("scheduler is stopped")
)
