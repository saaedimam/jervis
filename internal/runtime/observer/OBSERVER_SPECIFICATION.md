# Runtime Observer Specification (Phase 1.4.0 - FROZEN)

## 1. Overview & Responsibilities

The Runtime Observer subsystem is the canonical, read-only event observation layer of Project Jervis. It allows internal subsystems (such as Logger, Telemetry, and Metrics) to observe runtime events without modifying system state or intercepting execution flows.

### Core Responsibilities
1. **Read-Only Event Observation**: Inspect runtime events without modifying envelopes, payloads, or metadata.
2. **Deterministic FIFO Dispatch**: Notify registered observers in strict registration sequence (FIFO).
3. **Fault Isolation & Recovery**: Isolate observer failures and panics per handler using `recover()`, ensuring remaining observers are always notified (Continue-on-Error).
4. **Aggregate Error Reporting**: Collect observer handler errors and recovered panics into a composite `AggregateError`.
5. **Zero Side Effects**: Operates purely in-memory with zero persistence, zero permission checks, zero memory engine calls, zero domain service calls, and zero AI provider invocations.
6. **Pure Synchronous Execution**: Operates 100% synchronously without goroutines, background worker loops, channels, mutexes, atomic operations, or `context.Context`.

### Restricted Actions (Explicit Invariants)
An Observer MUST NOT under any circumstances:
- Publish or re-emit events to the EventBus.
- Mutate the observed `Notification` or underlying `Event`.
- Invoke the Permission Engine.
- Invoke the Memory Engine.
- Invoke Domain Services (Planner, Projects, Meetings, Habits, etc.).
- Invoke AI Provider services or external network endpoints.

---

## 2. Architecture & Subsystem Boundaries

The Observer subsystem is partitioned into modular, single-responsibility packages mirroring the EventBus architecture:

```
internal/runtime/observer/
├── contracts/        # Interfaces (Notification, Observer, Observable, Dispatcher, Registry)
├── notification/     # Immutable Notification implementation wrapping eventcontracts.Event
├── errors/           # Canonical error definitions & AggregateError
├── registry/         # Deterministic FIFO observer storage & lookup
└── dispatcher/       # Synchronous execution engine with panic isolation
```

### System Architecture Flow

```
Runtime Event Source
        │
        ▼
   Notification (Wraps canonical eventcontracts.Event interface + ObservedAt)
        │
        ▼
 Observer Facade / Dispatcher
        │
        ▼
  Registry Lookup (FIFO Order)
        │
        ▼
┌─────────────────────────────────────────────────────────────┐
│ Sequential Notification Pipeline                            │
│ ┌─────────────┐   ┌─────────────┐   ┌─────────────┐        │
│ │ Observer 1  │──>│ Observer 2  │──>│ Observer N  │        │
│ └──────┬──────┘   └──────┬──────┘   └──────┬──────┘        │
│        │                 │                 │               │
│  Panic / Error     Panic / Error     Panic / Error         │
│        │                 │                 │               │
│        ▼                 ▼                 ▼               │
│  Recover & Log    Recover & Log     Recover & Log          │
└────────┬─────────────────┬─────────────────┬───────────────┘
         │                 │                 │
         └─────────────────┼─────────────────┘
                           ▼
                    AggregateError
                           │
                           ▼
                 Return to Caller
```

---

## 3. Execution Pipeline

When an event is dispatched to registered observers, the execution pipeline executes in 5 deterministic stages:

1. **Stage 1 — Notification Construction**:
   - Construct immutable `Notification` wrapping the canonical `eventcontracts.Event` interface and `ObservedAt` timestamp.
2. **Stage 2 — Observer Lookup**:
   - Retrieve all registered `Observer` instances from the `Registry` in strict FIFO registration order.
3. **Stage 3 — Panic-Protected Sequential Loop**:
   - Iterate through the resolved observer slice.
   - Wrap each `observer.Handle(notification)` invocation in a `deferred` `recover()` block.
4. **Stage 4 — Error & Panic Aggregation**:
   - If `Handle()` returns an error, append it to the composite `AggregateError`.
   - If a panic is recovered, wrap it in `ErrObserverPanic` and append it to `AggregateError`.
   - **Continue-on-Error Policy**: Execution NEVER stops upon encountering an error or panic; subsequent observers MUST be notified.
5. **Stage 5 — Completion & Return**:
   - If no errors or panics occurred, return `nil`.
   - Otherwise, return the compiled `AggregateError`.

---

## 4. Ordering & Immutability Rules

### Deterministic FIFO Ordering
- Observers are executed strictly in the order they were registered.
- Direct Go map iteration is strictly forbidden during dispatch.
- Lookups and snapshots return slice copies sorted by registration sequence `Seq ASC`.

### Immutability Guarantees
- The `Notification` interface wraps `eventcontracts.Event` and provides read-only getters (`Event()`, `ObservedAt()`).
- Observers MUST treat all received `Notification` and `Event` objects as read-only.
- The underlying event implementation MUST be immutable.

---

## 5. Dependency Graph

```
internal/runtime/observer/dispatcher
    ↓
internal/runtime/observer/registry
    ↓
internal/runtime/observer/notification
    ↓
internal/runtime/observer/contracts
    ↓
internal/runtime/eventbus/contracts
    ↓
internal/runtime/types & errors
    ↓
Standard Library (fmt, sort, strings, time)
```

### Dependency Invariants
- Zero cyclic imports with `internal/runtime/eventbus`.
- Zero imports of `internal/memory`, `internal/services`, `internal/aiprovider`, or `internal/interfaces`.
- Zero third-party dependencies.

---

## 6. Architecture Invariants & Conflict Audit

| Target Document | Checked Invariants | Audit Result | Resolution |
| :--- | :--- | :--- | :--- |
| `ARCHITECTURE_INVARIANTS.md` | Invariant 2 (No AI calls), Invariant 6 (Interfaces), Invariant 7 (No cycles), Invariant 14 (Pure synchronous) | **COMPLIANT** | Observer uses pure interface contracts and synchronous execution without goroutines or AI dependencies. |
| `EVENT_BUS_SPECIFICATION.md` | Event model reuse & decoupled notification layer | **COMPLIANT** | Reuses canonical `eventcontracts.Event` interface via composition. No duplicate fields. |
| `PERMISSION_ENGINE_SPECIFICATION.md` | Permission Engine bypass prevention | **COMPLIANT** | Observer performs no actions requiring authorization; strictly read-only observation. |
| `PROJECT_CONTEXT.md` | Layer 1 Runtime ownership & non-blocking synchronous rules | **COMPLIANT** | Aligned with Phase 1 runtime baseline. |