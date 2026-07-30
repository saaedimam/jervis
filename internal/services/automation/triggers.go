package automation

import (
	events "github.com/ioriimasu/jervis/internal/runtime/eventbus/contracts"
)

// CronTrigger activates a workflow based on a cron expression.
// Note: Full cron support evaluation is deferred; time-based workflows
// will integrate via internal/runtime/scheduler.
type CronTrigger struct {
	Expression string
}

// IsTriggered checks if the event matches the cron schedule.
func (c *CronTrigger) IsTriggered(event events.Event) bool {
	// Deferred implementation for cron triggers
	return false
}

// EventTrigger activates a workflow based on a specific event type.
type EventTrigger struct {
	EventType string
}

// IsTriggered checks if the event matches the expected type/pattern.
func (e *EventTrigger) IsTriggered(event events.Event) bool {
	if event == nil {
		return false
	}
	return event.Type() == e.EventType
}
