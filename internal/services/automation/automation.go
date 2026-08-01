package automation

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	events "github.com/saaedimam/jervis/internal/runtime/eventbus/contracts"
	eventspb "github.com/saaedimam/jervis/internal/runtime/eventbus/events"
	"github.com/saaedimam/jervis/internal/runtime/types"
)

// Service represents the automation service.
type Service interface {
	Start(ctx context.Context) error
	Stop() error
	Registry() Registry
	HandleEvent(ctx context.Context, event events.Event) error
}

type service struct {
	mu        sync.Mutex
	publisher events.Publisher
	registry  Registry
	engine    Engine
	// Add scheduler or other dependencies here in the future
}

// NewService creates a new Automation Service instance.
func NewService(publisher events.Publisher) Service {
	return &service{
		publisher: publisher,
		registry:  NewRegistry(),
		engine:    NewEngine(),
	}
}

// Start begins processing automation workflows.
func (s *service) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. In a complete implementation, we would load workflows from persistence here.

	// 2. We would subscribe to the event bus for all registered EventTriggers.
	// (Deferred to a later phase or when the eventbus integration provides dynamic subscribing).

	return nil
}

// Stop halts the automation service.
func (s *service) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return nil
}

// Registry returns the workflow registry.
func (s *service) Registry() Registry {
	return s.registry
}

// HandleEvent evaluates all workflows against an incoming event.
// This is typically called by an event subscriber.
func (s *service) HandleEvent(ctx context.Context, event events.Event) error {
	workflows := s.registry.List()

	for _, w := range workflows {
		if w.Trigger.IsTriggered(event) {
			// In a real system, we'd spawn a goroutine to not block the event bus,
			// and track its lifecycle, but for now we execute directly or log.
			payload := map[string]any{
				"event_id":   event.ID(),
				"event_type": event.Type(),
			}

			if execErr := s.engine.Execute(ctx, w, payload); execErr != nil {
				// Build a deterministic failure event.
				evtID, idErr := types.NewEventID(fmt.Sprintf("automation-failed-%d", time.Now().UnixNano()))
				if idErr != nil {
					return fmt.Errorf("workflow execution failed: %w", execErr)
				}

				failEvt, buildErr := eventspb.NewBuilder().
					SetID(evtID).
					SetType(eventspb.EventType("automation.failed")).
					SetSource("automation.service").
					SetTimestamp(types.Now()).
					SetPayload(map[string]any{
						"workflow_id":      w.ID,
						"original_error":   execErr.Error(),
						"trigger_event_id": event.ID(),
					}).
					Build()
				if buildErr != nil {
					return fmt.Errorf("workflow execution failed: %w", execErr)
				}

				if pubErr := s.publisher.Publish(failEvt); pubErr != nil {
					// Preserve both the workflow error and the publish error.
					return errors.Join(execErr, pubErr)
				}

				// Return the original workflow error after a successful publish.
				return fmt.Errorf("workflow execution failed: %w", execErr)
			}
		}
	}

	return nil
}
