package errors_test

import (
	"errors"
	"strings"
	"testing"

	observererrors "github.com/ioriimasu/jervis/internal/runtime/observer/errors"
)

func TestErrObserverPanic(t *testing.T) {
	err := &observererrors.ErrObserverPanic{
		ObserverID: "obs-1",
		Recovered:  "something went wrong",
	}

	expected := "observer [obs-1] panicked: something went wrong"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestAggregateError(t *testing.T) {
	t.Run("nil errs", func(t *testing.T) {
		err := observererrors.NewAggregateError(nil)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		err := observererrors.NewAggregateError([]error{})
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("multiple errors", func(t *testing.T) {
		errs := []error{
			errors.New("error 1"),
			&observererrors.ErrObserverPanic{ObserverID: "obs-2", Recovered: "panic!"},
		}
		agg := observererrors.NewAggregateError(errs)
		if agg == nil {
			t.Fatal("expected aggregate error, got nil")
		}

		msg := agg.Error()
		if !strings.Contains(msg, "2 observer error(s)") {
			t.Errorf("missing error count in message: %q", msg)
		}
		if !strings.Contains(msg, "error 1") {
			t.Errorf("missing error 1 in message: %q", msg)
		}
		if !strings.Contains(msg, "obs-2") {
			t.Errorf("missing observer ID in message: %q", msg)
		}

		if len(agg.Errors()) != 2 {
			t.Errorf("expected 2 errors, got %d", len(agg.Errors()))
		}
	})
}
