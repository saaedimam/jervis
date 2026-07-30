# Engineering Principles

> Last Updated: 2026-07-31
> Owner: @saaedimam

## 1. Engineering Philosophy

Project Jervis is engineered as a local-first personal OS and developer automation runtime. The core engineering philosophy is built upon three pillars:

- **Determinism Over Probability**: The core runtime, event bus, memory engine, and service layer must execute deterministically. AI probabilistic models are strictly isolated downstream and never dictate system flow or control state.
- **Zero Runtime AI Lock-in**: The system must boot, manage state, record timeline events, execute schedules, and run domain services even if every AI provider is removed or offline.
- **Strict Single-Direction Layering**: Dependencies flow exclusively downward (`OS -> Runtime -> Memory Engine -> Service Layer -> AI Provider Layer -> Interfaces`). Upward calls, cyclic dependencies, or layer-bypassing are architectural defects.

## 2. Core Design Principles

1. **Local-First**: Prioritize operating entirely on the user's local machine for privacy, security, and low latency.
2. **Runtime Ownership**: The Jervis Runtime owns all system state, event dispatching, and process lifecycles. AI never owns the runtime.
3. **Transparency & Auditability**: System events, permissions, and tool execution logs are recorded on an immutable timeline.
4. **User in Control**: Actions with side effects require authorization enforced by the Runtime Permissions module.
5. **Modularity & Decoupling**: Components interact strictly through defined interfaces. Storage and AI providers are implementation details.
6. **Graceful Degradation**: The system operates deterministically even when AI providers are offline or removed.

## 3. Mandatory Design Rules

The following 15 rules must **NEVER** be violated under any circumstances:

1. **Observer never calls AI.**
2. **Event Bus never calls AI.**
3. **Memory never depends on AI.**
4. **Services never depend on a specific provider.**
5. **Runtime knows nothing about OpenAI, Claude, Gemini or any vendor.**
6. **Everything communicates through interfaces.**
7. **No cyclic dependency.**
8. **The runtime must continue working if every AI provider is removed.**
9. **Interfaces never contain business logic.**
10. **Business logic never depends on UI.**
11. **Storage is implementation detail.**
12. **Plugins cannot bypass Runtime.**
13. **Permissions are enforced before execution.**
14. **Runtime owns state.**
15. **AI only consumes context and produces responses.**

## 4. Deterministic Development

- **Predictable State Transitions**: All state mutations must occur via explicit events processed by the Runtime Event Bus and logged to the immutable Timeline ledger.
- **Idempotency**: Service operations and scheduling triggers must be idempotent where possible. Re-processing an event from the timeline must yield predictable, reproducible state.
- **Hermetic Execution**: Local automation scripts and services must operate within sandboxed boundaries enforced by the Runtime Permissions module.

## 5. Architecture Ownership

- **Runtime as System Owner**: The `Runtime` module (`internal/runtime`) owns process lifecycles, configuration, permissions, observer, session state, and event routing.
- **Non-Bypassable Security**: Neither client interfaces nor third-party plugins can bypass the Runtime Permissions check prior to action execution.
- **Storage Independence**: Storage layers (SQLite, flat files, key-value stores) are private implementation details of the Memory Engine and must be hidden behind abstract repository interfaces.

## 6. Dependency Policy Overview

- **Minimal External Footprint**: Prefer standard library implementations over third-party packages.
- **Zero Cgo Policy for Core Binary**: Ensure maximum portability and static cross-compilation by avoiding C dependencies in the runtime core.
- **Vetted Third-Party Code**: All external dependencies require explicit audit for security vulnerabilities, licensing compliance (permissive only), and active maintenance.

## 7. Versioning Policy

- **Semantic Versioning 2.0.0**: All release artifacts follow `MAJOR.MINOR.PATCH` formatting.
- **Public API Boundary**: The public contracts exposed in `pkg/plugin/` and CLI command signatures define the public API.
- **Breaking Changes**: Breaking changes to CLI contracts or plugin interfaces require an increment of the `MAJOR` version.

## 8. Backward Compatibility & Deprecation Policy

- **Schema Compatibility**: Memory engine schemas and persistent configuration files must maintain backward compatibility across minor versions. Schema migrations must be automatic and non-destructive.
- **Deprecation Cycle**: Any interface, configuration parameter, or command marked for deprecation must remain functional for at least one full `MINOR` release cycle with warning diagnostics emitted before removal in the next `MAJOR` release.

## 9. Documentation Policy

- **Self-Documenting Code**: Code must be clear, typed, and structured so that implementation intent is obvious.
- **Package & Interface Documentation**: Every public package, struct, and interface must be fully documented using standardized godoc comments.
- **Architecture Synchronization**: Any modification to component boundaries, interfaces, or system rules requires simultaneous updates to `ARCHITECTURE.md`, `PROJECT_STRUCTURE.md`, and relevant context documents.

## 10. Coding Standards

1. **Language & Style**: Strict adherence to the standard style guide of the chosen implementation language.
2. **Type Safety**: Enforce static typing and explicit interfaces across all layer boundaries.
3. **Documentation**: All public package contracts, layer interfaces, and services must be fully documented.
4. **Testing**: Unit tests are required for Runtime, Memory Engine, and Services. Integration tests for AI Provider adapters.
5. **Error Handling**: Use structured error types. Errors must be propagated through the Event Bus or returned cleanly; never swallowed silently.
