package automation

import (
	"context"
	"fmt"
)

// Engine defines the interface for the automation engine.
type Engine interface {
	// Execute runs a workflow's actions sequentially.
	Execute(ctx context.Context, workflow Workflow, initialPayload map[string]any) error
}

type engine struct{}

// NewEngine creates a new workflow execution engine.
func NewEngine() Engine {
	return &engine{}
}

// Execute runs the actions of a workflow in sequence.
// It stops and propagates any error encountered by an action.
func (e *engine) Execute(ctx context.Context, workflow Workflow, initialPayload map[string]any) error {
	// Defensive copy of the payload to pass down the chain
	payload := make(map[string]any)
	for k, v := range initialPayload {
		payload[k] = v
	}

	for i, action := range workflow.Actions {
		// Respect context cancellation
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("workflow %q cancelled at action %d: %w", workflow.ID, i, err)
		}

		if err := action.Execute(ctx, payload); err != nil {
			return fmt.Errorf("workflow %q failed at action %d: %w", workflow.ID, i, err)
		}
	}

	return nil
}
