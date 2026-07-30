package projects

import (
	"context"
	"fmt"

	"github.com/saaedimam/jervis/internal/memory/store/contracts"
)

// Schema defines the SQL statements for the projects service.
const Schema = `
CREATE TABLE IF NOT EXISTS projects (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status);
CREATE INDEX IF NOT EXISTS idx_projects_created_at ON projects(created_at);
`

// Initialize sets up the projects service schema in the provided store.
func Initialize(ctx context.Context, store contracts.Store) error {
	_, err := store.Exec(ctx, Schema)
	if err != nil {
		return fmt.Errorf("failed to initialize projects schema: %w", err)
	}
	return nil
}
