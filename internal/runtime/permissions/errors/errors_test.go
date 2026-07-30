package errors_test

import (
	"testing"

	pErrors "github.com/saaedimam/jervis/internal/runtime/permissions/errors"
)

func TestPermissionErrors(t *testing.T) {
	errList := []error{
		pErrors.ErrPermissionDenied,
		pErrors.ErrValidationFailed,
		pErrors.ErrInvalidSubject,
		pErrors.ErrInvalidResource,
		pErrors.ErrInvalidAction,
		pErrors.ErrDuplicatePolicy,
	}

	for _, err := range errList {
		if err == nil || err.Error() == "" {
			t.Fatalf("expected non-empty error message, got nil or empty")
		}
	}
}
