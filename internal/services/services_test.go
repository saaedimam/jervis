package services_test

import (
	"context"
	"testing"

	"github.com/saaedimam/jervis/internal/memory/store/sqlite"
	"github.com/saaedimam/jervis/internal/services"
)

func TestContainer(t *testing.T) {
	ctx := context.Background()
	store, err := sqlite.New(":memory:")
	if err != nil {
		t.Fatalf("failed to create sqlite store: %v", err)
	}
	defer func() { _ = store.Close() }()

	// Initialize container
	container, err := services.NewContainer(ctx, store, nil, "")
	if err != nil {
		t.Fatalf("failed to create container: %v", err)
	}

	if container.Planner == nil || container.Projects == nil || container.Habits == nil || container.Meetings == nil {
		t.Error("one or more services were not initialized")
	}

	// Verify one service works
	_, err = container.Planner.CreateTask(ctx, "task-1", "Test Task", "")
	if err != nil {
		t.Fatalf("failed to use planner service from container: %v", err)
	}

	fetched, err := container.Planner.GetTask(ctx, "task-1")
	if err != nil || fetched.ID != "task-1" {
		t.Error("failed to fetch task from container service")
	}
}
