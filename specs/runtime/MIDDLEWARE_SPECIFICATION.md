# Event Bus Middleware Architecture Specification

## 1. Overview
This document establishes the canonical specification for the **Middleware Component** (`internal/runtime/eventbus/middleware`) within the Jervis Event Bus.

Middleware provides a deterministic, composable interceptor pipeline around event publication and dispatch. Middleware allows cross-cutting concerns—such as structural validation, auditing, permission checks, and execution telemetry—to execute transparently before and after event handlers run.

---

## 2. Core Principles & Boundaries

### 2.1 Owned Responsibilities
1. **Event Interception**: Intercept event publication before and after handler dispatch.
2. **Deterministic Chain Execution**: Execute middleware hooks in strict, predictable FIFO (registration sequence) order.
3. **Short-Circuit Authorization**: Allow authorization or auditing middleware to abort event dispatch prior to handler invocation.
4. **Panic Recovery & Isolation**: Capture and convert middleware panics into canonical error structures.

### 2.2 Explicit Non-Responsibilities & Architectural Boundaries
- **No Handler Ownership**: Middleware MUST NOT own or directly invoke subscriber handlers. Handler invocation remains strictly owned by the `Dispatcher`.
- **No Event Mutation**: Middleware MUST NOT alter the immutable event envelope (`contracts.Event`).
- **No Persistence or Network Calls**: Middleware operates strictly in-process and in-memory.
- **No AI Provider Awareness**: Middleware MUST NOT import or invoke `internal/aiprovider` (Invariant 2).

---

## 3. Middleware Lifecycle

The Middleware lifecycle follows an "onion" wrapper pattern:

```
Publisher.Publish(event)
   │
   ▼
[M1 Pre-Hook] (Auditing / Telemetry)
   │
   ▼
[M2 Pre-Hook] (Permission Check)
   │
   ▼
[Terminal Hook] ---> Dispatcher.Dispatch(event, handlers)
   │                      │
   │                      ▼
   │                 Subscriber Handlers Executed
   │                      │
   │                      ▼
[M2 Post-Hook] <──────────┘
   │
   ▼
[M1 Post-Hook]
   │
   ▼
Return Result to Publisher
```

---

## 4. Execution Lifecycle & Next() Semantics

### 4.1 Middleware Contract
All middleware MUST implement the canonical `contracts.Middleware` interface:

```go
type Middleware interface {
    Execute(event Event, next func(event Event) error) error
}
```

### 4.2 Next() Function Semantics
- The `next` parameter is a closure representing the remaining downstream chain of middleware ending in `Dispatcher.Dispatch`.
- Calling `next(event)` passes control to the next middleware or terminal dispatcher.
- Returning the result of `next(event)` allows post-processing logic (Post-Hooks) to execute after downstream components return.

---

## 5. Short-Circuit Policy

1. **Short-Circuit Trigger**: A middleware MAY return an `error` without invoking `next(event)`.
2. **Effects of Short-Circuiting**:
   - Downstream middleware pre-hooks MUST NOT execute.
   - Subscriber handlers MUST NOT execute (`Dispatcher.Dispatch` is never reached).
   - Outer (upstream) middleware receive the returned short-circuit error in their post-hook phase and return it back up the call stack to the publisher.

---

## 6. Panic Policy & Recovery

1. **Middleware Panic Isolation**: The middleware chain runner wraps every `Middleware.Execute` call in a `defer recover()` block.
2. **Panic Conversion**: If a middleware panics during execution:
   - The panic is trapped via `recover()`.
   - The panic value is converted into a wrapped `errs.ErrHandlerFailure` error.
   - `next(event)` is NOT called if the panic occurs before `next()`.
   - The error is returned up the middleware chain.
3. **Dispatcher Panic Interaction**: If a panic occurs inside a subscriber handler during `Dispatcher.Dispatch`, the Dispatcher's own panic recovery converts it into an error and returns it to `next()`, allowing middleware post-hooks to observe the failure safely.

---

## 7. Performance & Future Compatibility

- **Pure Synchronous Execution**: Zero goroutines, zero channels, zero mutexes, zero `context.Context` parameters.
- **$O(M)$ Complexity**: For $M$ registered middleware, execution adds $O(M)$ stack frame allocations.
- **Future Async Compatibility**: As synchronous call chains preserve call stack frames, future worker pool dispatchers can integrate underneath the terminal `next` closure without modifying middleware signatures.
