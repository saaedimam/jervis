package openai

import (
	"context"
	"errors"

	"github.com/ioriimasu/jervis/internal/aiprovider"
)

var ErrInvalidPrompt = errors.New("openai: prompt cannot be empty")

type Provider struct {
	apiKey string
}

func New(apiKey string) aiprovider.Provider {
	return &Provider{apiKey: apiKey}
}

func (p *Provider) Name() string {
	return "openai"
}

func (p *Provider) Generate(ctx context.Context, prompt string, opts aiprovider.Options) (*aiprovider.Response, error) {
	if prompt == "" {
		return nil, ErrInvalidPrompt
	}

	model := opts.Model
	if model == "" {
		model = "gpt-4o"
	}

	return &aiprovider.Response{
		Content:      "Mock OpenAI response for prompt: " + prompt,
		ProviderName: "openai",
		Model:        model,
	}, nil
}
