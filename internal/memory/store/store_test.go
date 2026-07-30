//go:build !sqlite

package store_test

import (
	"testing"

	"github.com/saaedimam/jervis/internal/memory/store"
	"github.com/saaedimam/jervis/internal/memory/store/contracts"
)

// TestStoreNew verifies that Store.New returns a valid SQLite store.
func TestStoreNew(t *testing.T) {
	s, err := store.New(":memory:")
	if err != nil {
		t.Fatalf("Expected no error creating store, got: %v", err)
	}
	if s == nil {
		t.Fatal("Expected non-nil store")
	}

	// Verify it implements the Store interface
	var _ contracts.Store = s

	// Clean up
	if err := s.Close(); err != nil {
		t.Errorf("Error closing store: %v", err)
	}
}

// TestStoreNewMemory verifies that Store.NewMemory returns a valid in-memory store.
func TestStoreNewMemory(t *testing.T) {
	s, err := store.NewMemory()
	if err != nil {
		t.Fatalf("Expected no error creating memory store, got: %v", err)
	}
	if s == nil {
		t.Fatal("Expected non-nil store")
	}

	// Verify it implements the Store interface
	var _ contracts.Store = s

	// Clean up
	if err := s.Close(); err != nil {
		t.Errorf("Error closing store: %v", err)
	}
}
