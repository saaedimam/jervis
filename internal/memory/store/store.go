// Package store provides persistent storage drivers for the Memory Engine.
//
// The primary abstraction is the contracts.Store interface, which allows the Memory Engine
// to interact with different storage backends (defaulting to SQLite) without coupling
// to a specific database implementation.
//
// The default implementation uses modernc.org/sqlite, a pure Go SQLite implementation.
package store

import (
	"github.com/ioriimasu/jervis/internal/memory/store/contracts"
	"github.com/ioriimasu/jervis/internal/memory/store/sqlite"
)

// New returns a new Store implementation using SQLite with the provided data source name.
// For an in-memory database, use ":memory:" as the data source name.
func New(dsn string) (contracts.Store, error) {
	return sqlite.New(dsn)
}

// NewMemory returns a new in-memory Store for testing purposes.
func NewMemory() (contracts.Store, error) {
	return sqlite.New(":memory:")
}
