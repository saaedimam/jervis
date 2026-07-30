package engine_test

import (
	"testing"

	"github.com/saaedimam/jervis/internal/runtime/permissions/capability"
	"github.com/saaedimam/jervis/internal/runtime/permissions/contracts"
	"github.com/saaedimam/jervis/internal/runtime/permissions/engine"
	"github.com/saaedimam/jervis/internal/runtime/permissions/policy"
	"github.com/saaedimam/jervis/internal/runtime/permissions/registry"
	"github.com/saaedimam/jervis/internal/runtime/permissions/rule"
)

func TestEngineInvalidCapability(t *testing.T) {
	eng := engine.New(nil)

	// Nil capability
	dec := eng.Authorize(nil)
	if dec.IsAllowed() {
		t.Fatalf("expected nil capability to be denied")
	}
	if dec.Reason() != "invalid capability: permissions: validation failed: capability cannot be nil" {
		t.Errorf("unexpected reason: %q", dec.Reason())
	}

	// Invalid empty subject
	invalidCap := invalidCapMock{res: "fs:file", act: "read"}
	dec = eng.Authorize(invalidCap)
	if dec.IsAllowed() {
		t.Fatalf("expected invalid capability to be denied")
	}
}

type invalidCapMock struct {
	res string
	act string
}

func (i invalidCapMock) Subject() string  { return "" }
func (i invalidCapMock) Resource() string { return i.res }
func (i invalidCapMock) Action() string   { return i.act }

func TestEngineEmptyRegistryDefaultDeny(t *testing.T) {
	reg := registry.New()
	eng := engine.New(reg)

	cap, _ := capability.New("user:admin", "fs:file.txt", "read")
	dec := eng.Authorize(cap)

	if dec.IsAllowed() {
		t.Fatalf("expected authorization to be denied for empty registry")
	}
	if dec.Effect() != contracts.EffectDeny {
		t.Fatalf("expected EffectDeny, got %v", dec.Effect())
	}
	if dec.Reason() != "default deny policy enforced" {
		t.Fatalf("unexpected reason: %q", dec.Reason())
	}
}

func TestEngineExplicitAllowSinglePolicy(t *testing.T) {
	reg := registry.New()
	eng := engine.New(reg)

	r1, _ := rule.New("r-1", "user:admin", "fs:*", "read", contracts.EffectAllow, "Allow admin read")
	p1, _ := policy.New("p-1", "Admin Policy", "Desc", "1.0.0", []contracts.Rule{r1})

	if err := eng.Registry().Register(p1); err != nil {
		t.Fatalf("unexpected error registering policy: %v", err)
	}

	cap, _ := capability.New("user:admin", "fs:config.json", "read")
	dec := eng.Authorize(cap)

	if !dec.IsAllowed() {
		t.Fatalf("expected authorization to be allowed")
	}
	if dec.Effect() != contracts.EffectAllow {
		t.Fatalf("expected EffectAllow, got %v", dec.Effect())
	}
	if dec.Reason() != "explicit allow" {
		t.Fatalf("unexpected reason: %q", dec.Reason())
	}
}

func TestEngineDenyPrecedenceAcrossMultipleRulesAndPolicies(t *testing.T) {
	reg := registry.New()
	eng := engine.New(reg)

	// Policy 1: Allow read access to all files
	rAllow, _ := rule.New("r-allow", "user:*", "fs:*", "read", contracts.EffectAllow, "Allow all read")
	p1, _ := policy.New("p-1", "Allow Policy", "Desc", "1.0.0", []contracts.Rule{rAllow})

	// Policy 2: Deny read access to secret file
	rDeny, _ := rule.New("r-deny", "*", "fs:secret.key", "read", contracts.EffectDeny, "Deny secret read")
	p2, _ := policy.New("p-2", "Deny Policy", "Desc", "1.0.0", []contracts.Rule{rDeny})

	_ = reg.Register(p1)
	_ = reg.Register(p2)

	// Public file request should be allowed
	capPublic, _ := capability.New("user:guest", "fs:public.txt", "read")
	decPublic := eng.Authorize(capPublic)
	if !decPublic.IsAllowed() || decPublic.Reason() != "explicit allow" {
		t.Fatalf("expected public file read to be allowed")
	}

	// Secret file request should be denied due to deny precedence
	capSecret, _ := capability.New("user:guest", "fs:secret.key", "read")
	decSecret := eng.Authorize(capSecret)
	if decSecret.IsAllowed() {
		t.Fatalf("expected secret file read to be denied")
	}
	if decSecret.Reason() != "explicit deny" {
		t.Fatalf("unexpected reason: %q, want 'explicit deny'", decSecret.Reason())
	}
}

func TestEngineMultipleRulesInSinglePolicy(t *testing.T) {
	reg := registry.New()
	eng := engine.New(reg)

	r1, _ := rule.New("r-1", "user:guest", "fs:public.txt", "read", contracts.EffectAllow, "Allow guest public read")
	r2, _ := rule.New("r-2", "user:guest", "fs:secret.txt", "read", contracts.EffectDeny, "Deny guest secret read")

	p1, _ := policy.New("p-1", "Mixed Policy", "Desc", "1.0.0", []contracts.Rule{r1, r2})
	_ = reg.Register(p1)

	// Request for public file -> Allow
	cap1, _ := capability.New("user:guest", "fs:public.txt", "read")
	dec1 := eng.Authorize(cap1)
	if !dec1.IsAllowed() || dec1.Reason() != "explicit allow" {
		t.Fatalf("expected cap1 to be allowed")
	}

	// Request for secret file -> Deny
	cap2, _ := capability.New("user:guest", "fs:secret.txt", "read")
	dec2 := eng.Authorize(cap2)
	if dec2.IsAllowed() || dec2.Reason() != "explicit deny" {
		t.Fatalf("expected cap2 to be denied with 'explicit deny'")
	}
}

func TestEngineDeterministicEvaluationOrder(t *testing.T) {
	reg := registry.New()
	eng := engine.New(reg)

	r1, _ := rule.New("r-1", "user:admin", "fs:file", "read", contracts.EffectAllow, "")
	pB, _ := policy.New("p-b", "Policy B", "Desc", "1.0.0", []contracts.Rule{r1})
	pA, _ := policy.New("p-a", "Policy A", "Desc", "1.0.0", []contracts.Rule{r1})

	_ = reg.Register(pB)
	_ = reg.Register(pA)

	cap, _ := capability.New("user:admin", "fs:file", "read")

	// Perform 100 evaluations to verify strict deterministic behavior
	for i := 0; i < 100; i++ {
		dec := eng.Authorize(cap)
		if !dec.IsAllowed() || dec.Reason() != "explicit allow" {
			t.Fatalf("non-deterministic behavior detected on iteration %d", i)
		}
	}
}
