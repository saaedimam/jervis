package middleware_test

import (
	"errors"
	"testing"

	"github.com/ioriimasu/jervis/internal/runtime/eventbus/contracts"
	"github.com/ioriimasu/jervis/internal/runtime/eventbus/middleware"
)

func TestFuncMiddlewareAdapter(t *testing.T) {
	called := false
	mwFunc := middleware.Func(func(evt contracts.Event, next func(contracts.Event) error) error {
		called = true
		return next(evt)
	})

	var _ contracts.Middleware = mwFunc

	terminalCalled := false
	err := mwFunc.Execute(nil, func(evt contracts.Event) error {
		terminalCalled = true
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatalf("expected middleware function to be called")
	}
	if !terminalCalled {
		t.Fatalf("expected terminal function to be called")
	}

	// Test Func returning error
	errFunc := middleware.Func(func(evt contracts.Event, next func(contracts.Event) error) error {
		return errors.New("middleware failure")
	})
	if err := errFunc.Execute(nil, nil); err == nil || err.Error() != "middleware failure" {
		t.Fatalf("expected 'middleware failure', got %v", err)
	}
}
