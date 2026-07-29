package openai_test

import (
	"context"
	"testing"

	"github.com/ioriimasu/jervis/internal/aiprovider"
	"github.com/ioriimasu/jervis/internal/aiprovider/openai"
)

func TestOpenAIProvider(t *testing.T) {
	ctx := context.Background()
	p := openai.New("test-key")

	if p.Name() != "openai" {
		t.Errorf("expected provider name 'openai', got '%s'", p.Name())
	}

	_, err := p.Generate(ctx, "", aiprovider.Options{})
	if err != openai.ErrInvalidPrompt {
		t.Errorf("expected ErrInvalidPrompt, got %v", err)
	}

	res, err := p.Generate(ctx, "Hello AI", aiprovider.Options{Model: "gpt-4o"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.ProviderName != "openai" || res.Model != "gpt-4o" {
		t.Errorf("unexpected response fields: %+v", res)
	}

	// Test default model fallback
	resDefault, err := p.Generate(ctx, "Hello Default", aiprovider.Options{})
	if err != nil || resDefault.Model != "gpt-4o" {
		t.Errorf("expected default model gpt-4o, got %s", resDefault.Model)
	}
}
