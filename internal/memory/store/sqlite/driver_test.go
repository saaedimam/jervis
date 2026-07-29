package sqlite

import (
	"context"
	"testing"
)

func TestDriver(t *testing.T) {
	ctx := context.Background()

	// Use in-memory database for testing
	driver, err := New(":memory:")
	if err != nil {
		t.Fatalf("failed to create sqlite driver: %v", err)
	}
	defer driver.Close()

	// Test Initialize (Schema creation)
	if err := driver.Initialize(ctx); err != nil {
		t.Fatalf("failed to initialize schema: %v", err)
	}

	// Test Insert
	res, err := driver.Exec(ctx, "INSERT INTO events (id, type, source, timestamp, priority, payload) VALUES (?, ?, ?, ?, ?, ?)",
		"test-id", "test-type", "test-source", "2024-07-29T12:00:00Z", 1, []byte("test-payload"))
	if err != nil {
		t.Fatalf("failed to insert event: %v", err)
	}

	affected, _ := res.RowsAffected()
	if affected != 1 {
		t.Errorf("expected 1 row affected, got %d", affected)
	}

	// Test Query
	rows, err := driver.Query(ctx, "SELECT id, type FROM events WHERE id = ?", "test-id")
	if err != nil {
		t.Fatalf("failed to query event: %v", err)
	}

	if !rows.Next() {
		_ = rows.Close()
		t.Fatal("expected row to be returned")
	}

	var id, eventType string
	if err := rows.Scan(&id, &eventType); err != nil {
		_ = rows.Close()
		t.Fatalf("failed to scan row: %v", err)
	}
	_ = rows.Close()

	if id != "test-id" || eventType != "test-type" {
		t.Errorf("expected test-id and test-type, got %s and %s", id, eventType)
	}

	// Test QueryRow
	var scannedType string
	err = driver.QueryRow(ctx, "SELECT type FROM events WHERE id = ?", "test-id").Scan(&scannedType)
	if err != nil {
		t.Fatalf("failed to QueryRow event: %v", err)
	}
	if scannedType != "test-type" {
		t.Errorf("expected test-type, got %s", scannedType)
	}
}

func TestDriverErrors(t *testing.T) {
	ctx := context.Background()

	// Test invalid path opening error (Ping fails on invalid sqlite file uri)
	_, err := New("file:///invalid_path_does_not_exist/db.sqlite?mode=ro")
	if err == nil {
		t.Error("expected error opening read-only database in non-existent directory")
	}

	// Test Initialize failure on closed DB
	driver, err := New(":memory:")
	if err != nil {
		t.Fatalf("failed to create driver: %v", err)
	}
	_ = driver.Close()

	err = driver.Initialize(ctx)
	if err == nil {
		t.Error("expected error initializing schema on closed database")
	}
}
