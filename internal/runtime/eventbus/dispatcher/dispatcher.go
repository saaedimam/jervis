package dispatcher

import (
	"fmt"

	"github.com/saaedimam/jervis/internal/runtime/eventbus/contracts"
	errs "github.com/saaedimam/jervis/internal/runtime/eventbus/errors"
	"github.com/saaedimam/jervis/internal/runtime/eventbus/events"
)

// MaxDispatchDepth caps the maximum allowed recursive event dispatch call depth.
const MaxDispatchDepth = 16

// Dispatcher implements contracts.Dispatcher with synchronous execution, deterministic ordering, and panic recovery.
type Dispatcher struct {
	depth int
}

var _ contracts.Dispatcher = (*Dispatcher)(nil)

// NewDispatcher constructs a new Dispatcher instance.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{}
}

// Dispatch executes the 7-stage event dispatch pipeline.
func (d *Dispatcher) Dispatch(event contracts.Event, handlers []contracts.Handler) error {
	// Stage 1: Validate Event
	if err := events.ValidateEvent(event); err != nil {
		return err
	}

	// Stage 2: Validate Dispatch Depth
	if d.depth >= MaxDispatchDepth {
		return fmt.Errorf("%w: maximum dispatch depth %d exceeded", errs.ErrDispatchFailed, MaxDispatchDepth)
	}

	d.depth++
	defer func() {
		d.depth--
	}()

	// Stage 3: Resolve Handler List
	if len(handlers) == 0 {
		return nil
	}

	// Stage 4: Sort Handlers (Priority DESC, Sequence ASC, ID ASC)
	sortedHandlers := d.sortHandlers(handlers)

	// Stage 5 & Stage 6: Invoke Handlers with Panic Isolation & Error Aggregation
	agg := &AggregateError{}
	for _, h := range sortedHandlers {
		if err := d.invokeHandler(event, h); err != nil {
			agg.Add(err)
		}
	}

	// Stage 7: Return Result
	if agg.HasErrors() {
		return agg
	}
	return nil
}

func (d *Dispatcher) invokeHandler(event contracts.Event, handler contracts.Handler) (err error) {
	if handler == nil {
		return fmt.Errorf("%w: handler cannot be nil", errs.ErrHandlerFailure)
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%w: handler %q panicked: %v", errs.ErrHandlerFailure, handler.ID(), r)
		}
	}()
	if hErr := handler.Handle(event); hErr != nil {
		return fmt.Errorf("%w: handler %q failed: %w", errs.ErrHandlerFailure, handler.ID(), hErr)
	}
	return nil
}

type priorityProvider interface {
	Priority() uint8
}

type seqProvider interface {
	Seq() uint64
}

func (d *Dispatcher) sortHandlers(handlers []contracts.Handler) []contracts.Handler {
	sorted := make([]contracts.Handler, len(handlers))
	copy(sorted, handlers)

	// Stable sort matching priority DESC, seq ASC, ID ASC
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			pI, pJ := getPriority(sorted[i]), getPriority(sorted[j])
			sI, sJ := getSeq(sorted[i]), getSeq(sorted[j])

			swap := false
			if pI < pJ {
				swap = true
			} else if pI == pJ {
				if sI > sJ {
					swap = true
				} else if sI == sJ {
					if sorted[i].ID() > sorted[j].ID() {
						swap = true
					}
				}
			}

			if swap {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	return sorted
}

func getPriority(h contracts.Handler) uint8 {
	if p, ok := h.(priorityProvider); ok {
		return p.Priority()
	}
	return uint8(events.PriorityNormal)
}

func getSeq(h contracts.Handler) uint64 {
	if s, ok := h.(seqProvider); ok {
		return s.Seq()
	}
	return 0
}
