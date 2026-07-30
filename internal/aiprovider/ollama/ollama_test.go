package ollama

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ioriimasu/jervis/internal/aiprovider/contracts"
)

func TestOllamaAdapter_Name(t *testing.T) {
	provider := New("")
	if provider.Name() != "ollama" {
		t.Errorf("expected 'ollama', got %s", provider.Name())
	}
}

func TestOllamaAdapter_Chat(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"message": {"role": "assistant", "content": "hello ollama"},
			"prompt_eval_count": 10,
			"eval_count": 5
		}`))
	}))
	defer server.Close()

	provider := New(server.URL)
	ctx := context.Background()
	opts := contracts.ChatOptions{Temperature: 0.5, MaxTokens: 100}
	messages := []contracts.Message{
		{Role: contracts.RoleUser, Content: "hi"},
	}

	resp, err := provider.Chat(ctx, "llama3", messages, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Choices) != 1 {
		t.Fatalf("expected 1 choice, got %d", len(resp.Choices))
	}
	if resp.Choices[0].Message.Content != "hello ollama" {
		t.Errorf("expected 'hello ollama', got %s", resp.Choices[0].Message.Content)
	}
	if resp.Usage.TotalTokens != 15 {
		t.Errorf("expected 15 total tokens, got %d", resp.Usage.TotalTokens)
	}
}

func TestOllamaAdapter_Chat_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`internal server error`))
	}))
	defer server.Close()

	provider := New(server.URL)
	ctx := context.Background()
	_, err := provider.Chat(ctx, "llama3", nil, contracts.ChatOptions{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestOllamaAdapter_ChatStream(t *testing.T) {
	provider := New("")
	_, err := provider.ChatStream(context.Background(), "llama3", nil, contracts.ChatOptions{})
	if err == nil {
		t.Fatal("expected error for unimplemented ChatStream")
	}
}
