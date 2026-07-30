package registry

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ioriimasu/jervis/internal/runtime/permissions/contracts"
	errs "github.com/ioriimasu/jervis/internal/runtime/permissions/errors"
)

// Registry manages in-memory storage and retrieval of security policies.
type Registry struct {
	policies map[string]contracts.Policy
}

// New constructs an initialized Registry instance.
func New() *Registry {
	return &Registry{
		policies: make(map[string]contracts.Policy),
	}
}

// Register adds a policy to the registry after structural validation.
func (r *Registry) Register(policy contracts.Policy) error {
	if policy == nil {
		return fmt.Errorf("%w: policy cannot be nil", errs.ErrValidationFailed)
	}
	if err := policy.Validate(); err != nil {
		return fmt.Errorf("%w: invalid policy: %v", errs.ErrValidationFailed, err)
	}
	if _, exists := r.policies[policy.ID()]; exists {
		return fmt.Errorf("%w: policy ID %q is already registered", errs.ErrDuplicatePolicy, policy.ID())
	}
	r.policies[policy.ID()] = policy
	return nil
}

// Unregister removes a policy from the registry by its ID.
func (r *Registry) Unregister(id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("%w: policy ID cannot be empty", errs.ErrValidationFailed)
	}
	if _, exists := r.policies[id]; !exists {
		return fmt.Errorf("%w: policy ID %q not found", errs.ErrValidationFailed, id)
	}
	delete(r.policies, id)
	return nil
}

// Get retrieves a policy by ID if present.
func (r *Registry) Get(id string) (contracts.Policy, bool) {
	p, exists := r.policies[id]
	return p, exists
}

// Policies returns a defensive copy slice of all registered policies sorted deterministically by ID ascending.
func (r *Registry) Policies() []contracts.Policy {
	if len(r.policies) == 0 {
		return nil
	}
	list := make([]contracts.Policy, 0, len(r.policies))
	for _, p := range r.policies {
		list = append(list, p)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].ID() < list[j].ID()
	})
	return list
}

// Snapshot returns a defensive copy slice of all registered policies sorted deterministically by ID ascending.
func (r *Registry) Snapshot() []contracts.Policy {
	return r.Policies()
}

// Count returns the number of registered policies.
func (r *Registry) Count() int {
	return len(r.policies)
}

// Contains reports whether a policy with the specified ID exists in the registry.
func (r *Registry) Contains(id string) bool {
	_, exists := r.policies[id]
	return exists
}

// Clear removes all registered policies from the registry.
func (r *Registry) Clear() {
	r.policies = make(map[string]contracts.Policy)
}
