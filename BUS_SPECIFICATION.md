# Event Bus Facade Architecture Specification

## 1. Overview
The **Event Bus Facade** (`internal/runtime/eventbus`) serves as the single canonical entry point for all event-driven communications within Project Jervis.

The Facade encapsulates and orchestrates the **Subscription Registry** (`registry`), **Synchronous Dispatcher** (`dispatcher`), and **Middleware Chain** (`middleware`) behind a unified, simple, and non-bypassable facade interface.

---

## 2. Facade Responsibilities & Ownership Boundaries

### 2.1 Owned Responsibilities
1. **Orchestration**: Coordinate event validation, topic registry lookup, middleware chain execution, and handler dispatching.
2. **Subscription Management**: Expose high-level `Subscribe` and `Unsubscribe` methods.
3. **Middleware Pipeline Composition**: Allow callers to register middleware via `Use()`.
4. **Interface Isolation**: Hide internal component implementation details from external callers (`internal/runtime/...`).

### 2.2 Explicit Non-Responsibilities & Boundaries
- **No Direct Handler Invocation**: The Facade MUST NOT directly execute subscriber handlers; handler invocation is delegated to `dispatcher.Dispatcher`.
- **No Registry Data Ownership**: Topic lookup and subscription storage are strictly delegated to `registry.Registry`.
- **No AI Provider Awareness**: Zero imports from `internal/aiprovider` (Invariant 2).
- **No Storage or Logging**: Zero disk or database operations.
- **No Background Concurrency**: Execution is 100% synchronous on the caller's call stack.

---

## 3. Public API Specification

The Facade exposes the following frozen public API methods:

```go
// EventBus represents the canonical Event Bus facade instance.
type EventBus struct {
    registry   *registry.Registry
    dispatcher *dispatcher.Dispatcher
    chain      *middleware.Chain
}

// New constructs an initialized EventBus facade.
func New() *EventBus

// Publish validates, intercepts, and dispatches an event.
func (b *EventBus) Publish(event contracts.Event) error

// Subscribe registers a subscriber handler for an event type pattern.
func (b *EventBus) Subscribe(pattern string, handler contracts.Handler, priority uint8) (subscription.SubscriptionID, error)

// Unsubscribe removes a subscription by ID.
func (b *EventBus) Unsubscribe(subID subscription.SubscriptionID) error

// Use registers one or more middleware interceptors in FIFO order.
func (b *EventBus) Use(mw ...contracts.Middleware)

// Count returns the total active subscription count.
func (b *EventBus) Count() int
```

---

## 4. Execution Lifecycles

### 4.1 Publish Lifecycle
```
Publish(event)
  ├── Stage 1: Validate Event (events.ValidateEvent)
  ├── Stage 2: Registry Lookup (registry.Lookup) -> returns matching handlers
  ├── Stage 3: Middleware Chain Execution (chain.Execute)
  │             └── Pre-Hooks (M1 -> M2)
  │                   └── Dispatcher.Dispatch(event, handlers)
  │                         └── Panic-Isolated Handler Loop
  │                   └── Post-Hooks (M2 -> M1)
  └── Stage 4: Return Result
```

### 4.2 Subscribe Lifecycle
```
Subscribe(pattern, handler, priority)
  ├── Stage 1: Generate SubscriptionID & Construct Subscription
  ├── Stage 2: Validate Subscription (sub.Validate())
  ├── Stage 3: Validate Pattern (registry.ValidatePattern(pattern))
  ├── Stage 4: Register in Registry (registry.Register(sub))
  └── Stage 5: Return SubscriptionID
```

### 4.3 Unsubscribe Lifecycle
```
Unsubscribe(subID)
  ├── Stage 1: Validate SubscriptionID (subID.IsZero())
  ├── Stage 2: Remove from Registry (registry.Unregister(subID))
  └── Stage 3: Return Result
```

---

## 5. Error & Panic Policy

- **Validation Errors**: Validation failures during `Publish`, `Subscribe`, or `Unsubscribe` abort processing immediately and return `errs.ErrValidationFailed`.
- **Dispatcher & Handler Errors**: Handler errors and panics are aggregated into a composite `dispatcher.AggregateError` and returned to the publisher.
- **Middleware Panics**: Trapped by middleware chain panic recovery and returned as wrapped `errs.ErrHandlerFailure` errors.

---

## 6. Performance & Concurrency Guarantees

- **Synchronous Execution**: Zero goroutines, channels, or mutexes.
- **Thread Safety**: Single-threaded call stack execution guarantees deterministic FIFO/priority ordering.
- **Future Compatibility**: Interfaces allow wrapping background worker pools in future phases without breaking public facade signatures.
