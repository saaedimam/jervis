# Engineering Principles

## 1. Engineering Philosophy
Project Jervis is engineered as a local-first personal OS and developer automation runtime. The core engineering philosophy is built upon three pillars:
- **Determinism Over Probability**: The core runtime, event bus, memory engine, and service layer must execute deterministically. AI probabilistic models are strictly isolated downstream and never dictate system flow or control state.
- **Zero Runtime AI Lock-in**: The system must boot, manage state, record timeline events, execute schedules, and run domain services even if every AI provider is removed or offline.
- **Strict Single-Direction Layering**: Dependencies flow exclusively downward (`OS -> Runtime -> Memory Engine -> Service Layer -> AI Provider Layer -> Interfaces`). Upward calls, cyclic dependencies, or layer-bypassing are architectural defects.

## 2. Deterministic Development
- **Predictable State Transitions**: All state mutations must occur via explicit events processed by the Runtime Event Bus and logged to the immutable Timeline ledger.
- **Idempotency**: Service operations and scheduling triggers must be idempotent where possible. Re-processing an event from the timeline must yield predictable, reproducible state.
- **Hermetic Execution**: Local automation scripts and services must operate within sandboxed boundaries enforced by the Runtime Permissions module.

## 3. Architecture Ownership
- **Runtime as System Owner**: The `Runtime` module (`internal/runtime`) owns process lifecycles, configuration, permissions, observer, session state, and event routing.
- **Non-Bypassable Security**: Neither client interfaces nor third-party plugins can bypass the Runtime Permissions check prior to action execution.
- **Storage Independence**: Storage layers (SQLite, flat files, key-value stores) are private implementation details of the Memory Engine and must be hidden behind abstract repository interfaces.

## 4. Dependency Policy Overview
- **Minimal External Footprint**: Prefer standard library implementations over third-party packages.
- **Zero Cgo Policy for Core Binary**: Ensure maximum portability and static cross-compilation by avoiding C dependencies in the runtime core.
- **Vetted Third-Party Code**: All external dependencies require explicit audit for security vulnerabilities, licensing compliance (permissive only), and active maintenance.

## 5. Versioning Policy
- **Semantic Versioning 2.0.0**: All release artifacts follow `MAJOR.MINOR.PATCH` formatting.
- **Public API Boundary**: The public contracts exposed in `pkg/plugin/` and CLI command signatures define the public API.
- **Breaking Changes**: Breaking changes to CLI contracts or plugin interfaces require a increment of the `MAJOR` version.

## 6. Backward Compatibility & Deprecation Policy
- **Schema Compatibility**: Memory engine schemas and persistent configuration files must maintain backward compatibility across minor versions. Schema migrations must be automatic and non-destructive.
- **Deprecation Cycle**: Any interface, configuration parameter, or command marked for deprecation must remain functional for at least one full `MINOR` release cycle with warning diagnostics emitted before removal in the next `MAJOR` release.

## 7. Documentation Policy
- **Self-Documenting Code**: Code must be clear, typed, and structured so that implementation intent is obvious.
- **Package & Interface Documentation**: Every public package, struct, and interface must be fully documented using standardized godoc comments.
- **Architecture Synchronization**: Any modification to component boundaries, interfaces, or system rules requires simultaneous updates to `ARCHITECTURE.md`, `PROJECT_STRUCTURE.md`, and `AI_CONTEXT.md`.
