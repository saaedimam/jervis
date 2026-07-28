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

## Phase 1.2.1: Event Bus Foundation Packages
- **Packages**:
  - `internal/runtime/eventbus/contracts`
  - `internal/runtime/eventbus/events`
  - `internal/runtime/eventbus/errors`
- **Status**: Frozen
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
