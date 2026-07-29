package ollama

import (
	"context"
	"errors"

	"github.com/ioriimasu/jervis/internal/aiprovider"
)

var ErrInvalidPrompt = errors.New("ollama: prompt cannot be empty")

type Provider struct {
	baseURL string
}

func New(baseURL string) aiprovider.Provider {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	return &Provider{baseURL: baseURL}
}

func (p *Provider) Name() string {
	return "ollama"
}

func (p *Provider) Generate(ctx context.Context, prompt string, opts aiprovider.Options) (*aiprovider.Response, error) {
	if prompt == "" {
		return nil, ErrInvalidPrompt
	}

	model := opts.Model
	if model == "" {
		model = "llama3"
	}

	return &aiprovider.Response{
		Content:      "Mock Ollama response for prompt: " + prompt,
		ProviderName: "ollama",
		Model:        model,
	}, nil
}
