package decision_test

import (
	"testing"

	"github.com/ioriimasu/jervis/internal/runtime/permissions/contracts"
	"github.com/ioriimasu/jervis/internal/runtime/permissions/decision"
)

func TestDecisionAllow(t *testing.T) {
	d := decision.NewAllow("explicit allow rule matched")

	if !d.IsAllowed() {
		t.Errorf("expected IsAllowed() to be true for Allow decision")
	}
	if d.Effect() != contracts.EffectAllow {
		t.Errorf("Effect() = %v, want EffectAllow", d.Effect())
	}
	if d.Reason() != "explicit allow rule matched" {
		t.Errorf("Reason() = %q, want explicit allow rule matched", d.Reason())
	}
	if d.String() != "DECISION[ALLOW]: explicit allow rule matched" {
		t.Errorf("String() = %q, want DECISION[ALLOW]: explicit allow rule matched", d.String())
	}
}

func TestDecisionDeny(t *testing.T) {
	d := decision.NewDeny("default deny policy enforced")

	if d.IsAllowed() {
		t.Errorf("expected IsAllowed() to be false for Deny decision")
	}
	if d.Effect() != contracts.EffectDeny {
		t.Errorf("Effect() = %v, want EffectDeny", d.Effect())
	}
	if d.Reason() != "default deny policy enforced" {
		t.Errorf("Reason() = %q, want default deny policy enforced", d.Reason())
	}
	if d.String() != "DECISION[DENY]: default deny policy enforced" {
		t.Errorf("String() = %q, want DECISION[DENY]: default deny policy enforced", d.String())
	}
}
