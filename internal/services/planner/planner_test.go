package planner_test

import (
	"testing"

	"github.com/ioriimasu/jervis/internal/services/planner"
)

func TestPlannerService(t *testing.T) {
	svc := planner.New(nil)

	// Test invalid creation
	_, err := svc.CreateTask("", "title", "desc")
	if err != planner.ErrInvalidTask {
		t.Errorf("expected ErrInvalidTask, got %v", err)
	}

	// Test valid creation
	t1, err := svc.CreateTask("task-1", "First Task", "Description 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if t1.ID != "task-1" || t1.Status != planner.StatusPending {
		t.Errorf("unexpected task fields: %+v", t1)
	}

	// Test duplicate creation
	_, err = svc.CreateTask("task-1", "Duplicate", "desc")
	if err != planner.ErrDuplicateTask {
		t.Errorf("expected ErrDuplicateTask, got %v", err)
	}

	// Test GetTask
	fetched, err := svc.GetTask("task-1")
	if err != nil {
		t.Fatalf("unexpected error getting task: %v", err)
	}
	if fetched.Title != "First Task" {
		t.Errorf("expected title 'First Task', got '%s'", fetched.Title)
	}

	// Test GetTask non-existent
	_, err = svc.GetTask("non-existent")
	if err != planner.ErrTaskNotFound {
		t.Errorf("expected ErrTaskNotFound, got %v", err)
	}

	// Test UpdateTaskStatus
	updated, err := svc.UpdateTaskStatus("task-1", planner.StatusCompleted)
	if err != nil {
		t.Fatalf("unexpected error updating status: %v", err)
	}
	if updated.Status != planner.StatusCompleted {
		t.Errorf("expected StatusCompleted, got %v", updated.Status)
	}

	// Test UpdateTaskStatus non-existent
	_, err = svc.UpdateTaskStatus("non-existent", planner.StatusCompleted)
	if err != planner.ErrTaskNotFound {
		t.Errorf("expected ErrTaskNotFound, got %v", err)
	}

	// Test ListTasks
	_, _ = svc.CreateTask("task-2", "Second Task", "Description 2")
	list := svc.ListTasks()
	if len(list) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(list))
	}
	if list[0].ID != "task-1" || list[1].ID != "task-2" {
		t.Errorf("unexpected list order: %+v", list)
	}
}
