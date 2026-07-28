package errors

import "errors"

var (
	// ErrDuplicateSubscriber is returned when a handler is already registered for an event type.
	ErrDuplicateSubscriber = errors.New("eventbus: duplicate subscriber handler registered")

	// ErrHandlerFailure is returned when an event handler fails during execution.
	ErrHandlerFailure = errors.New("eventbus: handler execution failed")

	// ErrInvalidEvent is returned when an event envelope fails structural validation.
	ErrInvalidEvent = errors.New("eventbus: invalid event envelope")

	// ErrInvalidPriority is returned when an event or subscription priority is out of allowed bounds.
	ErrInvalidPriority = errors.New("eventbus: invalid priority value")

	// ErrDispatchFailed is returned when dispatching an event to subscribers fails.
	ErrDispatchFailed = errors.New("eventbus: event dispatch failed")

	// ErrValidationFailed is returned when event content or envelope validation fails.
	ErrValidationFailed = errors.New("eventbus: event validation failed")
)
