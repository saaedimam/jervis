# PROJECT_CONTEXT

## Jervis Project - Complete Engineering Context v2.0
### Generated: 2026-07-29T06:35:00Z
### Schema Version: 2.0.0

---

## TABLE OF CONTENTS

1. [System Overview](#system-overview)
2. [Repository Structure](#repository-structure)
3. [Architecture Layers](#architecture-layers)
4. [Modules & Packages](#modules--packages)
5. [Dependency Rules](#dependency-rules)
6. [Build Process](#build-process)
7. [Testing Strategy](#testing-strategy)
8. [Documentation Strategy](#documentation-strategy)
9. [Deployment Strategy](#deployment-strategy)
10. [Versioning Strategy](#versioning-strategy)
11. [Branch Strategy](#branch-strategy)
12. [Engineering Standards](#engineering-standards)
13. [Design Principles](#design-principles)
14. [Coding Standards](#coding-standards)
15. [Roadmap](#roadmap)

---

## SYSTEM OVERVIEW

### 1. Project Summary

**What this project is**: Jervis, a local-first personal OS and service platform designed for developer productivity and workflow automation.

**Primary goal**: Build a deterministic runtime-owned system that manages events, memory, domain services (planner, projects, habits, meetings), and security permissions, with optional AI provider integrations.

**Long-term vision**: A multi-interface platform (CLI, MCP Server, REST API, Desktop, Menu Bar) where the Jervis Runtime owns system execution and memory, and replaceable AI providers act strictly as context processing tools.

**Current development stage**: Phase 2 Memory Engine **IN PROGRESS** (66% complete - Phase 2.1 done, Phase 2.2 pending). Canonical tech stack: **Go (Golang 1.22+)**.

### 2. Architecture Baseline

**Canonical 5-Tier Hierarchy** (`OS -> Runtime -> Memory Engine -> Service Layer -> AI Provider Layer -> Interfaces`):
- Places system ownership in the deterministic Runtime rather than AI
- Ensures the runtime operates independently even if all AI providers are removed

**Runtime Ownership**:
- Ensures permissions, scheduling, event bus, observer, session management, and lifecycle are non-bypassable

**AI Provider Layer Decoupling**:
- Prevents vendor lock-in; treats OpenAI, Claude, Gemini, Ollama, and local models as replaceable utility engines

**Canonical Language Selection**: Go (Golang 1.22+)

**Runtime Foundation Concurrency Rule**: Phase 1 primitives use pure synchronous value semantics and deterministic state transitions without channels or background goroutines.

---

## REPOSITORY STRUCTURE

```
jervis/
├── 📁 .github/           # GitHub workflows and governance
│   ├── workflows/
│   ├── CODEOWNERS
│   └── GOVERNANCE.md
├── 📁 .jervis/           # Project metadata
│   ├── file_registry.txt
│   ├── file_hashes.txt
│   └── context/
├── 📁 cmd/               # CLI entry points
│   └── jervis/
├── 📁 context/           # Project documentation
│   ├── MASTER_CONTEXT.md
│   ├── PROJECT_CONTEXT.md (this file)
│   ├── CURRENT_CONTEXT.md
│   ├── API_FREEZE.md
│   └── MILESTONES.md
├── 📁 docs/              # Specifications and ADRs
│   ├── SPECIFICATIONS/
│   ├── ADR/
│   ├── EVENT_BUS_SPECIFICATION.md
│   ├── PERMISSION_ENGINE_SPECIFICATION.md
│   ├── OBSERVER_SPECIFICATION.md
│   └── ...
├── 📁 internal/          # Private implementation
│   ├── runtime/          # Core runtime
│   │   ├── contracts/    # ARCH-001
│   │   ├── types/        # ARCH-001
│   │   ├── eventbus/     # ARCH-002
│   │   ├── permissions/  # ARCH-003
│   │   ├── observer/     # ARCH-004
│   │   ├── storage/      # ARCH-005 (Phase 2)
│   │   ├── scheduler/    # ARCH-006
│   │   ├── planner/      # ARCH-007 (Phase 3)
│   │   ├── ai/           # ARCH-008 (Phase 4)
│   │   └── ...
│   ├── memory/           # Layer 2: Memory Engine
│   │   ├── working/      # ARCH-005
│   │   ├── timeline/     # ARCH-006
│   │   └── store/        # ARCH-007 (Phase 2.2)
│   ├── services/         # Layer 3: Service Layer
│   │   ├── planner/
│   │   ├── projects/
│   │   ├── habits/
│   │   └── meetings/
│   ├── aiprovider/       # Layer 4: AI Provider
│   │   ├── openai/
│   │   ├── claude/
│   │   ├── gemini/
│   │   ├── ollama/
│   │   └── local/
│   ├── interfaces/       # Layer 5: Interfaces
│   │   ├── cli/
│   │   ├── mcp/
│   │   ├── rest/
│   │   └── desktop/
│   └── desktop/          # Desktop client
├── 📁 pkg/               # Public packages
│   └── plugin/           # Plugin contracts
├── 📁 scripts/           # Automation scripts
│   ├── notion_*.sh
│   ├── jervis_*.sh
│   └── populate_*.sh
├── 📁 migrations/        # Database migrations
├── 📄 go.mod             # Go module definition
├── 📄 Makefile           # Build automation
└── 📄 README.md          # Project overview
```

**File Counts**:
- Total Files: 242
- Go Files: ~120
- Test Files: ~120
- Documentation: ~25
- Scripts: ~15

---

## ARCHITECTURE LAYERS

### Layer 1: OS Foundation

**Components**:
- OS Integration
- File System Abstraction
- Process Management
- Configuration Management

**Status**: ⚪ Foundation (infrastructure only)

### Layer 2: Runtime

**Components**:

#### ARCH-001: Runtime Core (6 packages)
- **contracts** (PKG-001): Interface definitions
- **types** (PKG-002): Domain types
- **errors** (PKG-003): Error handling
- **version** (PKG-004): Build versioning
- **buildinfo** (PKG-005): Version management
- **config** (PKG-006): Configuration

**Status**: 🟡 In Progress

#### ARCH-002: Event Bus Engine (8 packages)
- **contracts** (PKG-007): Event contracts
- **eventcontracts** (PKG-008): Event types
- **errors** (PKG-009): Event errors
- **registry** (PKG-010): Handler registry
- **dispatcher** (PKG-011): Event dispatch
- **middleware** (PKG-012): Middleware chain
- **eventbus** (PKG-013): Public API
- **internal** (PKG-014): Implementation

**Features**:
- In-process synchronous routing
- Priority-based dispatch (`PriorityLow`, `PriorityNormal`, `PriorityHigh`, `PriorityCritical`)
- Panic isolation per handler
- Continue-on-Error aggregation (`AggregateError`)
- MaxDispatchDepth = 16
- FIFO Middleware Chain
- Zero AI awareness
- Zero persistence

**Status**: ✅ COMPLETE (100% coverage, Frozen)

#### ARCH-003: Permission Engine (8 packages)
- **contracts** (PKG-015): Permission contracts
- **errors** (PKG-016): Permission errors
- **policy** (PKG-017): Policy engine
- **rules** (PKG-018): Rule definitions
- **resolution** (PKG-019): Rule resolution
- **permissions** (PKG-020): Public API
- **internal** (PKG-021): Implementation
- **types** (PKG-022): Domain types

**Features**:
- Capability-based access control (CBAC)
- Default Deny fallback
- Deny-First override rule precedence
- Wildcard evaluation (`*` and `prefix*`)
- 6-stage evaluation engine
- Zero AI awareness
- Zero persistence

**Status**: ✅ COMPLETE (100% coverage, Frozen)

#### ARCH-004: Observer (7 packages)
- **contracts** (PKG-023): Observer contracts
- **subscribe** (PKG-024): Subscription management
- **notify** (PKG-025): Notification
- **wait** (PKG-026): Synchronization
- **observer** (PKG-027): Public API
- **internal** (PKG-028): Implementation
- **types** (PKG-029): Domain types

**Features**:
- Compositional pattern
- Subscribe, Notify, Wait operations
- Thread-safe
- Defensive copies
- Zero AI awareness

**Status**: 🟡 Pending Implementation (Architecture Refrozen)

#### ARCH-005: Scheduler
- Background cron and interval task engine
- Supports Interval, Once (Deferred), and basic Cron schedules
- Deterministic `Tick(time.Time)` driving logic with panic isolation per job
- Thread-safe Registry with FIFO execution order
- Zero side effects on OS state

**Status**: ✅ COMPLETE (100% coverage)

#### ARCH-006: Session Management Engine
- Isolated state and context manager
- Thread-safe metadata storage per session
- Session lifecycle (`Created`, `Running`, `Stopped`)
- Defensive copies for state isolation

**Status**: ✅ COMPLETE (100% coverage)

**Overall Layer 2 Status**: 🟡 In Progress (86% complete)

### Layer 3: Memory Engine

**Components**:

#### ARCH-007: Working Memory
- In-memory sliding window FIFO store
- Thread-safe implementation
- Defensive copies
- Basic query support

**Status**: ✅ COMPLETE (100% coverage)

#### ARCH-008: Timeline Ledger
- Immutable append-only ledger
- Thread-safe
- Defensive copies
- Event sourcing pattern

**Status**: ✅ COMPLETE (100% coverage)

#### ARCH-009: Knowledge Store
- SQLite persistence layer
- Entity storage
- Query interface
- `modernc.org/sqlite` driver

**Status**: 🟡 IN PROGRESS (Phase 2.2)

**Overall Layer 3 Status**: 🟡 In Progress (66% complete)

### Layer 4: Service Layer (Phase 3)

**Components** (Future):
- **Planner Service**: Task management
- **Projects Service**: Project tracking
- **Habits Service**: Habit formation
- **Meetings Service**: Meeting management
- **Notion Service**: Notion integration
- **Calendar Service**: Calendar integration
- **Automation Service**: Workflow automation

**Status**: ⚪ Not Started (Q3 2026)

### Layer 5: AI Provider Layer (Phase 4)

**Components** (Future):
- **OpenAI Provider**: OpenAI API integration
- **Claude Provider**: Anthropic API integration
- **Gemini Provider**: Google API integration
- **Ollama Provider**: Local Ollama integration
- **Local Provider**: Direct model execution
- **Provider Abstractions**: Common interface

**Status**: ⚪ Not Started (Q4 2026)

### Layer 6: Interfaces (Phase 5)

**Components** (Future):
- **CLI**: Command-line interface
- **MCP Server**: Model Context Protocol
- **REST API**: HTTP API
- **Desktop**: Tauri + React application
- **Menu Bar**: System tray application

**Status**: ⚪ Not Started (Q4 2026)

---

## MODULES & PACKAGES

### Complete Package Inventory

#### ARCH-001: Runtime Core

| Package | ID | Path | Purpose | Status | Coverage |
|---------|-----|------|---------|--------|----------|
| contracts | PKG-001 | internal/runtime/contracts | Interface definitions | 🟡 Active | N/A |
| types | PKG-002 | internal/runtime/types | Domain types | 🟡 Active | N/A |
| errors | PKG-003 | internal/runtime/errors | Error handling | 🟡 Active | N/A |
| version | PKG-004 | internal/runtime/version | Build versioning | 🟡 Active | N/A |
| buildinfo | PKG-005 | internal/runtime/buildinfo | Version management | 🟡 Active | N/A |
| config | PKG-006 | internal/runtime/config | Configuration | 🟡 Active | N/A |

#### ARCH-002: Event Bus

| Package | ID | Path | Purpose | Status | Coverage |
|---------|-----|------|---------|--------|----------|
| contracts | PKG-007 | internal/runtime/eventbus/contracts | Event contracts | 🔒 Frozen | 100% |
| eventcontracts | PKG-008 | internal/runtime/eventbus/eventcontracts | Event types | 🔒 Frozen | 100% |
| errors | PKG-009 | internal/runtime/eventbus/errors | Event errors | 🔒 Frozen | 100% |
| registry | PKG-010 | internal/runtime/eventbus/registry | Handler registry | 🔒 Frozen | 100% |
| dispatcher | PKG-011 | internal/runtime/eventbus/dispatcher | Event dispatch | 🔒 Frozen | 100% |
| middleware | PKG-012 | internal/runtime/eventbus/middleware | Middleware chain | 🔒 Frozen | 100% |
| eventbus | PKG-013 | internal/runtime/eventbus | Public API | 🔒 Frozen | 100% |
| internal | PKG-014 | internal/runtime/eventbus/internal | Implementation | 🔒 Frozen | 100% |

**Exported APIs**: 8 (API-029..031 + 5 more)

#### ARCH-003: Permission Engine

| Package | ID | Path | Purpose | Status | Coverage |
|---------|-----|------|---------|--------|----------|
| contracts | PKG-015 | internal/runtime/permissions/contracts | Permission contracts | 🔒 Frozen | 100% |
| errors | PKG-016 | internal/runtime/permissions/errors | Permission errors | 🔒 Frozen | 100% |
| policy | PKG-017 | internal/runtime/permissions/policy | Policy engine | 🔒 Frozen | 100% |
| rules | PKG-018 | internal/runtime/permissions/rules | Rule definitions | 🔒 Frozen | 100% |
| resolution | PKG-019 | internal/runtime/permissions/resolution | Rule resolution | 🔒 Frozen | 100% |
| permissions | PKG-020 | internal/runtime/permissions | Public API | 🔒 Frozen | 100% |
| internal | PKG-021 | internal/runtime/permissions/internal | Implementation | 🔒 Frozen | 100% |
| types | PKG-022 | internal/runtime/permissions/types | Domain types | 🔒 Frozen | 100% |

**Exported APIs**: 12

#### ARCH-004: Observer

| Package | ID | Path | Purpose | Status | Coverage |
|---------|-----|------|---------|--------|----------|
| contracts | PKG-023 | internal/runtime/observer/contracts | Observer contracts | 🔒 Frozen | 100% |
| subscribe | PKG-024 | internal/runtime/observer/subscribe | Subscription mgmt | 🔒 Frozen | 100% |
| notify | PKG-025 | internal/runtime/observer/notify | Notification | 🔒 Frozen | 100% |
| wait | PKG-026 | internal/runtime/observer/wait | Synchronization | 🔒 Frozen | 100% |
| observer | PKG-027 | internal/runtime/observer | Public API | 🔒 Frozen | 100% |
| internal | PKG-028 | internal/runtime/observer/internal | Implementation | 🔒 Frozen | 100% |
| types | PKG-029 | internal/runtime/observer/types | Domain types | 🔒 Frozen | 100% |

**Exported APIs**: 11

### Summary Statistics

| Architecture | Packages | APIs | Status | Coverage |
|--------------|----------|------|--------|----------|
| ARCH-001 | 6 | TBD | 🟡 Active | - |
| ARCH-002 | 8 | 8 | 🔒 Frozen | 100% |
| ARCH-003 | 8 | 12 | 🔒 Frozen | 100% |
| ARCH-004 | 7 | 11 | 🔒 Frozen | 100% |
| **Total** | **29** | **31+** | **23 Frozen** | **100%** |

---

## DEPENDENCY RULES

### Import Hierarchy

```
Layer 6 (Interfaces) → Layer 5 (AI) → Layer 4 (Services)
  → Layer 3 (Memory) → Layer 2 (Runtime) → Layer 1 (OS)
```

**Rule**: Lower layers CANNOT import higher layers.

### Allowed Imports

```go
// Within Runtime
internal/runtime/eventbus → internal/runtime/contracts
internal/runtime/eventbus → internal/runtime/types
internal/runtime/eventbus → internal/runtime/errors

// Memory can import Runtime
internal/memory/working → internal/runtime/contracts
internal/memory/working → internal/runtime/types

// Services can import Runtime and Memory
internal/services/planner → internal/runtime/contracts
internal/services/planner → internal/memory/working
```

### Forbidden Imports

```go
// ❌ Cross-layer violation
internal/runtime/eventbus → internal/services/planner

// ❌ Circular dependency
internal/services/planner → internal/runtime/eventbus

// ❌ Interface → Implementation
internal/interfaces/cli → internal/services/planner/internal
```

### Dependency Direction Rules

1. **Contracts First**: All imports go through contracts
2. **No Implementation Leaks**: Never import `internal/` from other layers
3. **No Circles**: Dependency graph must be DAG
4. **Explicit Only**: All dependencies declared in go.mod

---

## BUILD PROCESS

### Makefile Targets

```makefile
.PHONY: test lint build clean ci-test ci-lint ci-build

# Development
test:                    ## Run all tests with coverage
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:                    ## Run golangci-lint
	golangci-lint run ./...

build:                   ## Build jervis binary
	go build -ldflags "-X main.Version=$(VERSION)" -o bin/jervis ./cmd/jervis

clean:                   ## Remove build artifacts
	rm -rf bin/ coverage.out coverage.html

# CI/CD
ci-test:                 ## CI test with race detection
	go test -race -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | grep total | awk '{print $$3}'

ci-lint:                 ## CI lint with strict config
	golangci-lint run --config=.golangci.yml ./...

ci-build:                ## CI build with version injection
	go build -ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.Date=$(DATE)" \
	  -o bin/jervis ./cmd/jervis

# Documentation
docs-serve:              ## Serve documentation locally
	mkdocs serve

docs-build:              ## Build documentation site
	mkdocs build

# Release
release-patch:           ## Bump patch version
	./scripts/bump_version.sh patch

release-minor:           ## Bump minor version
	./scripts/bump_version.sh minor

release-major:           ## Bump major version
	./scripts/bump_version.sh major
```

### Build Pipeline

```
1. Checkout Code
   ↓
2. Install Dependencies (go mod download)
   ↓
3. Lint (golangci-lint)
   ↓
4. Test (go test -race -cover)
   ↓
5. Build (go build)
   ↓
6. Package (Create artifacts)
   ↓
7. Release (GitHub release)
```

---

## TESTING STRATEGY

### Coverage Requirements

- **Minimum**: 100% for all packages
- **Critical Paths**: 100% + benchmarks
- **Race Detection**: Enabled for all tests
- **Fuzzing**: Required for parsers

### Test Structure

```go
// File: eventbus_test.go
// Package: eventbus_test (external package)

package eventbus_test

import (
    "testing"
    
    "github.com/jervis/internal/runtime/eventbus"
)

// Table-driven test
func TestEventBus_Publish(t *testing.T) {
    tests := []struct {
        name    string
        event   Event
        wantErr bool
    }{
        {"valid event", validEvent, false},
        {"nil event", nil, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            bus := eventbus.New()
            
            // Act
            err := bus.Publish(tt.event)
            
            // Assert
            if (err != nil) != tt.wantErr {
                t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

// Benchmark
func BenchmarkEventBus_Publish(b *testing.B) {
    bus := eventbus.New()
    event := createBenchmarkEvent()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        bus.Publish(event)
    }
}

// Fuzz test
func FuzzEventBus_Publish(f *testing.F) {
    f.Add([]byte("test event"))
    
    f.Fuzz(func(t *testing.T, data []byte) {
        bus := eventbus.New()
        event := parseEvent(data)
        bus.Publish(event)
    })
}
```

### Testing Rules

1. **Package Suffix**: Use `_test` suffix for test packages
2. **Table-Driven**: Prefer table-driven tests
3. **Subtests**: Use `t.Run()` for scenarios
4. **Parallel**: Mark safe tests with `t.Parallel()`
5. **Defensive Copies**: Always test defensive copies
6. **Panic Recovery**: Test panic isolation
7. **Error Cases**: Cover all error paths
8. **Edge Cases**: Test boundaries and nil values

### Coverage Commands

```bash
# Run all tests
go test -race ./...

# Generate coverage report
go test -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# View coverage summary
go tool cover -func=coverage.out

# Race detector
go test -race ./...

# Fuzz testing
go test -fuzz=FuzzPublish -fuzztime=30s ./...

# Benchmark
go test -bench=. -benchmem ./...
```

---

## DOCUMENTATION STRATEGY

### Documentation Hierarchy

```
Level 1: README.md                    (Quick start)
  ↓
Level 2: PROJECT_CONTEXT.md           (This file - Architecture)
  ↓
Level 3: MASTER_CONTEXT.md            (Current state)
  ↓
Level 4: docs/SPECIFICATIONS/*.md     (Technical details)
  ↓
Level 5: docs/ADR/*.md                (Decision records)
  ↓
Level 6: Code comments                (Implementation)
  ↓
Level 7: Notion Knowledge Graph       (Canonical source)
```

### Documentation Requirements

| Document | Purpose | Update Frequency | Owner |
|----------|---------|------------------|-------|
| README.md | Quick start, features | Per release | Engineering |
| PROJECT_CONTEXT.md | Architecture, structure | Per phase | Architecture |
| MASTER_CONTEXT.md | Current state, metrics | Per session | AI/Engineering |
| Specifications | Technical design | Per component | Component Lead |
| ADRs | Architecture decisions | Per decision | Decision Author |
| Notion | Knowledge graph | Continuous | Sync Agent |

### Documentation Standards

1. **Markdown**: All docs in Markdown
2. **Diagrams**: Use Mermaid for diagrams
3. **APIs**: Document in specifications
4. **Changes**: Track in ADRs
5. **Sync**: Automated sync to Notion

### Key Documents

**Vision & Architecture**:
- VISION.md
- ARCHITECTURE.md
- PRINCIPLES.md
- ROADMAP.md
- PROJECT_STRUCTURE.md
- DECISIONS.md

**Governance**:
- CONSTITUTION.md
- ARCHITECTURE_INVARIANTS.md
- ADR_GUIDE.md
- QUALITY_GATES.md
- CI_CD_SPECIFICATION.md
- PROJECT_METRICS.md

**Event Bus**:
- EVENT_BUS_SPECIFICATION.md
- EVENT_MODEL.md
- EVENT_CONTRACTS.md
- EVENT_IMPLEMENTATION_PLAN.md
- DISPATCHER_SPECIFICATION.md
- MIDDLEWARE_SPECIFICATION.md
- BUS_SPECIFICATION.md
- BUS_CONTRACTS.md
- BUS_PIPELINE.md
- BUS_SEQUENCE.md

**Permission Engine**:
- PERMISSION_ENGINE_SPECIFICATION.md
- PERMISSION_MODEL.md
- PERMISSION_CONTRACTS.md
- PERMISSION_IMPLEMENTATION_PLAN.md

**Observer**:
- OBSERVER_SPECIFICATION.md
- OBSERVER_MODEL.md
- OBSERVER_CONTRACTS.md
- OBSERVER_IMPLEMENTATION_PLAN.md

---

## DEPLOYMENT STRATEGY

### Environments

```
Development → Staging → Production
```

### Deployment Steps

1. **Create Release Branch**
   ```bash
   git checkout -b release/v1.2.3
   ```

2. **Update Version**
   ```bash
   ./scripts/bump_version.sh patch
   ```

3. **Run Tests**
   ```bash
   make ci-test
   make ci-lint
   ```

4. **Build Artifacts**
   ```bash
   make ci-build
   ```

5. **Create Release**
   ```bash
   git tag v1.2.3
   git push origin v1.2.3
   ```

6. **GitHub Release**
   - Automated via GitHub Actions
   - Artifacts uploaded
   - Release notes generated

7. **Update Notion**
   - Sync knowledge graph
   - Update milestone status
   - Record session

### Rollback Strategy

- **Code**: Revert to previous tag
- **Database**: Migrations are backward compatible
- **Config**: Feature flags enable gradual rollout
- **Notion**: Knowledge graph versioned

---

## VERSIONING STRATEGY

### Semantic Versioning

```
vMAJOR.MINOR.PATCH

MAJOR: Breaking changes (requires ADR)
MINOR: New features (backward compatible)
PATCH: Bug fixes (backward compatible)
```

### Version Sources

| Source | Format | Update Mechanism |
|--------|--------|------------------|
| Git Tag | vX.Y.Z | Manual (release) |
| go.mod | module version | `go mod edit` |
| buildinfo.Version | vX.Y.Z | ldflags at build |
| Notion | vX.Y.Z | Sync agent |

### Breaking Changes

**Require**:
1. ADR documenting the change
2. Migration guide
3. Deprecation period (if applicable)
4. Version bump to MAJOR

### Version Commands

```bash
# View current version
git describe --tags --abbrev=0

# Create new version
./scripts/bump_version.sh [patch|minor|major]

# Tag release
git tag v1.2.3
git push origin v1.2.3
```

---

## BRANCH STRATEGY

### Branch Types

```
main                    - Production code (protected)
  ↓
release/vX.Y            - Release preparation
  ↓
feature/description     - Feature development
  ↓
fix/description         - Bug fixes
  ↓
docs/description        - Documentation updates
  ↓
refactor/description    - Code refactoring
```

### Workflow

```bash
# Start new feature
git checkout -b feature/working-memory
git push -u origin feature/working-memory

# Make commits
git commit -m "feat: implement WorkingMemory"

# Create Pull Request
gh pr create --title "feat: Working Memory" --body "..."

# Merge after review
git checkout main
git pull origin main
```

### Branch Protection Rules

- [x] Require PR for main
- [x] Require 1 approving review
- [x] Require CI pass
- [x] Require up-to-date branch
- [x] Require linear history
- [x] Include administrators

### Naming Conventions

| Type | Pattern | Example |
|------|---------|---------|
| Feature | `feature/description` | `feature/timeline-ledger` |
| Fix | `fix/issue-description` | `fix/race-condition` |
| Docs | `docs/description` | `docs/api-reference` |
| Refactor | `refactor/description` | `refactor/event-bus` |
| Release | `release/vX.Y` | `release/v1.2` |

---

## ENGINEERING STANDARDS

### Code Quality

| Aspect | Requirement | Tool |
|--------|-------------|------|
| Test Coverage | 100% | go test -cover |
| Linting | Zero errors | golangci-lint |
| Race Detection | Zero races | go test -race |
| Formatting | gofmt compliant | gofmt |
| Imports | goimports compliant | goimports |

### Performance

| Aspect | Requirement | Tool |
|--------|-------------|------|
| Benchmarks | Critical paths | go test -bench |
| Memory Profile | When needed | pprof |
| CPU Profile | When needed | pprof |
| Optimization | No premature | - |

### Security

| Aspect | Requirement |
|--------|-------------|
| Secrets | No hardcoded |
| Input | Validate all |
| SQL | Injection prevention |
| XSS | Prevention |
| Auth | CBAC enforcement |

### Observability

| Aspect | Requirement |
|--------|-------------|
| Logging | Structured |
| Metrics | Key paths |
| Tracing | Future consideration |
| Health | Endpoints |

---

## DESIGN PRINCIPLES

### 1. Determinism

**Rule**: All execution must be deterministic.

**Applies To**: Event dispatch, scheduling, state transitions

**Enforcement**: 
- Synchronous execution
- No randomness
- No race conditions
- Ordered operations

### 2. Synchronicity

**Rule**: 100% synchronous execution for Phase 1.

**Applies To**: Runtime core, event bus, permissions

**Enforcement**:
- No goroutines in core
- No channels
- Blocking operations only
- Deterministic ordering

### 3. Immutability

**Rule**: Value objects are immutable.

**Applies To**: Events, decisions, configurations

**Enforcement**:
- Defensive copies on boundaries
- No mutable state in values
- Return new instances

### 4. Panic Isolation

**Rule**: Panics don't crash the system.

**Applies To**: Event handlers, schedulers, middleware

**Enforcement**:
- defer/recover in handlers
- Error aggregation
- Continue on error

### 5. Testability

**Rule**: 100% test coverage.

**Applies To**: All code

**Enforcement**:
- No untested code paths
- Defensive copies tested
- Error cases covered
- Edge cases tested

### 6. Composability

**Rule**: Components compose cleanly.

**Applies To**: All packages

**Enforcement**:
- No circular dependencies
- Clear interfaces
- Single responsibility

### 7. Replaceability

**Rule**: AI providers are replaceable.

**Applies To**: AI layer

**Enforcement**:
- Interface abstractions
- No hard dependencies
- Swappable implementations

---

## CODING STANDARDS

### Go Standards

```go
// Package documentation
// Package eventbus provides synchronous event dispatch with priority-based
// routing and panic isolation.
package eventbus

// Exported type with complete godoc
type EventBus struct {
    // registry holds all registered handlers
    registry *Registry
    
    // middleware holds the middleware chain
    middleware []Middleware
}

// New creates a new EventBus with default configuration.
func New(opts ...Option) *EventBus {
    // Implementation
}

// Publish sends an event to all registered handlers.
// Returns AggregateError if any handler fails.
func (eb *EventBus) Publish(event Event) error {
    // Implementation
}
```

### Naming Conventions

| Element | Convention | Example |
|---------|-----------|---------|
| Types | PascalCase | `EventBus`, `HandlerFunc` |
| Variables | camelCase | `eventBus`, `handler` |
| Constants | UPPER_SNAKE_CASE | `MaxDispatchDepth` |
| Interfaces | -er suffix | `Reader`, `Writer`, `Handler` |
| Tests | Test + FunctionName | `TestPublish`, `TestHandler` |
| Files | snake_case.go | `event_bus.go` |
| Packages | lowercase | `eventbus`, `permissions` |

### Formatting

```bash
# Format code
gofmt -w .

# Format and organize imports
goimports -w .

# Check formatting
gofmt -d .
```

### Comments

```go
// Godoc for exported items

// Package eventbus provides ...
package eventbus

// EventBus is the main ...
type EventBus struct{}

// Publish sends an event ...
func (eb *EventBus) Publish(event Event) error {}

// Inline for complex logic
func complexFunction() {
    // Validate inputs before processing to ensure
    // we don't waste resources on invalid data
    if err := validate(); err != nil {
        return err
    }
}

// TODO with reference
// TODO(SESSION-021): Refactor this after Phase 2.2
func needsRefactoring() {}
```

---

## ROADMAP

### Current Phase: Phase 2 Memory Engine

**Status**: 🟡 In Progress (66%)

#### Phase 2.1: Working Memory & Timeline ✅ COMPLETE
- [x] Working Memory FIFO store
- [x] Timeline Ledger immutable events
- [x] Thread-safe implementations
- [x] Defensive copies
- [x] 100% test coverage
- **Completion**: 2026-07-29

#### Phase 2.2: Knowledge Store 🟡 IN PROGRESS
- [ ] Design SQLite schema
- [ ] Implement `modernc.org/sqlite` driver
- [ ] Create persistence abstractions
- [ ] Implement query interface
- [ ] 100% test coverage
- **Target**: August 2026

#### Phase 2.3: Query Engine ⏳ PENDING
- [ ] Design query DSL
- [ ] Implement query parser
- [ ] Optimize query execution
- [ ] Add indexing
- **Target**: September 2026

### Phase 3: Domain Services ⏳ PENDING (Q3 2026)

**Services**:
- [ ] Planner Service
- [ ] Projects Service
- [ ] Habits Service
- [ ] Meetings Service
- [ ] Notion Integration
- [ ] Calendar Integration

### Phase 4: AI Provider Layer ⏳ PENDING (Q4 2026)

**Providers**:
- [ ] OpenAI integration
- [ ] Claude integration
- [ ] Gemini integration
- [ ] Ollama integration
- [ ] Local model support
- [ ] Provider abstraction

### Phase 5: Client Interfaces ⏳ PENDING (Q4 2026)

**Interfaces**:
- [ ] CLI improvements
- [ ] MCP Server
- [ ] REST API
- [ ] Desktop (Tauri + React)
- [ ] Menu Bar app

### v1.0.0 Release ⏳ PENDING (Dec 2026)

**Criteria**:
- [ ] All phases complete
- [ ] 100% test coverage
- [ ] Documentation complete
- [ ] Stable APIs
- [ ] Production ready

---

## QUICK REFERENCE

### Canonical IDs

| Entity | Pattern | Example | Count |
|--------|---------|---------|-------|
| Architecture | ARCH-### | ARCH-002 | 4 |
| Package | PKG-### | PKG-014 | 29 |
| File | FILE-#### | FILE-0008 | 242 |
| API | API-### | API-008 | 31+ |
| Spec | SPEC-### | SPEC-001 | 15 |
| ADR | ADR-#### | ADR-0002 | 4 |
| Session | SESSION-### | SESSION-021 | 21 |
| Milestone | PHASE-X.Y | PHASE-2.1 | 19 |

### Notion Databases

| Database | ID | Entries | Status |
|----------|-----|---------|--------|
| Architecture | d3dcb133-f96e-4e8e-944f-5825c2d1eee0 | 4 | 🟡 Schema |
| Package | 9c8bb7d5-5675-4cc5-b1b7-6a9c1ac3fe2f | 29 | 🟡 Schema |
| Specification | f30e0d51-a787-421a-ad6b-77935f7d2e53 | 15 | 🟡 Schema |
| File | d5b8d71a-c568-4288-9443-f3deb8b316bc | 23 | 🟡 Schema |
| API | 5e2dad61-5186-46f7-be6b-e7e5c3715f04 | 31 | 🟡 Schema |
| ADR | abc5d892-1299-4813-b8bf-a143d6c8c73c | 4 | 🟡 Schema |
| Milestone | 39ae6e23-2bc1-4e34-a7b0-a1da9410b081 | 19 | 🟡 Schema |
| Session | c1e36ebb-a3fc-4aea-a3d2-ac8214e1e40a | 21 | 🟡 Schema |
| Commit | 69c5145a-b84c-43e5-83b2-05d746a80e26 | 10+ | 🟡 Schema |
| Memory | 38a76b5b-b20e-498e-b6e9-e643c2ae7d8b | TBD | 🟡 Schema |
| Dependencies | 1de04b92-6fe3-4756-b85d-c9370f838a3b | TBD | 🟡 Schema |

### Key Commands

```bash
# Run tests
make test

# Run lint
make lint

# Build
make build

# CI pipeline
make ci-test ci-lint ci-build

# Sync to Notion
./scripts/jervis_compiler.sh

# View coverage
go tool cover -html=coverage.out
```

### Contact

- **Project**: Jervis
- **Repository**: github.com/jervis
- **Notion**: Jervis Knowledge Graph
- **Status**: Phase 2.1 Complete, Phase 2.2 In Progress

---

**PROJECT_CONTEXT v2.0.0 Complete**
**Last Updated**: 2026-07-29T06:35:00Z
