package registry

import (
	"context"
	"io"
	"testing"

	"github.com/ioriimasu/jervis/internal/aiprovider/contracts"
)

type mockProvider struct {
	name string
}

func (m *mockProvider) Name() string { return m.name }
func (m *mockProvider) Chat(ctx context.Context, model string, messages []contracts.Message, opts contracts.ChatOptions) (*contracts.Response, error) {
	return nil, nil
}
func (m *mockProvider) ChatStream(ctx context.Context, model string, messages []contracts.Message, opts contracts.ChatOptions) (io.ReadCloser, error) {
	return nil, nil
}

func TestRegistry(t *testing.T) {
	reg := New()

	// Initial default should be empty/error
	if _, err := reg.Default(); err != ErrProviderNotFound {
		t.Errorf("expected ErrProviderNotFound, got %v", err)
	}
	if _, err := reg.Get("foo"); err != ErrProviderNotFound {
		t.Errorf("expected ErrProviderNotFound, got %v", err)
	}

	p1 := &mockProvider{name: "p1"}
	reg.Register(p1)

	// First registered should be default
	d, err := reg.Default()
	if err != nil || d.Name() != "p1" {
		t.Errorf("expected default 'p1', got err=%v name=%v", err, d)
	}

	p2 := &mockProvider{name: "p2"}
	reg.Register(p2)

	// Get specific
	got, err := reg.Get("p2")
	if err != nil || got.Name() != "p2" {
		t.Errorf("expected 'p2', got err=%v name=%v", err, got)
	}

	// SetDefault
	if err := reg.SetDefault("p2"); err != nil {
		t.Errorf("unexpected error setting default: %v", err)
	}
	d2, _ := reg.Default()
	if d2.Name() != "p2" {
		t.Errorf("expected default 'p2', got %v", d2.Name())
	}

	// SetDefault invalid
	if err := reg.SetDefault("invalid"); err != ErrProviderNotFound {
		t.Errorf("expected ErrProviderNotFound, got %v", err)
	}
}
