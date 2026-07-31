# Event Dispatch Pipeline Specification

## 1. Overview
The **Event Dispatch Pipeline** defines the sequential stages executed whenever an event is published and dispatched through the Jervis Event Bus.

---

## 2. Pipeline Stages

```
   Publisher: Publish(event)
              │
              ▼
   ┌──────────────────────────────────────────┐
   │ Stage 1: Structural Event Validation     │
   │ - events.ValidateEvent(event)            │
   └──────────────────────────────────────────┘
              │
              ▼
   ┌──────────────────────────────────────────┐
   │ Stage 2: Dispatch Recursion Guard        │
   │ - Depth <= MaxDispatchDepth (16)         │
   └──────────────────────────────────────────┘
              │
              ▼
   ┌──────────────────────────────────────────┐
   │ Stage 3: Handler Ordering                │
   │ - Sort by Priority DESC                  │
   │ - Sort by Sequence ASC                   │
   └──────────────────────────────────────────┘
              │
              ▼
   ┌──────────────────────────────────────────┐
   │ Stage 4: Pre-Dispatch Middleware Hooks   │
   │ - Executed in order                      │
   └──────────────────────────────────────────┘
              │
              ▼
   ┌──────────────────────────────────────────┐
   │ Stage 5: Synchronous Handler Loop        │
   │ For each Handler h in Handlers:          │
   │   ┌──────────────────────────────────┐   │
   │   │ Panic Isolation Wrapper (defer)  │   │
   │   │ err := h.Handle(event)           │   │
   │   └──────────────────────────────────┘   │
   │   Record errors, Continue to next h      │
   └──────────────────────────────────────────┘
              │
              ▼
   ┌──────────────────────────────────────────┐
   │ Stage 6: Post-Dispatch Middleware Hooks  │
   │ - Executed in reverse order              │
   └──────────────────────────────────────────┘
              │
              ▼
   ┌──────────────────────────────────────────┐
   │ Stage 7: Composite Error Normalization   │
   │ - Return aggregated dispatch errors      │
   └──────────────────────────────────────────┘
```

---

## 3. Stage Details

### Stage 1: Structural Event Validation
- Validates that `event` is non-nil, `ID` is non-empty, `Type` is lowercase dot-separated, `Timestamp` is non-zero, `Priority` is valid (`0..3`), `Payload` is non-nil, and `Version` is non-empty.
- **Outcome**: If validation fails, dispatch stops immediately with `ErrValidationFailed`.

### Stage 2: Dispatch Recursion Guard
- Checks current stack depth counter.
- **Cap**: `MaxDispatchDepth = 16`.
- **Outcome**: If depth > 16, dispatch stops immediately with `ErrDispatchFailed`.

### Stage 3: Handler Ordering
- Sorts `handlers` deterministically:
  - Higher priority (`PriorityCritical`) executed before lower priority (`PriorityLow`).
  - Equal priority handlers executed in registration sequence order (FIFO).

### Stage 4: Pre-Dispatch Middleware Hooks
- Calls registered middleware pre-hooks sequentially.

### Stage 5: Synchronous Handler Loop
- Iterates over handlers. Each handler invocation is wrapped in a `defer recover()` block.
- **Policy**: Continue-on-Error. If a handler returns an error or panics, the error is recorded into the execution accumulator and the loop proceeds to the next handler.

### Stage 6: Post-Dispatch Middleware Hooks
- Calls registered middleware post-hooks.

### Stage 7: Composite Error Normalization
- If any handler failed or panicked, returns a composite error aggregating all failures.
