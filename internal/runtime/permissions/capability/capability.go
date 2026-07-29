package capability

import (
	"fmt"

	"github.com/ioriimasu/jervis/internal/runtime/permissions/contracts"
	errs "github.com/ioriimasu/jervis/internal/runtime/permissions/errors"
)

// Subject represents the security entity requesting permission.
type Subject string

// Resource represents the target resource identifier.
type Resource string

// Action represents the operation performed on the resource.
type Action string

// Capability is an immutable value object binding Subject, Resource, and Action.
type Capability struct {
	subject  Subject
	resource Resource
	action   Action
}

var _ contracts.Capability = Capability{}

// New constructs and validates an immutable Capability.
func New(sub Subject, res Resource, act Action) (Capability, error) {
	c := Capability{
		subject:  sub,
		resource: res,
		action:   act,
	}
	if err := c.validateSelf(); err != nil {
		return Capability{}, err
	}
	return c, nil
}

// Subject returns the string subject.
func (c Capability) Subject() string {
	return string(c.subject)
}

// Resource returns the string resource.
func (c Capability) Resource() string {
	return string(c.resource)
}

// Action returns the string action.
func (c Capability) Action() string {
	return string(c.action)
}

// String returns the formatted capability string "subject:resource:action".
func (c Capability) String() string {
	if c.IsZero() {
		return ""
	}
	return fmt.Sprintf("%s:%s:%s", c.subject, c.resource, c.action)
}

// IsZero checks whether the Capability is uninitialized or empty.
func (c Capability) IsZero() bool {
	return c.subject == "" && c.resource == "" && c.action == ""
}

func (c Capability) validateSelf() error {
	if c.subject == "" {
		return fmt.Errorf("%w: subject cannot be empty", errs.ErrInvalidSubject)
	}
	if c.resource == "" {
		return fmt.Errorf("%w: resource cannot be empty", errs.ErrInvalidResource)
	}
	if c.action == "" {
		return fmt.Errorf("%w: action cannot be empty", errs.ErrInvalidAction)
	}
	return nil
}
