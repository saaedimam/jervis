# Runtime Observer Data Model (Phase 1.4.0 - FROZEN)

## 1. Canonical Event Model Wrapping

To prevent model duplication and maintain strict consistency across the Jervis Runtime, the Observer subsystem does **NOT** duplicate event fields (`ID`, `Type`, `Source`, `Timestamp`, `CorrelationID`, `CausationID`, `Priority`, `Payload`, `Metadata`, `Version`).

Instead, `Notification` wraps the canonical EventBus interface `eventcontracts.Event` (from `internal/runtime/eventbus/contracts`).

---

## 2. Notification Model

The `Notification` struct (in `notification` package) provides an immutable implementation of the `contracts.Notification` interface.

```go
package notification

import (
	eventcontracts "github.com/saaedimam/jervis/internal/runtime/eventbus/contracts"
	"github.com/saaedimam/jervis/internal/runtime/types"
)

type Notification struct {
	event      eventcontracts.Event
	observedAt types.Timestamp
}
```

### Methods & Immutability Rules
- `New(event eventcontracts.Event, observedAt types.Timestamp) (Notification, error)`
  - Validates that `event` is non-nil and `observedAt` is not zero.
  - Wraps the event as-is (assuming the Event implementation is immutable).
- `Event() eventcontracts.Event`
  - Returns the wrapped canonical `Event` interface.
- `ObservedAt() types.Timestamp`
  - Returns the observation timestamp.
- `IsZero() bool`
  - Reports whether the `Notification` is uninitialized.
- `String() string`
  - Returns formatted string: `NOTIFICATION[<event_type>:<event_id>] at <observedAt>`.

---

## 3. Error Models & Aggregation

### AggregateError
To enforce Continue-on-Error and isolate observer failures, `dispatcher` uses `AggregateError`.

```go
type AggregateError struct {
	errors []error
}
```

- `NewAggregateError(errs []error) *AggregateError`
- `(a *AggregateError) Errors() []error` (Returns defensive copy of accumulated errors)
- `(a *AggregateError) Error() string` (Formatted composite error message)

### Canonical Errors
- `ErrInvalidNotification`: Notification event is nil or observedAt is invalid.
- `ErrDuplicateObserver`: Observer ID already registered.
- `ErrObserverNotFound`: Observer ID not found during unregistration.
- `ErrObserverFailure`: Observer handle execution returned an error.
- `ErrDispatchFailed`: Notification dispatch encountered one or more observer errors.
- `ErrObserverPanic`: Observer handle execution panicked (wrapped error).