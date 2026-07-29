package errors

import "errors"

var (
	// ErrPermissionDenied indicates that authorization was denied by policy.
	ErrPermissionDenied = errors.New("permissions: permission denied")

	// ErrValidationFailed indicates that a capability request or policy structure is invalid.
	ErrValidationFailed = errors.New("permissions: validation failed")

	// ErrInvalidSubject indicates that the subject attribute is invalid or empty.
	ErrInvalidSubject = errors.New("permissions: invalid subject")

	// ErrInvalidResource indicates that the resource attribute is invalid or empty.
	ErrInvalidResource = errors.New("permissions: invalid resource")

	// ErrInvalidAction indicates that the action attribute is invalid or empty.
	ErrInvalidAction = errors.New("permissions: invalid action")

	// ErrDuplicatePolicy indicates that a policy with the same ID is already registered.
	ErrDuplicatePolicy = errors.New("permissions: duplicate policy")
)
