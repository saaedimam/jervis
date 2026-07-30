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

const (
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
		"model":       model,
		"messages":    messages,
		"temperature": opts.Temperature,
		"max_tokens":  opts.MaxTokens,
		"stream":      false,
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
	req.Header.Set("Authorization", "Bearer "+a.apiKey)

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
		Choices []struct {
			Message contracts.Message `json:"message"`
		} `json:"choices"`
		Usage contracts.Usage `json:"usage"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return nil, err
	}

	res := &contracts.Response{
		Usage: openAIResp.Usage,
	}
	for _, c := range openAIResp.Choices {
		res.Choices = append(res.Choices, contracts.Choice{Message: c.Message})
	}

	return res, nil
}

func (a *adapter) ChatStream(ctx context.Context, model string, messages []contracts.Message, opts contracts.ChatOptions) (io.ReadCloser, error) {
	// TODO: Implement streaming
	return nil, fmt.Errorf("streaming not yet implemented for openai")
}
