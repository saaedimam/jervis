# PHASE_1_4_1_CHECKLIST.md
## Phase 1.4.1: Observer Foundation & Contracts Implementation Checklist

### Implementation Order
1. [x] Create `internal/runtime/observer/contracts/interfaces.go`
   - Define `Notification` interface (Event(), ObservedAt())
   - Define `Observer` interface (ID(), Handle(notification) error)
   - Define `Observable` interface (Notify(notification) error)
   - Define `Registry` interface (Register(), Unregister(), Observers(), Count(), Contains(), Clear())
   - Define `Dispatcher` interface (Dispatch(notification) error)

2. [x] Create `internal/runtime/observer/notification/notification.go`
   - Implement `Notification` struct wrapping `eventcontracts.Event`
   - Implement `New(event, observedAt)` constructor with validation
   - Implement `Event()` method returning wrapped event
   - Implement `ObservedAt()` method returning timestamp
   - Implement `IsZero()` and `String()` methods

3. [x] Create `internal/runtime/observer/errors/errors.go`
   - Define canonical error variables:
     - `ErrInvalidNotification`
     - `ErrDuplicateObserver`
     - `ErrObserverNotFound`
     - `ErrObserverFailure`
     - `ErrDispatchFailed`
     - `ErrObserverPanic`

4. [x] Create `internal/runtime/observer/registry/registry.go`
   - Implement `registry` struct with:
     - `observers []Observer` (FIFO order)
     - `observerMap map[string]Observer` (ID lookup)
   - Implement `NewRegistry()` constructor
   - Implement `Register()` with duplicate ID checking
   - Implement `Unregister()` by ID with FIFO preservation
   - Implement `Observers()` returning defensive copy slice
   - Implement `Count()`, `Contains()`, `Clear()` methods

5. [x] Create `internal/runtime/observer/dispatcher/dispatcher.go`
   - Implement `dispatcher` struct with embedded `Registry`
   - Implement `NewDispatcher(registry)` constructor
   - Implement `Dispatch(notification)` method:
     - Iterate through observers in FIFO order
     - Wrap each `Handle()` call in `deferred recover()`
     - Collect errors and panics into `AggregateError`
     - Continue-on-error policy (never stop early)
     - Return `nil` if no errors, otherwise `AggregateError`

6. [x] Create internal documentation
   - Update `doc.go` with package documentation

### Acceptance Criteria
- [x] All files compile: `go build ./internal/runtime/observer/...`
- [x] No lint errors: `golangci-lint run ./internal/runtime/observer/...`
- [x] 100% statement coverage: `go test -cover -race ./internal/runtime/observer/...`
- [x] No race conditions: `go test -race ./internal/runtime/observer/...` passes
- [x] All interfaces properly implemented and usable
- [x] Notification is immutable (no setter methods)
- [x] Observer registration maintains FIFO order
- [x] Dispatcher handles panics and continues notification
- [x] AggregateError properly collects and returns errors
- [x] Defensive copies returned by Registry.Observers()
- [ ] No dependencies on: context, goroutines, channels, mutexes, reflection, generics
- [ ] No imports from: internal/memory, internal/services, internal/aiprovider, internal/interfaces
- [ ] No EventBus dependency cycles (checks via `go list)

### Required Tests
1. **ObserverRegistration order of [ReportPackageObserver (from, name] }
  /* Observer ID type tests*/// Test Observer ID uniqueness
- [x] Test Observer Handle returns error handling
- [x] Test Registry Unregister by ID
- [x] Test Registry Observers returns defensive copy
- [x] Test Dispatcher Dispatch with multiple observers
- [x] Test Dispatcher continues after observer error
- [x] Test Dispatcher continues after observer panic
- [x] Test AggregateError error aggregation
- [x] Test Notification immutability (attempted mutation fails to compile)
- [x] Test Notification Event() returns wrapped event
- [x] Test Notification ObservedAt() returns timestamp
- [x] Test Notification IsZero() and String() methods
- [x] Test zero-value handling in Notification New()
- [x] Test nil event rejection in Notification New()
- [x] Test zero timestamp rejection in Notification New()

### Frozen APIs (Must Not Change After This Point)
#### Contracts Package (`internal/runtime/observer/contracts`)
```go
type Notification interface {
    Event() eventcontracts.Event
    ObservedAt() types.Timestamp
}

type Observer interface {
    ID() string
    Handle(notification Notification) error
}

type Observable interface {
    Notify(notification Notification) error
}

type Registry interface {
    Register(obs Observer) error
    Unregister(id string) error
    Observers() []Observer
    Count() int
    Contains(id string) bool
    Clear()
}

type Dispatcher interface {
    Dispatch(notification Notification) error
}
```

#### Notification Package (`internal/runtime/observer/notification`)
```go
func New(event eventcontracts.Event, observedAt types.Timestamp) (Notification, error)

type Notification interface {
    Event() eventcontracts.Event
    ObservedAt() types.Timestamp
    IsZero() bool
    String() string
}
```

#### Errors Package (`internal/runtime/observer/errors`)
```go
var (
    ErrInvalidNotification      = errors.New("observer: invalid notification")
    ErrDuplicateObserver        = errors.New("observer: duplicate observer ID")
    ErrObserverNotFound         = errors.New("observer: observer not found")
    ErrObserverFailure          = errors.New("observer: handler execution failed")
    ErrDispatchFailed           = errors.New("observer: dispatch failed")
    ErrObserverPanic            = errors.New("observer: handler execution panicked")
)

type AggregateError struct {
    errors []error
}

func NewAggregateError(errs []error) *AggregateError
func (a *AggregateError) Errors() []error
func (a *AggregateError) Error() string
```

### Risks and Mitigations
| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Observer modifies notification | Low | High | Documentation + code review + type system (interface only provides getters) |
| Deadlock or blocking in Handler | Medium | High | Documentation: Handler must return quickly; consider timeout pattern in future phases |
| Observer registration during dispatch | Low | Medium | Documented as undefined behavior; Phase 1 assumes single-threaded registration |
| Memory leak from uncleared observers | Low | Low | Registry.Clear() provided; documentation emphasizes cleanup |
| Incorrect error aggregation | Low | High | Comprehensive unit tests for Dispatcher error handling |
| Performance degradation with many observers | Low | Medium | FIFO O(n) notification is acceptable for Phase 1; optimize later if needed |
| Circular import with EventBus | Low | High | Import validation in CI; explicit dependency rules in CONTRIBUTING.md |

### Dependencies
**Allowed Imports:**
- `github.com/ioriimasu/jervis/internal/runtime/eventbus/contracts`
- `github.com/ioriimasu/jervis/internal/runtime/types`
- `github.com/ioriimasu/jervis/internal/runtime/errors`
- Standard library: `fmt`, `sort`, `strings`, `sync/atomic` (only for testing)

**Forbidden Imports:**
- `context`
- `sync` (mutexes/rwmutexes)
- `container/*`
- `reflect`
- `unsafe`
- `internal/memory/*`
- `internal/services/*`
- `internal/aiprovider/*`
- `internal/interfaces/*`
- Any third-party libraries

**Dependency Flow:**
```
Dispatcher → Registry → Notification → Contracts → EventBus Contracts → Types/Errors
```

### Definition of Done
All checklist items completed, all tests passing with 100% coverage and race detector clean, specifications reviewed and frozen, and architecture sign-off obtained.