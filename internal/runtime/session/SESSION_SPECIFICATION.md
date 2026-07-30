# Runtime Session Specification (Phase 1.6.0 - DRAFT)

## 1. Overview & Responsibilities

The Runtime Session Engine is the canonical state and context management layer of Project Jervis. It isolates workspace variables, user context, and session-specific metadata within a deterministic execution environment.

### Core Responsibilities
1. **State Isolation**: Ensure each session has its own private metadata store.
2. **Metadata Management**: Provide thread-safe key-value storage for session-specific information.
3. **Session Lifecycle**: Manage the transition of sessions through canonical states (`Created`, `Running`, `Stopped`).
4. **Deterministic Registry**: Maintain a strict order of sessions and provide predictable lookup capabilities.
5. **Zero Side Effects (on system state)**: The engine coordinates session lifecycle without modifying underlying OS resources directly.

---

## 2. Architecture & Subsystem Boundaries

```
internal/runtime/session/
├── contracts/        # Interfaces (Session, Registry, Manager)
├── model/            # Concrete Session implementation
├── errors/           # Canonical session errors
├── registry/         # Thread-safe session storage
└── session.go        # Facade orchestrating sessions and registry
```

---

## 3. Implementation Invariants

- **Thread-Safety**: All metadata operations and registry lookups MUST use proper synchronization (Mutex/RWMutex).
- **Defensive Copies**: `Metadata()` and `All()` methods MUST return defensive copies to prevent external mutation of internal state.
- **NoAI Policy**: The session engine must not depend on any AI Provider services.
- **Layer 1 Ownership**: The session engine must not depend on Memory, Service, or AI layers.

---

## 4. Exit Criteria (Phase 1.6)

- [x] Implementation of core contracts and registry.
- [x] Thread-safe metadata management with isolation.
- [x] 100% statement coverage across registry and session logic.
- [x] Facade for session lifecycle (Create, Get, Close).
