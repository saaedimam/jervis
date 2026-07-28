# Project Context

## 1. Project Summary
- **What this project is**: Jervis, a local-first personal OS and service platform designed for developer productivity and workflow automation.
- **Primary goal**: Build a deterministic runtime-owned system that manages events, memory, domain services (planner, projects, habits, meetings), and security permissions, with optional AI provider integrations.
- **Long-term vision**: A multi-interface platform (CLI, MCP Server, REST API, Desktop, Menu Bar) where the Jervis Runtime owns system execution and memory, and replaceable AI providers act strictly as context processing tools.
- **Current development stage**: Phase 1.1 Foundation **COMPLETED** (100% test coverage). Phase 1.2.1 Event Bus Foundation **COMPLETED** (100% test coverage). Phase 1.2.2 Event Bus Subscription Registry **COMPLETED** (100% test coverage). Phase 1.2.4 Synchronous Dispatcher **COMPLETED** (100% test coverage). Phase 1.2.5 Middleware Architecture Specifications **FROZEN**. Canonical tech stack: **Go (Golang 1.22+)**.

## 2. Architecture Baseline
- **Canonical 5-Tier Hierarchy (`OS -> Runtime -> Memory Engine -> Service Layer -> AI Provider Layer -> Interfaces`)**:
  - Places system ownership in the deterministic Runtime rather than AI. Ensures the runtime operates independently even if all AI providers are removed.
- **Runtime Ownership**:
  - Ensures permissions, scheduling, event bus, observer, session management, and lifecycle are non-bypassable.
- **AI Provider Layer Decoupling**:
  - Prevents vendor lock-in; treats OpenAI, Claude, Gemini, Ollama, and local models as replaceable utility engines.
- **Canonical Language Selection**: Go (Golang 1.22+).
- **Runtime Foundation Concurrency Rule**: Phase 1 primitives use pure synchronous value semantics and deterministic state transitions without channels or background goroutines.
- **Event Bus Architecture**: In-process synchronous routing pipeline with `Priority uint8` (`PriorityLow`, `PriorityNormal`, `PriorityHigh`, `PriorityCritical`), pure Go pattern matching, deterministic priority-based dispatch ordering, panic isolation per handler, Continue-on-Error error aggregation (`AggregateError`), `MaxDispatchDepth = 16`, FIFO Middleware Chain, zero AI awareness, and zero persistence.

## 3. Repository Structure
- `cmd/`: Entrypoints for Interfaces (`jervis` CLI, `mcp` server, `daemon` for desktop/menubar).
- `internal/runtime/`: Layer 1: Runtime (Observer, Event Bus, Scheduler, Lifecycle, Session, Permissions, Config).
- `internal/memory/`: Layer 2: Memory Engine (Working, Episodic, Semantic, Timeline, Retrieval, Compression, Store).
- `internal/services/`: Layer 3: Service Layer (Planner, Projects, Meetings, Habits, Notion, Calendar, Automation).
- `internal/aiprovider/`: Layer 4: AI Provider Layer (OpenAI, Claude, Gemini, Ollama, Local, Abstractions).
- `internal/interfaces/`: Layer 5: Interfaces (CLI, MCP, REST, Desktop, Menu Bar).
- `pkg/plugin/`: Public contracts and dynamic plugin interfaces.
- `docs/`: Architectural specifications, diagrams, ADRs.
- `context/`: Canonical context architecture v1.0 (long-term state, session logs, milestones, freeze logs, ADR index).
- `prompts/`: Reusable agent prompts and synchronization workflows.

## 4. Key Documents
- `VISION.md`, `ARCHITECTURE.md`, `PRINCIPLES.md`, `ROADMAP.md`, `PROJECT_STRUCTURE.md`, `DECISIONS.md`.
- Governance: `CONSTITUTION.md`, `ARCHITECTURE_INVARIANTS.md`, `ADR_GUIDE.md`, `QUALITY_GATES.md`, `CI_CD_SPECIFICATION.md`, `PROJECT_METRICS.md`.
- Specifications: `EVENT_BUS_SPECIFICATION.md`, `EVENT_MODEL.md`, `EVENT_CONTRACTS.md`, `EVENT_IMPLEMENTATION_PLAN.md`, `DISPATCHER_SPECIFICATION.md`, `DISPATCHER_CONTRACTS.md`, `DISPATCH_PIPELINE.md`, `DISPATCH_SEQUENCE.md`, `MIDDLEWARE_SPECIFICATION.md`, `MIDDLEWARE_CONTRACTS.md`, `MIDDLEWARE_PIPELINE.md`, `MIDDLEWARE_ORDERING.md`.

## 5. Pending Tasks
1. Implement Phase 1.2.5: Event Bus Middleware & Bus Facade implementation code (`internal/runtime/eventbus/{middleware, bus.go}`).
2. Implement Phase 1.3: Observer Component (`internal/runtime/observer`).
3. Implement Phase 1.4: Scheduler Component (`internal/runtime/scheduler`).
4. Implement Phase 1.5: Session & Permissions Engine (`internal/runtime/{session, permissions}`).
5. Implement Phase 2 (Memory Engine), Phase 3 (Domain Services), Phase 4 (AI Providers), Phase 5 (Client Interfaces).
