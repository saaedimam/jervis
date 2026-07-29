package dispatcher

import (
	"fmt"

	"github.com/ioriimasu/jervis/internal/runtime/observer/contracts"
	obserrors "github.com/ioriimasu/jervis/internal/runtime/observer/errors"
)

type dispatcher struct {
	registry contracts.Registry
}

// NewDispatcher creates a new Observer Dispatcher.
func NewDispatcher(registry contracts.Registry) contracts.Dispatcher {
	return &dispatcher{
		registry: registry,
	}
}

func (d *dispatcher) Dispatch(n contracts.Notification) error {
	if n == nil {
		return obserrors.ErrInvalidNotification
	}

	if d.registry == nil {
		return nil
	}

	observers := d.registry.Observers()
	if len(observers) == 0 {
		return nil
	}

	var collectedErrors []error

	for _, obs := range observers {
		if obs == nil {
			continue
		}

		err := d.safeHandle(obs, n)
		if err != nil {
			collectedErrors = append(collectedErrors, err)
		}
	}

	if len(collectedErrors) == 0 {
		return nil
	}

	return obserrors.NewAggregateError(collectedErrors)
}

func (d *dispatcher) safeHandle(obs contracts.Observer, n contracts.Notification) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%w: observer [%s] panicked: %v", obserrors.ErrObserverPanic, obs.ID(), r)
		}
	}()

	return obs.Handle(n)
}
