package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/saaedimam/jervis/internal/aiprovider/contracts"
)

type adapter struct {
	baseURL string
	client  *http.Client
}

// New creates a new Ollama (Local) provider adapter.
func New(baseURL string) contracts.Provider {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	return &adapter{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 120 * time.Second, // Local models can be slow
		},
	}
}

func (a *adapter) Name() string {
	return "ollama"
}

func (a *adapter) Chat(ctx context.Context, model string, messages []contracts.Message, opts contracts.ChatOptions) (*contracts.Response, error) {
	reqBody := map[string]any{
		"model":    model,
		"messages": messages,
		"stream":   false,
		"options": map[string]any{
			"temperature": opts.Temperature,
			"seed":        opts.Seed,
			"num_predict": opts.MaxTokens,
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", a.baseURL+"/api/chat", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama error (%d): %s", resp.StatusCode, string(body))
	}

	var ollamaResp struct {
		Message         contracts.Message `json:"message"`
		TotalDuration   int64             `json:"total_duration"`
		PromptEvalCount int               `json:"prompt_eval_count"`
		EvalCount       int               `json:"eval_count"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, err
	}

	return &contracts.Response{
		Choices: []contracts.Choice{
			{
				Message: ollamaResp.Message,
			},
		},
		Usage: contracts.Usage{
			PromptTokens:     ollamaResp.PromptEvalCount,
			CompletionTokens: ollamaResp.EvalCount,
			TotalTokens:      ollamaResp.PromptEvalCount + ollamaResp.EvalCount,
		},
	}, nil
}

func (a *adapter) ChatStream(ctx context.Context, model string, messages []contracts.Message, opts contracts.ChatOptions) (io.ReadCloser, error) {
	return nil, fmt.Errorf("streaming not yet implemented for ollama")
}
