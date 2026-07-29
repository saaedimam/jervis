package aiprovider

import (
	"context"
)

// Request options for AI generation.
type Options struct {
	Model       string  `json:"model"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
}

// Response represents normalized output from an AI Provider.
type Response struct {
	Content      string `json:"content"`
	ProviderName string `json:"provider_name"`
	Model        string `json:"model"`
}

// Provider defines the clean, decoupled interface for AI providers.
type Provider interface {
	Name() string
	Generate(ctx context.Context, prompt string, opts Options) (*Response, error)
}
