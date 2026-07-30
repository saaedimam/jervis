package observer

import (
	"github.com/ioriimasu/jervis/internal/runtime/eventbus/contracts"
	obscontracts "github.com/ioriimasu/jervis/internal/runtime/observer/contracts"
	"github.com/ioriimasu/jervis/internal/runtime/observer/dispatcher"
	"github.com/ioriimasu/jervis/internal/runtime/observer/notification"
	"github.com/ioriimasu/jervis/internal/runtime/observer/registry"
	"github.com/ioriimasu/jervis/internal/runtime/types"
)

// RuntimeObserver provides the high-level API for the Observer subsystem.
type RuntimeObserver struct {
	registry   obscontracts.Registry
	dispatcher obscontracts.Dispatcher
}

// New creates a new RuntimeObserver instance.
func New() *RuntimeObserver {
	reg := registry.New()
	return &RuntimeObserver{
		registry:   reg,
		dispatcher: dispatcher.New(reg),
	}
}

// Register adds an observer to the runtime.
func (o *RuntimeObserver) Register(obs obscontracts.Observer) error {
	return o.registry.Register(obs)
}

// Unregister removes an observer from the runtime.
func (o *RuntimeObserver) Unregister(id string) error {
	return o.registry.Unregister(id)
}

// Notify creates and dispatches a notification for the given event.
func (o *RuntimeObserver) Notify(event contracts.Event) error {
	n := notification.New(event, types.Now())
	return o.dispatcher.Dispatch(n)
}

// Registry returns the underlying observer registry.
func (o *RuntimeObserver) Registry() obscontracts.Registry {
	return o.registry
}
