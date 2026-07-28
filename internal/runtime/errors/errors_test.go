package errors_test

import (
	"errors"
	"testing"

	errs "github.com/ioriimasu/jervis/internal/runtime/errors"
)

func TestCanonicalErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "ErrInvalidState",
			err:      errs.ErrInvalidState,
			expected: "invalid runtime state transition",
		},
		{
			name:     "ErrConfiguration",
			err:      errs.ErrConfiguration,
			expected: "invalid runtime configuration",
		},
		{
			name:     "ErrPermission",
			err:      errs.ErrPermission,
			expected: "runtime permission denied",
		},
		{
			name:     "ErrAlreadyRunning",
			err:      errs.ErrAlreadyRunning,
			expected: "runtime is already running",
		},
		{
			name:     "ErrNotRunning",
			err:      errs.ErrNotRunning,
			expected: "runtime is not running",
		},
		{
			name:     "ErrShutdown",
			err:      errs.ErrShutdown,
			expected: "runtime is shutdown",
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
