package planner

import (
	"context"
	"fmt"

	"github.com/saaedimam/jervis/internal/memory/store/contracts"
)

// Schema defines the SQL statements for the planner service.
const Schema = `
CREATE TABLE IF NOT EXISTS planner_tasks (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_planner_tasks_status ON planner_tasks(status);
CREATE INDEX IF NOT EXISTS idx_planner_tasks_created_at ON planner_tasks(created_at);
`

// Initialize sets up the planner service schema in the provided store.
func Initialize(ctx context.Context, store contracts.Store) error {
	_, err := store.Exec(ctx, Schema)
	if err != nil {
		return fmt.Errorf("failed to initialize planner schema: %w", err)
	}
	return nil
}
