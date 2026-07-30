package errors_test

import (
	"errors"
	"testing"

	errs "github.com/saaedimam/jervis/internal/runtime/eventbus/errors"
)

func TestEventBusErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "ErrDuplicateSubscriber",
			err:      errs.ErrDuplicateSubscriber,
			expected: "eventbus: duplicate subscriber handler registered",
		},
		{
			name:     "ErrHandlerFailure",
			err:      errs.ErrHandlerFailure,
			expected: "eventbus: handler execution failed",
		},
		{
			name:     "ErrInvalidEvent",
			err:      errs.ErrInvalidEvent,
			expected: "eventbus: invalid event envelope",
		},
		{
			name:     "ErrInvalidPriority",
			err:      errs.ErrInvalidPriority,
			expected: "eventbus: invalid priority value",
		},
		{
			name:     "ErrDispatchFailed",
			err:      errs.ErrDispatchFailed,
			expected: "eventbus: event dispatch failed",
		},
		{
			name:     "ErrValidationFailed",
			err:      errs.ErrValidationFailed,
			expected: "eventbus: event validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.expected {
				t.Errorf("expected error string %q, got %q", tt.expected, tt.err.Error())
			}
			if !errors.Is(tt.err, tt.err) {
				t.Errorf("errors.Is failed for %v", tt.err)
			}
		})
	}
}
