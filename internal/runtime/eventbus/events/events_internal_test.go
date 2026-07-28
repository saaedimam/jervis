package events

import (
	"testing"
)

func TestEnvelopeDirectNilMetadata(t *testing.T) {
	env := &Envelope{metadata: nil}
	m := env.Metadata()
	if m == nil {
		t.Fatalf("expected non-nil map")
	}
	if len(m) != 0 {
		t.Fatalf("expected empty map")
	}
}
