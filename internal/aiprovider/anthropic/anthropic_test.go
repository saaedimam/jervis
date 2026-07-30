package anthropic

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ioriimasu/jervis/internal/aiprovider/contracts"
)

func TestAnthropicAdapter_Name(t *testing.T) {
	provider := New("test-key")
	if provider.Name() != "anthropic" {
		t.Errorf("expected 'anthropic', got %s", provider.Name())
	}
}

func TestAnthropicAdapter_Chat(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-api-key") != "test-key" {
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
			"content": [{"type": "text", "text": "hello anthropic"}],
			"usage": {"input_tokens": 10, "output_tokens": 5}
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
	}

	resp, err := provider.Chat(ctx, "claude-3-opus", messages, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Choices) != 1 {
		t.Fatalf("expected 1 choice, got %d", len(resp.Choices))
	}
	if resp.Choices[0].Message.Content != "hello anthropic" {
		t.Errorf("expected 'hello anthropic', got %s", resp.Choices[0].Message.Content)
	}
	if resp.Usage.TotalTokens != 15 {
		t.Errorf("expected 15 total tokens, got %d", resp.Usage.TotalTokens)
	}
}

func TestAnthropicAdapter_Chat_Error(t *testing.T) {
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
	_, err := provider.Chat(ctx, "claude-3", nil, contracts.ChatOptions{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAnthropicAdapter_ChatStream(t *testing.T) {
	provider := New("test-key")
	_, err := provider.ChatStream(context.Background(), "claude-3", nil, contracts.ChatOptions{})
	if err == nil {
		t.Fatal("expected error for unimplemented ChatStream")
	}
}
