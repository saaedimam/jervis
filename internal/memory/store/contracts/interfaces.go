package contracts

import (
	"context"
	"database/sql"
)

// Store defines the interface for the persistent storage driver.
// It abstracts the underlying database (SQLite) to allow for easier testing and potential future migrations.
type Store interface {
	// Exec executes a query without returning any rows.
	// The args are for any placeholder parameters in the query.
	Exec(ctx context.Context, query string, args ...any) (sql.Result, error)

	// Query executes a query that returns rows, typically a SELECT.
	// The args are for any placeholder parameters in the query.
	Query(ctx context.Context, query string, args ...any) (*sql.Rows, error)

	// QueryRow executes a query that is expected to return at most one row.
	// QueryRow always returns a non-nil value. Errors are deferred until Row's Scan method is called.
	QueryRow(ctx context.Context, query string, args ...any) *sql.Row

	// Close closes the database and prevents new queries from starting.
	Close() error
}
