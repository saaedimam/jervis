package policy_test

import (
	"errors"
	"testing"

	"github.com/saaedimam/jervis/internal/runtime/permissions/contracts"
	errs "github.com/saaedimam/jervis/internal/runtime/permissions/errors"
	"github.com/saaedimam/jervis/internal/runtime/permissions/policy"
	"github.com/saaedimam/jervis/internal/runtime/permissions/rule"
)

type invalidMockRule struct{}

func (i invalidMockRule) ID() string               { return "" }
func (i invalidMockRule) Subject() string          { return "" }
func (i invalidMockRule) Resource() string         { return "" }
func (i invalidMockRule) Action() string           { return "" }
func (i invalidMockRule) Effect() contracts.Effect { return contracts.EffectNeutral }
func (i invalidMockRule) Description() string      { return "" }
func (i invalidMockRule) Evaluate(c contracts.Capability) contracts.Effect {
	return contracts.EffectNeutral
}
func (i invalidMockRule) Validate() error { return errors.New("invalid rule") }

func TestPolicyValidCreationAndAccessors(t *testing.T) {
	r1, _ := rule.New("r-1", "user:admin", "fs:*", "read", contracts.EffectAllow, "Allow read")
	r2, _ := rule.New("r-2", "user:guest", "fs:secret", "*", contracts.EffectDeny, "Deny secret")

	rules := []contracts.Rule{r1, r2}
	pol, err := policy.New("p-1", "FileSystem Policy", "Default FS policy", "1.0.0", rules)
	if err != nil {
		t.Fatalf("unexpected error creating policy: %v", err)
	}

	if pol.ID() != "p-1" {
		t.Errorf("ID() = %q, want p-1", pol.ID())
	}
	if pol.Name() != "FileSystem Policy" {
		t.Errorf("Name() = %q, want FileSystem Policy", pol.Name())
	}
	if pol.Description() != "Default FS policy" {
		t.Errorf("Description() = %q", pol.Description())
	}
	if pol.Version() != "1.0.0" {
		t.Errorf("Version() = %q, want 1.0.0", pol.Version())
	}
	if pol.Count() != 2 {
		t.Errorf("Count() = %d, want 2", pol.Count())
	}
	if pol.IsZero() {
		t.Errorf("expected IsZero() to be false")
	}
	if pol.String() != "POLICY[p-1 v1.0.0]: FileSystem Policy (2 rules)" {
		t.Errorf("String() = %q", pol.String())
	}

	// Test defensive copy of input rules
	rules[0] = nil
	if pol.Rules()[0] == nil {
		t.Fatalf("Policy constructor did not make a defensive copy of input rules")
	}

	// Test defensive copy returned by Rules()
	outRules := pol.Rules()
	outRules[0] = nil
	if pol.Rules()[0] == nil {
		t.Fatalf("Rules() did not return a defensive copy of policy rules")
	}
}

func TestPolicyZeroValue(t *testing.T) {
	var p policy.Policy
	if !p.IsZero() {
		t.Errorf("expected IsZero() to be true for zero value")
	}
	if p.String() != "" {
		t.Errorf("expected String() to be empty for zero value")
	}
	if p.Rules() != nil {
		t.Errorf("expected nil rules for zero value")
	}
}

func TestPolicyValidationErrors(t *testing.T) {
	r1, _ := rule.New("r-1", "sub", "res", "act", contracts.EffectAllow, "")

	// Empty ID
	_, err := policy.New("", "Name", "Desc", "1.0.0", []contracts.Rule{r1})
	if !errors.Is(err, errs.ErrValidationFailed) {
		t.Fatalf("expected ErrValidationFailed for empty ID, got %v", err)
	}

	// Empty Name
	_, err = policy.New("p-1", "", "Desc", "1.0.0", []contracts.Rule{r1})
	if !errors.Is(err, errs.ErrValidationFailed) {
		t.Fatalf("expected ErrValidationFailed for empty Name, got %v", err)
	}

	// Empty Version
	_, err = policy.New("p-1", "Name", "Desc", "", []contracts.Rule{r1})
	if !errors.Is(err, errs.ErrValidationFailed) {
		t.Fatalf("expected ErrValidationFailed for empty Version, got %v", err)
	}

	// Empty Rules
	_, err = policy.New("p-1", "Name", "Desc", "1.0.0", nil)
	if !errors.Is(err, errs.ErrValidationFailed) {
		t.Fatalf("expected ErrValidationFailed for empty rules slice, got %v", err)
	}

	// Nil Rule in Slice
	_, err = policy.New("p-1", "Name", "Desc", "1.0.0", []contracts.Rule{nil})
	if !errors.Is(err, errs.ErrValidationFailed) {
		t.Fatalf("expected ErrValidationFailed for nil rule in slice, got %v", err)
	}

	// Invalid Rule in Slice
	_, err = policy.New("p-1", "Name", "Desc", "1.0.0", []contracts.Rule{invalidMockRule{}})
	if !errors.Is(err, errs.ErrValidationFailed) {
		t.Fatalf("expected ErrValidationFailed for invalid rule in slice, got %v", err)
	}
}
