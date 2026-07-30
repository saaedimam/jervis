package habits

import (
	"context"
	"fmt"

	"github.com/saaedimam/jervis/internal/memory/store/contracts"
)

// Schema defines the SQL statements for the habits service.
const Schema = `
CREATE TABLE IF NOT EXISTS habits (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    frequency TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS habit_logs (
    habit_id TEXT NOT NULL,
    logged_date DATE NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT 0,
    PRIMARY KEY (habit_id, logged_date),
    FOREIGN KEY (habit_id) REFERENCES habits(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_habits_status ON habits(status);
CREATE INDEX IF NOT EXISTS idx_habit_logs_date ON habit_logs(logged_date);
`

// Initialize sets up the habits service schema in the provided store.
func Initialize(ctx context.Context, store contracts.Store) error {
	_, err := store.Exec(ctx, Schema)
	if err != nil {
		return fmt.Errorf("failed to initialize habits schema: %w", err)
	}
	return nil
}
