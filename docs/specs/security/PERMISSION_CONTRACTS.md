# Runtime Permission Contracts Specification

## 1. Overview
This document specifies the frozen interface contracts for the Jervis Permission Engine (`internal/runtime/permissions`).

---

## 2. Core Contracts

```go
package contracts

// Effect represents the rule evaluation outcome (EffectAllow, EffectDeny, EffectNeutral).
type Effect uint8

const (
	EffectNeutral Effect = iota
	EffectAllow
	EffectDeny
)

// Capability defines the interface for an immutable access capability.
type Capability interface {
	Subject() string
	Resource() string
	Action() string
}

// Decision represents the final authorization decision.
type Decision interface {
	IsAllowed() bool
	Effect() Effect
	Reason() string
}

// Rule defines a policy rule evaluation statement.
type Rule interface {
	ID() string
	Evaluate(cap Capability) Effect
}

// Policy defines a collection of rules identified by ID.
type Policy interface {
	ID() string
	Rules() []Rule
}

// Authorizer evaluates requests against active policies.
type Authorizer interface {
	Authorize(cap Capability) (Decision, error)
}

// PermissionEngine defines the top-level facade contract for the permission subsystem.
type PermissionEngine interface {
	Authorizer
	RegisterPolicy(policy Policy) error
	UnregisterPolicy(id string) error
	Policies() []Policy
}
```

---

## 3. Invariants
- **No `context.Context`**: Methods execute synchronously on the caller stack.
- **Immutable Returns**: Slices and data structures returned from `Policies()` or `Rules()` MUST be defensive copies.
