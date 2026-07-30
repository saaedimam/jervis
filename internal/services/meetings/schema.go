package meetings

import (
	"context"
	"fmt"

	"github.com/saaedimam/jervis/internal/memory/store/contracts"
)

// Schema defines the SQL statements for the meetings service.
const Schema = `
CREATE TABLE IF NOT EXISTS meetings (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL,
    location TEXT,
    status TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_meetings_status ON meetings(status);
CREATE INDEX IF NOT EXISTS idx_meetings_time ON meetings(start_time, end_time);
`

// Initialize sets up the meetings service schema in the provided store.
func Initialize(ctx context.Context, store contracts.Store) error {
	_, err := store.Exec(ctx, Schema)
	if err != nil {
		return fmt.Errorf("failed to initialize meetings schema: %w", err)
	}
	return nil
}
