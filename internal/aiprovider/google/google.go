package google

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
	baseURL = "https://generativelanguage.googleapis.com/v1beta/models"
)

type adapter struct {
	apiKey string
	client *http.Client
}

// New creates a new Google (Gemini) provider adapter.
func New(apiKey string) contracts.Provider {
	return &adapter{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (a *adapter) Name() string {
	return "google"
}

func (a *adapter) Chat(ctx context.Context, model string, messages []contracts.Message, opts contracts.ChatOptions) (*contracts.Response, error) {
	// Gemini uses a different structure: {"contents": [{"role": "user", "parts": [{"text": "..."}]}]}
	// System prompt is handled separately in some versions, or as a special role.

	type Part struct {
		Text string `json:"text"`
	}
	type Content struct {
		Role  string `json:"role"`
		Parts []Part `json:"parts"`
	}

	var contents []Content
	var systemInstruction *Content

	for _, msg := range messages {
		role := string(msg.Role)
		if msg.Role == contracts.RoleAssistant {
			role = "model"
		}

		if msg.Role == contracts.RoleSystem {
			systemInstruction = &Content{
				Parts: []Part{{Text: msg.Content}},
			}
		} else {
			contents = append(contents, Content{
				Role:  role,
				Parts: []Part{{Text: msg.Content}},
			})
		}
	}

	reqBody := map[string]any{
		"contents": contents,
		"generationConfig": map[string]any{
			"temperature":     opts.Temperature,
			"maxOutputTokens": opts.MaxTokens,
		},
	}
	if systemInstruction != nil {
		reqBody["system_instruction"] = systemInstruction
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/%s:generateContent?key=%s", baseURL, model, a.apiKey)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("google error (%d): %s", resp.StatusCode, string(body))
	}

	var geminiResp struct {
		Candidates []struct {
			Content Content `json:"content"`
		} `json:"candidates"`
		UsageMetadata struct {
			PromptTokenCount     int `json:"promptTokenCount"`
			CandidatesTokenCount int `json:"candidatesTokenCount"`
			TotalTokenCount      int `json:"totalTokenCount"`
		} `json:"usageMetadata"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return nil, err
	}

	content := ""
	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		content = geminiResp.Candidates[0].Content.Parts[0].Text
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
			PromptTokens:     geminiResp.UsageMetadata.PromptTokenCount,
			CompletionTokens: geminiResp.UsageMetadata.CandidatesTokenCount,
			TotalTokens:      geminiResp.UsageMetadata.TotalTokenCount,
		},
	}, nil
}

func (a *adapter) ChatStream(ctx context.Context, model string, messages []contracts.Message, opts contracts.ChatOptions) (io.ReadCloser, error) {
	return nil, fmt.Errorf("streaming not yet implemented for google")
}
