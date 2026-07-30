package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ioriimasu/jervis/internal/app"
)

// Server wraps the REST API server implementation.
type Server struct {
	app    *app.App
	server *http.Server
}

// NewServer creates a new Jervis REST API server.
func NewServer(a *app.App, port int) *Server {
	mux := http.NewServeMux()

	s := &Server{
		app: a,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
	}

	// Register routes
	mux.HandleFunc("/api/v1/planner/tasks", s.handleTasks)
	mux.HandleFunc("/api/v1/notion/sync", s.handleNotionSync)

	return s
}

// Start starts the REST API server.
func (s *Server) Start() error {
	fmt.Printf("REST API server starting on %s...\n", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *Server) handleTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	switch r.Method {
	case http.MethodGet:
		tasks, err := s.app.Services.Planner.ListTasks(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(w).Encode(tasks)
	case http.MethodPost:
		var req struct {
			ID          string `json:"id"`
			Title       string `json:"title"`
			Description string `json:"description"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		task, err := s.app.Services.Planner.CreateTask(ctx, req.ID, req.Title, req.Description)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(task)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleNotionSync(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Type string `json:"type"`
		Name string `json:"name"`
		File string `json:"file"`
		ID   string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	var err error
	switch req.Type {
	case "context":
		err = s.app.Services.Notion.SyncContext(ctx, req.Name, req.File, req.ID)
	case "tasks":
		err = s.app.Services.Notion.SyncTasks(ctx, req.ID, s.app.Services.Planner)
	default:
		http.Error(w, "Unsupported sync type", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "Sync successful")
}
