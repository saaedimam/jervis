package projects_test

import (
	"context"
	"testing"

	"github.com/ioriimasu/jervis/internal/memory/store/sqlite"
	"github.com/ioriimasu/jervis/internal/services/projects"
)

func setupTestProjects(t *testing.T) (projects.Service, func()) {
	ctx := context.Background()
	store, err := sqlite.New(":memory:")
	if err != nil {
		t.Fatalf("failed to create sqlite store: %v", err)
	}

	if err := projects.Initialize(ctx, store); err != nil {
		t.Fatalf("failed to initialize projects schema: %v", err)
	}

	svc := projects.New(store, nil)
	return svc, func() { _ = store.Close() }
}

func TestProjectsService(t *testing.T) {
	svc, cleanup := setupTestProjects(t)
	defer cleanup()

	ctx := context.Background()

	// Test invalid creation
	_, err := svc.CreateProject(ctx, "", "name", "desc")
	if err != projects.ErrInvalidProject {
		t.Errorf("expected ErrInvalidProject, got %v", err)
	}

	// Test valid creation
	p1, err := svc.CreateProject(ctx, "proj-1", "First Project", "Description 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p1.ID != "proj-1" || p1.Status != projects.StatusActive {
		t.Errorf("unexpected project fields: %+v", p1)
	}

	// Test duplicate creation
	_, err = svc.CreateProject(ctx, "proj-1", "Duplicate", "desc")
	if err == nil {
		t.Errorf("expected error for duplicate project ID, got nil")
	}

	// Test GetProject
	fetched, err := svc.GetProject(ctx, "proj-1")
	if err != nil {
		t.Fatalf("unexpected error getting project: %v", err)
	}
	if fetched.Name != "First Project" {
		t.Errorf("expected name 'First Project', got '%s'", fetched.Name)
	}

	// Test GetProject non-existent
	_, err = svc.GetProject(ctx, "non-existent")
	if err != projects.ErrProjectNotFound {
		t.Errorf("expected ErrProjectNotFound, got %v", err)
	}

	// Test UpdateProjectStatus
	updated, err := svc.UpdateProjectStatus(ctx, "proj-1", projects.StatusCompleted)
	if err != nil {
		t.Fatalf("unexpected error updating status: %v", err)
	}
	if updated.Status != projects.StatusCompleted {
		t.Errorf("expected StatusCompleted, got %v", updated.Status)
	}

	// Test UpdateProjectStatus non-existent
	_, err = svc.UpdateProjectStatus(ctx, "non-existent", projects.StatusCompleted)
	if err != projects.ErrProjectNotFound {
		t.Errorf("expected ErrProjectNotFound, got %v", err)
	}

	// Test ListProjects
	_, _ = svc.CreateProject(ctx, "proj-2", "Second Project", "Description 2")
	list, err := svc.ListProjects(ctx)
	if err != nil {
		t.Fatalf("unexpected error listing projects: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 projects, got %d", len(list))
	}
	if list[0].ID != "proj-1" || list[1].ID != "proj-2" {
		t.Errorf("unexpected list order: %+v", list)
	}
}
