package policy

import (
	"fmt"
	"strings"

	"github.com/ioriimasu/jervis/internal/runtime/permissions/contracts"
	errs "github.com/ioriimasu/jervis/internal/runtime/permissions/errors"
)

// Policy represents an immutable collection of security rules.
type Policy struct {
	id          string
	name        string
	description string
	version     string
	rules       []contracts.Rule
}

var _ contracts.Policy = Policy{}

// New constructs and validates an immutable Policy instance with defensive rule copies.
func New(id, name, description, version string, rules []contracts.Rule) (Policy, error) {
	rulesCopy := make([]contracts.Rule, len(rules))
	copy(rulesCopy, rules)

	p := Policy{
		id:          id,
		name:        name,
		description: description,
		version:     version,
		rules:       rulesCopy,
	}

	if err := p.Validate(); err != nil {
		return Policy{}, err
	}

	return p, nil
}

// ID returns the policy identifier.
func (p Policy) ID() string {
	return p.id
}

// Name returns the human-readable policy name.
func (p Policy) Name() string {
	return p.name
}

// Description returns the policy description.
func (p Policy) Description() string {
	return p.description
}

// Version returns the policy semver string.
func (p Policy) Version() string {
	return p.version
}

// Rules returns a defensive copy of the registered rules in this policy.
func (p Policy) Rules() []contracts.Rule {
	if len(p.rules) == 0 {
		return nil
	}
	cp := make([]contracts.Rule, len(p.rules))
	copy(cp, p.rules)
	return cp
}

// Count returns the number of rules contained in this policy.
func (p Policy) Count() int {
	return len(p.rules)
}

// IsZero reports whether the Policy is uninitialized.
func (p Policy) IsZero() bool {
	return p.id == "" && p.name == "" && p.version == "" && len(p.rules) == 0
}

// String returns a formatted string representation of the policy.
func (p Policy) String() string {
	if p.IsZero() {
		return ""
	}
	return fmt.Sprintf("POLICY[%s v%s]: %s (%d rules)", p.id, p.version, p.name, len(p.rules))
}

// Validate verifies structural invariants of the policy and its rules.
func (p Policy) Validate() error {
	if strings.TrimSpace(p.id) == "" {
		return fmt.Errorf("%w: policy ID cannot be empty", errs.ErrValidationFailed)
	}
	if strings.TrimSpace(p.name) == "" {
		return fmt.Errorf("%w: policy name cannot be empty", errs.ErrValidationFailed)
	}
	if strings.TrimSpace(p.version) == "" {
		return fmt.Errorf("%w: policy version cannot be empty", errs.ErrValidationFailed)
	}
	if len(p.rules) == 0 {
		return fmt.Errorf("%w: policy must contain at least one rule", errs.ErrValidationFailed)
	}
	for i, r := range p.rules {
		if r == nil {
			return fmt.Errorf("%w: policy rule at index %d cannot be nil", errs.ErrValidationFailed, i)
		}
		if err := r.Validate(); err != nil {
			return fmt.Errorf("%w: invalid rule at index %d: %v", errs.ErrValidationFailed, i, err)
		}
	}
	return nil
}
