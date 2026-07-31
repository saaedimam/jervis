# Dispatcher Architectural Specification

## 1. Overview
The **Dispatcher** (`internal/runtime/eventbus/dispatcher`) is the central runtime engine responsible for executing matching subscriber handlers when an event is published to the Jervis Event Bus.

The Dispatcher enforces **deterministic, synchronous execution**, **panic isolation**, and **strict priority ordering** across local process boundaries without depending on external services, background goroutines, or AI models.

---

## 2. Core Responsibilities

### 2.1 Owned Responsibilities
1. **Synchronous Execution**: Invoke registered `contracts.Handler` callback functions sequentially on the caller's call stack.
2. **Deterministic Priority Dispatch**: Ensure handlers execute strictly according to `Priority` (DESCENDING: `PriorityCritical` -> `PriorityHigh` -> `PriorityNormal` -> `PriorityLow`) and registration sequence (`Seq` ASCENDING).
3. **Panic Recovery & Isolation**: Trap handler panics via `recover()`, wrapping them into canonical `ErrHandlerFailure` errors without crashing the calling thread or process.
4. **Error Aggregation & Continuing Execution**: Execute all matching handlers for an event even if individual handlers return errors or panic (Continue-on-Error policy). Aggregate all failures into a composite dispatch error.
5. **Recursion & Depth Guard**: Enforce a maximum dispatch depth (`MaxDispatchDepth = 16`) to prevent infinite recursive event loops.
6. **Middleware Pipeline Hook Support**: Execute pre-dispatch and post-dispatch middleware hooks.

### 2.2 Explicit Non-Responsibilities
- **No Persistence or Logging**: The Dispatcher does not persist events or write diagnostic logs to disk.
- **No Goroutines / Concurrency (Phase 1)**: Execution occurs entirely on the publishing thread.
- **No Network Transport**: Operations are local process memory calls only.
- **No AI Provider Dependencies**: Zero imports from `internal/aiprovider`.

---

## 3. Dispatch Lifecycle

Every call to `Dispatch(event, handlers)` follows a rigid 7-stage lifecycle:

```
+-------------------------------------------------------------------+
| 1. Validation Stage (events.ValidateEvent)                        |
+-------------------------------------------------------------------+
                                  |
                                  v
+-------------------------------------------------------------------+
| 2. Recursion Depth Guard (Depth <= MaxDispatchDepth = 16)         |
+-------------------------------------------------------------------+
                                  |
                                  v
+-------------------------------------------------------------------+
| 3. Handler Ordering (Priority DESC, Registration Seq ASC)          |
+-------------------------------------------------------------------+
                                  |
                                  v
+-------------------------------------------------------------------+
| 4. Middleware Chain Pre-Hooks Execution                           |
+-------------------------------------------------------------------+
                                  |
                                  v
+-------------------------------------------------------------------+
| 5. Handler Loop Execution (Panic Isolation per Handler)           |
+-------------------------------------------------------------------+
                                  |
                                  v
+-------------------------------------------------------------------+
| 6. Error Aggregation & Normalization                              |
+-------------------------------------------------------------------+
                                  |
                                  v
+-------------------------------------------------------------------+
| 7. Middleware Chain Post-Hooks Execution                          |
+-------------------------------------------------------------------+
```

---

## 4. Architectural Invariants & Policies

### 4.1 Handler Execution Order
- Handlers MUST be sorted before invocation:
  1. `Priority` DESCENDING (`PriorityCritical` > `PriorityHigh` > `PriorityNormal` > `PriorityLow`).
  2. `Seq` ASCENDING (registration sequence FIFO for handlers of identical priority).

### 4.2 Validation Order
- **Stage 1 Validation**: `events.ValidateEvent(event)` is executed before looking up handlers or invoking middleware. If structural validation fails, dispatch aborts immediately with `errs.ErrValidationFailed`.
- **Stage 2 Handler Check**: If the handler list is empty, `Dispatch` returns `nil` immediately.

### 4.3 Panic Recovery Policy
- Each handler invocation MUST be wrapped in a panic isolation block:
```go
defer func() {
    if r := recover(); r != nil {
        err = fmt.Errorf("%w: handler %s panicked: %v", errs.ErrHandlerFailure, handler.ID(), r)
    }
}()
```
- A handler panic MUST NOT crash the process or interrupt the execution of remaining handlers.

### 4.4 Error Propagation & Stop-on-Error vs. Continue Policy
- **Policy**: **Continue-on-Error**.
- When a handler returns an `error` or panics, the Dispatcher records the failure and proceeds immediately to the next handler in the priority list.
- After all handlers complete, if any handler failed or panicked, a composite error wrapping all handler errors is returned to the caller.

### 4.5 Maximum Dispatch Depth
- To guard against infinite event loops (`EventA` triggers `EventB` which triggers `EventA`), the Dispatcher maintains a call stack depth counter.
- **Maximum Depth Cap**: `MaxDispatchDepth = 16`.
- Exceeding `MaxDispatchDepth` causes `Dispatch` to abort immediately with `errs.ErrDispatchFailed`.

### 4.6 Middleware Hook Points
- **Pre-Dispatch Hooks**: Executed before the handler loop begins. Can inspect event state or inject context attributes.
- **Post-Dispatch Hooks**: Executed after all handlers complete. Can inspect collected execution errors or record overall latency.

### 4.7 Future Async Compatibility (Spec Only)
- In Phase 1, `Dispatch` is strictly synchronous.
- In future phases (Phase 2+), worker pool dispatchers may implement `contracts.Dispatcher` to execute handlers asynchronously on background goroutines without changing subscriber or handler contracts.

---

## 5. Summary of Frozen Specification Rules
- **Pure Synchronous Execution**: Zero goroutines or channels.
- **Panic Guard**: Process crash impossible from subscriber panics.
- **Zero AI Imports**: Strict adherence to Invariant 2.
- **100% Deterministic**: Explicit sorting guarantees identical execution order across runs.
