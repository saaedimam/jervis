package automation

import (
	"context"
	"fmt"
	"sync"

	events "github.com/ioriimasu/jervis/internal/runtime/eventbus/contracts"
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

			if err := s.engine.Execute(ctx, w, payload); err != nil {
				// Record failure, perhaps publish an automation.failed event
				_ = s.publisher.Publish(nil) // Placeholder for failure event
				return fmt.Errorf("workflow execution failed: %w", err)
			}
		}
	}

	return nil
}
