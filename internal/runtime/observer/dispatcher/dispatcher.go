package dispatcher

import (
	"github.com/ioriimasu/jervis/internal/runtime/observer/contracts"
	"github.com/ioriimasu/jervis/internal/runtime/observer/errors"
)

// dispatcher implements contracts.Dispatcher interface.
type dispatcher struct {
	registry contracts.Registry
}

// New creates a new observer dispatcher using the provided registry.
func New(registry contracts.Registry) contracts.Dispatcher {
	return &dispatcher{
		registry: registry,
	}
}

// Dispatch executes all registered observers sequentially in FIFO order with panic isolation.
func (d *dispatcher) Dispatch(n contracts.Notification) error {
	if n == nil {
		return errors.ErrInvalidNotification
	}

	observers := d.registry.Observers()
	if len(observers) == 0 {
		return nil
	}

	var errs []error
	for _, obs := range observers {
		if err := d.safeHandle(obs, n); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.NewAggregateError(errs)
	}

	return nil
}

// safeHandle wraps observer execution in a recover block.
func (d *dispatcher) safeHandle(obs contracts.Observer, n contracts.Notification) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = &errors.ErrObserverPanic{
				ObserverID: obs.ID(),
				Recovered:  r,
			}
		}
	}()

	return obs.Handle(n)
}
