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

## Phase 1.2.6: Event Bus Middleware Package
- **Packages**:
  - `internal/runtime/eventbus/middleware` (`Chain`, `NewChain()`, `Use()`, `Middlewares()`, `Count()`, `Execute()`, `Func`)
- **Status**: Frozen
- **Design Invariants**:
  - FIFO registration sequence execution order (FIFO entering, LIFO exiting).
  - Explicit `next()` closure semantics with short-circuit error return support.
  - Panic recovery isolation per middleware block converting panics into `errs.ErrHandlerFailure`.
  - Dispatcher remains sole owner of subscriber handler execution.
- **Breaking Changes**: Formal ADR required before modifying public exported middleware interfaces.

---

## Phase 1.2.8: Event Bus Facade Package
- **Packages**:
  - `internal/runtime/eventbus` (`EventBus`, `New()`, `Publish()`, `Subscribe()`, `Unsubscribe()`, `Use()`, `Count()`)
- **Status**: Frozen (Complete Event Bus Engine Implementation)
- **Design Invariants**:
  - Canonical single entry point facade orchestrating `events.ValidateEvent`, `registry.Registry`, `middleware.Chain`, and `dispatcher.Dispatcher`.
  - 100% synchronous execution without goroutines, channels, mutexes, or `context.Context`.
  - 100.0% statement coverage achieved across all 8 eventbus packages with `-race`.
- **Breaking Changes**: Formal ADR required before modifying public facade method signatures.

---

## Phase 1.3.1: Runtime Permission Engine Core Foundation Packages
- **Packages**:
  - `internal/runtime/permissions/contracts` (`Capability`, `Decision`, `Effect`, `Validator`, `Rule`, `Policy`)
  - `internal/runtime/permissions/capability` (`Capability`, `New()`, `Subject`, `Resource`, `Action`, `IsZero()`, `String()`)
  - `internal/runtime/permissions/decision` (`Decision`, `NewAllow()`, `NewDeny()`, `IsAllowed()`, `Effect()`, `Reason()`, `String()`)
  - `internal/runtime/permissions/validator` (`Validator`, `New()`, `Validate()`)
  - `internal/runtime/permissions/errors` (Canonical error variables)
- **Status**: Frozen
- **Design Invariants**:
  - Pure value objects and structural validation without state or policy storage.
  - 100.0% statement coverage achieved across all 5 foundation packages with `-race`.
- **Breaking Changes**: Formal ADR required before modifying permission core types or contracts.

---

## Phase 1.3.2: Runtime Permission Engine Rule & Policy Domain Models
- **Packages**:
  - `internal/runtime/permissions/rule` (`Rule`, `New()`, `ID()`, `Subject()`, `Resource()`, `Action()`, `Effect()`, `Description()`, `Evaluate()`, `Validate()`, `IsZero()`, `String()`)
  - `internal/runtime/permissions/policy` (`Policy`, `New()`, `ID()`, `Name()`, `Description()`, `Version()`, `Rules()`, `Count()`, `Validate()`, `IsZero()`, `String()`)
- **Status**: Frozen
- **Design Invariants**:
  - Immutable domain models with defensive rules slice copies in constructors and getters.
  - Wildcard pattern matching (`*` and `prefix*`) in `Rule.Evaluate`.
  - 100.0% statement coverage achieved across all permission domain packages with `-race`.
- **Breaking Changes**: Formal ADR required before modifying Rule or Policy exported interfaces.

---

## Phase 1.3.3: Runtime Permission Engine Policy Registry Package
- **Packages**:
  - `internal/runtime/permissions/registry` (`Registry`, `New()`, `Register()`, `Unregister()`, `Get()`, `Policies()`, `Snapshot()`, `Count()`, `Contains()`, `Clear()`)
- **Status**: Frozen
- **Design Invariants**:
  - In-memory policy storage component strictly decoupled from evaluation logic.
  - Deterministic Policy ID ascending sort order returned by `Policies()` and `Snapshot()`.
  - Defensive slice copies returned for all policy query methods.
  - 100.0% statement coverage achieved across all 8 permissions packages with `-race`.
- **Breaking Changes**: Formal ADR required before modifying Registry exported method signatures.

---

## Phase 1.3.4: Runtime Permission Engine Evaluator Package
- **Packages**:
  - `internal/runtime/permissions/engine` (`Engine`, `New()`, `Authorize()`, `Registry()`)
- **Status**: Frozen (Complete Permission Evaluation Engine Implementation)
- **Design Invariants**:
  - Strict 6-stage evaluation pipeline (Capability Validation -> Policy Retrieval -> Rule Evaluation -> Deny Short-Circuit -> Allow Accumulation -> Default Deny Fallback).
  - Explicit Deny-First override rule precedence ("explicit deny" reason).
  - Default Deny fallback policy ("default deny policy enforced" reason).
  - 100% synchronous execution without goroutines, channels, mutexes, or `context.Context`.
  - 100.0% statement coverage achieved across all 9 permissions packages with `-race`.
- **Breaking Changes**: Formal ADR required before modifying Engine exported method signatures.

---

## Phase 1.4.0: Runtime Observer Subsystem (Architectural Freeze)
- **Status**: Frozen (Architecture & Interface Contracts)
- **Design Invariants**:
  - **Read-Only Observation**: Passive event monitoring without state mutation or EventBus feedback loops.
  - **Compositional Wrapping**: `Notification` wraps the canonical `eventcontracts.Event` interface (Composition over Duplication).
  - **Deterministic FIFO Dispatch**: Strict registration sequence execution order (FIFO).
  - **Panic Isolation**: Individual handler recovery (`recover()`) ensuring system-wide notification integrity (Continue-on-Error).
  - **Error Aggregation**: Composite `AggregateError` for multi-handler failure reporting.
  - **Pure Synchronous Execution**: 100% synchronous logic without goroutines, channels, or mutexes.
- **Architectural Packages (Drafted)**:
  - `internal/runtime/observer/contracts`
  - `internal/runtime/observer/notification`
  - `internal/runtime/observer/errors`
  - `internal/runtime/observer/registry`
  - `internal/runtime/observer/dispatcher`
- **Breaking Changes**: Formal ADR required before modifying frozen Observer specifications or interface contracts.
