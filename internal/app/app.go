package app

import (
	"context"
	"fmt"

	"github.com/ioriimasu/jervis/internal/aiprovider/anthropic"
	"github.com/ioriimasu/jervis/internal/aiprovider/google"
	"github.com/ioriimasu/jervis/internal/aiprovider/ollama"
	"github.com/ioriimasu/jervis/internal/aiprovider/openai"
	"github.com/ioriimasu/jervis/internal/aiprovider/prompts"
	"github.com/ioriimasu/jervis/internal/aiprovider/registry"
	"github.com/ioriimasu/jervis/internal/memory/store/sqlite"
	"github.com/ioriimasu/jervis/internal/memory/timeline"
	"github.com/ioriimasu/jervis/internal/memory/working"
	"github.com/ioriimasu/jervis/internal/runtime/eventbus"
	"github.com/ioriimasu/jervis/internal/services"
)

// App is the central container for the Jervis Personal OS.
type App struct {
	EventBus    *eventbus.EventBus
	Store       *sqlite.Driver
	Memory      *working.WorkingMemory
	Timeline    *timeline.Timeline
	Services    *services.Container
	AIProviders *registry.Registry
	Prompts     *prompts.Engine
}

// Config defines the application configuration.
type Config struct {
	DatabasePath  string
	MemoryLimit   int
	NotionToken   string
	OpenAIKey     string
	AnthropicKey  string
	GoogleKey     string
	OllamaBaseURL string
	DefaultAIProv string
	PromptDir     string
}

// New initializes the entire Jervis stack.
func New(ctx context.Context, cfg Config) (*App, error) {
	// 1. Initialize EventBus
	bus := eventbus.New()

	// 2. Initialize Store
	store, err := sqlite.New(cfg.DatabasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to init store: %w", err)
	}

	// 3. Initialize Memory Engine
	mem := working.New(cfg.MemoryLimit)

	// 4. Initialize Timeline (Ledger)
	tl := timeline.New(store)

	// 5. Initialize Services
	container, err := services.NewContainer(ctx, store, bus, cfg.NotionToken)
	if err != nil {
		_ = store.Close()
		return nil, fmt.Errorf("failed to init services: %w", err)
	}

	// 6. Initialize AI Providers
	aiReg := registry.New()
	if cfg.OpenAIKey != "" {
		aiReg.Register(openai.New(cfg.OpenAIKey))
	}
	if cfg.AnthropicKey != "" {
		aiReg.Register(anthropic.New(cfg.AnthropicKey))
	}
	if cfg.GoogleKey != "" {
		aiReg.Register(google.New(cfg.GoogleKey))
	}
	// Always register Ollama if it's running locally
	aiReg.Register(ollama.New(cfg.OllamaBaseURL))

	if cfg.DefaultAIProv != "" {
		_ = aiReg.SetDefault(cfg.DefaultAIProv)
	}

	// 7. Initialize Prompt Engine
	promptDir := cfg.PromptDir
	if promptDir == "" {
		promptDir = "prompts"
	}
	pe := prompts.NewEngine(promptDir)

	return &App{
		EventBus:    bus,
		Store:       store,
		Memory:      mem,
		Timeline:    tl,
		Services:    container,
		AIProviders: aiReg,
		Prompts:     pe,
	}, nil
}

// Close gracefully shuts down all components.
func (a *App) Close() error {
	if a.Store != nil {
		return a.Store.Close()
	}
	return nil
}

// DefaultConfig returns the default configuration.
func DefaultConfig() Config {
	return Config{
		DatabasePath: "jervis.db",
		MemoryLimit:  1000,
	}
}
