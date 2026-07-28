# Event Bus Implementation Plan (Phase 1.2)

## 1. Overview
This document specifies the exact sequential implementation sub-phases for Phase 1.2 (Event Bus). Implementation must proceed strictly step-by-step. No phase may commence until all prior phase exit criteria are fully satisfied.

---

## 2. Phase Breakdown

### 2.1 Phase 1.2.1: Core Contracts & Event Envelope
- **Goal**: Implement the event envelope data structure, event builder, canonical errors, and contract interfaces.
- **Files to Create**:
  - `internal/runtime/eventbus/contracts/contracts.go`
  - `internal/runtime/eventbus/events/envelope.go`
  - `internal/runtime/eventbus/events/envelope_test.go`
  - `internal/runtime/eventbus/errors/errors.go`
  - `internal/runtime/eventbus/errors/errors_test.go`
- **Exit Criteria**:
  - 100% test coverage for envelope construction, immutability, and validation.
  - Zero dependencies outside `internal/runtime/{types, errors, version}` and standard library.
- **Tests**: Table-driven tests for envelope creation, header validation, and string formatters.

---

### 2.2 Phase 1.2.2: Registry & Synchronous Dispatcher
- **Goal**: Implement thread-safe subscription registry and synchronous priority-ordered dispatcher.
- **Files to Create**:
  - `internal/runtime/eventbus/registry/registry.go`
  - `internal/runtime/eventbus/registry/registry_test.go`
  - `internal/runtime/eventbus/dispatcher/dispatcher.go`
  - `internal/runtime/eventbus/dispatcher/dispatcher_test.go`
- **Exit Criteria**:
  - 100% test coverage for subscription registration, unregistration, duplicate handling, and priority execution.
  - Handlers execute strictly in priority order.
  - Panic recovery isolation verified via synthetic panicking handlers.
- **Tests**: Priority ordering tests, panic recovery tests, concurrency-safe registration tests (`-race`).

---

### 2.3 Phase 1.2.3: Validator, Middleware Chain, & Bus Integration
- **Goal**: Implement event validation rules, middleware execution pipeline, and top-level `Bus` facade.
- **Files to Create**:
  - `internal/runtime/eventbus/validation/validator.go`
  - `internal/runtime/eventbus/validation/validator_test.go`
  - `internal/runtime/eventbus/middleware/chain.go`
  - `internal/runtime/eventbus/middleware/chain_test.go`
  - `internal/runtime/eventbus/bus.go`
  - `internal/runtime/eventbus/bus_test.go`
- **Exit Criteria**:
  - 100% test coverage for complete bus publish and subscribe workflow.
  - Middleware chain executes pre-and-post hooks deterministically.
  - All quality gates pass (`go fmt`, `go test -race`, `golangci-lint`).
- **Tests**: End-to-end event publishing tests, middleware chain order tests, benchmark tests (>100,000 ops/sec).
