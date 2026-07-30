package middleware_test

import (
	"errors"
	"testing"

	"github.com/saaedimam/jervis/internal/runtime/eventbus/contracts"
	errs "github.com/saaedimam/jervis/internal/runtime/eventbus/errors"
	"github.com/saaedimam/jervis/internal/runtime/eventbus/events"
	"github.com/saaedimam/jervis/internal/runtime/eventbus/middleware"
	"github.com/saaedimam/jervis/internal/runtime/types"
)

func createTestEvent(t *testing.T) contracts.Event {
	evtID, _ := types.NewEventID("evt-001")
	env, err := events.NewBuilder().
		SetID(evtID).
		SetType("runtime.test.event").
		SetSource("test").
		SetPayload("payload").
		Build()
	if err != nil {
		t.Fatalf("failed to build test event: %v", err)
	}
	return env
}

func TestChainBasicOperations(t *testing.T) {
	m1 := middleware.Func(func(evt contracts.Event, next func(contracts.Event) error) error { return next(evt) })
	m2 := middleware.Func(func(evt contracts.Event, next func(contracts.Event) error) error { return next(evt) })

	chain := middleware.NewChain(m1, nil, m2)

	if chain.Count() != 2 {
		t.Fatalf("expected count 2, got %d", chain.Count())
	}

	mws := chain.Middlewares()
	if len(mws) != 2 {
		t.Fatalf("expected 2 middlewares in slice")
	}

	// Test defensive copy
	mws[0] = nil
	if chain.Middlewares()[0] == nil {
		t.Fatalf("Middlewares() did not return a defensive copy")
	}

	emptyChain := middleware.NewChain()
	if emptyChain.Count() != 0 {
		t.Fatalf("expected 0 count for empty chain")
	}
	if emptyChain.Middlewares() != nil {
		t.Fatalf("expected nil slice for empty chain Middlewares()")
	}
}

func TestChainFIFOPreLIFOPostExecution(t *testing.T) {
	var sequence []string

	m1 := middleware.Func(func(evt contracts.Event, next func(contracts.Event) error) error {
		sequence = append(sequence, "M1-Pre")
		err := next(evt)
		sequence = append(sequence, "M1-Post")
		return err
	})

	m2 := middleware.Func(func(evt contracts.Event, next func(contracts.Event) error) error {
		sequence = append(sequence, "M2-Pre")
		err := next(evt)
		sequence = append(sequence, "M2-Post")
		return err
	})

	chain := middleware.NewChain(m1, m2)
	evt := createTestEvent(t)

	err := chain.Execute(evt, func(e contracts.Event) error {
		sequence = append(sequence, "Terminal")
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedSeq := []string{"M1-Pre", "M2-Pre", "Terminal", "M2-Post", "M1-Post"}
	if len(sequence) != len(expectedSeq) {
		t.Fatalf("sequence length mismatch: got %v, want %v", sequence, expectedSeq)
	}
	for i, s := range expectedSeq {
		if sequence[i] != s {
			t.Errorf("position %d: got %s, want %s", i, sequence[i], s)
		}
	}
}

func TestChainShortCircuit(t *testing.T) {
	var sequence []string
	shortCircuitErr := errors.New("unauthorized")

	m1 := middleware.Func(func(evt contracts.Event, next func(contracts.Event) error) error {
		sequence = append(sequence, "M1-Pre")
		err := next(evt)
		sequence = append(sequence, "M1-Post")
		return err
	})

	m2ShortCircuit := middleware.Func(func(evt contracts.Event, next func(contracts.Event) error) error {
		sequence = append(sequence, "M2-ShortCircuit")
		return shortCircuitErr
	})

	m3 := middleware.Func(func(evt contracts.Event, next func(contracts.Event) error) error {
		sequence = append(sequence, "M3-Pre")
		return next(evt)
	})

	chain := middleware.NewChain(m1, m2ShortCircuit, m3)
	evt := createTestEvent(t)

	err := chain.Execute(evt, func(e contracts.Event) error {
		sequence = append(sequence, "Terminal")
		return nil
	})

	if !errors.Is(err, shortCircuitErr) {
		t.Fatalf("expected shortCircuitErr, got %v", err)
	}

	expectedSeq := []string{"M1-Pre", "M2-ShortCircuit", "M1-Post"}
	if len(sequence) != len(expectedSeq) {
		t.Fatalf("sequence length mismatch: got %v, want %v", sequence, expectedSeq)
	}
	for i, s := range expectedSeq {
		if sequence[i] != s {
			t.Errorf("position %d: got %s, want %s", i, sequence[i], s)
		}
	}
}

func TestChainPanicRecovery(t *testing.T) {
	mPanic := middleware.Func(func(evt contracts.Event, next func(contracts.Event) error) error {
		panic("middleware boom")
	})

	chain := middleware.NewChain(mPanic)
	evt := createTestEvent(t)

	err := chain.Execute(evt, func(e contracts.Event) error {
		t.Fatalf("terminal should not be called on middleware panic")
		return nil
	})

	if err == nil {
		t.Fatalf("expected error from middleware panic")
	}
	if !errors.Is(err, errs.ErrHandlerFailure) {
		t.Fatalf("expected ErrHandlerFailure for middleware panic, got %v", err)
	}
}

func TestChainNilTerminalValidation(t *testing.T) {
	chain := middleware.NewChain()
	evt := createTestEvent(t)

	err := chain.Execute(evt, nil)
	if !errors.Is(err, errs.ErrValidationFailed) {
		t.Fatalf("expected ErrValidationFailed for nil terminal function")
	}
}
