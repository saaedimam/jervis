# Coding Standards

## 1. Project Layout & Package Boundaries
Package layout strictly follows the canonical 5-tier directory layout defined in `PROJECT_STRUCTURE.md`:
- `cmd/`: Application binaries (`jervis`, `mcp`, `daemon`). No business logic.
- `internal/runtime/`: Layer 1: Runtime core (Observer, Event Bus, Scheduler, Lifecycle, Session, Permissions, Config).
- `internal/memory/`: Layer 2: Memory Engine (Working, Episodic, Semantic, Timeline, Retrieval, Compression, Store).
- `internal/services/`: Layer 3: Service Layer (Planner, Projects, Meetings, Habits, Notion, Calendar, Automation).
- `internal/aiprovider/`: Layer 4: AI Provider Layer (OpenAI, Claude, Gemini, Ollama, Local).
- `internal/interfaces/`: Layer 5: Interface Handlers (CLI, MCP, REST, Desktop, Menu Bar).
- `pkg/plugin/`: Public contracts for external plugins.

Cross-layer calls must flow strictly downward. Upward imports (e.g. `runtime` importing `services` or `aiprovider`) are strictly forbidden.

## 2. Naming Conventions
- **Packages**: Lowercase, single-word names (`eventbus`, `session`, `planner`, `aiprovider`).
- **Interfaces**: Concise noun or suffix `-er` (`Broker`, `Retriever`, `Provider`, `Subscriber`).
- **Structs**: PascalCase nouns describing concrete implementations (`MemoryStore`, `OpenAIAdapter`).
- **Functions & Methods**: MixedCaps / camelCase for unexported, PascalCase for exported (`ExecuteAction`, `validatePermission`).
- **Constants**: PascalCase for exported, camelCase for unexported (`DefaultTimeout`, `maxContextWindow`).

## 3. Formatting
- **Standard Tooling**: Code must pass `gofmt -s` and `goimports` cleanly.
- **Line Length**: Soft target 100 characters; max 120 characters.
- **Imports**: Grouped into three distinct blocks separated by blank lines:
  1. Standard Library packages
  2. Third-party external packages
  3. Internal `jervis` project packages

## 4. Logging
- **Standard Logger**: Use standard library `log/slog` for structured JSON logging.
- **No AI Calls in Logger**: The Observer and Event Bus loggers must NEVER invoke AI inference.
- **Log Levels**:
  - `DEBUG`: Verbose internal state transitions, event bus routing details.
  - `INFO`: Normal operational events (service execution, CLI invocation, lifecycle changes).
  - `WARN`: Recoverable errors, permission checks denied, provider fallbacks.
  - `ERROR`: Unrecoverable service failures, storage access errors, process crashes.

## 5. Error Handling
- **Typed Errors**: Define explicit sentinel errors and custom error types per package (`ErrPermissionDenied`, `ErrContextExhausted`).
- **Error Wrapping**: Wrap errors using `fmt.Errorf("context: %w", err)` to preserve stack traces.
- **No Swallowed Errors**: Errors must never be discarded with `_`. Every error must be logged, returned, or handled explicitly.

## 6. Testing Conventions
- **File Placement**: Unit test files placed in the same package as code under test (`eventbus_test.go`).
- **Table-Driven Tests**: Mandatory for testing multiple inputs/outputs across services and providers.
- **Mock Interfaces**: Use interface mocks; no reliance on live network connections during unit tests.

## 7. Comments & Documentation
- **Godoc Compliance**: Every exported package, type, constant, variable, and function must have a Godoc comment beginning with the symbol's name.
- **Why over What**: Comments must explain the architectural rationale, security considerations, or boundary constraints rather than restating code logic.

## 8. API Contracts
- **Interface Definitions**: All public package boundaries defined using Go `interface` types in `pkg/` or package roots.
- **JSON Contracts**: REST API and MCP protocol payloads defined using explicit JSON struct tags (`json:"action_id"`).

## 9. Configuration Management
- **Centralized Config**: Loaded strictly inside Layer 1 `internal/runtime/config`.
- **Precedence**: Command-line flags > Environment variables (`JERVIS_*`) > Configuration file (`~/.jervis/config.yaml`) > Hardcoded defaults.
- **Immutability**: Config structs are immutable once initialized during runtime boot.
