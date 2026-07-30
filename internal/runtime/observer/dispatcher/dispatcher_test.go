package dispatcher_test

import (
	"errors"
	"testing"

	"github.com/saaedimam/jervis/internal/runtime/observer/contracts"
	"github.com/saaedimam/jervis/internal/runtime/observer/dispatcher"
	observererrors "github.com/saaedimam/jervis/internal/runtime/observer/errors"
	"github.com/saaedimam/jervis/internal/runtime/observer/notification"
	"github.com/saaedimam/jervis/internal/runtime/observer/registry"
	"github.com/saaedimam/jervis/internal/runtime/observer/testutils"
	"github.com/saaedimam/jervis/internal/runtime/types"
)

func TestDispatcher(t *testing.T) {
	reg := registry.New()
	disp := dispatcher.New(reg)

	t.Run("nil notification", func(t *testing.T) {
		if err := disp.Dispatch(nil); err != observererrors.ErrInvalidNotification {
			t.Errorf("expected ErrInvalidNotification, got %v", err)
		}
	})

	t.Run("empty registry", func(t *testing.T) {
		evt := &testutils.MockEvent{}
		n := notification.New(evt, types.Now())
		if err := disp.Dispatch(n); err != nil {
			t.Errorf("expected nil error for empty registry, got %v", err)
		}
	})

	t.Run("sequential execution and aggregation", func(t *testing.T) {
		var executed []string
		obs1 := &testutils.MockObserver{
			IDVal: "obs-1",
			HandleFunc: func(n contracts.Notification) error {
				executed = append(executed, "obs-1")
				return nil
			},
		}
		obs2 := &testutils.MockObserver{
			IDVal: "obs-2",
			HandleFunc: func(n contracts.Notification) error {
				executed = append(executed, "obs-2")
				return errors.New("err-2")
			},
		}
		obs3 := &testutils.MockObserver{
			IDVal: "obs-3",
			HandleFunc: func(n contracts.Notification) error {
				executed = append(executed, "obs-3")
				panic("panic-3")
			},
		}

		_ = reg.Register(obs1)
		_ = reg.Register(obs2)
		_ = reg.Register(obs3)

		evt := &testutils.MockEvent{}
		n := notification.New(evt, types.Now())
		err := disp.Dispatch(n)

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		agg, ok := err.(*observererrors.AggregateError)
		if !ok {
			t.Fatalf("expected AggregateError, got %T", err)
		}

		if len(agg.Errors()) != 2 {
			t.Errorf("expected 2 errors, got %d", len(agg.Errors()))
		}

		if len(executed) != 3 {
			t.Errorf("expected 3 executions, got %d", len(executed))
		}

		if executed[0] != "obs-1" || executed[1] != "obs-2" || executed[2] != "obs-3" {
			t.Errorf("execution order mismatch: %v", executed)
		}
	})
}
