package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/saaedimam/jervis/internal/app"
)

// Server wraps the REST API server implementation.
type Server struct {
	app           *app.App
	server        *http.Server
	authKey       string
	deviceAuthKey string
}

// ServerConfig defines parameters for configuring the REST API server.
type ServerConfig struct {
	Host          string
	Port          int
	APIAuthKey    string
	DeviceAuthKey string
}

// NewServer creates a new Jervis REST API server.
// Default binds to localhost (127.0.0.1) for security.
// Set bindAddress to "0.0.0.0" to expose externally (not recommended).
func NewServer(a *app.App, port int) *Server {
	return NewServerWithAuth(a, port, "")
}

// NewServerWithAuth creates a new Jervis REST API server with optional auth.
// If authKey is empty, no authentication is required.
// If authKey is set, requests must include "Authorization: Bearer <authKey>".
func NewServerWithAuth(a *app.App, port int, authKey string) *Server {
	return NewServerWithConfig(a, ServerConfig{
		Host:       "127.0.0.1",
		Port:       port,
		APIAuthKey: authKey,
	})
}

// NewServerWithConfig creates a new Jervis REST API server with explicit configuration.
func NewServerWithConfig(a *app.App, cfg ServerConfig) *Server {
	mux := http.NewServeMux()

	host := cfg.Host
	if host == "" {
		host = "127.0.0.1"
	}

	s := &Server{
		app:           a,
		authKey:       cfg.APIAuthKey,
		deviceAuthKey: cfg.DeviceAuthKey,
		server: &http.Server{
			Addr:              fmt.Sprintf("%s:%d", host, cfg.Port),
			Handler:           mux,
			ReadTimeout:       45 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      45 * time.Second,
			IdleTimeout:       120 * time.Second,
		},
	}

	// Register routes with optional auth middleware
	mux.HandleFunc("/api/v1/planner/tasks", s.withAuth(s.handleTasks))
	mux.HandleFunc("/api/v1/notion/sync", s.withAuth(s.handleNotionSync))

	// Register device control plane routes
	mux.HandleFunc("/v1/device/events", s.withDeviceAuth(s.handleDeviceEvents))
	mux.HandleFunc("/v1/device/commands", s.withDeviceOrAPIAuth(s.handleDeviceCommands))
	mux.HandleFunc("/v1/device/results", s.withDeviceAuth(s.handleDeviceResults))
	mux.HandleFunc("/v1/device/status", s.withDeviceOrAPIAuth(s.handleDeviceStatus))

	return s
}

// maxRequestBodySize defines the maximum number of bytes accepted for any
// request body that is decoded by the REST API. The current value mirrors the
// historic 1 MiB limit introduced in Phase 1. It is expressed as a constant so
// the invariant is shared across the package and can be adjusted centrally if
// a different policy is adopted.
const maxRequestBodySize = 1 << 20 // 1 MiB

// withAuth wraps a handler with optional authentication.
func (s *Server) withAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.authKey != "" {
			auth := r.Header.Get("Authorization")
			if auth == "" || auth != "Bearer "+s.authKey {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	}
}

// Start starts the REST API server.
func (s *Server) Start() error {
	fmt.Printf("REST API server starting on %s...\n", s.server.Addr)
	return s.server.ListenAndServe()
}

// Handler returns the underlying http.Handler (useful for testing).
func (s *Server) Handler() http.Handler {
	return s.server.Handler
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

		r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize) // 1 MB limit
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

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize) // 1 MB limit
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
