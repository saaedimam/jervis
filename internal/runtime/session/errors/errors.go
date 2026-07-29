package errors

import "fmt"

var (
	// ErrSessionAlreadyExists indicates a session with the same ID is already registered.
	ErrSessionAlreadyExists = fmt.Errorf("session already exists: ID must be unique")

	// ErrSessionNotFound indicates a session ID was not found in the registry.
	ErrSessionNotFound = fmt.Errorf("session not found: ID does not exist")

	// ErrInvalidSessionID indicates a session ID format is invalid.
	ErrInvalidSessionID = fmt.Errorf("invalid session ID")
)
