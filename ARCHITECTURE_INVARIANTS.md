# Architecture Invariants

This document enumerates the 15 immutable Architectural Invariants for Project Jervis. Every code contribution MUST comply with these invariants. Any pull request introducing a violation MUST be rejected.

---

## Invariant 1: Observer Never Calls AI
- **Requirement**: The Observer component (`internal/runtime/observer`) MUST NOT invoke, depend on, or route messages to any AI provider or LLM adapter.
- **Reason**: System health monitoring, metrics, and diagnostics must execute deterministically without external API latency, token costs, or network failure risks.
- **Violation Impact**: Cascading failure of system monitoring if AI providers fail or experience network latency.
- **Detection Method**: Automated AST package import analysis checking that `internal/runtime/observer` imports zero packages from `internal/aiprovider`.

---

## Invariant 2: Event Bus Never Calls AI
- **Requirement**: The Event Bus (`internal/runtime/eventbus`) MUST NOT invoke AI providers or depend on LLM inference for message routing.
- **Reason**: The asynchronous event broker must deliver messages with low, predictable latency (<1ms).
- **Violation Impact**: Severe event processing bottlenecks, thread starvation, and non-deterministic event handling.
- **Detection Method**: AST linter verification confirming zero imports of `internal/aiprovider` within `internal/runtime/eventbus`.

---

## Invariant 3: Memory Never Depends on AI
- **Requirement**: The Memory Engine (`internal/memory`) MUST manage working context, timeline logging, storage, retrieval, and compression without requiring AI inference.
- **Reason**: Memory persistence and context window management are foundational runtime services that must function offline.
- **Violation Impact**: Inability to record events, access context history, or load state when AI providers are unavailable.
- **Detection Method**: Static import linter ensuring `internal/memory` packages have zero dependency on `internal/aiprovider`.

---

## Invariant 4: Services Never Depend on a Specific Provider
- **Requirement**: Domain services (`internal/services`) MUST NOT import or reference vendor-specific AI packages (e.g. `openai`, `claude`, `gemini`). Services MUST interact strictly through the abstract `Provider` interface.
- **Reason**: Prevents vendor lock-in and allows seamless provider swapping.
- **Violation Impact**: Tight coupling to proprietary vendor SDKs, breaking modularity.
- **Detection Method**: AST linter checking that `internal/services` imports only the top-level `internal/aiprovider` interface contract.

---

## Invariant 5: Runtime Knows Nothing About AI Vendors
- **Requirement**: The Runtime layer (`internal/runtime`) MUST NOT contain references to specific AI vendor names or APIs.
- **Reason**: Runtime is the system owner and must remain completely decoupled from AI processing implementations.
- **Violation Impact**: Compromised architectural separation; runtime failure if vendor APIs change.
- **Detection Method**: AST import verification enforcing zero dependency from `internal/runtime` to `internal/aiprovider`.

---

## Invariant 6: Everything Communicates Through Interfaces
- **Requirement**: Inter-component and inter-layer communication MUST occur exclusively across typed Go interface contracts.
- **Reason**: Enables strict component decoupling, unit test mocking, and modular replacement.
- **Violation Impact**: High coupling, fragile codebase, and inability to unit test components in isolation.
- **Detection Method**: Static analysis validating that all cross-package fields use interface types.

---

## Invariant 7: No Cyclic Dependencies
- **Requirement**: The package dependency graph MUST remain strictly directed and acyclic (DAG). Dependencies MUST flow exclusively downward.
- **Reason**: Cyclic dependencies create initialization deadlocks and prevent clean package isolation.
- **Violation Impact**: Build failure, package initialization cycles, and architectural degradation.
- **Detection Method**: Automated `go vet` and import graph DAG verification in CI.

---

## Invariant 8: Runtime Works Without AI Providers
- **Requirement**: The Runtime, Memory Engine, and Service Layer MUST boot and operate completely if all AI providers are removed.
- **Reason**: Guarantees offline capability and deterministic local workflow automation.
- **Violation Impact**: System unusable when offline or when API limits/credentials expire.
- **Detection Method**: Integration test suite executed with zero AI providers configured.

---

## Invariant 9: Interfaces Contain Zero Business Logic
- **Requirement**: Client interface packages (`internal/interfaces`) MUST only format inputs, parse flags, and render responses. Business logic MUST reside in `internal/services`.
- **Reason**: Ensures consistent behavior across CLI, MCP, REST, and Desktop interfaces.
- **Violation Impact**: Duplicated logic, inconsistent behavior across different user interfaces.
- **Detection Method**: Code review and AST structural inspection.

---

## Invariant 10: Business Logic Never Depends on UI
- **Requirement**: Domain services (`internal/services`) and core runtime MUST NOT import UI or interface packages.
- **Reason**: Ensures business logic can run headlessly in background daemons or CLI scripts.
- **Violation Impact**: Inability to run Jervis in headless environments or CLI tools without GUI dependencies.
- **Detection Method**: AST linter confirming zero imports from `internal/services` to `internal/interfaces`.

---

## Invariant 11: Storage is Implementation Detail
- **Requirement**: Persistent storage drivers (SQLite, flat files) MUST be hidden behind abstract repository interfaces in `internal/memory/store`.
- **Reason**: Allows swapping underlying storage drivers without altering Memory Engine or Service logic.
- **Violation Impact**: Storage technology leak into business logic, preventing database migrations.
- **Detection Method**: Interface abstraction checking in package code reviews.

---

## Invariant 12: Plugins Cannot Bypass Runtime
- **Requirement**: Third-party plugins MUST NOT execute system operations or access memory without passing through Runtime Permission checks.
- **Reason**: Prevents malicious or buggy plugins from compromising host security.
- **Violation Impact**: Arbitrary code execution, unauthorized data exfiltration, system compromise.
- **Detection Method**: WASM sandbox runtime checks and AST verification of plugin host bindings.

---

## Invariant 13: Permissions Enforced Before Execution
- **Requirement**: Capability authorization MUST be evaluated by `internal/runtime/permissions` prior to executing any action or tool call.
- **Reason**: Protects against unauthorized file writes, network requests, or shell command execution.
- **Violation Impact**: Security vulnerabilities and unapproved destructive side effects.
- **Detection Method**: Unit and integration test suites verifying authorization checks prior to action execution.

---

## Invariant 14: Runtime Owns State
- **Requirement**: System session state, execution context, and configuration MUST be owned and mutated exclusively by the Runtime (`internal/runtime`).
- **Reason**: Prevents race conditions, inconsistent state, and unauthorized state mutation.
- **Violation Impact**: State corruption and non-deterministic application behavior.
- **Detection Method**: Thread-safety analysis and package mutation boundary verification.

---

## Invariant 15: AI Only Consumes Context and Produces Responses
- **Requirement**: AI Providers MUST act purely as stateless transformers consuming context payloads and returning structured text/tool responses.
- **Reason**: Keeps AI isolated downstream as a context consumer rather than system manager.
- **Violation Impact**: Loss of deterministic control and vulnerability to prompt injection hijacking.
- **Detection Method**: Interface contract inspection verifying AI adapter inputs and outputs.
