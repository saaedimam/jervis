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
  - Defensive slice copies returned for all lookup and snapshot APIs.
- **Breaking Changes**: Formal ADR required before modifying public exported registry surfaces.

---

## Phase 1.2.4: Event Bus Dispatcher Package
- **Packages**:
  - `internal/runtime/eventbus/dispatcher` (`Dispatcher`, `NewDispatcher()`, `AggregateError`, `NewAggregateError()`, `MaxDispatchDepth = 16`)
- **Status**: Frozen
- **Design Invariants**:
  - 7-stage execution pipeline (Validation -> Recursion Guard -> Handler Lookup -> Priority Sort -> Panic-Protected Handler Loop -> Error Aggregation -> Return).
  - Deterministic priority ordering (`Priority` DESC, `Seq` ASC, `ID` ASC).
  - Panic recovery isolation per handler (`recover()`) wrapping panics in `errs.ErrHandlerFailure`.
  - Continue-on-Error policy with composite `AggregateError` returned to caller.
  - Recursion depth capped at `MaxDispatchDepth = 16`.
- **Breaking Changes**: Formal ADR required before modifying public exported dispatcher types.

---

## Phase 1.2.5 Specifications: Event Bus Middleware Architecture
- **Specs**:
  - `MIDDLEWARE_SPECIFICATION.md`
  - `MIDDLEWARE_CONTRACTS.md`
  - `MIDDLEWARE_PIPELINE.md`
  - `MIDDLEWARE_ORDERING.md`
- **Status**: Specification Frozen (Zero Runtime Implementation Code)
- **Design Invariants**:
  - Canonical FIFO registration order.
  - Explicit `next()` closure semantics with short-circuit authorization support.
  - Panic recovery isolation per middleware block.
  - Dispatcher remains sole owner of subscriber handler execution.
- **Breaking Changes**: Formal ADR required for contract or middleware pipeline modifications.
