package capability_test

import (
	"errors"
	"testing"

	"github.com/saaedimam/jervis/internal/runtime/permissions/capability"
	errs "github.com/saaedimam/jervis/internal/runtime/permissions/errors"
)

func TestCapabilityValidCreation(t *testing.T) {
	capObj, err := capability.New("user:admin", "fs:config.json", "read")
	if err != nil {
		t.Fatalf("unexpected error creating capability: %v", err)
	}

	if capObj.Subject() != "user:admin" {
		t.Errorf("Subject() = %q, want user:admin", capObj.Subject())
	}
	if capObj.Resource() != "fs:config.json" {
		t.Errorf("Resource() = %q, want fs:config.json", capObj.Resource())
	}
	if capObj.Action() != "read" {
		t.Errorf("Action() = %q, want read", capObj.Action())
	}
	if capObj.String() != "user:admin:fs:config.json:read" {
		t.Errorf("String() = %q, want user:admin:fs:config.json:read", capObj.String())
	}
	if capObj.IsZero() {
		t.Errorf("expected IsZero() to be false for initialized capability")
	}
}

func TestCapabilityZeroValue(t *testing.T) {
	var zeroCap capability.Capability
	if !zeroCap.IsZero() {
		t.Errorf("expected IsZero() to be true for uninitialized capability")
	}
	if zeroCap.String() != "" {
		t.Errorf("expected String() to be empty for zero value, got %q", zeroCap.String())
	}
	if zeroCap.Subject() != "" || zeroCap.Resource() != "" || zeroCap.Action() != "" {
		t.Errorf("accessors should return empty strings for zero value")
	}
}

func TestCapabilityValidationErrors(t *testing.T) {
	// Empty Subject
	_, err := capability.New("", "fs:file", "read")
	if !errors.Is(err, errs.ErrInvalidSubject) {
		t.Fatalf("expected ErrInvalidSubject, got %v", err)
	}

	// Empty Resource
	_, err = capability.New("user:1", "", "read")
	if !errors.Is(err, errs.ErrInvalidResource) {
		t.Fatalf("expected ErrInvalidResource, got %v", err)
	}

	// Empty Action
	_, err = capability.New("user:1", "fs:file", "")
	if !errors.Is(err, errs.ErrInvalidAction) {
		t.Fatalf("expected ErrInvalidAction, got %v", err)
	}
}
