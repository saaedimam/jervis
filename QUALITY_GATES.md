# Quality Gates Specification

## 1. Overview
This document defines the mandatory Quality Gates that MUST be satisfied before any pull request or code contribution can be merged into the main branch of Project Jervis. All quality gates are automated and enforced through continuous integration pipelines.

---

## 2. Mandatory Quality Gates

### Gate 1: Architecture Invariant Gate
- **Requirement**: The codebase MUST NOT violate any of the 15 Architectural Invariants defined in `ARCHITECTURE_INVARIANTS.md`.
- **Validation**: Automated AST import analysis MUST verify zero upward or cyclic imports between layers.

### Gate 2: Code Formatting & Lint Gate
- **Requirement**: Code MUST pass standard `gofmt -s` and `goimports` formatting checks without modifications. `golangci-lint` MUST pass with zero warnings.

### Gate 3: Unit Testing Gate
- **Requirement**: 100% of unit tests in `internal/runtime`, `internal/memory`, and `internal/services` MUST pass without failures or race conditions (`go test -race`).

### Gate 4: Integration Testing Gate
- **Requirement**: Integration test suites in `tests/integration/` MUST pass using mock HTTP servers. Zero live network calls SHALL occur.

### Gate 5: Contract Testing Gate
- **Requirement**: MCP protocol and REST API JSON contracts MUST pass schema validation tests in `tests/contract/`.

### Gate 6: Benchmark & Performance Regression Gate
- **Requirement**: Benchmark tests in `tests/benchmark/` MUST verify that performance metrics remain within allowed budgets.

### Gate 7: Security Vulnerability Scan Gate
- **Requirement**: `govulncheck` MUST return zero known security vulnerabilities across all direct and indirect dependencies.

### Gate 8: Dependency & License Audit Gate
- **Requirement**: All dependencies MUST be verified against `03_DEPENDENCY_POLICY.md`. Copyleft licenses (GPL, AGPL) MUST NOT be present.

### Gate 9: Architecture Invariant Verification Gate
- **Requirement**: Package dependency graphs MUST be verified as directed acyclic graphs (DAGs). Upward package imports MUST trigger immediate failure.

### Gate 10: Code Coverage Threshold Gate
- **Requirement**: Overall project statement code coverage MUST meet or exceed **85%**. Layer 1 (Runtime) coverage MUST meet or exceed **90%**.

### Gate 11: Binary Size Budget Gate
- **Requirement**: The compiled single static CLI executable binary size MUST NOT exceed **15 Megabytes** (uncompressed).

### Gate 12: Startup Latency Budget Gate
- **Requirement**: Cold execution CLI startup latency MUST NOT exceed **15 Milliseconds**.

### Gate 13: Memory Footprint Budget Gate
- **Requirement**: Idle process memory usage for the background daemon/runtime MUST NOT exceed **25 Megabytes** RSS.

---

## 3. Enforcement Policy
- Quality gate checks SHALL be strictly blocking.
- Overriding a failed quality gate SHALL NOT be permitted under any circumstances.
- Pull requests with failing quality gates MUST be automatically blocked from merging by repository rule sets.
