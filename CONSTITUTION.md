# Project Jervis Governance Constitution

## 1. Preamble & Scope
This Constitution establishes the supreme governance framework for Project Jervis. All development activities, architectural decisions, code contributions, and release operations MUST strictly comply with the provisions set forth herein.

---

## 2. Mission & Vision
- **Mission**: To build a deterministic, local-first personal OS and automation runtime that empowers developers while preserving data privacy and execution integrity.
- **Vision**: To provide a unified multi-interface runtime where system state, event distribution, and execution security are owned exclusively by the deterministic Runtime, and external AI providers act strictly as pluggable processing tools.

---

## 3. Core Values & Engineering Philosophy
- **Determinism First**: The system MUST execute predictably. Probabilistic models MUST NOT control system lifecycles, event dispatching, or permission evaluations.
- **Runtime Sovereignty**: The Jervis Runtime MUST maintain absolute ownership over state, sessions, and capabilities. AI MUST NEVER own the runtime loop.
- **Local Sovereignty & Privacy**: Data MUST remain local to the user's system by default. Cloud access MUST be opt-in and explicitly authorized.
- **Zero AI Dependency for Core Functions**: The Runtime, Memory Engine, and Service Layer MUST remain 100% operational even if all AI providers are disabled or removed.

---

## 4. Architecture Ownership
- The **Principal Software Architect** SHALL hold ultimate authority over system architecture, component boundaries, and layer definitions.
- Components MUST adhere strictly to the 5-tier single-direction hierarchy (`OS -> Runtime -> Memory Engine -> Service Layer -> AI Provider Layer -> Interfaces`).
- Any modification to layer boundaries, security models, or runtime ownership MUST be approved through a formal Architecture Decision Record (ADR).

---

## 5. Decision Process & ADR Governance
- Technical changes affecting component boundaries, interfaces, dependencies, or security models MUST be proposed via an ADR.
- An ADR MUST be reviewed by the Architecture Board and SHALL require unanimous approval before any code implementation begins.
- Minor refactorings that preserve existing contracts MAY proceed with standard peer code review.

---

## 6. Breaking Change Policy
- Breaking changes to public interfaces (`pkg/plugin/`), storage schemas, or CLI command flags MUST increment the `MAJOR` version under Semantic Versioning 2.0.0.
- Deprecated capabilities MUST be maintained for at least one full `MINOR` release cycle before removal.
- Migration scripts MUST be provided for persistent schema changes to guarantee non-destructive state updates.

---

## 7. Definition of Done & Acceptance Criteria
A task or pull request SHALL NOT be marked as complete or merged into the main branch unless it satisfies the following criteria:
1. **Rule Compliance**: The code MUST NOT violate any of the 15 Architecture Invariants.
2. **Quality Gates Passed**: All mandatory quality gates (linting, tests, coverage, security, architecture validation) MUST pass without warnings.
3. **Test Coverage**: Unit and integration test coverage MUST meet or exceed the 85% project threshold.
4. **Documentation**: All public APIs, structs, interfaces, and CLI flags MUST be fully documented.
5. **No Regression**: Performance benchmarks MUST verify that CLI cold startup remains under 15ms and memory footprint remains under 25MB.

---

## 8. Project Lifecycle & Invariants
- The project lifecycle SHALL consist of strictly ordered implementation phases as defined in `08_IMPLEMENTATION_ORDER.md`.
- No implementation phase MAY commence until the exit criteria of all preceding phases are satisfied.
- The 15 Architecture Invariants defined in `ARCHITECTURE_INVARIANTS.md` MUST remain permanently enforced throughout the lifetime of the project.
