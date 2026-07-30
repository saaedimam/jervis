// Package contracts contains interfaces for the query engine subsystem.
// It defines the contract between the query engine and other layers.
package contracts

import (
	"context"

	storecontracts "github.com/ioriimasu/jervis/internal/memory/store/contracts"
	eventsContracts "github.com/ioriimasu/jervis/internal/runtime/eventbus/contracts"
)

// QueryEngine defines the interface for the query engine component.
// It provides methods for querying stored events using a domain-specific language.
type QueryEngine interface {
	// Query executes a query string against the underlying storage and returns matching events.
	// The query string follows the defined DSL syntax.
	Query(ctx context.Context, queryString string) ([]eventsContracts.Event, error)

	// Prepare parses and validates a query string without executing it.
	// This allows for query validation and optimization planning.
	Prepare(queryString string) (PreparedQuery, error)

	// Close releases any resources held by the query engine.
	Close() error
}

// PreparedQuery represents a parsed and validated query ready for execution.
// It can be executed multiple times with different parameters or contexts.
type PreparedQuery interface {
	// Execute runs the prepared query against the provided store.
	Execute(ctx context.Context, store Store) ([]eventsContracts.Event, error)

	// String returns the original query string.
	String() string
}

// Store represents the storage interface used by the query engine.
// This is the same store interface used by the timeline engine for consistency.
type Store = storecontracts.Store
