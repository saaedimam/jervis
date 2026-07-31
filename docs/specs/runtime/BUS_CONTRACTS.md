# Event Bus Facade Contracts Specification

## 1. Overview
This document specifies the frozen interface contracts satisfied by the Event Bus Facade (`internal/runtime/eventbus`).

---

## 2. Core Interfaces

```go
package contracts

// Publisher defines the contract for publishing events to the Event Bus.
type Publisher interface {
	Publish(event Event) error
}

// Subscriber defines the contract for registering and unregistering event handlers.
type Subscriber interface {
	Subscribe(eventType string, handler Handler, priority uint8) error
	Unsubscribe(eventType string, handlerID string) error
}

// Bus defines the composite Event Bus facade contract.
type Bus interface {
	Publisher
	Subscribe(pattern string, handler Handler, priority uint8) (string, error)
	Unsubscribe(subscriptionID string) error
	Use(middleware ...Middleware)
	Count() int
}
```

---

## 3. Invariants & Rules

1. **No `context.Context`**: Method signatures receive parameters directly without `context.Context`.
2. **Immutable Input**: Published `Event` instances MUST NOT be mutated by the Event Bus.
3. **Synchronous Execution**: `Publish` returns only after validation, middleware execution, and handler dispatch complete.
4. **Deterministic Signatures**: Exported method names, parameter types, and return values are permanently frozen.
