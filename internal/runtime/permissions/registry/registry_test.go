package registry_test

import (
	"errors"
	"testing"

	"github.com/saaedimam/jervis/internal/runtime/permissions/contracts"
	errs "github.com/saaedimam/jervis/internal/runtime/permissions/errors"
	"github.com/saaedimam/jervis/internal/runtime/permissions/policy"
	"github.com/saaedimam/jervis/internal/runtime/permissions/registry"
	"github.com/saaedimam/jervis/internal/runtime/permissions/rule"
)

type invalidPolicy struct{}

func (i invalidPolicy) ID() string              { return "invalid" }
func (i invalidPolicy) Name() string            { return "Invalid" }
func (i invalidPolicy) Description() string     { return "" }
func (i invalidPolicy) Version() string         { return "1.0.0" }
func (i invalidPolicy) Rules() []contracts.Rule { return nil }
func (i invalidPolicy) Count() int              { return 0 }
func (i invalidPolicy) Validate() error         { return errors.New("invalid policy rules") }

func TestRegistryBasicOperations(t *testing.T) {
	reg := registry.New()

	if reg.Count() != 0 {
		t.Fatalf("expected initial count 0, got %d", reg.Count())
	}
	if reg.Policies() != nil {
		t.Fatalf("expected nil for empty Policies()")
	}

	r1, _ := rule.New("r-1", "user:*", "fs:*", "read", contracts.EffectAllow, "Allow read")
	pB, _ := policy.New("p-b", "Policy B", "Desc", "1.0.0", []contracts.Rule{r1})
	pA, _ := policy.New("p-a", "Policy A", "Desc", "1.0.0", []contracts.Rule{r1})

	// Register
	if err := reg.Register(pB); err != nil {
		t.Fatalf("unexpected error registering pB: %v", err)
	}
	if err := reg.Register(pA); err != nil {
		t.Fatalf("unexpected error registering pA: %v", err)
	}

	if reg.Count() != 2 {
		t.Fatalf("expected count 2, got %d", reg.Count())
	}

	if !reg.Contains("p-a") || !reg.Contains("p-b") {
		t.Fatalf("Contains failed for registered policies")
	}

	// Get existing
	gotP, exists := reg.Get("p-a")
	if !exists || gotP.ID() != "p-a" {
		t.Fatalf("Get('p-a') failed")
	}

	// Get missing
	_, exists = reg.Get("missing")
	if exists {
		t.Fatalf("Get('missing') returned true")
	}

	// Test deterministic ID ascending ordering (Policies and Snapshot)
	pols := reg.Policies()
	snaps := reg.Snapshot()
	if len(pols) != 2 || pols[0].ID() != "p-a" || pols[1].ID() != "p-b" {
		t.Fatalf("Policies() ordering mismatch: got [%s, %s]", pols[0].ID(), pols[1].ID())
	}
	if len(snaps) != 2 || snaps[0].ID() != "p-a" || snaps[1].ID() != "p-b" {
		t.Fatalf("Snapshot() ordering mismatch")
	}

	// Test defensive copy
	pols[0] = nil
	if reg.Policies()[0] == nil {
		t.Fatalf("Policies() did not return a defensive copy")
	}

	// Unregister
	if err := reg.Unregister("p-a"); err != nil {
		t.Fatalf("unexpected error unregistering p-a: %v", err)
	}
	if reg.Count() != 1 {
		t.Fatalf("expected count 1 after unregister, got %d", reg.Count())
	}

	// Clear
	reg.Clear()
	if reg.Count() != 0 || reg.Contains("p-b") {
		t.Fatalf("Clear() failed to empty registry")
	}
}

func TestRegistryValidationErrors(t *testing.T) {
	reg := registry.New()

	// Register nil policy
	if err := reg.Register(nil); !errors.Is(err, errs.ErrValidationFailed) {
		t.Fatalf("expected ErrValidationFailed for nil policy")
	}

	// Register invalid policy
	if err := reg.Register(invalidPolicy{}); !errors.Is(err, errs.ErrValidationFailed) {
		t.Fatalf("expected ErrValidationFailed for invalid policy")
	}

	// Register duplicate policy
	r1, _ := rule.New("r-1", "user:*", "fs:*", "read", contracts.EffectAllow, "")
	pA, _ := policy.New("p-a", "Policy A", "Desc", "1.0.0", []contracts.Rule{r1})
	_ = reg.Register(pA)

	if err := reg.Register(pA); !errors.Is(err, errs.ErrDuplicatePolicy) {
		t.Fatalf("expected ErrDuplicatePolicy for duplicate registration, got %v", err)
	}

	// Unregister empty ID
	if err := reg.Unregister(""); !errors.Is(err, errs.ErrValidationFailed) {
		t.Fatalf("expected ErrValidationFailed for empty unregister ID")
	}

	// Unregister non-existent ID
	if err := reg.Unregister("missing"); !errors.Is(err, errs.ErrValidationFailed) {
		t.Fatalf("expected ErrValidationFailed for missing unregister ID")
	}
}
