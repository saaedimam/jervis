# Testing Strategy

## 1. Testing Pyramid & Overview
Project Jervis enforces a strict testing methodology to ensure that the core runtime, event bus, memory engine, and service layer remain 100% deterministic, robust, and secure.

```
       / \
      / E2E \           (CLI / MCP / REST Contracts)
     /-------\
    / Integration \     (Memory + Services + Provider Adapters)
   /---------------\
  /   Architecture  \   (AST Layer Boundary Linting)
 /-------------------\
/      Unit Tests      \ (Runtime / Memory / Services Core Logic)
-------------------------
```

---

## 2. Testing Levels

### Unit Testing
- **Scope**: Isolated testing of pure functions, data structures, event routing logic, working memory context sliding, timeline appending, and permission evaluation.
- **Location**: `*_test.go` files alongside code under test inside `internal/runtime`, `internal/memory`, `internal/services`.
- **Execution**: Fast, in-memory execution (<1 second total suite runtime). Zero network or disk I/O dependency.

### Integration Testing
- **Scope**: Verifying interaction between adjacent layers (e.g., Service Layer querying Memory Engine, Memory Engine flushing to Storage Driver, Service calling AI Provider Adapter with mock HTTP responses).
- **Location**: `tests/integration/`.
- **Mocking Policy**: Mock HTTP server (`net/http/httptest`) used for external LLM API endpoints (OpenAI, Claude, Gemini). No live API tokens consumed during automated CI test execution.

### Contract Testing
- **Scope**: Ensuring API interfaces, MCP protocol schemas, REST JSON schemas, and CLI flags conform to frozen public contracts.
- **Location**: `tests/contract/`.
- **Verification**: Schema validation against JSON specs and standard MCP protocol payloads.

### Architecture Testing (Layer Linter)
- **Scope**: Automated Go AST linting to enforce the mandatory 15 design rules and single-direction dependency flow (`OS -> Runtime -> Memory Engine -> Service Layer -> AI Provider Layer -> Interfaces`).
- **Tooling**: Custom AST linter or `golangci-lint` import rules (`depguard`).
- **Enforcement Rule**: Fails build if any package in `internal/runtime` or `internal/memory` imports `internal/aiprovider` or if cyclic dependencies are detected.

### Performance Testing
- **Scope**: Benchmarking CLI startup time, memory consumption, Event Bus event throughput, and Timeline append latency.
- **Location**: `tests/benchmark/`.
- **Key Metrics & Thresholds**:
  - CLI Cold Startup Latency: `< 15ms`
  - Memory Footprint (Idle Daemon): `< 25MB`
  - Event Bus Dispatch Latency: `< 1ms` per event batch
  - Memory Timeline Write Latency: `< 5ms` per record

### Regression Testing
- **Scope**: Verifying state recovery from local storage after unexpected daemon termination or process crashes.
- **Location**: `tests/regression/`.

---

## 3. Required Coverage Goals

| Module / Component | Minimum Code Coverage | Target Coverage |
| :--- | :---: | :---: |
| Layer 1: Runtime (`internal/runtime`) | **90%** | 95% |
| Layer 2: Memory Engine (`internal/memory`) | **85%** | 90% |
| Layer 3: Service Layer (`internal/services`) | **85%** | 90% |
| Layer 4: AI Provider Adapters (`internal/aiprovider`) | **80%** | 85% |
| Layer 5: Client Interfaces (`internal/interfaces`) | **75%** | 80% |
| **TOTAL PROJECT AVERAGE** | **85%** | **90%** |
