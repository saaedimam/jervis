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
  - `internal/runtime/eventbus/contracts` (Interface contracts: `Publisher`, `Subscriber`, `Handler`, `Dispatcher`, `Validator`, `Middleware`, `EventFilter`)
  - `internal/runtime/eventbus/events` (`Envelope`, `Header`, `Priority` uint8, `EventType` string, `Builder` returning `*Envelope`, `Envelope.Clone()`)
  - `internal/runtime/eventbus/errors` (Canonical error variables)
- **Status**: Frozen
- **Breaking Changes**: Formal ADR required before modifying public exported interfaces or envelope structures.

---

## Phase 1.2.2: Event Bus Subscription Registry Packages
- **Packages**:
  - `internal/runtime/eventbus/subscription` (`SubscriptionID`, `Subscription`, `New()`, `NewWithSeq()`)
  - `internal/runtime/eventbus/registry` (`Registry`, `NewRegistry()`, `Register()`, `Unregister()`, `Lookup()`, `LookupExact()`, `LookupPattern()`, `Contains()`, `Count()`, `Clear()`, `Snapshot()`, `ValidatePattern()`, `MatchesPattern()`)
- **Status**: Frozen
- **Design Invariants**:
  - Deterministic priority ordering (Priority DESC then registration sequence ASC).
  - Pure Go pattern matching (`*`, `prefix.*`, `prefix*`, exact string). No regex or glob libraries.
  - Zero mutexes, channels, goroutines, or `context.Context` dependencies.
  - Defensive slice copies returned for all lookup and snapshot methods.
- **Breaking Changes**: Formal ADR required before modifying public exported registry surfaces.

---

## Phase 1.2 Specifications: Event Bus Core Architecture
- **Specs**:
  - `EVENT_BUS_SPECIFICATION.md`
  - `EVENT_MODEL.md`
  - `EVENT_CONTRACTS.md`
  - `EVENT_IMPLEMENTATION_PLAN.md`
- **Status**: Specification Frozen
- **Breaking Changes**: Formal ADR required for contract or envelope schema modifications.
