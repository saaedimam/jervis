package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ioriimasu/jervis/internal/memory/store/contracts"
	_ "modernc.org/sqlite"
)

// Driver implements the Store interface using modernc.org/sqlite.
type Driver struct {
	db *sql.DB
}

var _ contracts.Store = (*Driver)(nil)

// New constructs a new SQLite driver and opens a connection to the specified path.
// If the path is ":memory:", an in-memory database is created.
func New(path string) (*Driver, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping sqlite database: %w", err)
	}

	return &Driver{db: db}, nil
}

// Exec executes a query without returning any rows.
func (d *Driver) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return d.db.ExecContext(ctx, query, args...)
}

// Query executes a query that returns rows.
func (d *Driver) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return d.db.QueryContext(ctx, query, args...)
}

// QueryRow executes a query that is expected to return at most one row.
func (d *Driver) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return d.db.QueryRowContext(ctx, query, args...)
}

// Close closes the database connection.
func (d *Driver) Close() error {
	return d.db.Close()
}
