package registry_test

import (
	"errors"
	"testing"

	"github.com/ioriimasu/jervis/internal/runtime/eventbus/contracts"
	errs "github.com/ioriimasu/jervis/internal/runtime/eventbus/errors"
	"github.com/ioriimasu/jervis/internal/runtime/eventbus/events"
	"github.com/ioriimasu/jervis/internal/runtime/eventbus/registry"
	"github.com/ioriimasu/jervis/internal/runtime/eventbus/subscription"
)

type mockHandler struct {
	id string
}

func (m *mockHandler) ID() string                         { return m.id }
func (m *mockHandler) Handle(event contracts.Event) error { return nil }

func TestValidatePattern(t *testing.T) {
	validPatterns := []string{
		"*",
		"runtime.started",
		"runtime.*",
		"memory.timeline.*",
		"a.b",
	}
	for _, p := range validPatterns {
		if err := registry.ValidatePattern(p); err != nil {
			t.Errorf("expected pattern %q to be valid, got %v", p, err)
		}
	}

	invalidPatterns := []string{
		"",
		"INVALID",
		"runtime started",
		"runtime..started",
		".runtime.started",
	}
	for _, p := range invalidPatterns {
		if err := registry.ValidatePattern(p); !errors.Is(err, errs.ErrValidationFailed) {
			t.Errorf("expected ErrValidationFailed for pattern %q, got %v", p, err)
		}
	}
}

func TestMatchesPattern(t *testing.T) {
	tests := []struct {
		eventType string
		pattern   string
		match     bool
	}{
		{"runtime.started", "runtime.started", true},
		{"runtime.started", "runtime.stopped", false},
		{"runtime.started", "runtime.*", true},
		{"runtime.lifecycle.booted", "runtime.*", true},
		{"memory.appended", "runtime.*", false},
		{"any.event.type", "*", true},
		{"runtime.started", "runtime*", true},
	}

	for _, tt := range tests {
		got := registry.MatchesPattern(tt.eventType, tt.pattern)
		if got != tt.match {
			t.Errorf("MatchesPattern(%q, %q) = %v, want %v", tt.eventType, tt.pattern, got, tt.match)
		}
	}
}

func TestRegistryRegisterAndLookup(t *testing.T) {
	reg := registry.NewRegistry()
	h1 := &mockHandler{id: "h-1"}
	h2 := &mockHandler{id: "h-2"}
	h3 := &mockHandler{id: "h-3"}

	sub1, _ := subscription.New("sub-1", "runtime.started", events.PriorityNormal, h1)
	sub2, _ := subscription.New("sub-2", "runtime.*", events.PriorityHigh, h2)
	sub3, _ := subscription.New("sub-3", "*", events.PriorityCritical, h3)

	if err := reg.Register(sub1); err != nil {
		t.Fatalf("register sub1 failed: %v", err)
	}
	if err := reg.Register(sub2); err != nil {
		t.Fatalf("register sub2 failed: %v", err)
	}
	if err := reg.Register(sub3); err != nil {
		t.Fatalf("register sub3 failed: %v", err)
	}

	if reg.Count() != 3 {
		t.Fatalf("expected count 3, got %d", reg.Count())
	}
	if !reg.Contains("sub-1") || !reg.Contains("sub-2") || !reg.Contains("sub-3") {
		t.Fatalf("expected registry to contain all subscriptions")
	}

	// Lookup matches: sub3 (PriorityCritical), sub2 (PriorityHigh), sub1 (PriorityNormal)
	matches := reg.Lookup("runtime.started")
	if len(matches) != 3 {
		t.Fatalf("expected 3 matches for runtime.started, got %d", len(matches))
	}
	if matches[0].ID() != "sub-3" || matches[1].ID() != "sub-2" || matches[2].ID() != "sub-1" {
		t.Errorf("expected deterministic priority ordering [sub-3, sub-2, sub-1], got [%s, %s, %s]",
			matches[0].ID(), matches[1].ID(), matches[2].ID())
	}

	// LookupExact
	exactMatches := reg.LookupExact("runtime.started")
	if len(exactMatches) != 1 || exactMatches[0].ID() != "sub-1" {
		t.Fatalf("LookupExact mismatch")
	}

	// LookupPattern
	patternMatches := reg.LookupPattern("runtime.*")
	if len(patternMatches) != 1 || patternMatches[0].ID() != "sub-2" {
		t.Fatalf("LookupPattern mismatch")
	}
}

func TestRegistryOrderingStabilitySamePriority(t *testing.T) {
	reg := registry.NewRegistry()
	h1 := &mockHandler{id: "h-1"}
	h2 := &mockHandler{id: "h-2"}

	sub1, _ := subscription.New("sub-1", "runtime.started", events.PriorityNormal, h1)
	sub2, _ := subscription.New("sub-2", "runtime.started", events.PriorityNormal, h2)

	_ = reg.Register(sub1)
	_ = reg.Register(sub2)

	matches := reg.Lookup("runtime.started")
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches")
	}
	if matches[0].ID() != "sub-1" || matches[1].ID() != "sub-2" {
		t.Errorf("expected stable FIFO ordering for same priority, got [%s, %s]", matches[0].ID(), matches[1].ID())
	}
}

func TestRegistryDuplicateDetectionsAndValidation(t *testing.T) {
	reg := registry.NewRegistry()
	h1 := &mockHandler{id: "h-1"}

	sub1, _ := subscription.New("sub-1", "runtime.started", events.PriorityNormal, h1)
	if err := reg.Register(sub1); err != nil {
		t.Fatalf("initial register failed: %v", err)
	}

	// Duplicate SubscriptionID
	subDupID, _ := subscription.New("sub-1", "other.event", events.PriorityNormal, &mockHandler{id: "h-other"})
	if err := reg.Register(subDupID); !errors.Is(err, errs.ErrDuplicateSubscriber) {
		t.Fatalf("expected ErrDuplicateSubscriber for duplicate sub ID")
	}

	// Duplicate Handler ID for same pattern
	subDupHandler, _ := subscription.New("sub-2", "runtime.started", events.PriorityHigh, h1)
	if err := reg.Register(subDupHandler); !errors.Is(err, errs.ErrDuplicateSubscriber) {
		t.Fatalf("expected ErrDuplicateSubscriber for duplicate handler on same pattern")
	}

	// Invalid pattern in sub
	subInvalidPattern, _ := subscription.New("sub-invalid", "INVALID PATTERN", events.PriorityNormal, h1)
	if err := reg.Register(subInvalidPattern); !errors.Is(err, errs.ErrValidationFailed) {
		t.Fatalf("expected ErrValidationFailed for invalid pattern")
	}
}

func TestRegistryUnregisterAndClear(t *testing.T) {
	reg := registry.NewRegistry()
	h1 := &mockHandler{id: "h-1"}
	sub1, _ := subscription.New("sub-1", "runtime.started", events.PriorityNormal, h1)

	_ = reg.Register(sub1)

	if err := reg.Unregister(""); !errors.Is(err, errs.ErrValidationFailed) {
		t.Fatalf("expected ErrValidationFailed for empty subID on unregister")
	}
	if err := reg.Unregister("non-existent"); !errors.Is(err, errs.ErrValidationFailed) {
		t.Fatalf("expected ErrValidationFailed for non-existent subID")
	}

	if err := reg.Unregister("sub-1"); err != nil {
		t.Fatalf("unregister sub-1 failed: %v", err)
	}
	if reg.Contains("sub-1") || reg.Count() != 0 {
		t.Fatalf("expected count 0 after unregister")
	}

	// Snapshot and Clear
	_ = reg.Register(sub1)
	snap := reg.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected snapshot size 1")
	}

	reg.Clear()
	if reg.Count() != 0 {
		t.Fatalf("expected count 0 after clear")
	}
}
