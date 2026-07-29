package claude

import (
	"context"
	"errors"

	"github.com/ioriimasu/jervis/internal/aiprovider"
)

var ErrInvalidPrompt = errors.New("claude: prompt cannot be empty")

type Provider struct {
	apiKey string
}

func New(apiKey string) aiprovider.Provider {
	return &Provider{apiKey: apiKey}
}

func (p *Provider) Name() string {
	return "claude"
}

func (p *Provider) Generate(ctx context.Context, prompt string, opts aiprovider.Options) (*aiprovider.Response, error) {
	if prompt == "" {
		return nil, ErrInvalidPrompt
	}

	model := opts.Model
	if model == "" {
		model = "claude-3-5-sonnet"
	}

	return &aiprovider.Response{
		Content:      "Mock Claude response for prompt: " + prompt,
		ProviderName: "claude",
		Model:        model,
	}, nil
}
