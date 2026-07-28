package errors

import "errors"

var (
	// ErrInvalidState indicates an invalid runtime state or state transition.
	ErrInvalidState = errors.New("invalid runtime state transition")

	// ErrConfiguration indicates invalid or incomplete runtime configuration.
	ErrConfiguration = errors.New("invalid runtime configuration")

	// ErrPermission indicates a capability or resource permission violation.
	ErrPermission = errors.New("runtime permission denied")

	// ErrAlreadyRunning indicates an attempt to start a runtime component that is already running.
	ErrAlreadyRunning = errors.New("runtime is already running")

	// ErrNotRunning indicates an operation was attempted on a component that is not running.
	ErrNotRunning = errors.New("runtime is not running")

	// ErrShutdown indicates the runtime component is shutting down or already shutdown.
	ErrShutdown = errors.New("runtime is shutdown")
)
