package dispatcher_test

import (
	"errors"
	"fmt"
	"testing"

	errs "github.com/saaedimam/jervis/internal/runtime/eventbus/errors"

	"github.com/saaedimam/jervis/internal/runtime/eventbus/dispatcher"
)

func TestAggregateErrorEmpty(t *testing.T) {
	agg := dispatcher.NewAggregateError(nil)
	if agg.HasErrors() {
		t.Fatalf("expected HasErrors() to be false")
	}
	if agg.Count() != 0 {
		t.Fatalf("expected Count() to be 0")
	}
	if agg.Errors() != nil {
		t.Fatalf("expected Errors() to be nil")
	}
	if agg.Error() != "no errors" {
		t.Fatalf("expected 'no errors', got %q", agg.Error())
	}
	if agg.Unwrap() != nil {
		t.Fatalf("expected Unwrap() to be nil")
	}

	agg.Add(nil)
	if agg.HasErrors() {
		t.Fatalf("expected HasErrors() to remain false after adding nil error")
	}
}

func TestAggregateErrorSingle(t *testing.T) {
	err1 := errors.New("handler 1 failed")
	agg := dispatcher.NewAggregateError([]error{err1, nil})

	if !agg.HasErrors() {
		t.Fatalf("expected HasErrors() to be true")
	}
	if agg.Count() != 1 {
		t.Fatalf("expected Count() to be 1")
	}
	if agg.Error() != "handler 1 failed" {
		t.Fatalf("expected 'handler 1 failed', got %q", agg.Error())
	}

	unwrapped := agg.Unwrap()
	if len(unwrapped) != 1 || unwrapped[0] != err1 {
		t.Fatalf("Unwrap() mismatch")
	}
}

func TestAggregateErrorMultiple(t *testing.T) {
	err1 := fmt.Errorf("%w: handler h1 failed", errs.ErrHandlerFailure)
	err2 := fmt.Errorf("%w: handler h2 panicked", errs.ErrHandlerFailure)

	agg := &dispatcher.AggregateError{}
	agg.Add(err1)
	agg.Add(err2)

	if !agg.HasErrors() {
		t.Fatalf("expected HasErrors() to be true")
	}
	if agg.Count() != 2 {
		t.Fatalf("expected Count() to be 2")
	}

	errStr := agg.Error()
	expected := "dispatch failed with 2 error(s): eventbus: handler execution failed: handler h1 failed; eventbus: handler execution failed: handler h2 panicked"
	if errStr != expected {
		t.Fatalf("unexpected formatted error string:\n got: %q\nwant: %q", errStr, expected)
	}

	if !errors.Is(agg, errs.ErrHandlerFailure) {
		t.Fatalf("errors.Is failed to find ErrHandlerFailure in unwrapped aggregate error")
	}

	// Test defensive copy
	errsSlice := agg.Errors()
	errsSlice[0] = nil
	if agg.Errors()[0] == nil {
		t.Fatalf("Errors() did not return a defensive copy")
	}
}
