package validator

import (
	"fmt"
	"strings"

	"github.com/saaedimam/jervis/internal/runtime/permissions/contracts"
	errs "github.com/saaedimam/jervis/internal/runtime/permissions/errors"
)

// Validator implements structural validation for capabilities.
type Validator struct{}

var _ contracts.Validator = (*Validator)(nil)

// New constructs a Validator instance.
func New() *Validator {
	return &Validator{}
}

// Validate performs structural validation on a Capability interface.
func (v *Validator) Validate(cap contracts.Capability) error {
	if cap == nil {
		return fmt.Errorf("%w: capability cannot be nil", errs.ErrValidationFailed)
	}

	sub := strings.TrimSpace(cap.Subject())
	if sub == "" {
		return fmt.Errorf("%w: subject cannot be empty", errs.ErrInvalidSubject)
	}

	res := strings.TrimSpace(cap.Resource())
	if res == "" {
		return fmt.Errorf("%w: resource cannot be empty", errs.ErrInvalidResource)
	}

	act := strings.TrimSpace(cap.Action())
	if act == "" {
		return fmt.Errorf("%w: action cannot be empty", errs.ErrInvalidAction)
	}

	// Validate character boundaries (no space or control chars in subject or action)
	if strings.ContainsAny(sub, " \t\r\n") {
		return fmt.Errorf("%w: subject contains whitespace", errs.ErrInvalidSubject)
	}
	if strings.ContainsAny(act, " \t\r\n") {
		return fmt.Errorf("%w: action contains whitespace", errs.ErrInvalidAction)
	}

	return nil
}
