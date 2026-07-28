package types_test

import (
	"testing"
	"time"

	"github.com/ioriimasu/jervis/internal/runtime/types"
)

func TestRuntimeID(t *testing.T) {
	_, err := types.NewRuntimeID("")
	if err == nil {
		t.Fatalf("expected error for empty RuntimeID")
	}

	id, err := types.NewRuntimeID("rt-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id.String() != "rt-123" {
		t.Fatalf("expected string %q, got %q", "rt-123", id.String())
	}
	if id.IsZero() {
		t.Fatalf("expected IsZero false")
	}

	var zeroID types.RuntimeID
	if !zeroID.IsZero() {
		t.Fatalf("expected IsZero true for uninitialized RuntimeID")
	}
}

func TestSessionID(t *testing.T) {
	_, err := types.NewSessionID("")
	if err == nil {
		t.Fatalf("expected error for empty SessionID")
	}

	id, err := types.NewSessionID("sess-456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id.String() != "sess-456" {
		t.Fatalf("expected string %q, got %q", "sess-456", id.String())
	}
	if id.IsZero() {
		t.Fatalf("expected IsZero false")
	}

	var zeroID types.SessionID
	if !zeroID.IsZero() {
		t.Fatalf("expected IsZero true for uninitialized SessionID")
	}
}

func TestEventID(t *testing.T) {
	_, err := types.NewEventID("")
	if err == nil {
		t.Fatalf("expected error for empty EventID")
	}

	id, err := types.NewEventID("evt-789")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id.String() != "evt-789" {
		t.Fatalf("expected string %q, got %q", "evt-789", id.String())
	}
	if id.IsZero() {
		t.Fatalf("expected IsZero false")
	}

	var zeroID types.EventID
	if !zeroID.IsZero() {
		t.Fatalf("expected IsZero true for uninitialized EventID")
	}
}

func TestTimestamp(t *testing.T) {
	var zeroTS types.Timestamp
	if !zeroTS.IsZero() {
		t.Fatalf("expected IsZero true for zero timestamp")
	}
	if zeroTS.String() != "" {
		t.Fatalf("expected empty string for zero timestamp string")
	}

	now := time.Now()
	ts := types.NewTimestamp(now)
	if ts.IsZero() {
		t.Fatalf("expected IsZero false")
	}
	if ts.Time() != now.UTC() {
		t.Fatalf("expected UTC time %v, got %v", now.UTC(), ts.Time())
	}
	if ts.UnixNano() != now.UTC().UnixNano() {
		t.Fatalf("expected unix nano %d, got %d", now.UTC().UnixNano(), ts.UnixNano())
	}
	if ts.String() == "" {
		t.Fatalf("expected non-empty string format")
	}

	nowTS := types.Now()
	if nowTS.IsZero() {
		t.Fatalf("expected Now() to produce non-zero timestamp")
	}
}

func TestState(t *testing.T) {
	validStates := []types.State{
		types.StateCreated,
		types.StateInitializing,
		types.StateRunning,
		types.StateStopping,
		types.StateStopped,
		types.StateFailed,
	}

	for _, s := range validStates {
		if !s.IsValid() {
			t.Errorf("expected state %s to be valid", s)
		}
		if s.String() == "" {
			t.Errorf("expected non-empty string representation")
		}
	}

	invalidState := types.State("Unknown")
	if invalidState.IsValid() {
		t.Errorf("expected Unknown state to be invalid")
	}
}
