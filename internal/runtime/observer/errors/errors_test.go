package errors_test

import (
	"errors"
	"testing"

	obserrors "github.com/ioriimasu/jervis/internal/runtime/observer/errors"
)

func TestErrors(t *testing.T) {
	if obserrors.ErrInvalidNotification == nil {
		t.Fatal("ErrInvalidNotification is nil")
	}
	if obserrors.ErrDuplicateObserver == nil {
		t.Fatal("ErrDuplicateObserver is nil")
	}
	if obserrors.ErrObserverNotFound == nil {
		t.Fatal("ErrObserverNotFound is nil")
	}
	if obserrors.ErrObserverFailure == nil {
		t.Fatal("ErrObserverFailure is nil")
	}
	if obserrors.ErrDispatchFailed == nil {
		t.Fatal("ErrDispatchFailed is nil")
	}
	if obserrors.ErrObserverPanic == nil {
		t.Fatal("ErrObserverPanic is nil")
	}
}

func TestAggregateError(t *testing.T) {
	var nilAgg *obserrors.AggregateError
	if nilAgg.Error() != "no observer errors" {
		t.Errorf("Expected 'no observer errors', got %s", nilAgg.Error())
	}
	if nilAgg.Errors() != nil {
		t.Errorf("Expected nil errors slice")
	}

	emptyAgg := obserrors.NewAggregateError([]error{})
	if emptyAgg != nil {
		t.Errorf("Expected nil for empty errors slice")
	}

	err1 := errors.New("err 1")
	err2 := errors.New("err 2")
	agg := obserrors.NewAggregateError([]error{err1, nil, err2})

	if agg == nil {
		t.Fatal("Expected non-nil AggregateError")
	}

	if len(agg.Errors()) != 2 {
		t.Fatalf("Expected 2 errors, got %d", len(agg.Errors()))
	}

	str := agg.Error()
	if str == "" {
		t.Error("Expected non-empty string representation")
	}
}
