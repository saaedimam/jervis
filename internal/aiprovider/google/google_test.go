package google

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ioriimasu/jervis/internal/aiprovider/contracts"
)

func TestGoogleAdapter_Name(t *testing.T) {
	provider := New("test-key")
	if provider.Name() != "google" {
		t.Errorf("expected 'google', got %s", provider.Name())
	}
}

func TestGoogleAdapter_Chat(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("key") != "test-key" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"candidates": [{"content": {"role": "model", "parts": [{"text": "hello google"}]}}],
			"usageMetadata": {"promptTokenCount": 10, "candidatesTokenCount": 5, "totalTokenCount": 15}
		}`))
	}))
	defer server.Close()

	originalBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = originalBaseURL }()

	provider := New("test-key")
	ctx := context.Background()
	opts := contracts.ChatOptions{Temperature: 0.5, MaxTokens: 100}
	messages := []contracts.Message{
		{Role: contracts.RoleSystem, Content: "sys prompt"},
		{Role: contracts.RoleUser, Content: "hi"},
		{Role: contracts.RoleAssistant, Content: "hello"},
	}

	resp, err := provider.Chat(ctx, "gemini-1.5", messages, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Choices) != 1 {
		t.Fatalf("expected 1 choice, got %d", len(resp.Choices))
	}
	if resp.Choices[0].Message.Content != "hello google" {
		t.Errorf("expected 'hello google', got %s", resp.Choices[0].Message.Content)
	}
	if resp.Usage.TotalTokens != 15 {
		t.Errorf("expected 15 total tokens, got %d", resp.Usage.TotalTokens)
	}
}

func TestGoogleAdapter_Chat_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`internal server error`))
	}))
	defer server.Close()

	originalBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = originalBaseURL }()

	provider := New("test-key")
	ctx := context.Background()
	_, err := provider.Chat(ctx, "gemini-1.5", nil, contracts.ChatOptions{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGoogleAdapter_ChatStream(t *testing.T) {
	provider := New("test-key")
	_, err := provider.ChatStream(context.Background(), "gemini-1.5", nil, contracts.ChatOptions{})
	if err == nil {
		t.Fatal("expected error for unimplemented ChatStream")
	}
}
