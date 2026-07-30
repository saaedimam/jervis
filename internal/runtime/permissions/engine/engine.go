package engine

import (
	"fmt"

	"github.com/ioriimasu/jervis/internal/runtime/permissions/contracts"
	"github.com/ioriimasu/jervis/internal/runtime/permissions/decision"
	"github.com/ioriimasu/jervis/internal/runtime/permissions/registry"
	"github.com/ioriimasu/jervis/internal/runtime/permissions/validator"
)

// Engine orchestrates synchronous capability validation and policy evaluation.
type Engine struct {
	registry  *registry.Registry
	validator *validator.Validator
}

// New constructs a new Engine backed by the provided Policy Registry.
// If reg is nil, a new empty Registry is assigned.
func New(reg *registry.Registry) *Engine {
	if reg == nil {
		reg = registry.New()
	}
	return &Engine{
		registry:  reg,
		validator: validator.New(),
	}
}

// Registry returns the underlying Policy Registry.
func (e *Engine) Registry() *registry.Registry {
	return e.registry
}

// Authorize evaluates an authorization request against all registered policies.
// The evaluation pipeline executes in 6 deterministic stages:
// Stage 1: Validate capability structure
// Stage 2: Retrieve policies from registry
// Stage 3 & 4: Evaluate rules; immediately return DecisionDeny on first EffectDeny ("explicit deny")
// Stage 5: Track EffectAllow match; if no EffectDeny encountered, return DecisionAllow ("explicit allow")
// Stage 6: Fall through to default deny ("default deny policy enforced")
func (e *Engine) Authorize(cap contracts.Capability) contracts.Decision {
	// Stage 1: Validate capability
	if err := e.validator.Validate(cap); err != nil {
		return decision.NewDeny(fmt.Sprintf("invalid capability: %v", err))
	}

	// Stage 2: Retrieve registered policies
	policies := e.registry.Policies()

	hasAllow := false

	// Stage 3: Evaluate rules across all policies
	for _, pol := range policies {
		for _, r := range pol.Rules() {
			eff := r.Evaluate(cap)
			switch eff {
			case contracts.EffectDeny:
				// Stage 4: Short-circuit on explicit deny
				return decision.NewDeny("explicit deny")
			case contracts.EffectAllow:
				// Stage 5: Record explicit allow
				hasAllow = true
			}
		}
	}

	// Stage 5: Explicit allow return
	if hasAllow {
		return decision.NewAllow("explicit allow")
	}

	// Stage 6: Default deny fallback
	return decision.NewDeny("default deny policy enforced")
}
