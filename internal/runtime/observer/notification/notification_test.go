package notification_test

import (
	"testing"

	"github.com/ioriimasu/jervis/internal/runtime/observer/notification"
	"github.com/ioriimasu/jervis/internal/runtime/observer/testutils"
	"github.com/ioriimasu/jervis/internal/runtime/types"
)

func TestNotification(t *testing.T) {
	evtID, _ := types.NewEventID("evt-1")
	mockEvt := &testutils.MockEvent{
		IDVal: evtID,
	}

	observedAt := types.Now()
	n := notification.New(mockEvt, observedAt)

	if n == nil {
		t.Fatal("expected non-nil notification")
	}
	if n.Event() != mockEvt {
		t.Errorf("Event() mismatch")
	}
	if n.ObservedAt() != observedAt {
		t.Errorf("ObservedAt() mismatch")
	}
	if n.Event().ID() != evtID {
		t.Errorf("Event ID mismatch")
	}

	t.Run("nil event", func(t *testing.T) {
		if notification.New(nil, observedAt) != nil {
			t.Errorf("expected nil notification for nil event")
		}
	})

	t.Run("zero timestamp", func(t *testing.T) {
		if notification.New(mockEvt, types.Timestamp{}) != nil {
			t.Errorf("expected nil notification for zero timestamp")
		}
	})
}
