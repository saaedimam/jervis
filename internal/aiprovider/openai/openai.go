package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ioriimasu/jervis/internal/aiprovider/contracts"
)

var (
	baseURL = "https://api.openai.com/v1"
)

type adapter struct {
	apiKey string
	client *http.Client
}

// New creates a new OpenAI provider adapter.
func New(apiKey string) contracts.Provider {
	return &adapter{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (a *adapter) Name() string {
	return "openai"
}

func (a *adapter) Chat(ctx context.Context, model string, messages []contracts.Message, opts contracts.ChatOptions) (*contracts.Response, error) {
	reqBody := map[string]any{
		"model":    model,
		"messages": messages,
	}

	if opts.Temperature != 0 {
		reqBody["temperature"] = opts.Temperature
	}
	if opts.MaxTokens != 0 {
		reqBody["max_tokens"] = opts.MaxTokens
	}
	if opts.Stream {
		reqBody["stream"] = true
	}
	if opts.Seed != nil {
		reqBody["seed"] = *opts.Seed
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/chat/completions", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.apiKey))

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("openai error (%d): %s", resp.StatusCode, string(body))
	}

	var openAIResp struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Created int64  `json:"created"`
		Model   string `json:"model"`
		Choices []struct {
			Index   int `json:"index"`
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return nil, err
	}

	res := &contracts.Response{
		Usage: contracts.Usage{
			PromptTokens:     openAIResp.Usage.PromptTokens,
			CompletionTokens: openAIResp.Usage.CompletionTokens,
			TotalTokens:      openAIResp.Usage.TotalTokens,
		},
	}
	for _, c := range openAIResp.Choices {
		res.Choices = append(res.Choices, contracts.Choice{Message: contracts.Message{
			Role:    contracts.Role(c.Message.Role),
			Content: c.Message.Content,
		}})
	}

	return res, nil
}

func (a *adapter) ChatStream(ctx context.Context, model string, messages []contracts.Message, opts contracts.ChatOptions) (io.ReadCloser, error) {
	// TODO: Implement streaming
	return nil, fmt.Errorf("streaming not yet implemented for openai")
}
