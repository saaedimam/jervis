package projects_test

import (
	"testing"

	"github.com/ioriimasu/jervis/internal/services/projects"
)

func TestProjectsService(t *testing.T) {
	svc := projects.New(nil)

	// Test invalid creation
	_, err := svc.CreateProject("", "Name", "Desc")
	if err != projects.ErrInvalidProject {
		t.Errorf("expected ErrInvalidProject, got %v", err)
	}

	// Test valid creation
	p1, err := svc.CreateProject("proj-1", "Project 1", "Desc 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p1.ID != "proj-1" || p1.Status != projects.StatusActive {
		t.Errorf("unexpected project fields: %+v", p1)
	}

	// Test duplicate creation
	_, err = svc.CreateProject("proj-1", "Duplicate", "desc")
	if err != projects.ErrDuplicateProject {
		t.Errorf("expected ErrDuplicateProject, got %v", err)
	}

	// Test GetProject
	fetched, err := svc.GetProject("proj-1")
	if err != nil {
		t.Fatalf("unexpected error getting project: %v", err)
	}
	if fetched.Name != "Project 1" {
		t.Errorf("expected Name 'Project 1', got '%s'", fetched.Name)
	}

	// Test GetProject non-existent
	_, err = svc.GetProject("non-existent")
	if err != projects.ErrProjectNotFound {
		t.Errorf("expected ErrProjectNotFound, got %v", err)
	}

	// Test UpdateProjectStatus
	updated, err := svc.UpdateProjectStatus("proj-1", projects.StatusCompleted)
	if err != nil {
		t.Fatalf("unexpected error updating status: %v", err)
	}
	if updated.Status != projects.StatusCompleted {
		t.Errorf("expected StatusCompleted, got %v", updated.Status)
	}

	// Test UpdateProjectStatus non-existent
	_, err = svc.UpdateProjectStatus("non-existent", projects.StatusCompleted)
	if err != projects.ErrProjectNotFound {
		t.Errorf("expected ErrProjectNotFound, got %v", err)
	}

	// Test ListProjects
	_, _ = svc.CreateProject("proj-2", "Project 2", "Desc 2")
	list := svc.ListProjects()
	if len(list) != 2 {
		t.Fatalf("expected 2 projects, got %d", len(list))
	}
	if list[0].ID != "proj-1" || list[1].ID != "proj-2" {
		t.Errorf("unexpected list order: %+v", list)
	}
}
