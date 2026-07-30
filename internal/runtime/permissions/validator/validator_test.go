package validator_test

import (
	"errors"
	"testing"

	"github.com/saaedimam/jervis/internal/runtime/permissions/capability"
	errs "github.com/saaedimam/jervis/internal/runtime/permissions/errors"
	"github.com/saaedimam/jervis/internal/runtime/permissions/validator"
)

type mockCap struct {
	subject  string
	resource string
	action   string
}

func (m mockCap) Subject() string  { return m.subject }
func (m mockCap) Resource() string { return m.resource }
func (m mockCap) Action() string   { return m.action }

func TestValidatorSuccess(t *testing.T) {
	v := validator.New()
	capObj, err := capability.New("user:admin", "fs:config.json", "read")
	if err != nil {
		t.Fatalf("unexpected error constructing capability: %v", err)
	}

	if err := v.Validate(capObj); err != nil {
		t.Fatalf("expected valid capability, got error: %v", err)
	}
}

func TestValidatorNilCapability(t *testing.T) {
	v := validator.New()
	if err := v.Validate(nil); !errors.Is(err, errs.ErrValidationFailed) {
		t.Fatalf("expected ErrValidationFailed for nil capability, got %v", err)
	}
}

func TestValidatorValidationErrors(t *testing.T) {
	v := validator.New()

	// Empty Subject
	if err := v.Validate(mockCap{subject: "", resource: "res", action: "act"}); !errors.Is(err, errs.ErrInvalidSubject) {
		t.Fatalf("expected ErrInvalidSubject for empty subject, got %v", err)
	}

	// Whitespace Subject
	if err := v.Validate(mockCap{subject: "user admin", resource: "res", action: "act"}); !errors.Is(err, errs.ErrInvalidSubject) {
		t.Fatalf("expected ErrInvalidSubject for whitespace subject, got %v", err)
	}

	// Empty Resource
	if err := v.Validate(mockCap{subject: "sub", resource: "", action: "act"}); !errors.Is(err, errs.ErrInvalidResource) {
		t.Fatalf("expected ErrInvalidResource for empty resource, got %v", err)
	}

	// Empty Action
	if err := v.Validate(mockCap{subject: "sub", resource: "res", action: ""}); !errors.Is(err, errs.ErrInvalidAction) {
		t.Fatalf("expected ErrInvalidAction for empty action, got %v", err)
	}

	// Whitespace Action
	if err := v.Validate(mockCap{subject: "sub", resource: "res", action: "read write"}); !errors.Is(err, errs.ErrInvalidAction) {
		t.Fatalf("expected ErrInvalidAction for whitespace action, got %v", err)
	}
}
