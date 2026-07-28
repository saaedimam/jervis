package events

import (
	"testing"

	"github.com/ioriimasu/jervis/internal/runtime/types"
)

func TestEnvelopeHeaderAndNilMetadata(t *testing.T) {
	evtID, _ := types.NewEventID("evt-hdr")
	env := &Envelope{
		header: Header{
			ID:       evtID,
			Type:     "runtime.test.event",
			Source:   "test",
			Priority: Normal,
			Version:  DefaultVersion,
		},
		payload:  "data",
		metadata: nil,
	}

	hdr := env.Header()
	if hdr.ID != evtID {
		t.Fatalf("expected header ID %v, got %v", evtID, hdr.ID)
	}

	m := env.Metadata()
	if m == nil || len(m) != 0 {
		t.Fatalf("expected empty non-nil metadata map for nil internal metadata")
	}
}
