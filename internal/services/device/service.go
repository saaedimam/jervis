package device

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	storecontracts "github.com/saaedimam/jervis/internal/memory/store/contracts"
	eventcontracts "github.com/saaedimam/jervis/internal/runtime/eventbus/contracts"
	"github.com/saaedimam/jervis/internal/runtime/eventbus/events"
	"github.com/saaedimam/jervis/internal/runtime/types"
)

type serviceImpl struct {
	store       storecontracts.Store
	publisher   eventcontracts.Publisher
	queue       *CommandQueue
	idempotency *IdempotencyManager

	mu       sync.RWMutex
	lastSeen map[string]time.Time
	lastEvt  map[string]string
	lastCmd  map[string]string
	lastRes  map[string]string
}

// New constructs a new DeviceService instance backed by SQLite and the Jervis event bus.
func New(store storecontracts.Store, publisher eventcontracts.Publisher) Service {
	return &serviceImpl{
		store:       store,
		publisher:   publisher,
		queue:       NewCommandQueue(),
		idempotency: NewIdempotencyManager(),
		lastSeen:    make(map[string]time.Time),
		lastEvt:     make(map[string]string),
		lastCmd:     make(map[string]string),
		lastRes:     make(map[string]string),
	}
}

func (s *serviceImpl) IngestEvent(ctx context.Context, env *Envelope) (bool, error) {
	if err := env.ValidateEnvelope(); err != nil {
		return false, err
	}

	s.recordActivity(env.Device)

	// 1. Fast memory idempotency check
	if s.idempotency.MarkEventSeen(env.ID, env.Timestamp) {
		s.audit("device.event_duplicate", map[string]any{
			"device_id": env.Device,
			"event_id":  env.ID,
			"type":      env.Type,
		})
		return true, nil
	}

	// 2. Durable SQLite check
	var existingID string
	err := s.store.QueryRow(ctx, "SELECT id FROM device_events WHERE id = ?", env.ID).Scan(&existingID)
	if err == nil {
		s.audit("device.event_duplicate", map[string]any{
			"device_id": env.Device,
			"event_id":  env.ID,
			"type":      env.Type,
		})
		return true, nil
	} else if err != sql.ErrNoRows {
		return false, fmt.Errorf("failed checking event idempotency: %w", err)
	}

	// 3. Persist to SQLite
	now := time.Now().UTC()
	payloadStr := string(env.Payload)
	_, err = s.store.Exec(ctx,
		"INSERT INTO device_events (id, device_id, event_type, timestamp, payload, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		env.ID, env.Device, env.Type, env.Timestamp, payloadStr, now,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			s.audit("device.event_duplicate", map[string]any{
				"device_id": env.Device,
				"event_id":  env.ID,
				"type":      env.Type,
			})
			return true, nil
		}
		return false, fmt.Errorf("failed to persist device event: %w", err)
	}

	// 4. Update status tracking
	s.mu.Lock()
	s.lastEvt[env.Device] = fmt.Sprintf("%s (%s)", env.ID, env.Type)
	s.mu.Unlock()

	// 5. Audit log
	s.audit("device.event_received", map[string]any{
		"device_id": env.Device,
		"event_id":  env.ID,
		"type":      env.Type,
		"timestamp": env.Timestamp,
		"bytes":     len(env.Payload),
	})

	// 6. Publish to Jervis EventBus for Hermes/Jervis orchestration
	if s.publisher != nil {
		evtID, idErr := types.NewEventID(fmt.Sprintf("dev-%s", env.ID))
		if idErr == nil {
			var unmarshaledPayload any
			_ = json.Unmarshal(env.Payload, &unmarshaledPayload)

			eventObj, buildErr := events.NewBuilder().
				SetID(evtID).
				SetType(events.EventType("device." + env.Type)).
				SetSource(env.Device).
				SetPayload(unmarshaledPayload).
				Build()
			if buildErr == nil {
				_ = s.publisher.Publish(eventObj)
			}
		}
	}

	return false, nil
}

func (s *serviceImpl) QueueCommand(ctx context.Context, cmd *Envelope) error {
	if err := cmd.ValidateEnvelope(); err != nil {
		return err
	}

	// Strictly decode and validate the command payload
	var cp CommandPayload
	if err := DecodeStrict(cmd.Payload, &cp); err != nil {
		s.audit("device.command_rejected", map[string]any{
			"device_id": cmd.Device,
			"reason":    "invalid_or_unknown_fields",
			"error":     err.Error(),
		})
		return err
	}

	if err := cp.Validate(); err != nil {
		s.audit("device.command_rejected", map[string]any{
			"device_id":  cmd.Device,
			"command_id": cp.CommandID,
			"action":     cp.Action,
			"reason":     "validation_failed",
			"error":      err.Error(),
		})
		return err
	}

	// Check durable database to prevent duplicate command IDs
	var existingID string
	err := s.store.QueryRow(ctx, "SELECT id FROM device_commands WHERE id = ?", cp.CommandID).Scan(&existingID)
	if err == nil {
		return fmt.Errorf("device: duplicate command_id %q", cp.CommandID)
	} else if err != sql.ErrNoRows {
		return fmt.Errorf("failed checking command uniqueness: %w", err)
	}

	// Guard queue capacity before persisting
	if s.queue.Depth(cmd.Device) >= MaxQueueSize {
		s.audit("device.command_rejected", map[string]any{
			"device_id":  cmd.Device,
			"command_id": cp.CommandID,
			"action":     cp.Action,
			"reason":     "queue_full",
		})
		return ErrQueueFull
	}

	// Persist pending command into SQLite
	now := time.Now().UTC()
	_, err = s.store.Exec(ctx,
		"INSERT INTO device_commands (id, device_id, action, payload, status, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		cp.CommandID, cmd.Device, cp.Action, string(cmd.Payload), "PENDING", now,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return fmt.Errorf("device: duplicate command_id %q", cp.CommandID)
		}
		return fmt.Errorf("failed to persist device command: %w", err)
	}

	// Enqueue in memory FIFO queue (max 100)
	if err := s.queue.Enqueue(cmd); err != nil {
		// Roll back the persisted command if memory queue failed
		_, _ = s.store.Exec(ctx, "DELETE FROM device_commands WHERE id = ?", cp.CommandID)
		s.audit("device.command_rejected", map[string]any{
			"device_id":  cmd.Device,
			"command_id": cp.CommandID,
			"action":     cp.Action,
			"reason":     "queue_full",
		})
		return err
	}

	s.audit("device.command_queued", map[string]any{
		"device_id":   cmd.Device,
		"command_id":  cp.CommandID,
		"action":      cp.Action,
		"queue_depth": s.queue.Depth(cmd.Device),
	})

	return nil
}

func (s *serviceImpl) WaitForCommand(ctx context.Context, deviceID string, timeout time.Duration) (*Envelope, error) {
	s.recordActivity(deviceID)

	cmd, err := s.queue.WaitForCommand(ctx, deviceID, timeout)
	if err != nil {
		return nil, err
	}

	if cmd == nil {
		s.audit("device.long_poll_timeout", map[string]any{
			"device_id": deviceID,
			"timeout":   timeout.String(),
		})
		return nil, nil
	}

	// Extract command_id for auditing and status update
	var cp CommandPayload
	_ = json.Unmarshal(cmd.Payload, &cp)

	now := time.Now().UTC()
	_, _ = s.store.Exec(ctx,
		"UPDATE device_commands SET status = 'DELIVERED', delivered_at = ? WHERE id = ?",
		now, cp.CommandID,
	)

	s.mu.Lock()
	s.lastCmd[deviceID] = fmt.Sprintf("%s (%s)", cp.CommandID, cp.Action)
	s.mu.Unlock()

	s.audit("device.command_delivered", map[string]any{
		"device_id":   deviceID,
		"command_id":  cp.CommandID,
		"action":      cp.Action,
		"queue_depth": s.queue.Depth(deviceID),
	})

	return cmd, nil
}

func (s *serviceImpl) IngestResult(ctx context.Context, res *Envelope) (bool, error) {
	if err := res.ValidateEnvelope(); err != nil {
		return false, err
	}

	s.recordActivity(res.Device)

	var rp ResultPayload
	if err := DecodeStrict(res.Payload, &rp); err != nil {
		return false, err
	}

	if err := rp.Validate(); err != nil {
		return false, err
	}

	// Fast in-memory idempotency check
	if s.idempotency.MarkResultSeen(rp.CommandID, res.Timestamp) {
		s.audit("device.result_duplicate", map[string]any{
			"device_id":  res.Device,
			"command_id": rp.CommandID,
			"action":     rp.Action,
		})
		return true, nil
	}

	// Durable SQLite idempotency check
	var existingID string
	err := s.store.QueryRow(ctx, "SELECT command_id FROM device_results WHERE command_id = ?", rp.CommandID).Scan(&existingID)
	if err == nil {
		s.audit("device.result_duplicate", map[string]any{
			"device_id":  res.Device,
			"command_id": rp.CommandID,
			"action":     rp.Action,
		})
		return true, nil
	} else if err != sql.ErrNoRows {
		return false, fmt.Errorf("failed checking result idempotency: %w", err)
	}

	now := time.Now().UTC()
	var errStr *string
	if rp.Error != nil {
		errStr = rp.Error
	}

	dataStr := string(rp.Data)
	_, err = s.store.Exec(ctx,
		"INSERT INTO device_results (command_id, device_id, action, status, data, error, timestamp, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		rp.CommandID, res.Device, rp.Action, rp.Status, dataStr, errStr, res.Timestamp, now,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			s.audit("device.result_duplicate", map[string]any{
				"device_id":  res.Device,
				"command_id": rp.CommandID,
				"action":     rp.Action,
			})
			return true, nil
		}
		return false, fmt.Errorf("failed to persist device result: %w", err)
	}

	// Transition command status from DELIVERED/PENDING to COMPLETED or FAILED
	finalStatus := "COMPLETED"
	if rp.Status != StatusSuccess {
		finalStatus = "FAILED"
	}
	_, _ = s.store.Exec(ctx, "UPDATE device_commands SET status = ? WHERE id = ?", finalStatus, rp.CommandID)

	s.mu.Lock()
	s.lastRes[res.Device] = fmt.Sprintf("%s (%s - %s)", rp.CommandID, rp.Action, rp.Status)
	s.mu.Unlock()

	s.audit("device.result_received", map[string]any{
		"device_id":  res.Device,
		"command_id": rp.CommandID,
		"action":     rp.Action,
		"status":     rp.Status,
	})

	// Publish result to Jervis EventBus
	if s.publisher != nil {
		evtID, idErr := types.NewEventID(fmt.Sprintf("res-%s", rp.CommandID))
		if idErr == nil {
			eventObj, buildErr := events.NewBuilder().
				SetID(evtID).
				SetType(events.EventType("device.result." + rp.Action)).
				SetSource(res.Device).
				SetPayload(rp).
				Build()
			if buildErr == nil {
				_ = s.publisher.Publish(eventObj)
			}
		}
	}

	return false, nil
}

func (s *serviceImpl) GetDeviceStatus(ctx context.Context, deviceID string) (*DeviceStatus, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	lastSeenTime, hasSeen := s.lastSeen[deviceID]
	connected := hasSeen && time.Since(lastSeenTime) < 60*time.Second

	var lastSeenUnix int64
	if hasSeen {
		lastSeenUnix = lastSeenTime.Unix()
	}

	return &DeviceStatus{
		DeviceID:    deviceID,
		Connected:   connected,
		LastEvent:   s.lastEvt[deviceID],
		LastCommand: s.lastCmd[deviceID],
		LastResult:  s.lastRes[deviceID],
		QueueDepth:  s.queue.Depth(deviceID),
		LastSeen:    lastSeenUnix,
	}, nil
}

func (s *serviceImpl) recordActivity(deviceID string) {
	s.mu.Lock()
	s.lastSeen[deviceID] = time.Now().UTC()
	s.mu.Unlock()
}

func (s *serviceImpl) audit(action string, meta map[string]any) {
	// Structured audit log: metadata only, never sensitive payload contents
	builder := fmt.Sprintf("[AUDIT] action=%s", action)
	for k, v := range meta {
		builder += fmt.Sprintf(" %s=%v", k, v)
	}
	log.Println(builder)
}
