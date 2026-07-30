package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/ioriimasu/jervis/internal/memory/contracts"
	storecontracts "github.com/ioriimasu/jervis/internal/memory/store/contracts"
	runtimecontracts "github.com/ioriimasu/jervis/internal/runtime/eventbus/contracts"
	"github.com/ioriimasu/jervis/internal/runtime/eventbus/events"
	"github.com/ioriimasu/jervis/internal/runtime/types"
)

// Engine implements Timeline using a persistent Store.
type Engine struct {
	mu    sync.RWMutex
	store storecontracts.Store
}

// New constructs a new Timeline Engine with the specified persistent store.
func New(store storecontracts.Store) *Engine {
	return &Engine{
		store: store,
	}
}

// Append adds an event to the persistent ledger.
func (e *Engine) Append(ctx context.Context, event runtimecontracts.Event) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	payload, err := json.Marshal(event.Payload())
	if err != nil {
		return fmt.Errorf("failed to marshal event payload: %w", err)
	}

	metadata, err := json.Marshal(event.Metadata())
	if err != nil {
		return fmt.Errorf("failed to marshal event metadata: %w", err)
	}

	_, err = e.store.Exec(ctx,
		"INSERT INTO events (id, type, source, timestamp, priority, payload, metadata) VALUES (?, ?, ?, ?, ?, ?, ?)",
		event.ID().String(),
		event.Type(),
		event.Source(),
		event.Timestamp().Time(),
		event.Priority(),
		payload,
		metadata,
	)
	if err != nil {
		return fmt.Errorf("failed to persist event to store: %w", err)
	}

	return nil
}

// Query retrieves events based on filter criteria from the store.
func (e *Engine) Query(ctx context.Context, filter contracts.Filter) ([]runtimecontracts.Event, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	query := "SELECT id, type, source, timestamp, priority, payload, metadata FROM events WHERE 1=1"
	args := []any{}

	if filter.Type != "" {
		query += " AND type = ?"
		args = append(args, filter.Type)
	}
	if !filter.From.IsZero() {
		query += " AND timestamp >= ?"
		args = append(args, filter.From)
	}
	if !filter.To.IsZero() {
		query += " AND timestamp <= ?"
		args = append(args, filter.To)
	}

	query += " ORDER BY timestamp ASC"

	if filter.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, filter.Limit)
	}

	rows, err := e.store.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query events from store: %w", err)
	}
	defer rows.Close()

	results := []runtimecontracts.Event{}
	for rows.Next() {
		var id, evType, source string
		var timestamp time.Time
		var priority uint8
		var payloadData, metadataData []byte

		if err := rows.Scan(&id, &evType, &source, &timestamp, &priority, &payloadData, &metadataData); err != nil {
			return nil, fmt.Errorf("failed to scan event row: %w", err)
		}

		var payload any
		if err := json.Unmarshal(payloadData, &payload); err != nil {
			return nil, fmt.Errorf("failed to unmarshal event payload: %w", err)
		}

		var metadata map[string]string
		if err := json.Unmarshal(metadataData, &metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal event metadata: %w", err)
		}
		eid, err := types.NewEventID(id)
		if err != nil {
			return nil, fmt.Errorf("failed to create event ID: %w", err)
		}

		// Reconstruct event using the events builder
		builder := events.NewBuilder().
			SetID(eid).
			SetType(events.EventType(evType)).
			SetSource(source).
			SetTimestamp(types.NewTimestamp(timestamp)).
			SetPriority(events.Priority(priority)).
			SetPayload(payload)

		for k, v := range metadata {
			builder.SetMetadata(k, v)
		}

		event, err := builder.Build()
		if err != nil {
			return nil, fmt.Errorf("failed to rebuild event: %w", err)
		}

		results = append(results, event)
	}

	return results, nil
}
