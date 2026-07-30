package contracts_test

import (
	"testing"

	"github.com/saaedimam/jervis/internal/runtime/permissions/contracts"
)

func TestEffectString(t *testing.T) {
	tests := []struct {
		effect   contracts.Effect
		expected string
	}{
		{contracts.EffectAllow, "ALLOW"},
		{contracts.EffectDeny, "DENY"},
		{contracts.EffectNeutral, "NEUTRAL"},
		{contracts.Effect(99), "NEUTRAL"},
	}

	for _, tt := range tests {
		if got := tt.effect.String(); got != tt.expected {
			t.Errorf("Effect(%d).String() = %q, want %q", tt.effect, got, tt.expected)
		}
	}
}
