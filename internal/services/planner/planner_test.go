package planner_test

import (
	"context"
	"testing"

	"github.com/ioriimasu/jervis/internal/memory/store/sqlite"
	"github.com/ioriimasu/jervis/internal/services/planner"
)

func setupTestPlanner(t *testing.T) (planner.Service, func()) {
	ctx := context.Background()
	store, err := sqlite.New(":memory:")
	if err != nil {
		t.Fatalf("failed to create sqlite store: %v", err)
	}

	if err := planner.Initialize(ctx, store); err != nil {
		t.Fatalf("failed to initialize planner schema: %v", err)
	}

	svc := planner.New(store, nil)
	return svc, func() { store.Close() }
}

func TestPlannerService(t *testing.T) {
	svc, cleanup := setupTestPlanner(t)
	defer cleanup()

	ctx := context.Background()

	// Test invalid creation
	_, err := svc.CreateTask(ctx, "", "title", "desc")
	if err != planner.ErrInvalidTask {
		t.Errorf("expected ErrInvalidTask, got %v", err)
	}

	// Test valid creation
	t1, err := svc.CreateTask(ctx, "task-1", "First Task", "Description 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if t1.ID != "task-1" || t1.Status != planner.StatusPending {
		t.Errorf("unexpected task fields: %+v", t1)
	}

	// Test duplicate creation (SQLite will return error on PK violation)
	_, err = svc.CreateTask(ctx, "task-1", "Duplicate", "desc")
	if err == nil {
		t.Errorf("expected error for duplicate task ID, got nil")
	}

	// Test GetTask
	fetched, err := svc.GetTask(ctx, "task-1")
	if err != nil {
		t.Fatalf("unexpected error getting task: %v", err)
	}
	if fetched.Title != "First Task" {
		t.Errorf("expected title 'First Task', got '%s'", fetched.Title)
	}

	// Test GetTask non-existent
	_, err = svc.GetTask(ctx, "non-existent")
	if err != planner.ErrTaskNotFound {
		t.Errorf("expected ErrTaskNotFound, got %v", err)
	}

	// Test UpdateTaskStatus
	updated, err := svc.UpdateTaskStatus(ctx, "task-1", planner.StatusCompleted)
	if err != nil {
		t.Fatalf("unexpected error updating status: %v", err)
	}
	if updated.Status != planner.StatusCompleted {
		t.Errorf("expected StatusCompleted, got %v", updated.Status)
	}

	// Test UpdateTaskStatus non-existent
	_, err = svc.UpdateTaskStatus(ctx, "non-existent", planner.StatusCompleted)
	if err != planner.ErrTaskNotFound {
		t.Errorf("expected ErrTaskNotFound, got %v", err)
	}

	// Test ListTasks
	_, _ = svc.CreateTask(ctx, "task-2", "Second Task", "Description 2")
	list, err := svc.ListTasks(ctx)
	if err != nil {
		t.Fatalf("unexpected error listing tasks: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(list))
	}
	// Note: We use created_at ASC for ordering in ListTasks
	if list[0].ID != "task-1" || list[1].ID != "task-2" {
		t.Errorf("unexpected list order: %+v", list)
	}
}
