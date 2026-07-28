package registry

import (
	"fmt"
	"sort"
	"strings"
	"unicode"

	errs "github.com/ioriimasu/jervis/internal/runtime/eventbus/errors"
	"github.com/ioriimasu/jervis/internal/runtime/eventbus/subscription"
)

// Registry manages synchronous event subscriptions deterministically.
type Registry struct {
	subscriptions map[subscription.SubscriptionID]subscription.Subscription
	nextSeq       uint64
}

// NewRegistry creates a new initialized Registry.
func NewRegistry() *Registry {
	return &Registry{
		subscriptions: make(map[subscription.SubscriptionID]subscription.Subscription),
		nextSeq:       1,
	}
}

// ValidatePattern checks if a subscription pattern string is syntactically valid.
func ValidatePattern(p string) error {
	if p == "" {
		return fmt.Errorf("%w: pattern cannot be empty", errs.ErrValidationFailed)
	}
	if p == "*" {
		return nil
	}
	for _, r := range p {
		if unicode.IsUpper(r) || unicode.IsSpace(r) {
			return fmt.Errorf("%w: pattern %q must be lowercase and contain no spaces", errs.ErrValidationFailed, p)
		}
	}
	if strings.Contains(p, "..") || strings.HasPrefix(p, ".") {
		return fmt.Errorf("%w: pattern %q contains invalid dot sequence", errs.ErrValidationFailed, p)
	}

	cleanPattern := p
	if strings.HasSuffix(cleanPattern, ".*") {
		cleanPattern = strings.TrimSuffix(cleanPattern, ".*")
	} else if strings.HasSuffix(cleanPattern, "*") {
		cleanPattern = strings.TrimSuffix(cleanPattern, "*")
		if strings.HasSuffix(cleanPattern, ".") {
			cleanPattern = strings.TrimSuffix(cleanPattern, ".")
		}
	}

	if cleanPattern == "" {
		return fmt.Errorf("%w: pattern %q resolves to empty namespace", errs.ErrValidationFailed, p)
	}

	parts := strings.Split(cleanPattern, ".")
	for _, part := range parts {
		if part == "" {
			return fmt.Errorf("%w: pattern %q contains empty segment", errs.ErrValidationFailed, p)
		}
	}

	return nil
}

// MatchesPattern tests whether an eventType matches a subscription pattern.
func MatchesPattern(eventType string, pattern string) bool {
	if pattern == "*" {
		return true
	}
	if pattern == eventType {
		return true
	}
	if strings.HasSuffix(pattern, ".*") {
		prefix := strings.TrimSuffix(pattern, "*") // e.g. "runtime."
		return strings.HasPrefix(eventType, prefix) && len(eventType) > len(prefix)
	}
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(eventType, prefix)
	}
	return false
}

// Register adds a new subscription to the registry.
func (r *Registry) Register(sub subscription.Subscription) error {
	if err := sub.Validate(); err != nil {
		return err
	}
	if err := ValidatePattern(sub.Pattern()); err != nil {
		return err
	}
	if _, exists := r.subscriptions[sub.ID()]; exists {
		return fmt.Errorf("%w: subscription ID %q already registered", errs.ErrDuplicateSubscriber, sub.ID())
	}

	// Check duplicate handler ID registered for the exact same pattern
	for _, existing := range r.subscriptions {
		if existing.Pattern() == sub.Pattern() && existing.Handler().ID() == sub.Handler().ID() {
			return fmt.Errorf("%w: handler %q already registered for pattern %q", errs.ErrDuplicateSubscriber, sub.Handler().ID(), sub.Pattern())
		}
	}

	subWithSeq, err := subscription.NewWithSeq(sub.ID(), sub.Pattern(), sub.Priority(), sub.Handler(), r.nextSeq)
	if err != nil {
		return err
	}
	r.nextSeq++
	r.subscriptions[subWithSeq.ID()] = subWithSeq
	return nil
}

// Unregister removes a subscription by SubscriptionID.
func (r *Registry) Unregister(subID subscription.SubscriptionID) error {
	if subID.IsZero() {
		return fmt.Errorf("%w: subscription ID cannot be empty", errs.ErrValidationFailed)
	}
	if _, exists := r.subscriptions[subID]; !exists {
		return fmt.Errorf("%w: subscription ID %q not found", errs.ErrValidationFailed, subID)
	}
	delete(r.subscriptions, subID)
	return nil
}

// Lookup finds all matching subscriptions for an eventType, ordered by priority DESC then sequence ASC.
func (r *Registry) Lookup(eventType string) []subscription.Subscription {
	var matches []subscription.Subscription
	for _, sub := range r.subscriptions {
		if MatchesPattern(eventType, sub.Pattern()) {
			matches = append(matches, sub)
		}
	}
	r.sortSubscriptions(matches)
	return matches
}

// LookupExact returns subscriptions matching the exact eventType string.
func (r *Registry) LookupExact(eventType string) []subscription.Subscription {
	var matches []subscription.Subscription
	for _, sub := range r.subscriptions {
		if sub.Pattern() == eventType {
			matches = append(matches, sub)
		}
	}
	r.sortSubscriptions(matches)
	return matches
}

// LookupPattern returns subscriptions matching the given subscription pattern string.
func (r *Registry) LookupPattern(pattern string) []subscription.Subscription {
	var matches []subscription.Subscription
	for _, sub := range r.subscriptions {
		if sub.Pattern() == pattern {
			matches = append(matches, sub)
		}
	}
	r.sortSubscriptions(matches)
	return matches
}

// Contains reports whether a subscription with subID exists in the registry.
func (r *Registry) Contains(subID subscription.SubscriptionID) bool {
	_, exists := r.subscriptions[subID]
	return exists
}

// Count returns the total number of registered subscriptions.
func (r *Registry) Count() int {
	return len(r.subscriptions)
}

// Clear removes all registered subscriptions and resets the sequence counter.
func (r *Registry) Clear() {
	r.subscriptions = make(map[subscription.SubscriptionID]subscription.Subscription)
	r.nextSeq = 1
}

// Snapshot returns a defensive copy slice of all registered subscriptions, sorted deterministically.
func (r *Registry) Snapshot() []subscription.Subscription {
	snapshot := make([]subscription.Subscription, 0, len(r.subscriptions))
	for _, sub := range r.subscriptions {
		snapshot = append(snapshot, sub)
	}
	r.sortSubscriptions(snapshot)
	return snapshot
}

func (r *Registry) sortSubscriptions(subs []subscription.Subscription) {
	sort.SliceStable(subs, func(i, j int) bool {
		if subs[i].Priority() != subs[j].Priority() {
			return subs[i].Priority() > subs[j].Priority()
		}
		return subs[i].Seq() < subs[j].Seq()
	})
}
