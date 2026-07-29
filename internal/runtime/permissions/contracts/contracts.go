package contracts

// Effect represents the outcome of a rule evaluation.
type Effect uint8

const (
	// EffectNeutral indicates that a rule does not apply to the capability.
	EffectNeutral Effect = iota

	// EffectAllow indicates that a rule explicitly allows the capability.
	EffectAllow

	// EffectDeny indicates that a rule explicitly denies the capability.
	EffectDeny
)

// String returns the string representation of Effect.
func (e Effect) String() string {
	switch e {
	case EffectAllow:
		return "ALLOW"
	case EffectDeny:
		return "DENY"
	default:
		return "NEUTRAL"
	}
}

// Capability defines the interface for an immutable access capability.
type Capability interface {
	Subject() string
	Resource() string
	Action() string
}

// Decision defines the interface for an authorization decision.
type Decision interface {
	IsAllowed() bool
	Effect() Effect
	Reason() string
}

// Rule defines a policy rule statement contract.
type Rule interface {
	ID() string
	Subject() string
	Resource() string
	Action() string
	Effect() Effect
	Description() string
	Evaluate(cap Capability) Effect
	Validate() error
}

// Policy defines a collection of rules identified by ID.
type Policy interface {
	ID() string
	Name() string
	Description() string
	Version() string
	Rules() []Rule
	Count() int
	Validate() error
}

// Validator defines structural validation for capabilities.
type Validator interface {
	Validate(cap Capability) error
}
