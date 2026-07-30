package openai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ioriimasu/jervis/internal/aiprovider/contracts"
)

func TestOpenAIAdapter_Name(t *testing.T) {
	provider := New("test-key")
	if provider.Name() != "openai" {
		t.Errorf("expected 'openai', got %s", provider.Name())
	}
}

func TestOpenAIAdapter_Chat(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-key" {
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
			"choices": [{"message": {"role": "assistant", "content": "hello"}}],
			"usage": {"prompt_tokens": 10, "completion_tokens": 5, "total_tokens": 15}
		}`))
	}))
	defer server.Close()

	originalBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = originalBaseURL }()

	provider := New("test-key")
	ctx := context.Background()
	opts := contracts.ChatOptions{Temperature: 0.5, MaxTokens: 100}
	messages := []contracts.Message{{Role: contracts.RoleUser, Content: "hi"}}

	resp, err := provider.Chat(ctx, "gpt-4", messages, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Choices) != 1 {
		t.Fatalf("expected 1 choice, got %d", len(resp.Choices))
	}
	if resp.Choices[0].Message.Content != "hello" {
		t.Errorf("expected 'hello', got %s", resp.Choices[0].Message.Content)
	}
	if resp.Usage.TotalTokens != 15 {
		t.Errorf("expected 15 total tokens, got %d", resp.Usage.TotalTokens)
	}
}

func TestOpenAIAdapter_Chat_Error(t *testing.T) {
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
	_, err := provider.Chat(ctx, "gpt-4", nil, contracts.ChatOptions{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestOpenAIAdapter_ChatStream(t *testing.T) {
	provider := New("test-key")
	_, err := provider.ChatStream(context.Background(), "gpt-4", nil, contracts.ChatOptions{})
	if err == nil {
		t.Fatal("expected error for unimplemented ChatStream")
	}
}
