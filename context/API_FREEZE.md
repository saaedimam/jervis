# API Freeze Log

## Phase 1.1: Runtime Foundation Packages
- **Packages**:
  - `internal/runtime/contracts`
  - `internal/runtime/types`
  - `internal/runtime/errors`
  - `internal/runtime/version`
  - `internal/runtime/buildinfo`
  - `internal/runtime/config`
  - `internal/runtime/lifecycle`
- **Status**: Frozen
- **Breaking Changes**: Formal ADR required before modifying public exported APIs.

---

## Phase 1.2.1: Event Bus Foundation Packages (Refactored)
- **Packages**:
  - `internal/runtime/eventbus/contracts` (Interface contracts: `Publisher`, `Subscriber`, `Handler`, `Dispatcher`, `Validator`, `Middleware`, `EventFilter`)
  - `internal/runtime/eventbus/events` (`Envelope`, `Header`, `Priority` uint8, `EventType` string, `Builder` returning `*Envelope`, `Envelope.Clone()`)
  - `internal/runtime/eventbus/errors` (Canonical error variables)
- **Status**: Frozen (Final Refactor Complete)
- **Design Invariants**:
  - Zero `context.Context` usage in Event Bus public interfaces.
  - `Priority` represented as `type Priority uint8` (`PriorityLow`, `PriorityNormal`, `PriorityHigh`, `PriorityCritical`).
  - `EventType` defined as `type EventType string`.
  - `Builder.Build()` returns `(*Envelope, error)`.
  - `Envelope` enforces metadata immutability via defensive copies and supports `Clone()`.
- **Breaking Changes**: Formal ADR required before modifying public exported interfaces or envelope structures.

---

## Phase 1.2 Specifications: Event Bus Core Architecture
- **Specs**:
  - `EVENT_BUS_SPECIFICATION.md`
  - `EVENT_MODEL.md`
  - `EVENT_CONTRACTS.md`
  - `EVENT_IMPLEMENTATION_PLAN.md`
- **Status**: Specification Frozen
- **Breaking Changes**: Formal ADR required for contract or envelope schema modifications.
