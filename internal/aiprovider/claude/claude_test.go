package claude_test

import (
	"context"
	"testing"

	"github.com/ioriimasu/jervis/internal/aiprovider"
	"github.com/ioriimasu/jervis/internal/aiprovider/claude"
)

func TestClaudeProvider(t *testing.T) {
	ctx := context.Background()
	p := claude.New("test-key")

	if p.Name() != "claude" {
		t.Errorf("expected provider name 'claude', got '%s'", p.Name())
	}

	_, err := p.Generate(ctx, "", aiprovider.Options{})
	if err != claude.ErrInvalidPrompt {
		t.Errorf("expected ErrInvalidPrompt, got %v", err)
	}

	res, err := p.Generate(ctx, "Hello Claude", aiprovider.Options{Model: "claude-3-5-sonnet"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.ProviderName != "claude" || res.Model != "claude-3-5-sonnet" {
		t.Errorf("unexpected response fields: %+v", res)
	}

	// Test default model fallback
	resDefault, err := p.Generate(ctx, "Hello Default", aiprovider.Options{})
	if err != nil || resDefault.Model != "claude-3-5-sonnet" {
		t.Errorf("expected default model claude-3-5-sonnet, got %s", resDefault.Model)
	}
}
