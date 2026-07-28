# Continuous Integration & Continuous Delivery (CI/CD) Specification

## 1. Scope & Principles
This document defines the normative specification for the automated CI/CD automation pipeline for Project Jervis. All pull requests, main branch commits, and tag releases MUST trigger execution of this pipeline.

---

## 2. Pipeline Execution Stages

### Stage 1: Code Quality & Formatting Validation
- **Requirement**: The pipeline MUST execute static formatting validation using `gofmt -s` and `goimports`. Any formatting discrepancy SHALL fail the pipeline stage immediately.
- **Linter Enforcement**: The pipeline MUST execute `golangci-lint` configured with zero warning tolerance.

### Stage 2: License & Dependency Compliance Audit
- **Requirement**: The pipeline MUST inspect `go.mod` and vendor trees using automated license scanners.
- **License Policy**: Any dependency using GPL, AGPL, or non-permissive licenses MUST trigger immediate pipeline abortion.

### Stage 3: Security & Vulnerability Audit
- **Requirement**: The pipeline MUST execute `govulncheck` against the codebase. Any detected vulnerability rated High or Critical MUST fail the pipeline.

### Stage 4: Architecture Invariant Verification
- **Requirement**: The pipeline MUST run AST import graph inspection to verify single-direction layer flow (`OS -> Runtime -> Memory Engine -> Service Layer -> AI Provider Layer -> Interfaces`).
- **Enforcement**: Any import from `internal/runtime` or `internal/memory` pointing to `internal/aiprovider` MUST fail the pipeline.

### Stage 5: Automated Testing Execution
- **Unit Tests**: The pipeline MUST execute `go test -v -race` across all packages in `internal/`.
- **Integration Tests**: The pipeline MUST execute integration suites in `tests/integration/` using mock servers.
- **Contract Tests**: The pipeline MUST execute MCP protocol and REST API schema validation tests in `tests/contract/`.

### Stage 6: Code Coverage Reporting
- **Requirement**: The pipeline MUST calculate statement coverage. If total coverage falls below 85% or Runtime coverage falls below 90%, the stage MUST fail.

### Stage 7: Performance Benchmark & Budget Validation
- **Requirement**: The pipeline MUST execute performance benchmarks in `tests/benchmark/`.
- **Threshold Checks**: The pipeline MUST measure binary size (budget: 15MB) and startup latency (budget: 15ms). Violations MUST fail the stage.

### Stage 8: Cross-Compilation Matrix
- **Requirement**: The pipeline MUST execute static cross-compilation (`CGO_ENABLED=0`) across target platforms: `darwin/arm64`, `darwin/amd64`, `linux/amd64`, `linux/arm64`, `windows/amd64`.

### Stage 9: Release Artifact Generation & Signing
- **Trigger**: Tag push matching version pattern `v*`.
- **Requirement**: The pipeline MUST generate tarball archives, compute SHA256 checksums, sign manifests with GPG/Cosign, and create GitHub Releases.

### Stage 10: Homebrew Formula Automation
- **Trigger**: Successful release artifact generation.
- **Requirement**: The pipeline MUST update the Homebrew formula repository (`ioriimasu/homebrew-jervis`) with updated binary download URLs and SHA256 hashes.
