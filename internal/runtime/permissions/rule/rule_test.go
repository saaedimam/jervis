package rule_test

import (
	"errors"
	"testing"

	"github.com/saaedimam/jervis/internal/runtime/permissions/capability"
	"github.com/saaedimam/jervis/internal/runtime/permissions/contracts"
	errs "github.com/saaedimam/jervis/internal/runtime/permissions/errors"
	"github.com/saaedimam/jervis/internal/runtime/permissions/rule"
)

func TestRuleValidCreationAndAccessors(t *testing.T) {
	r, err := rule.New("r-1", "user:admin", "fs:config.json", "read", contracts.EffectAllow, "Allow admin read config")
	if err != nil {
		t.Fatalf("unexpected error creating rule: %v", err)
	}

	if r.ID() != "r-1" {
		t.Errorf("ID() = %q, want r-1", r.ID())
	}
	if r.Subject() != "user:admin" {
		t.Errorf("Subject() = %q, want user:admin", r.Subject())
	}
	if r.Resource() != "fs:config.json" {
		t.Errorf("Resource() = %q, want fs:config.json", r.Resource())
	}
	if r.Action() != "read" {
		t.Errorf("Action() = %q, want read", r.Action())
	}
	if r.Effect() != contracts.EffectAllow {
		t.Errorf("Effect() = %v, want EffectAllow", r.Effect())
	}
	if r.Description() != "Allow admin read config" {
		t.Errorf("Description() = %q, want description", r.Description())
	}
	if r.IsZero() {
		t.Errorf("expected IsZero() to be false")
	}
	if r.String() != "RULE[r-1]: ALLOW user:admin:fs:config.json:read" {
		t.Errorf("String() = %q", r.String())
	}
}

func TestRuleZeroValue(t *testing.T) {
	var r rule.Rule
	if !r.IsZero() {
		t.Errorf("expected IsZero() to be true for zero value")
	}
	if r.String() != "" {
		t.Errorf("expected String() to be empty for zero value")
	}
}

func TestRuleValidationErrors(t *testing.T) {
	// Empty ID
	_, err := rule.New("", "sub", "res", "act", contracts.EffectAllow, "")
	if !errors.Is(err, errs.ErrValidationFailed) {
		t.Fatalf("expected ErrValidationFailed for empty ID, got %v", err)
	}

	// Empty Subject
	_, err = rule.New("r-1", "", "res", "act", contracts.EffectAllow, "")
	if !errors.Is(err, errs.ErrInvalidSubject) {
		t.Fatalf("expected ErrInvalidSubject for empty subject, got %v", err)
	}

	// Empty Resource
	_, err = rule.New("r-1", "sub", "", "act", contracts.EffectAllow, "")
	if !errors.Is(err, errs.ErrInvalidResource) {
		t.Fatalf("expected ErrInvalidResource for empty resource, got %v", err)
	}

	// Empty Action
	_, err = rule.New("r-1", "sub", "res", "", contracts.EffectAllow, "")
	if !errors.Is(err, errs.ErrInvalidAction) {
		t.Fatalf("expected ErrInvalidAction for empty action, got %v", err)
	}

	// Invalid Effect (Neutral)
	_, err = rule.New("r-1", "sub", "res", "act", contracts.EffectNeutral, "")
	if !errors.Is(err, errs.ErrValidationFailed) {
		t.Fatalf("expected ErrValidationFailed for EffectNeutral, got %v", err)
	}
}

func TestRuleEvaluation(t *testing.T) {
	rDeny, _ := rule.New("r-deny", "user:*", "fs:*", "write", contracts.EffectDeny, "Deny write")
	rAllow, _ := rule.New("r-allow", "user:guest", "fs:file.txt", "read", contracts.EffectAllow, "Allow read")
	rActionMismatch, _ := rule.New("r-act", "user:guest", "fs:file.txt", "delete", contracts.EffectDeny, "Deny delete")

	cap1, _ := capability.New("user:admin", "fs:file.txt", "write")
	cap2, _ := capability.New("user:guest", "fs:file.txt", "read")
	cap3, _ := capability.New("user:guest", "net:api", "read")
	capSubjectMismatch, _ := capability.New("other:role", "fs:file.txt", "read")
	capResourceMismatch, _ := capability.New("user:guest", "db:table", "read")

	if eff := rDeny.Evaluate(cap1); eff != contracts.EffectDeny {
		t.Errorf("cap1 evaluate = %v, want EffectDeny", eff)
	}
	if eff := rAllow.Evaluate(cap2); eff != contracts.EffectAllow {
		t.Errorf("cap2 evaluate = %v, want EffectAllow", eff)
	}
	if eff := rAllow.Evaluate(cap3); eff != contracts.EffectNeutral {
		t.Errorf("cap3 evaluate = %v, want EffectNeutral", eff)
	}
	if eff := rAllow.Evaluate(capSubjectMismatch); eff != contracts.EffectNeutral {
		t.Errorf("capSubjectMismatch evaluate = %v, want EffectNeutral", eff)
	}
	if eff := rAllow.Evaluate(capResourceMismatch); eff != contracts.EffectNeutral {
		t.Errorf("capResourceMismatch evaluate = %v, want EffectNeutral", eff)
	}
	if eff := rActionMismatch.Evaluate(cap2); eff != contracts.EffectNeutral {
		t.Errorf("capActionMismatch evaluate = %v, want EffectNeutral", eff)
	}
	if eff := rAllow.Evaluate(nil); eff != contracts.EffectNeutral {
		t.Errorf("nil cap evaluate = %v, want EffectNeutral", eff)
	}
}
