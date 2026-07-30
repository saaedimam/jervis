package notion

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestServer(t *testing.T, expectedPath string, method string, response string, status int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if r.Method != method {
			t.Errorf("expected %s, got %s", method, r.Method)
		}
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(status)
		_, _ = w.Write([]byte(response))
	}))
}

func TestClient_GetPage(t *testing.T) {
	server := setupTestServer(t, "/pages/test-page-123", "GET", `{"id": "test-page-123"}`, http.StatusOK)
	defer server.Close()

	originalBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = originalBaseURL }()

	client := NewClient("test-token")
	page, err := client.GetPage(context.Background(), "test-page-123")
	if err != nil {
		t.Fatal(err)
	}
	if page["id"] != "test-page-123" {
		t.Errorf("expected id test-page-123, got %v", page["id"])
	}
}

func TestClient_UpdatePageBlocks(t *testing.T) {
	server := setupTestServer(t, "/blocks/page-id/children", "PATCH", `{}`, http.StatusOK)
	defer server.Close()

	originalBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = originalBaseURL }()

	client := NewClient("test-token")
	err := client.UpdatePageBlocks(context.Background(), "page-id", []any{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetBlockChildren(t *testing.T) {
	server := setupTestServer(t, "/blocks/block-id/children", "GET", `{"results": [{"id": "b1"}]}`, http.StatusOK)
	defer server.Close()

	originalBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = originalBaseURL }()

	client := NewClient("test-token")
	results, err := client.GetBlockChildren(context.Background(), "block-id")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0]["id"] != "b1" {
		t.Errorf("unexpected results: %v", results)
	}
}

func TestClient_DeleteBlock(t *testing.T) {
	server := setupTestServer(t, "/blocks/block-id", "DELETE", `{}`, http.StatusOK)
	defer server.Close()

	originalBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = originalBaseURL }()

	client := NewClient("test-token")
	err := client.DeleteBlock(context.Background(), "block-id")
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_CreatePageInDatabase(t *testing.T) {
	server := setupTestServer(t, "/pages", "POST", `{"id": "new-page"}`, http.StatusOK)
	defer server.Close()

	originalBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = originalBaseURL }()

	client := NewClient("test-token")
	id, err := client.CreatePageInDatabase(context.Background(), "db-id", map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if id != "new-page" {
		t.Errorf("expected new-page, got %s", id)
	}
}

func TestClient_UpdateDatabaseProperties(t *testing.T) {
	server := setupTestServer(t, "/databases/db-id", "PATCH", `{}`, http.StatusOK)
	defer server.Close()

	originalBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = originalBaseURL }()

	client := NewClient("test-token")
	err := client.UpdateDatabaseProperties(context.Background(), "db-id", map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_QueryDatabase(t *testing.T) {
	server := setupTestServer(t, "/databases/db-id/query", "POST", `{"results": [{"id": "r1"}]}`, http.StatusOK)
	defer server.Close()

	originalBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = originalBaseURL }()

	client := NewClient("test-token")
	results, err := client.QueryDatabase(context.Background(), "db-id", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatal("expected 1 result")
	}
}

func TestClient_UpdatePageProperties(t *testing.T) {
	server := setupTestServer(t, "/pages/page-id", "PATCH", `{}`, http.StatusOK)
	defer server.Close()

	originalBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = originalBaseURL }()

	client := NewClient("test-token")
	err := client.UpdatePageProperties(context.Background(), "page-id", map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_Do_Error(t *testing.T) {
	server := setupTestServer(t, "/pages/fail", "GET", `error`, http.StatusInternalServerError)
	defer server.Close()

	originalBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = originalBaseURL }()

	client := NewClient("test-token")
	_, err := client.GetPage(context.Background(), "fail")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
