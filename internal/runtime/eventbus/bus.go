package eventbus

import (
	"fmt"

	"github.com/saaedimam/jervis/internal/runtime/eventbus/contracts"
	"github.com/saaedimam/jervis/internal/runtime/eventbus/dispatcher"
	errs "github.com/saaedimam/jervis/internal/runtime/eventbus/errors"
	"github.com/saaedimam/jervis/internal/runtime/eventbus/events"
	"github.com/saaedimam/jervis/internal/runtime/eventbus/middleware"
	"github.com/saaedimam/jervis/internal/runtime/eventbus/registry"
	"github.com/saaedimam/jervis/internal/runtime/eventbus/subscription"
)

// EventBus is the canonical Event Bus facade orchestrating validation, registry, middleware, and dispatcher.
type EventBus struct {
	reg   *registry.Registry
	disp  *dispatcher.Dispatcher
	chain *middleware.Chain
}

var _ contracts.Publisher = (*EventBus)(nil)

// New constructs an initialized EventBus facade.
func New() *EventBus {
	return &EventBus{
		reg:   registry.NewRegistry(),
		disp:  dispatcher.NewDispatcher(),
		chain: middleware.NewChain(),
	}
}

// Publish validates, intercepts, and dispatches an event envelope synchronously.
func (b *EventBus) Publish(event contracts.Event) error {
	// Stage 1: Validate Event
	if err := events.ValidateEvent(event); err != nil {
		return err
	}

	// Stage 2: Registry Lookup
	subs := b.reg.Lookup(event.Type())

	// Stage 3: Extract Handlers
	handlers := make([]contracts.Handler, len(subs))
	for i, sub := range subs {
		handlers[i] = sub.Handler()
	}

	// Stage 4: Execute Middleware Chain & Dispatcher
	return b.chain.Execute(event, func(evt contracts.Event) error {
		return b.disp.Dispatch(evt, handlers)
	})
}

// Subscribe registers a subscriber handler for a topic pattern.
func (b *EventBus) Subscribe(
	pattern string,
	handler contracts.Handler,
	priority events.Priority,
) (subscription.SubscriptionID, error) {
	if handler == nil {
		return subscription.SubscriptionID(""), fmt.Errorf("%w: subscription handler cannot be nil", errs.ErrValidationFailed)
	}

	subID := subscription.SubscriptionID(fmt.Sprintf("%s:%s", pattern, handler.ID()))
	sub, err := subscription.New(subID, pattern, priority, handler)
	if err != nil {
		return subscription.SubscriptionID(""), err
	}

	if err := b.reg.Register(sub); err != nil {
		return subscription.SubscriptionID(""), err
	}

	return sub.ID(), nil
}

// Unsubscribe removes a subscription from the registry by its SubscriptionID.
func (b *EventBus) Unsubscribe(id subscription.SubscriptionID) error {
	if id.IsZero() {
		return fmt.Errorf("%w: subscription ID cannot be zero", errs.ErrValidationFailed)
	}
	return b.reg.Unregister(id)
}

// Use registers one or more middleware interceptors into the chain in FIFO order.
func (b *EventBus) Use(mw ...contracts.Middleware) {
	b.chain.Use(mw...)
}

// Count returns the current number of active subscriptions registered in the bus.
func (b *EventBus) Count() int {
	return b.reg.Count()
}
