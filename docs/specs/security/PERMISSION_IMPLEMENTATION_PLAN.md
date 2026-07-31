# Runtime Permission Engine Implementation Plan

## 1. Overview
This plan defines the step-by-step implementation strategy for the Jervis Permission Engine (`internal/runtime/permissions`).

---

## 2. Implementation Sub-Phases

### Sub-Phase 1.3.1: Foundation Package
- **Package**: `internal/runtime/permissions/{contracts, types, errors}`
- **Deliverables**:
  - `contracts/contracts.go`: Interfaces (`Capability`, `Decision`, `Rule`, `Policy`, `Authorizer`, `PermissionEngine`).
  - `types/types.go`: Value types (`Subject`, `Resource`, `Action`, `Capability`, `Decision`, `Effect`).
  - `errors/errors.go`: Canonical error variables (`ErrPermissionDenied`, `ErrInvalidSubject`, `ErrInvalidResource`, `ErrInvalidAction`, `ErrDuplicatePolicy`).

### Sub-Phase 1.3.2: Rule Evaluator Engine
- **Package**: `internal/runtime/permissions/evaluator`
- **Deliverables**:
  - Pattern matching algorithm for subject/resource/action wildcards.
  - Evaluation of rules with Deny-First override logic.

### Sub-Phase 1.3.3: Authorizer & Policy Store
- **Package**: `internal/runtime/permissions/authorizer`
- **Deliverables**:
  - In-memory policy store (`RegisterPolicy`, `UnregisterPolicy`, `Policies`).
  - `Authorize(cap Capability)` implementation with Default Deny fallback.

### Sub-Phase 1.3.4: Permission Engine Facade & Test Suite
- **Package**: `internal/runtime/permissions`
- **Deliverables**:
  - Top-level `Engine` facade connecting policy store and evaluator.
  - Comprehensive unit test suite achieving **100.0% statement coverage** with `-race`.

---

## 3. Quality & Verification Strategy
- **Static Analysis**: `go fmt ./...`, `go vet ./...`.
- **Test Coverage**: 100.0% statement coverage.
- **Race Detector**: `go test -v -cover -race ./internal/runtime/permissions/...`.
