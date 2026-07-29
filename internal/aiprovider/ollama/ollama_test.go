package ollama_test

import (
	"context"
	"testing"

	"github.com/ioriimasu/jervis/internal/aiprovider"
	"github.com/ioriimasu/jervis/internal/aiprovider/ollama"
)

func TestOllamaProvider(t *testing.T) {
	ctx := context.Background()
	p := ollama.New("")

	if p.Name() != "ollama" {
		t.Errorf("expected provider name 'ollama', got '%s'", p.Name())
	}

	pCustom := ollama.New("http://localhost:11434")
	if pCustom.Name() != "ollama" {
		t.Errorf("expected provider name 'ollama'")
	}

	_, err := p.Generate(ctx, "", aiprovider.Options{})
	if err != ollama.ErrInvalidPrompt {
		t.Errorf("expected ErrInvalidPrompt, got %v", err)
	}

	res, err := p.Generate(ctx, "Hello Ollama", aiprovider.Options{Model: "llama3"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.ProviderName != "ollama" || res.Model != "llama3" {
		t.Errorf("unexpected response fields: %+v", res)
	}

	// Test default model fallback
	resDefault, err := p.Generate(ctx, "Hello Default", aiprovider.Options{})
	if err != nil || resDefault.Model != "llama3" {
		t.Errorf("expected default model llama3, got %s", resDefault.Model)
	}
}
