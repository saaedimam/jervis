package device

import (
	"context"
	"fmt"

	"github.com/saaedimam/jervis/internal/memory/store/contracts"
)

// Schema defines the SQLite DDL for the device subsystem.
const Schema = `
CREATE TABLE IF NOT EXISTS device_events (
    id TEXT PRIMARY KEY,
    device_id TEXT NOT NULL,
    event_type TEXT NOT NULL,
    timestamp INTEGER NOT NULL,
    payload TEXT NOT NULL,
    created_at DATETIME NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_device_events_device_time ON device_events(device_id, timestamp);

CREATE TABLE IF NOT EXISTS device_commands (
    id TEXT PRIMARY KEY,
    device_id TEXT NOT NULL,
    action TEXT NOT NULL,
    payload TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    delivered_at DATETIME
);
CREATE INDEX IF NOT EXISTS idx_device_commands_device_status ON device_commands(device_id, status);

CREATE TABLE IF NOT EXISTS device_results (
    command_id TEXT PRIMARY KEY,
    device_id TEXT NOT NULL,
    action TEXT NOT NULL,
    status TEXT NOT NULL,
    data TEXT,
    error TEXT,
    timestamp INTEGER NOT NULL,
    created_at DATETIME NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_device_results_device ON device_results(device_id);
`

// Initialize creates the device subsystem tables and indexes if they do not exist.
func Initialize(ctx context.Context, store contracts.Store) error {
	_, err := store.Exec(ctx, Schema)
	if err != nil {
		return fmt.Errorf("failed to initialize device schema: %w", err)
	}
	return nil
}
