package subscription

import (
	"fmt"

	"github.com/saaedimam/jervis/internal/runtime/eventbus/contracts"
	errs "github.com/saaedimam/jervis/internal/runtime/eventbus/errors"
	"github.com/saaedimam/jervis/internal/runtime/eventbus/events"
)

// SubscriptionID represents a unique subscription identifier.
type SubscriptionID string

// String returns the string representation of SubscriptionID.
func (id SubscriptionID) String() string {
	return string(id)
}

// IsZero reports whether the SubscriptionID is empty.
func (id SubscriptionID) IsZero() bool {
	return id == ""
}

// Subscription represents an immutable event subscription container.
type Subscription struct {
	id       SubscriptionID
	pattern  string
	priority events.Priority
	handler  contracts.Handler
	seq      uint64
}

// New constructs a validated immutable Subscription instance.
func New(id SubscriptionID, pattern string, priority events.Priority, handler contracts.Handler) (Subscription, error) {
	sub := Subscription{
		id:       id,
		pattern:  pattern,
		priority: priority,
		handler:  handler,
	}
	if err := sub.Validate(); err != nil {
		return Subscription{}, err
	}
	return sub, nil
}

// NewWithSeq constructs a Subscription with an explicit registration sequence index for stable sorting.
func NewWithSeq(id SubscriptionID, pattern string, priority events.Priority, handler contracts.Handler, seq uint64) (Subscription, error) {
	sub, err := New(id, pattern, priority, handler)
	if err != nil {
		return Subscription{}, err
	}
	sub.seq = seq
	return sub, nil
}

// ID returns the SubscriptionID.
func (s Subscription) ID() SubscriptionID {
	return s.id
}

// Pattern returns the subscribed event topic pattern.
func (s Subscription) Pattern() string {
	return s.pattern
}

// Priority returns the priority level of the subscription.
func (s Subscription) Priority() events.Priority {
	return s.priority
}

// Handler returns the subscriber Handler interface.
func (s Subscription) Handler() contracts.Handler {
	return s.handler
}

// Seq returns the registration sequence number.
func (s Subscription) Seq() uint64 {
	return s.seq
}

// Validate verifies that the subscription parameters satisfy all structural invariants.
func (s Subscription) Validate() error {
	if s.id.IsZero() {
		return fmt.Errorf("%w: subscription ID cannot be empty", errs.ErrValidationFailed)
	}
	if s.handler == nil {
		return fmt.Errorf("%w: subscription handler cannot be nil", errs.ErrValidationFailed)
	}
	if s.handler.ID() == "" {
		return fmt.Errorf("%w: subscription handler ID cannot be empty", errs.ErrValidationFailed)
	}
	if err := events.ValidatePriority(s.priority); err != nil {
		return err
	}
	if s.pattern == "" {
		return fmt.Errorf("%w: subscription pattern cannot be empty", errs.ErrValidationFailed)
	}
	return nil
}
