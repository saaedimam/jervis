package errors

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrInvalidNotification indicates an invalid or zero-value notification was provided.
	ErrInvalidNotification = errors.New("observer: invalid notification")
	// ErrDuplicateObserver indicates an observer with the same ID already exists in the registry.
	ErrDuplicateObserver = errors.New("observer: duplicate observer ID")
	// ErrObserverNotFound indicates the specified observer ID was not found in the registry.
	ErrObserverNotFound = errors.New("observer: observer not found")
	// ErrObserverFailure indicates handler execution returned an error.
	ErrObserverFailure = errors.New("observer: handler execution failed")
	// ErrDispatchFailed indicates dispatch execution failed.
	ErrDispatchFailed = errors.New("observer: dispatch failed")
	// ErrObserverPanic indicates an observer handler panicked during execution.
	ErrObserverPanic = errors.New("observer: handler execution panicked")
)

// AggregateError collects multiple observer errors and panics.
type AggregateError struct {
	errors []error
}

// Errors returns a defensive copy of the slice of errors.
func (a *AggregateError) Errors() []error {
	if a == nil || len(a.errors) == 0 {
		return nil
	}
	cp := make([]error, len(a.errors))
	copy(cp, a.errors)
	return cp
}

// Error returns a formatted multiline representation of aggregated errors.
func (a *AggregateError) Error() string {
	if a == nil || len(a.errors) == 0 {
		return "no observer errors"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d observer error(s) occurred:", len(a.errors)))
	for _, err := range a.errors {
		sb.WriteString("\n - ")
		sb.WriteString(err.Error())
	}
	return sb.String()
}

// NewAggregateError constructs an AggregateError from a slice of errors.
func NewAggregateError(errs []error) *AggregateError {
	var validErrs []error
	for _, err := range errs {
		if err != nil {
			validErrs = append(validErrs, err)
		}
	}
	if len(validErrs) == 0 {
		return nil
	}
	return &AggregateError{errors: validErrs}
}