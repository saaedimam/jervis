package model

import (
	"testing"
)

func TestEntry(t *testing.T) {
	meta := map[string]string{"foo": "bar"}
	e := NewEntry("id1", "content", meta)

	if e.ID() != "id1" {
		t.Error("ID mismatch")
	}
	if e.Content() != "content" {
		t.Error("Content mismatch")
	}

	m := e.Metadata()
	if m["foo"] != "bar" {
		t.Error("Metadata mismatch")
	}

	// Ensure defensive copy
	m["foo"] = "changed"
	if e.Metadata()["foo"] != "bar" {
		t.Error("Defensive copy failed")
	}

	if e.Timestamp().IsZero() {
		t.Error("Timestamp should not be zero")
	}
}
