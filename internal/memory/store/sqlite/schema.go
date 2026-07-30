package sqlite

import (
	"context"
	"fmt"
)

// Schema defines the SQL statements for initializing the database.
const Schema = `
CREATE TABLE IF NOT EXISTS events (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL,
    source TEXT NOT NULL,
    timestamp DATETIME NOT NULL,
    priority INTEGER NOT NULL,
    payload BLOB NOT NULL,
    metadata BLOB
);

CREATE INDEX IF NOT EXISTS idx_events_timestamp ON events(timestamp);
CREATE INDEX IF NOT EXISTS idx_events_type ON events(type);
`

// Initialize sets up the database schema.
func (d *Driver) Initialize(ctx context.Context) error {
	_, err := d.Exec(ctx, Schema)
	if err != nil {
		return fmt.Errorf("failed to initialize sqlite schema: %w", err)
	}
	return nil
}
