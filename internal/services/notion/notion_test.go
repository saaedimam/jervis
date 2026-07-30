package notion

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/ioriimasu/jervis/internal/services/planner"
	"github.com/ioriimasu/jervis/internal/services/projects"
)

type mockPlanner struct{}

func (m *mockPlanner) CreateTask(ctx context.Context, id, title, description string) (*planner.Task, error) {
	return nil, nil
}

func (m *mockPlanner) GetTask(ctx context.Context, id string) (*planner.Task, error) { return nil, nil }

func (m *mockPlanner) UpdateTaskStatus(ctx context.Context, id string, status planner.TaskStatus) (*planner.Task, error) {
	return nil, nil
}

func (m *mockPlanner) ListTasks(ctx context.Context) ([]*planner.Task, error) {
	return []*planner.Task{
		{ID: "t1", Title: "Task 1", Status: planner.StatusPending},
		{ID: "t2", Title: "Task 2", Status: planner.StatusCompleted},
	}, nil
}

type mockProjects struct{}

func (m *mockProjects) CreateProject(ctx context.Context, id, name, desc string) (*projects.Project, error) {
	return nil, nil
}

func (m *mockProjects) GetProject(ctx context.Context, id string) (*projects.Project, error) {
	return nil, nil
}

func (m *mockProjects) UpdateProjectStatus(ctx context.Context, id string, status projects.ProjectStatus) (*projects.Project, error) {
	return nil, nil
}

func (m *mockProjects) ListProjects(ctx context.Context) ([]*projects.Project, error) {
	return []*projects.Project{
		{ID: "p1", Name: "Proj 1", Status: projects.StatusActive},
		{ID: "p2", Name: "Proj 2", Status: projects.StatusCompleted},
	}, nil
}

func TestService_SyncContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method == "GET" {
			_, _ = w.Write([]byte(`{"results": [{"id": "block1"}]}`))
		} else {
			_, _ = w.Write([]byte(`{}`))
		}
	}))
	defer server.Close()
	originalBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = originalBaseURL }()

	s := New("token", nil)

	tmpFile := filepath.Join(t.TempDir(), "ctx.md")
	_ = os.WriteFile(tmpFile, []byte("hello"), 0o644)

	err := s.SyncContext(context.Background(), "Ctx", tmpFile, "page1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestService_SyncTasks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method == "POST" && r.URL.Path == "/databases/db1/query" {
			_, _ = w.Write([]byte(`{"results": [{"id": "page1"}]}`))
		} else {
			_, _ = w.Write([]byte(`{"id": "new-page"}`))
		}
	}))
	defer server.Close()
	originalBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = originalBaseURL }()

	s := New("token", nil)
	err := s.SyncTasks(context.Background(), "db1", &mockPlanner{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestService_SyncProjects(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method == "POST" && r.URL.Path == "/databases/db1/query" {
			if callCount == 0 {
				_, _ = w.Write([]byte(`{"results": []}`))
			} else {
				_, _ = w.Write([]byte(`{"results": [{"id": "page1"}]}`))
			}
			callCount++
		} else {
			_, _ = w.Write([]byte(`{"id": "new-page"}`))
		}
	}))
	defer server.Close()
	originalBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = originalBaseURL }()

	s := New("token", nil)
	err := s.SyncProjects(context.Background(), "db1", &mockProjects{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestService_SyncMilestones(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id": "new-page"}`))
	}))
	defer server.Close()
	originalBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = originalBaseURL }()

	s := New("token", nil)
	tmpFile := filepath.Join(t.TempDir(), "mile.md")
	_ = os.WriteFile(tmpFile, []byte("hello"), 0o644)

	err := s.SyncMilestones(context.Background(), "db1", tmpFile)
	if err != nil {
		t.Fatal(err)
	}
}

func TestService_SyncADRs_And_Specs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id": "new-page"}`))
	}))
	defer server.Close()
	originalBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = originalBaseURL }()

	s := New("token", nil)

	tmpDir := t.TempDir()
	_ = os.WriteFile(filepath.Join(tmpDir, "adr1.md"), []byte("hello"), 0o644)

	err := s.SyncADRs(context.Background(), "db1", tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	err = s.SyncSpecifications(context.Background(), "db1", []string{"spec1"})
	if err != nil {
		t.Fatal(err)
	}

	err = s.SyncDashboard(context.Background(), "p1", "ok")
	if err != nil {
		t.Fatal(err)
	}
}
