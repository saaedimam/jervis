package automation

import (
	"context"

	events "github.com/ioriimasu/jervis/internal/runtime/eventbus/contracts"
)

// Action defines a contract for executable actions in the Automation Service.
type Action interface {
	Execute(ctx context.Context, payload map[string]any) error
}

// Trigger defines how a workflow should be activated.
// This can be a cron expression or an event pattern.
type Trigger interface {
	IsTriggered(event events.Event) bool
}

// Workflow represents a sequence of actions with a trigger mechanism.
type Workflow struct {
	ID      string
	Name    string
	Actions []Action
	Trigger Trigger
}

// Registry manages workflow definitions and their registration.
type Registry interface {
	Register(workflow Workflow) error
	Unregister(workflowID string) error
	List() []Workflow
	Get(workflowID string) (Workflow, bool)
}
