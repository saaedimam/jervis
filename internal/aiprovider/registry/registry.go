package registry

import (
	"errors"
	"sync"

	"github.com/ioriimasu/jervis/internal/aiprovider/contracts"
)

var ErrProviderNotFound = errors.New("aiprovider: provider not found")

// Registry manages a set of AI providers.
type Registry struct {
	mu        sync.RWMutex
	providers map[string]contracts.Provider
	defaultPr string
}

// New creates a new provider registry.
func New() *Registry {
	return &Registry{
		providers: make(map[string]contracts.Provider),
	}
}

// Register adds a provider to the registry.
func (r *Registry) Register(p contracts.Provider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[p.Name()] = p
	if r.defaultPr == "" {
		r.defaultPr = p.Name()
	}
}

// Get retrieves a provider by name.
func (r *Registry) Get(name string) (contracts.Provider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.providers[name]
	if !ok {
		return nil, ErrProviderNotFound
	}
	return p, nil
}

// SetDefault sets the default provider.
func (r *Registry) SetDefault(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.providers[name]; !ok {
		return ErrProviderNotFound
	}
	r.defaultPr = name
	return nil
}

// Default returns the default provider.
func (r *Registry) Default() (contracts.Provider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.defaultPr == "" {
		return nil, ErrProviderNotFound
	}
	return r.providers[r.defaultPr], nil
}
