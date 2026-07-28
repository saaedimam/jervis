package dispatcher

import (
	"fmt"
	"strings"
)

// AggregateError represents a collection of errors encountered during event dispatching.
type AggregateError struct {
	errors []error
}

// NewAggregateError constructs an AggregateError initialized with the provided error slice.
func NewAggregateError(errs []error) *AggregateError {
	agg := &AggregateError{
		errors: make([]error, 0, len(errs)),
	}
	for _, err := range errs {
		if err != nil {
			agg.errors = append(agg.errors, err)
		}
	}
	return agg
}

// Add appends a non-nil error to the collection.
func (a *AggregateError) Add(err error) {
	if err != nil {
		if a.errors == nil {
			a.errors = make([]error, 0, 2)
		}
		a.errors = append(a.errors, err)
	}
}

// Error returns a formatted string representation of all collected errors.
func (a *AggregateError) Error() string {
	if len(a.errors) == 0 {
		return "no errors"
	}
	if len(a.errors) == 1 {
		return a.errors[0].Error()
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "dispatch failed with %d error(s): ", len(a.errors))
	for i, err := range a.errors {
		if i > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(err.Error())
	}
	return sb.String()
}

// Errors returns a defensive copy slice of all collected errors.
func (a *AggregateError) Errors() []error {
	if len(a.errors) == 0 {
		return nil
	}
	cp := make([]error, len(a.errors))
	copy(cp, a.errors)
	return cp
}

// HasErrors reports whether any errors have been collected.
func (a *AggregateError) HasErrors() bool {
	return len(a.errors) > 0
}

// Count returns the total number of collected errors.
func (a *AggregateError) Count() int {
	return len(a.errors)
}

// Unwrap returns the underlying error slice for standard library errors.Is / errors.As support.
func (a *AggregateError) Unwrap() []error {
	return a.Errors()
}
