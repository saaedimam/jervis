package sqlite_test

import (
	"context"
	"os"
	"testing"

	"github.com/ioriimasu/jervis/internal/memory/store/sqlite"
)

func TestPersistence(t *testing.T) {
	tmpFile := "test_persistence.db"
	defer func() { _ = os.Remove(tmpFile) }()

	ctx := context.Background()

	// Instance 1: Create and initialize
	d1, err := sqlite.New(tmpFile)
	if err != nil {
		t.Fatalf("failed to create d1: %v", err)
	}
	if err := d1.Initialize(ctx); err != nil {
		t.Fatalf("failed to init d1: %v", err)
	}

	_, err = d1.Exec(ctx, "INSERT INTO events (id, type, source, timestamp, priority, payload) VALUES (?, ?, ?, ?, ?, ?)",
		"persist-1", "test.type", "test.source", "2026-01-01 00:00:00", 1, []byte("{}"))
	if err != nil {
		t.Fatalf("failed to insert in d1: %v", err)
	}
	_ = d1.Close()

	// Instance 2: Reopen and check
	d2, err := sqlite.New(tmpFile)
	if err != nil {
		t.Fatalf("failed to create d2: %v", err)
	}
	defer func() { _ = d2.Close() }()

	rows, err := d2.Query(ctx, "SELECT id FROM events WHERE id = ?", "persist-1")
	if err != nil {
		t.Fatalf("failed to query d2: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		t.Errorf("expected to find record in d2")
	}
}
