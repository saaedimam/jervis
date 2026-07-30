package rule

import (
	"fmt"
	"strings"

	"github.com/saaedimam/jervis/internal/runtime/permissions/contracts"
	errs "github.com/saaedimam/jervis/internal/runtime/permissions/errors"
)

// Rule represents an immutable authorization statement within a security policy.
type Rule struct {
	id          string
	subject     string
	resource    string
	action      string
	effect      contracts.Effect
	description string
}

var _ contracts.Rule = Rule{}

// New constructs and validates an immutable Rule instance.
func New(id, subject, resource, action string, effect contracts.Effect, description string) (Rule, error) {
	r := Rule{
		id:          id,
		subject:     subject,
		resource:    resource,
		action:      action,
		effect:      effect,
		description: description,
	}
	if err := r.Validate(); err != nil {
		return Rule{}, err
	}
	return r, nil
}

// ID returns the rule identifier.
func (r Rule) ID() string {
	return r.id
}

// Subject returns the target subject pattern.
func (r Rule) Subject() string {
	return r.subject
}

// Resource returns the target resource pattern.
func (r Rule) Resource() string {
	return r.resource
}

// Action returns the target action pattern.
func (r Rule) Action() string {
	return r.action
}

// Effect returns the intended evaluation outcome (EffectAllow or EffectDeny).
func (r Rule) Effect() contracts.Effect {
	return r.effect
}

// Description returns an optional description of the rule.
func (r Rule) Description() string {
	return r.description
}

// IsZero reports whether the Rule is uninitialized.
func (r Rule) IsZero() bool {
	return r.id == "" && r.subject == "" && r.resource == "" && r.action == ""
}

// String returns a formatted string representation of the rule.
func (r Rule) String() string {
	if r.IsZero() {
		return ""
	}
	return fmt.Sprintf("RULE[%s]: %s %s:%s:%s", r.id, r.effect.String(), r.subject, r.resource, r.action)
}

// Validate checks structural invariants of the rule.
func (r Rule) Validate() error {
	if strings.TrimSpace(r.id) == "" {
		return fmt.Errorf("%w: rule ID cannot be empty", errs.ErrValidationFailed)
	}
	if strings.TrimSpace(r.subject) == "" {
		return fmt.Errorf("%w: rule subject cannot be empty", errs.ErrInvalidSubject)
	}
	if strings.TrimSpace(r.resource) == "" {
		return fmt.Errorf("%w: rule resource cannot be empty", errs.ErrInvalidResource)
	}
	if strings.TrimSpace(r.action) == "" {
		return fmt.Errorf("%w: rule action cannot be empty", errs.ErrInvalidAction)
	}
	if r.effect != contracts.EffectAllow && r.effect != contracts.EffectDeny {
		return fmt.Errorf("%w: rule effect must be ALLOW or DENY", errs.ErrValidationFailed)
	}
	return nil
}

// Evaluate checks whether a capability matches this rule and returns its Effect or EffectNeutral.
func (r Rule) Evaluate(cap contracts.Capability) contracts.Effect {
	if cap == nil {
		return contracts.EffectNeutral
	}

	if !matchPattern(r.subject, cap.Subject()) {
		return contracts.EffectNeutral
	}
	if !matchPattern(r.resource, cap.Resource()) {
		return contracts.EffectNeutral
	}
	if !matchPattern(r.action, cap.Action()) {
		return contracts.EffectNeutral
	}

	return r.effect
}

func matchPattern(pattern, target string) bool {
	if pattern == "*" || pattern == target {
		return true
	}
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(target, prefix)
	}
	return false
}
