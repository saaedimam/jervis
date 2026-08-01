package sqlite

import (
	"context"
	"os"
	"testing"
)

// TestDriverPragmas verifies that the SQLite driver configures the expected pragmas
// and connection limits. The test uses a file‑backed temporary database to allow
// inspection of persistent PRAGMA values.
func TestDriverPragmas(t *testing.T) {
	// Create a temporary file‑backed SQLite DB.
	tmp, err := os.CreateTemp(t.TempDir(), "jervis_sqlite_*.db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpPath := tmp.Name()
	_ = tmp.Close()

	drv, err := New(tmpPath)
	if err != nil {
		t.Fatalf("failed to create sqlite driver: %v", err)
	}
	defer func() { _ = drv.Close() }()

	// Verify journal mode is WAL.
	var journal string
	if err := drv.QueryRow(context.Background(), "PRAGMA journal_mode;").Scan(&journal); err != nil {
		t.Fatalf("failed to query journal_mode pragma: %v", err)
	}
	if journal != "wal" && journal != "WAL" {
		t.Errorf("expected journal_mode WAL, got %s", journal)
	}

	// Verify busy timeout is set to 5000 ms.
	var timeout int
	if err := drv.QueryRow(context.Background(), "PRAGMA busy_timeout;").Scan(&timeout); err != nil {
		t.Fatalf("failed to query busy_timeout pragma: %v", err)
	}
	if timeout != 5000 {
		t.Errorf("expected busy_timeout 5000, got %d", timeout)
	}

	// Verify foreign keys are enabled.
	var foreign int
	if err := drv.QueryRow(context.Background(), "PRAGMA foreign_keys;").Scan(&foreign); err != nil {
		t.Fatalf("failed to query foreign_keys pragma: %v", err)
	}
	if foreign != 1 {
		t.Errorf("expected foreign_keys ON (1), got %d", foreign)
	}

	// Verify the driver limits open connections to a single connection.
	if max := drv.db.Stats().MaxOpenConnections; max != 1 {
		t.Errorf("expected MaxOpenConns 1, got %d", max)
	}

	// Concurrency sanity check: two goroutines performing a simple query should
	// both succeed without error, confirming the connection pool does not exceed
	// the configured limit.
	done := make(chan error, 2)
	queryFn := func() {
		_, err := drv.Exec(context.Background(), "SELECT 1")
		done <- err
	}
	go queryFn()
	go queryFn()
	for i := 0; i < 2; i++ {
		if err := <-done; err != nil {
			t.Fatalf("concurrent query failed: %v", err)
		}
	}
}
