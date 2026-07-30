package anthropic

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

var baseURL = "https://api.anthropic.com/v1"

const (
	apiVersion = "2023-06-01"
)

type adapter struct {
	apiKey string
	client *http.Client
}

// New creates a new Anthropic (Claude) provider adapter.
func New(apiKey string) contracts.Provider {
	return &adapter{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (a *adapter) Name() string {
	return "anthropic"
}

func (a *adapter) Chat(ctx context.Context, model string, messages []contracts.Message, opts contracts.ChatOptions) (*contracts.Response, error) {
	// Anthropic uses 'system' as a separate parameter, not part of the messages array.
	var systemPrompt string
	var anthropicMessages []map[string]string

	for _, msg := range messages {
		if msg.Role == contracts.RoleSystem {
			systemPrompt = msg.Content
		} else {
			anthropicMessages = append(anthropicMessages, map[string]string{
				"role":    string(msg.Role),
				"content": msg.Content,
			})
		}
	}

	reqBody := map[string]any{
		"model":      model,
		"messages":   anthropicMessages,
		"max_tokens": opts.MaxTokens,
		"stream":     false,
	}
	if systemPrompt != "" {
		reqBody["system"] = systemPrompt
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/messages", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", a.apiKey)
	req.Header.Set("anthropic-version", apiVersion)

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("anthropic error (%d): %s", resp.StatusCode, string(body))
	}

	var anthResp struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&anthResp); err != nil {
		return nil, err
	}

	content := ""
	if len(anthResp.Content) > 0 {
		content = anthResp.Content[0].Text
	}

	return &contracts.Response{
		Choices: []contracts.Choice{
			{
				Message: contracts.Message{
					Role:    contracts.RoleAssistant,
					Content: content,
				},
			},
		},
		Usage: contracts.Usage{
			PromptTokens:     anthResp.Usage.InputTokens,
			CompletionTokens: anthResp.Usage.OutputTokens,
			TotalTokens:      anthResp.Usage.InputTokens + anthResp.Usage.OutputTokens,
		},
	}, nil
}

func (a *adapter) ChatStream(ctx context.Context, model string, messages []contracts.Message, opts contracts.ChatOptions) (io.ReadCloser, error) {
	return nil, fmt.Errorf("streaming not yet implemented for anthropic")
}
