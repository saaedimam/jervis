package automation

import (
	"errors"
	"fmt"
	"sort"
	"sync"
)

var (
	ErrWorkflowExists   = errors.New("workflow already exists")
	ErrWorkflowNotFound = errors.New("workflow not found")
	ErrInvalidWorkflow  = errors.New("invalid workflow")
)

// workflowRegistry implementation using sync.RWMutex for thread safety
type workflowRegistry struct {
	mu        sync.RWMutex
	workflows map[string]Workflow
}

// NewRegistry creates a new instance of the Registry interface.
func NewRegistry() Registry {
	return &workflowRegistry{
		workflows: make(map[string]Workflow),
	}
}

func (r *workflowRegistry) Register(workflow Workflow) error {
	if workflow.ID == "" {
		return fmt.Errorf("%w: missing ID", ErrInvalidWorkflow)
	}
	if workflow.Name == "" {
		return fmt.Errorf("%w: missing Name", ErrInvalidWorkflow)
	}
	if len(workflow.Actions) == 0 {
		return fmt.Errorf("%w: no actions defined", ErrInvalidWorkflow)
	}
	if workflow.Trigger == nil {
		return fmt.Errorf("%w: missing trigger", ErrInvalidWorkflow)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.workflows[workflow.ID]; exists {
		return fmt.Errorf("%w: %s", ErrWorkflowExists, workflow.ID)
	}

	// Store a defensive copy (shallow copy of struct is fine here)
	r.workflows[workflow.ID] = workflow
	return nil
}

func (r *workflowRegistry) Unregister(workflowID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.workflows[workflowID]; !exists {
		return fmt.Errorf("%w: %s", ErrWorkflowNotFound, workflowID)
	}

	delete(r.workflows, workflowID)
	return nil
}

func (r *workflowRegistry) List() []Workflow {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]Workflow, 0, len(r.workflows))
	for _, w := range r.workflows {
		result = append(result, w)
	}

	// Deterministic ordering by ID
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})

	return result
}

func (r *workflowRegistry) Get(workflowID string) (Workflow, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	w, exists := r.workflows[workflowID]
	// Defensive copy returned inherently because Workflow is returned by value.
	return w, exists
}
