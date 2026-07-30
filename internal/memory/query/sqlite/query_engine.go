// Package sqlite implements the SQLite query engine.
package sqlite

import (
	"context"
	"fmt"
	"strings"
	"time"

	memContracts "github.com/ioriimasu/jervis/internal/memory/contracts"
	qContracts "github.com/ioriimasu/jervis/internal/memory/query/contracts"
	storecontracts "github.com/ioriimasu/jervis/internal/memory/store/contracts"
	"github.com/ioriimasu/jervis/internal/memory/timeline/engine"
	eventsContracts "github.com/ioriimasu/jervis/internal/runtime/eventbus/contracts"
)

// sqliteQueryEngine implements the QueryEngine interface using SQLite and the timeline engine.
type sqliteQueryEngine struct {
	timelineEngine *engine.Engine
}

// NewSQLiteQueryEngine creates a new query engine backed by the given SQLite store.
func NewSQLiteQueryEngine(store storecontracts.Store) qContracts.QueryEngine {
	te := engine.New(store)
	return &sqliteQueryEngine{timelineEngine: te}
}

// Query executes a query string against the SQLite database and returns matching events.
// The query string follows a simple DSL:
//
//	SELECT * FROM events [WHERE <condition> [AND <condition> ...]] [ORDER BY timestamp [ASC|DESC]] [LIMIT <n>]
//
// Supported conditions:
//
//	TYPE = <string>
//	TIMESTAMP >= <timestamp>
//	TIMESTAMP <= <timestamp>
//	LIMIT <integer>
//
// Note: This is a basic implementation and will be expanded in future iterations.
func (e *sqliteQueryEngine) Query(ctx context.Context, queryString string) ([]eventsContracts.Event, error) {
	// Parse the query string into a filter.
	filter, err := parseQueryToFilter(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse query: %w", err)
	}

	// Delegate to the timeline engine's query method.
	return e.timelineEngine.Query(ctx, filter)
}

// Prepare parses and validates a query string without executing it.
// This allows for query validation and optimization planning.
func (e *sqliteQueryEngine) Prepare(queryString string) (qContracts.PreparedQuery, error) {
	// For now, we just validate the query by attempting to parse it.
	_, err := parseQueryToFilter(queryString)
	if err != nil {
		return nil, err
	}
	return &preparedQuery{query: queryString}, nil
}

// Close releases any resources held by the query engine.
func (e *sqliteQueryEngine) Close() error {
	// The timeline engine does not have a close method, and the store is owned by the caller.
	return nil
}

// preparedQuery is a simple implementation of PreparedQuery.
type preparedQuery struct {
	query string
}

func (p *preparedQuery) Execute(ctx context.Context, store storecontracts.Store) ([]eventsContracts.Event, error) {
	// Reuse the Query method by creating a temporary query engine.
	qe := NewSQLiteQueryEngine(store)
	return qe.Query(ctx, p.query)
}

func (p *preparedQuery) String() string {
	return p.query
}

// parseQueryToFilter converts a simple query string into a filter.
// This is a basic parser for demonstration purposes.
// It supports: WHERE type = 'type', timestamp >= 'time', timestamp <= 'time', LIMIT N
func parseQueryToFilter(query string) (memContracts.Filter, error) {
	f := memContracts.Filter{}
	query = strings.TrimSpace(query)

	// Convert to uppercase for case-insensitive comparison of keywords
	upperQuery := strings.ToUpper(query)

	// Extract WHERE clause
	whereIdx := strings.Index(upperQuery, "WHERE")
	if whereIdx == -1 {
		// No WHERE clause, return empty filter
		return f, nil
	}

	// Extract the part after WHERE
	wherePart := query[whereIdx+5:]

	// Extract ORDER BY clause (we'll ignore for now but could implement later)
	orderByIdx := strings.Index(strings.ToUpper(wherePart), "ORDER BY")
	if orderByIdx != -1 {
		wherePart = wherePart[:orderByIdx]
	}

	// Extract LIMIT clause
	limitIdx := strings.Index(strings.ToUpper(wherePart), "LIMIT")
	if limitIdx != -1 {
		limitPart := strings.TrimSpace(wherePart[limitIdx+5:])
		// Parse the limit value
		var limitVal int
		_, err := fmt.Sscanf(limitPart, "%d", &limitVal)
		if err != nil {
			return memContracts.Filter{}, fmt.Errorf("invalid LIMIT value: %w", err)
		}
		f.Limit = limitVal
		// Remove the LIMIT clause from wherePart
		wherePart = strings.TrimSpace(wherePart[:limitIdx])
	}

	// Split conditions by AND
	conditions := strings.Split(strings.TrimSpace(wherePart), "AND")
	for _, condition := range conditions {
		condition = strings.TrimSpace(condition)
		if condition == "" {
			continue
		}

		// Parse TYPE = 'value'
		if strings.HasPrefix(strings.ToUpper(condition), "TYPE") {
			// Find the equals sign
			eqIdx := strings.Index(condition, "=")
			if eqIdx == -1 {
				return memContracts.Filter{}, fmt.Errorf("invalid TYPE condition: missing '='")
			}
			valueStr := strings.TrimSpace(condition[eqIdx+1:])
			// Remove quotes if present
			valueStr = strings.Trim(valueStr, "\"'")
			f.Type = valueStr
			continue
		}

		// Parse TIMESTAMP >= 'value'
		if strings.HasPrefix(strings.ToUpper(condition), "TIMESTAMP>=") {
			eqIdx := strings.Index(condition, ">=")
			if eqIdx == -1 {
				return memContracts.Filter{}, fmt.Errorf("invalid TIMESTAMP>= condition: missing '>='")
			}
			valueStr := strings.TrimSpace(condition[eqIdx+2:])
			valueStr = strings.Trim(valueStr, "\"'")
			timestamp, err := time.Parse(time.RFC3339, valueStr)
			if err != nil {
				return memContracts.Filter{}, fmt.Errorf("invalid timestamp format: %w", err)
			}
			f.From = timestamp
			continue
		}

		// Parse TIMESTAMP <= 'value'
		if strings.HasPrefix(strings.ToUpper(condition), "TIMESTAMP<=") {
			eqIdx := strings.Index(condition, "<=")
			if eqIdx == -1 {
				return memContracts.Filter{}, fmt.Errorf("invalid TIMESTAMP<= condition: missing '<='")
			}
			valueStr := strings.TrimSpace(condition[eqIdx+2:])
			valueStr = strings.Trim(valueStr, "\"'")
			timestamp, err := time.Parse(time.RFC3339, valueStr)
			if err != nil {
				return memContracts.Filter{}, fmt.Errorf("invalid timestamp format: %w", err)
			}
			f.To = timestamp
			continue
		}
	}

	return f, nil
}
