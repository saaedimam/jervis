package rest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/saaedimam/jervis/internal/app"
)

func TestServerAuthMiddleware(t *testing.T) {
	// Minimal app (in‑memory DB) – we only need the Services.Planner stub.
	cfg := app.DefaultConfig()
	cfg.DatabasePath = ":memory:"
	a, err := app.New(context.Background(), cfg)
	if err != nil {
		t.Fatalf("app.New: %v", err)
	}
	defer a.Close()

	const secret = "super‑secret"
	srv := NewServerWithAuth(a, 0, secret) // port 0 – we never call ListenAndServe

	// Helper to execute a request against the server's mux.
	do := func(authHeader string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/planner/tasks", nil)
		if authHeader != "" {
			req.Header.Set("Authorization", authHeader)
		}
		rec := httptest.NewRecorder()
		srv.server.Handler.ServeHTTP(rec, req)
		return rec
	}

	// 1️⃣ No Authorization header → 401
	if rec := do(""); rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without auth, got %d", rec.Code)
	}

	// 2️⃣ Wrong token → 401
	if rec := do("Bearer wrong-token"); rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 with wrong token, got %d", rec.Code)
	}

	// 3️⃣ Correct token → 200 (empty task list)
	if rec := do("Bearer " + secret); rec.Code != http.StatusOK {
		t.Fatalf("expected 200 with correct token, got %d", rec.Code)
	}
}
