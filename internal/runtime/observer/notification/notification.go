package notification

import (
	"fmt"

	eventcontracts "github.com/ioriimasu/jervis/internal/runtime/eventbus/contracts"
	"github.com/ioriimasu/jervis/internal/runtime/observer/contracts"
	"github.com/ioriimasu/jervis/internal/runtime/types"
)

// notification implements contracts.Notification interface.
type notification struct {
	event      eventcontracts.Event
	observedAt types.Timestamp
}

// New creates a new immutable notification wrapping an event.
func New(event eventcontracts.Event, observedAt types.Timestamp) contracts.Notification {
	if event == nil || observedAt.IsZero() {
		return nil
	}
	return &notification{
		event:      event,
		observedAt: observedAt,
	}
}

// Event returns the wrapped canonical Event interface.
func (n *notification) Event() eventcontracts.Event {
	if n == nil {
		return nil
	}
	return n.event
}

// ObservedAt returns the observation timestamp.
func (n *notification) ObservedAt() types.Timestamp {
	if n == nil {
		return types.Timestamp{}
	}
	return n.observedAt
}

// String returns a formatted representation of the notification.
func (n *notification) String() string {
	if n == nil || n.event == nil {
		return "Notification{nil}"
	}
	return fmt.Sprintf("Notification{event=%s, observedAt=%v}", n.event.ID(), n.observedAt)
}
