# Implementation Order

This document defines the exact, canonical sequence of implementation phases for Project Jervis. Development must proceed strictly bottom-up through the runtime hierarchy. No phase may begin until the exit criteria of all preceding phases are satisfied.

---

## Phase 1: Core Runtime Foundation Layer

### Goal
Implement the core deterministic Runtime module (`internal/runtime/`) owning system boot, configuration, session state, event distribution, process scheduling, and capability permission authorization.

### Deliverables
- Centralized configuration manager (`config`) loading flags, environment variables, and YAML files.
- Process lifecycle controller (`lifecycle`) managing boot, signal trapping, and graceful shutdown.
- Thread-safe Session State Manager (`session`) isolating session variables.
- Capability-based Permissions Enforcement engine (`permissions`) evaluating action requests.
- Asynchronous publish-subscribe Event Bus (`eventbus`) routing system events.
- Health monitoring and metric Observer (`observer`).
- Background cron and interval Scheduler (`scheduler`).

### Target Files
- `internal/runtime/config/config.go`
- `internal/runtime/lifecycle/lifecycle.go`
- `internal/runtime/session/session.go`
- `internal/runtime/permissions/permissions.go`
- `internal/runtime/eventbus/eventbus.go`
- `internal/runtime/observer/observer.go`
- `internal/runtime/scheduler/scheduler.go`

### Exit Criteria
- 100% of Runtime unit tests passing with `>90%` code coverage.
- AST Architecture Linter confirms zero imports of memory, services, AI providers, or UIs.
- Event Bus handles 10,000 events/second in-memory with zero lock contention.
- Permission engine correctly authorizes or denies capability requests.

### Dependencies
- Go standard library (`os`, `sync`, `context`, `log/slog`, `time`).

### Estimated Complexity
- **Medium** (Core concurrency, event routing, thread safety).

### Risk
- Misdesigned Event Bus interfaces could cause channel deadlocks or blocking backpressure under heavy load.

---

## Phase 2: Memory Engine Foundation Layer

### Goal
Implement the Memory Engine (`internal/memory/`) providing Working Memory, an immutable Timeline ledger, Episodic and Semantic storage, keyword/vector Retrieval, and heuristic Compression.

### Deliverables
- In-memory sliding window Working Memory store (`working`).
- Immutable append-only Timeline ledger (`timeline`) recording system events.
- Document and entity Knowledge Store driver (`store`) using Cgo-free SQLite (`modernc.org/sqlite`).
- Epidsodic (`episodic`) and Semantic (`semantic`) memory indexing.
- Context Retrieval pipeline (`retrieval`) with keyword and vector distance search.
- Rule-based context Compression engine (`compression`).

### Target Files
- `internal/memory/working/working.go`
- `internal/memory/timeline/timeline.go`
- `internal/memory/store/sqlite.go`
- `internal/memory/episodic/episodic.go`
- `internal/memory/semantic/semantic.go`
- `internal/memory/retrieval/retrieval.go`
- `internal/memory/compression/compression.go`

### Exit Criteria
- 100% of Memory unit tests passing with `>85%` code coverage.
- Timeline correctly records all Runtime events without loss or corruption.
- Context Compression truncates tokens deterministically without calling AI APIs.
- Memory engine operates 100% independently without any AI Provider configured.

### Dependencies
- Phase 1 (Runtime Layer).
- `modernc.org/sqlite` (Pure Go SQLite driver).

### Estimated Complexity
- **High** (Storage indexing, vector retrieval primitives, concurrency locking).

### Risk
- SQLite lock contention on concurrent writes if Timeline and Knowledge Store share database handles. (Mitigated by Write-Ahead Logging WAL mode).

---

## Phase 3: Domain Service Layer

### Goal
Implement domain business services (`internal/services/`) managing planning, local projects, meetings, habits, integrations (Notion, Calendar), and shell automation workflows.

### Deliverables
- Deterministic Planner service (`planner`) for task decomposition and execution mapping.
- Local repository and project tracking service (`projects`).
- Calendar and meeting preparation service (`meetings`, `calendar`).
- Recurring routine and habit tracking service (`habits`).
- Notion API integration client service (`notion`).
- Local shell script and automation workflow service (`automation`).

### Target Files
- `internal/services/planner/planner.go`
- `internal/services/projects/projects.go`
- `internal/services/meetings/meetings.go`
- `internal/services/habits/habits.go`
- `internal/services/notion/notion.go`
- `internal/services/calendar/calendar.go`
- `internal/services/automation/automation.go`

### Exit Criteria
- All domain services execute deterministically via Runtime Permission checks.
- Service unit and integration tests passing with `>85%` coverage.
- Automation service executes local tasks safely within sandbox limits.

### Dependencies
- Phase 1 (Runtime Layer), Phase 2 (Memory Engine).

### Estimated Complexity
- **Medium** (Business logic implementation, third-party API integration).

### Risk
- Notion / Google Calendar API rate limits or network failures. (Mitigated by retry logic and Memory Engine caching).

---

## Phase 4: AI Provider Abstraction Layer

### Goal
Implement the pluggable AI Provider Abstraction (`internal/aiprovider/`) supporting replaceable inference engines (OpenAI, Claude, Gemini, Ollama, local models).

### Deliverables
- Standard `Provider` and `Model` interface contracts (`provider.go`).
- OpenAI REST API adapter (`openai`).
- Anthropic Claude REST API adapter (`claude`).
- Google Gemini REST API adapter (`gemini`).
- Local Ollama HTTP API adapter (`ollama`).
- Local llama.cpp GGUF model binding wrapper (`local`).

### Target Files
- `internal/aiprovider/provider.go`
- `internal/aiprovider/openai/openai.go`
- `internal/aiprovider/claude/claude.go`
- `internal/aiprovider/gemini/gemini.go`
- `internal/aiprovider/ollama/ollama.go`
- `internal/aiprovider/local/local.go`

### Exit Criteria
- All AI Provider adapters conform strictly to `Provider` interface contract.
- Integration tests using mock HTTP servers (`net/http/httptest`) pass with `>80%` coverage.
- System switches between providers dynamically via config without runtime restart.
- Removing all providers leaves Runtime, Memory, and Services 100% functional.

### Dependencies
- Phase 1 (Runtime Layer), Phase 2 (Memory Engine), Phase 3 (Service Layer).

### Estimated Complexity
- **Medium** (REST JSON payload translation, streaming response handling).

### Risk
- Provider API schema drift or breaking vendor changes. (Mitigated by contract tests).

---

## Phase 5: Client Interfaces & Public Plugin SDK

### Goal
Implement external interfaces (`cmd/`, `internal/interfaces/`) exposing Jervis via CLI, MCP Server, REST API daemon, Desktop GUI, and Menu Bar widget, alongside public Plugin contracts (`pkg/plugin/`).

### Deliverables
- CLI entrypoint (`cmd/jervis/main.go`, `internal/interfaces/cli/`).
- MCP Server entrypoint (`cmd/mcp/main.go`, `internal/interfaces/mcp/`).
- REST Daemon entrypoint (`cmd/daemon/main.go`, `internal/interfaces/rest/`).
- Desktop and Menu Bar interface handlers (`internal/interfaces/desktop/`, `internal/interfaces/menubar/`).
- Public WASM/RPC Plugin interface contracts (`pkg/plugin/`).

### Target Files
- `cmd/jervis/main.go`
- `cmd/mcp/main.go`
- `cmd/daemon/main.go`
- `internal/interfaces/cli/cli.go`
- `internal/interfaces/mcp/mcp.go`
- `internal/interfaces/rest/rest.go`
- `internal/interfaces/desktop/desktop.go`
- `internal/interfaces/menubar/menubar.go`
- `pkg/plugin/plugin.go`

### Exit Criteria
- CLI binary compiles to single static executable `<15MB`.
- CLI cold startup latency `< 15ms`.
- MCP Server fully compliant with Model Context Protocol specification.
- WASM plugin sandbox loads and executes external test plugins safely.

### Dependencies
- Phases 1-4.
- `github.com/spf13/cobra`, `github.com/tetratelabs/wazero`.

### Estimated Complexity
- **High** (Multi-interface handling, WASM plugin sandboxing, CLI UX).

### Risk
- Inter-process communication latency between Desktop/Menu Bar UIs and REST daemon. (Mitigated by lightweight Unix domain sockets or HTTP/2 localhost connection pooling).
