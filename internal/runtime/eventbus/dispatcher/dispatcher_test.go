package dispatcher_test

import (
	"errors"
	"testing"

	"github.com/saaedimam/jervis/internal/runtime/eventbus/contracts"
	"github.com/saaedimam/jervis/internal/runtime/eventbus/dispatcher"
	errs "github.com/saaedimam/jervis/internal/runtime/eventbus/errors"
	"github.com/saaedimam/jervis/internal/runtime/eventbus/events"
	"github.com/saaedimam/jervis/internal/runtime/types"
)

type mockHandler struct {
	id       string
	priority uint8
	seq      uint64
	failErr  error
	panicVal any
	called   *[]string
}

func (m *mockHandler) ID() string      { return m.id }
func (m *mockHandler) Priority() uint8 { return m.priority }
func (m *mockHandler) Seq() uint64     { return m.seq }
func (m *mockHandler) Handle(event contracts.Event) error {
	if m.called != nil {
		*m.called = append(*m.called, m.id)
	}
	if m.panicVal != nil {
		panic(m.panicVal)
	}
	return m.failErr
}

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

func TestDispatchSuccessAndOrdering(t *testing.T) {
	disp := dispatcher.NewDispatcher()
	evt := createTestEvent(t)

	var executionOrder []string
	hLow := &mockHandler{id: "h-low", priority: uint8(events.PriorityLow), seq: 1, called: &executionOrder}
	hNorm1 := &mockHandler{id: "h-norm-1", priority: uint8(events.PriorityNormal), seq: 1, called: &executionOrder}
	hNorm2 := &mockHandler{id: "h-norm-2", priority: uint8(events.PriorityNormal), seq: 2, called: &executionOrder}
	hCrit := &mockHandler{id: "h-crit", priority: uint8(events.PriorityCritical), seq: 1, called: &executionOrder}

	// Pass handlers out of order
	handlers := []contracts.Handler{hLow, hNorm2, hCrit, hNorm1}

	err := disp.Dispatch(evt, handlers)
	if err != nil {
		t.Fatalf("expected dispatch success, got: %v", err)
	}

	expectedOrder := []string{"h-crit", "h-norm-1", "h-norm-2", "h-low"}
	if len(executionOrder) != len(expectedOrder) {
		t.Fatalf("execution count mismatch: got %v, want %v", executionOrder, expectedOrder)
	}
	for i, id := range expectedOrder {
		if executionOrder[i] != id {
			t.Errorf("position %d: got %s, want %s", i, executionOrder[i], id)
		}
	}
}

func TestDispatchEmptyHandlers(t *testing.T) {
	disp := dispatcher.NewDispatcher()
	evt := createTestEvent(t)

	if err := disp.Dispatch(evt, nil); err != nil {
		t.Fatalf("expected nil error for empty handlers, got %v", err)
	}
	if err := disp.Dispatch(evt, []contracts.Handler{}); err != nil {
		t.Fatalf("expected nil error for empty slice handlers, got %v", err)
	}
}

func TestDispatchInvalidEvent(t *testing.T) {
	disp := dispatcher.NewDispatcher()
	h := &mockHandler{id: "h1", priority: uint8(events.PriorityNormal)}

	// Invalid event (nil)
	if err := disp.Dispatch(nil, []contracts.Handler{h}); !errors.Is(err, errs.ErrInvalidEvent) {
		t.Fatalf("expected ErrInvalidEvent, got %v", err)
	}
}

func TestDispatchContinueOnErrorAndPanicRecovery(t *testing.T) {
	disp := dispatcher.NewDispatcher()
	evt := createTestEvent(t)

	var executionOrder []string
	h1 := &mockHandler{id: "h1", priority: uint8(events.PriorityNormal), seq: 1, called: &executionOrder, failErr: errors.New("h1 failed")}
	h2 := &mockHandler{id: "h2", priority: uint8(events.PriorityNormal), seq: 2, called: &executionOrder, panicVal: "h2 boom"}
	h3 := &mockHandler{id: "h3", priority: uint8(events.PriorityNormal), seq: 3, called: &executionOrder}

	err := disp.Dispatch(evt, []contracts.Handler{h1, h2, h3})
	if err == nil {
		t.Fatalf("expected error from aggregate dispatch")
	}

	// Verify all 3 handlers were invoked (Continue-on-error policy)
	if len(executionOrder) != 3 {
		t.Fatalf("expected 3 handlers invoked, got %d (%v)", len(executionOrder), executionOrder)
	}

	var aggErr *dispatcher.AggregateError
	if !errors.As(err, &aggErr) {
		t.Fatalf("expected error to be *dispatcher.AggregateError")
	}

	if aggErr.Count() != 2 {
		t.Fatalf("expected 2 aggregated errors, got %d", aggErr.Count())
	}

	if !errors.Is(err, errs.ErrHandlerFailure) {
		t.Fatalf("expected errors.Is(err, ErrHandlerFailure) to be true")
	}
}

func TestDispatchNilHandlerInList(t *testing.T) {
	disp := dispatcher.NewDispatcher()
	evt := createTestEvent(t)

	err := disp.Dispatch(evt, []contracts.Handler{nil})
	if err == nil {
		t.Fatalf("expected error for nil handler")
	}
	if !errors.Is(err, errs.ErrHandlerFailure) {
		t.Fatalf("expected ErrHandlerFailure for nil handler, got %v", err)
	}
}

type recursiveHandler struct {
	disp *dispatcher.Dispatcher
}

func (r *recursiveHandler) ID() string      { return "recursive" }
func (r *recursiveHandler) Priority() uint8 { return uint8(events.PriorityNormal) }
func (r *recursiveHandler) Seq() uint64     { return 1 }
func (r *recursiveHandler) Handle(event contracts.Event) error {
	return r.disp.Dispatch(event, []contracts.Handler{r})
}

func TestDispatchMaxDepthGuard(t *testing.T) {
	disp := dispatcher.NewDispatcher()
	evt := createTestEvent(t)

	recH := &recursiveHandler{disp: disp}

	err := disp.Dispatch(evt, []contracts.Handler{recH})
	if err == nil {
		t.Fatalf("expected max depth exceeded error")
	}

	if !errors.Is(err, errs.ErrDispatchFailed) {
		t.Fatalf("expected ErrDispatchFailed for depth exceeded, got %v", err)
	}
}

type plainHandler struct {
	id string
}

func (p *plainHandler) ID() string                         { return p.id }
func (p *plainHandler) Handle(event contracts.Event) error { return nil }

func TestDispatchPlainHandlersWithoutPriorityInterface(t *testing.T) {
	disp := dispatcher.NewDispatcher()
	evt := createTestEvent(t)

	h1 := &plainHandler{id: "b"}
	h2 := &plainHandler{id: "a"}

	err := disp.Dispatch(evt, []contracts.Handler{h1, h2})
	if err != nil {
		t.Fatalf("unexpected dispatch error: %v", err)
	}
}
