package decision

import (
	"fmt"

	"github.com/ioriimasu/jervis/internal/runtime/permissions/contracts"
)

// Decision represents an immutable authorization decision result.
type Decision struct {
	allowed bool
	effect  contracts.Effect
	reason  string
}

var _ contracts.Decision = Decision{}

// NewAllow constructs an ALLOW decision with a reason.
func NewAllow(reason string) Decision {
	return Decision{
		allowed: true,
		effect:  contracts.EffectAllow,
		reason:  reason,
	}
}

// NewDeny constructs a DENY decision with a reason.
func NewDeny(reason string) Decision {
	return Decision{
		allowed: false,
		effect:  contracts.EffectDeny,
		reason:  reason,
	}
}

// IsAllowed returns true if the decision permits access.
func (d Decision) IsAllowed() bool {
	return d.allowed
}

// Effect returns the contracts.Effect enum (EffectAllow or EffectDeny).
func (d Decision) Effect() contracts.Effect {
	return d.effect
}

// Reason returns the explanation string for the decision.
func (d Decision) Reason() string {
	return d.reason
}

// String returns a formatted decision string "DECISION[ALLOW/DENY]: reason".
func (d Decision) String() string {
	return fmt.Sprintf("DECISION[%s]: %s", d.effect.String(), d.reason)
}
