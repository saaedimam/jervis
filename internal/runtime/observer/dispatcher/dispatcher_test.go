package dispatcher_test

import (
	"errors"
	"testing"
	"time"

	eventcontracts "github.com/ioriimasu/jervis/internal/runtime/eventbus/contracts"
	"github.com/ioriimasu/jervis/internal/runtime/observer/contracts"
	obsdispatcher "github.com/ioriimasu/jervis/internal/runtime/observer/dispatcher"
	obserrors "github.com/ioriimasu/jervis/internal/runtime/observer/errors"
	obsnotification "github.com/ioriimasu/jervis/internal/runtime/observer/notification"
	obsregistry "github.com/ioriimasu/jervis/internal/runtime/observer/registry"
	"github.com/ioriimasu/jervis/internal/runtime/types"
)

type mockEvent struct{}

func (m *mockEvent) ID() types.EventID {
	id, _ := types.NewEventID("evt-1")
	return id
}
func (m *mockEvent) Type() string                           { return "test.event" }
func (m *mockEvent) Source() string                         { return "test" }
func (m *mockEvent) Timestamp() types.Timestamp             { return types.NewTimestamp(time.Now()) }
func (m *mockEvent) CorrelationID() string                  { return "" }
func (m *mockEvent) CausationID() string                    { return "" }
func (m *mockEvent) Priority() uint8                        { return 1 }
func (m *mockEvent) Payload() any                           { return nil }
func (m *mockEvent) Metadata() map[string]string            { return nil }
func (m *mockEvent) Version() string                        { return "1.0" }

var _ eventcontracts.Event = (*mockEvent)(nil)

type testObserver struct {
	id          string
	shouldError bool
	shouldPanic bool
	handled     bool
}

func (o *testObserver) ID() string {
	return o.id
}

func (o *testObserver) Handle(n contracts.Notification) error {
	o.handled = true
	if o.shouldPanic {
		panic("boom")
	}
	if o.shouldError {
		return errors.New("handler error")
	}
	return nil
}

func TestDispatcher(t *testing.T) {
	reg := obsregistry.NewRegistry()
	disp := obsdispatcher.NewDispatcher(reg)

	n, _ := obsnotification.New(&mockEvent{}, types.NewTimestamp(time.Now()))

	// Nil notification
	err := disp.Dispatch(nil)
	if err != obserrors.ErrInvalidNotification {
		t.Errorf("Expected ErrInvalidNotification, got %v", err)
	}

	// Empty registry
	err = disp.Dispatch(n)
	if err != nil {
		t.Errorf("Expected nil error for empty registry, got %v", err)
	}

	// Normal execution
	obs1 := &testObserver{id: "obs-1"}
	obs2 := &testObserver{id: "obs-2"}
	_ = reg.Register(obs1)
	_ = reg.Register(obs2)

	err = disp.Dispatch(n)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !obs1.handled || !obs2.handled {
		t.Error("Observers were not handled")
	}
}

func TestDispatcherErrorAndPanicIsolation(t *testing.T) {
	reg := obsregistry.NewRegistry()
	disp := obsdispatcher.NewDispatcher(reg)

	n, _ := obsnotification.New(&mockEvent{}, types.NewTimestamp(time.Now()))

	obsNormal := &testObserver{id: "obs-normal"}
	obsErr := &testObserver{id: "obs-err", shouldError: true}
	obsPanic := &testObserver{id: "obs-panic", shouldPanic: true}
	obsAfter := &testObserver{id: "obs-after"}

	_ = reg.Register(obsNormal)
	_ = reg.Register(obsErr)
	_ = reg.Register(obsPanic)
	_ = reg.Register(obsAfter)

	err := disp.Dispatch(n)
	if err == nil {
		t.Fatal("Expected error from dispatcher containing panic and error")
	}

	// Verify continue-on-error: obsAfter MUST have been executed despite previous error and panic!
	if !obsAfter.handled {
		t.Errorf("Dispatcher stopped execution early on error/panic; obsAfter was not handled")
	}
}

func TestNilRegistryDispatcher(t *testing.T) {
	disp := obsdispatcher.NewDispatcher(nil)
	n, _ := obsnotification.New(&mockEvent{}, types.NewTimestamp(time.Now()))
	err := disp.Dispatch(n)
	if err != nil {
		t.Errorf("Expected nil error for nil registry, got %v", err)
	}
}
