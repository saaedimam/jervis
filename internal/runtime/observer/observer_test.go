package observer_test

import (
	"testing"

	"github.com/ioriimasu/jervis/internal/runtime/observer"
	"github.com/ioriimasu/jervis/internal/runtime/observer/contracts"
	"github.com/ioriimasu/jervis/internal/runtime/observer/testutils"
)

func TestRuntimeObserver(t *testing.T) {
	ro := observer.New()

	var notified bool
	obs := &testutils.MockObserver{
		IDVal: "obs-facade",
		HandleFunc: func(n contracts.Notification) error {
			notified = true
			return nil
		},
	}

	_ = ro.Register(obs)

	evt := &testutils.MockEvent{}
	err := ro.Notify(evt)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !notified {
		t.Errorf("observer was not notified")
	}

	if err := ro.Unregister("obs-facade"); err != nil {
		t.Errorf("unexpected unregister error: %v", err)
	}

	notified = false
	_ = ro.Notify(evt)
	if notified {
		t.Errorf("unregistered observer was notified")
	}
}
