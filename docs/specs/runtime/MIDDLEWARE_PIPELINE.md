# Event Bus Middleware Pipeline Specification

## 1. Pipeline Overview
The **Middleware Pipeline** defines the exact stage-by-stage flow of control during event publication.

```
   Publisher.Publish(event)
              │
              ▼
   ┌──────────────────────────────────────────┐
   │ Stage 1: Structural Event Validation     │
   │ - events.ValidateEvent(event)            │
   └──────────────────────────────────────────┘
              │
              ▼
   ┌──────────────────────────────────────────┐
   │ Stage 2: Registry Lookup                 │
   │ - registry.Lookup(event.Type())          │
   └──────────────────────────────────────────┘
              │
              ▼
   ┌──────────────────────────────────────────┐
   │ Stage 3: Build Middleware Chain (FIFO)   │
   │ - Chain: M1 -> M2 -> ... -> Mn -> Terminal│
   └──────────────────────────────────────────┘
              │
              ▼
   ┌──────────────────────────────────────────┐
   │ Stage 4: Pre-Dispatch Middleware Hooks   │
   │ - M1 Pre-Hook                            │
   │ - M2 Pre-Hook                            │
   └──────────────────────────────────────────┘
              │ (Short-Circuit Check)
              ▼
   ┌──────────────────────────────────────────┐
   │ Stage 5: Terminal Dispatcher Stage       │
   │ - Dispatcher.Dispatch(event, handlers)   │
   └──────────────────────────────────────────┘
              │
              ▼
   ┌──────────────────────────────────────────┐
   │ Stage 6: Post-Dispatch Middleware Hooks  │
   │ - M2 Post-Hook                           │
   │ - M1 Post-Hook                           │
   └──────────────────────────────────────────┘
              │
              ▼
   ┌──────────────────────────────────────────┐
   │ Stage 7: Return Execution Result         │
   └──────────────────────────────────────────┘
```

---

## 2. Stage Detailed Description

1. **Stage 1 (Validation)**: Verifies event structural integrity (`events.ValidateEvent`). Aborts if invalid.
2. **Stage 2 (Registry Lookup)**: Resolves and orders matching handlers from `Registry`.
3. **Stage 3 (Chain Construction)**: Wraps `Dispatcher.Dispatch` at the end of the registered middleware array.
4. **Stage 4 (Pre-Hooks)**: Invokes middleware pre-hook logic in FIFO order.
5. **Stage 5 (Terminal Dispatcher)**: Invokes `Dispatcher.Dispatch(event, handlers)` if no middleware short-circuited.
6. **Stage 6 (Post-Hooks)**: Invokes middleware post-hook logic in reverse (LIFO) order as stack frames return.
7. **Stage 7 (Return)**: Returns final aggregated error or `nil` back to publisher.
