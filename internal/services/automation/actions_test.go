package automation

import (
	"context"
	"github.com/ioriimasu/jervis/internal/aiprovider/registry"
	"github.com/ioriimasu/jervis/internal/services/planner"
	"github.com/ioriimasu/jervis/internal/services/projects"
	"testing"
)

type mockNotion struct{}

func (m *mockNotion) SyncTasks(ctx context.Context, id string, p planner.Service) error { return nil }
func (m *mockNotion) SyncProjects(ctx context.Context, id string, p projects.Service) error {
	return nil
}
func (m *mockNotion) SyncContext(ctx context.Context, name, file, id string) error { return nil }
func (m *mockNotion) SyncMilestones(ctx context.Context, databaseID string, filePath string) error {
	return nil
}
func (m *mockNotion) SyncADRs(ctx context.Context, databaseID string, adrDir string) error {
	return nil
}
func (m *mockNotion) SyncSpecifications(ctx context.Context, databaseID string, specFiles []string) error {
	return nil
}
func (m *mockNotion) SyncDashboard(ctx context.Context, pageID string, status string) error {
	return nil
}

type mockPlanner struct{}

func (m *mockPlanner) CreateTask(ctx context.Context, id, title, desc string) (*planner.Task, error) {
	return nil, nil
}
func (m *mockPlanner) GetTask(ctx context.Context, id string) (*planner.Task, error) { return nil, nil }
func (m *mockPlanner) UpdateTaskStatus(ctx context.Context, id string, status planner.TaskStatus) (*planner.Task, error) {
	return nil, nil
}
func (m *mockPlanner) ListTasks(ctx context.Context) ([]*planner.Task, error) { return nil, nil }

func TestActions_SyncNotion(t *testing.T) {
	action := &SyncNotionAction{
		Notion:   &mockNotion{},
		Planner:  &mockPlanner{},
		Projects: nil,
	}

	err := action.Execute(context.Background(), map[string]any{"sync_type": "tasks", "target_id": "id"})
	if err != nil {
		t.Fatal(err)
	}

	err = action.Execute(context.Background(), map[string]any{"sync_type": "projects", "target_id": "id"})
	if err != nil {
		t.Fatal(err)
	}

	err = action.Execute(context.Background(), map[string]any{"sync_type": "context", "target_id": "id", "name": "n", "file": "f"})
	if err != nil {
		t.Fatal(err)
	}

	err = action.Execute(context.Background(), map[string]any{"sync_type": "unknown"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestActions_CreateTask(t *testing.T) {
	action := &CreateTaskAction{
		Planner: &mockPlanner{},
	}

	err := action.Execute(context.Background(), map[string]any{"id": "1", "title": "t"})
	if err != nil {
		t.Fatal(err)
	}

	err = action.Execute(context.Background(), map[string]any{"id": "1"})
	if err == nil {
		t.Fatal("expected error missing title")
	}
}

func TestActions_Chat(t *testing.T) {
	reg := registry.New()
	action := &ChatAction{
		Registry: reg,
	}

	err := action.Execute(context.Background(), map[string]any{"provider": "p", "model": "m"})
	if err == nil {
		t.Fatal("expected error missing prompt")
	}

	err = action.Execute(context.Background(), map[string]any{"provider": "missing", "model": "m", "prompt": "hi"})
	if err == nil {
		t.Fatal("expected error missing provider")
	}
}
