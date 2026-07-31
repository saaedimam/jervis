# Event Bus Middleware Contracts Specification

## 1. Overview
This document specifies the frozen interface contracts and method signatures for the Jervis Event Bus Middleware component (`internal/runtime/eventbus/middleware`).

---

## 2. Interface Contracts

```go
package contracts

// Middleware defines pipeline hooks executed before and after event handling.
type Middleware interface {
	Execute(event Event, next func(event Event) error) error
}

// Chain defines the composite middleware chain interface.
type Chain interface {
	Use(middleware ...Middleware)
	Execute(event Event, terminal func(event Event) error) error
}
```

---

## 3. Method Semantics & Invariants

### 3.1 `Middleware.Execute`
- **Signature**: `Execute(event Event, next func(event Event) error) error`
- **Parameters**:
  - `event`: Immutable `contracts.Event` envelope.
  - `next`: Chain continuation closure representing the next middleware or terminal dispatcher call.
- **Invariants**:
  - `Execute` MUST NOT modify `event`.
  - `Execute` MUST return an `error` if processing fails or short-circuits.
  - `Execute` MUST NOT accept `context.Context`.

### 3.2 `Chain.Use`
- **Signature**: `Use(middleware ...Middleware)`
- **Behavior**: Appends middleware to the end of the registration list in FIFO order.

### 3.3 `Chain.Execute`
- **Signature**: `Execute(event Event, terminal func(event Event) error) error`
- **Behavior**: Constructs the recursive closure chain and invokes the first middleware in FIFO sequence, passing `terminal` (which invokes `Dispatcher.Dispatch`) as the final closure.
