commit a0b4bfa386bd1fc8b08b072050711f6ea171b690
Author: IoriImasu <sayedimam.fahim@gmail.com>
Date:   Sat Aug 1 19:52:56 2026 +0600

    docs: streamline README for developer onboarding

diff --git a/README.md b/README.md
index ce944dc..5e130ad 100644
--- a/README.md
+++ b/README.md
@@ -6,59 +6,41 @@
 [![Codecov](https://codecov.io/gh/saaedimam/jervis/branch/main/graph/badge.svg)](https://codecov.io/gh/saaedimam/jervis)
 [![License](https://img.shields.io/github/license/saaedimam/jervis)](LICENSE)
 
-> **Status**
->
-> v1.0.0 released. See [releases](https://github.com/saaedimam/jervis/releases).
+Jervis is a local-first runtime for deterministic automation, persistent memory, and optional AI-assisted workflows.
 
-Jervis is a local-first runtime and context operating system for deterministic automation, persistent memory, and AI-assisted workflows. The runtime does not require an AI provider; AI integrations are enabled only when the corresponding provider is configured.
+It provides a single runtime for local services and interfaces such as the CLI, REST API, and MCP server. Core execution does not require an AI provider; external integrations are enabled explicitly through configuration.
 
 ---
 
 ## Why Jervis?
 
-- **Determinism over probability.** The runtime executes deterministically. AI is isolated downstream and never controls system state.
-- **Local sovereignty.** All data stays on your machine by default. Cloud access is opt-in and explicitly authorized per service.
-- **Single-direction architecture.** Dependencies flow strictly downward. No upward calls. No cycles.
-- **Extensible interfaces.** CLI, REST API, MCP server, and daemon mode over the same runtime core.
-
-> **Design Principle**
->
-> Runtime owns execution. AI consumes context. AI never owns state.
-
----
-
-## Who is Jervis for?
-
-Jervis is intended for developers building:
-
-- Local AI agents
-- MCP tools
-- Personal automation
-- Knowledge systems
-- Deterministic workflows
-
-It is not intended to be a cloud orchestration platform or a hosted AI service.
+- **Deterministic core** — runtime execution and state management do not depend on probabilistic AI output.
+- **Local-first operation** — core data and services run locally by default.
+- **AI as an optional dependency** — providers consume context without owning runtime state.
+- **Explicit boundaries** — runtime, memory, services, providers, and interfaces have defined dependency directions.
+- **Multiple interfaces, one runtime** — CLI, REST, and MCP expose the same underlying system.
 
 ---
 
 ## Features
 
-### Core Features (No External Dependencies)
+**Runtime and local services**
 
-- Local-first runtime for long-lived agent workflows
+- Persistent memory
 - Task and project management
-- Persistent memory store
 - Automation services
+- Calendar import and export
+- Background daemon
 
-### Optional Integrations
+**Optional integrations**
 
-Requires configuration to enable:
+- Notion synchronization
+- OpenAI
+- Anthropic
+- Google Gemini
+- Ollama
 
-- **Notion sync** — `NOTION_TOKEN` + database IDs
-- **Calendar** — iCal import/export
-- **AI chat** — OpenAI, Anthropic, Gemini, or Ollama API keys
-
-### Interfaces
+**Interfaces**
 
 - CLI
 - REST API
@@ -68,12 +50,9 @@ Requires configuration to enable:
 
 ## Quick Start
 
-**Prerequisites**
-
-- Go (see `go.mod`)
-- macOS (Linux and Windows support planned)
+Requires Go as specified by [`go.mod`](go.mod).
 
-**Build from source**
+### Installation
 
 ```bash
 git clone https://github.com/saaedimam/jervis.git
@@ -81,64 +60,72 @@ cd jervis
 make build
 ```
 
-**Run**
+### Run
 
 ```bash
-./bin/jervis daemon
+./bin/jervis version
 ```
 
-**Verify**
+### Verify
 
 ```bash
-./bin/jervis version
+./bin/jervis planner -create -id task-001 -title "First Jervis task"
+./bin/jervis planner -list
 ```
 
+No AI provider or external service is required for these commands.
+
 ---
 
 ## Minimal Example
 
+Create and retrieve a local task:
+
 ```bash
 ./bin/jervis planner -create -id task-001 -title "Review architecture invariants"
 ./bin/jervis planner -list
 ```
 
+Jervis stores planner state locally in `jervis.db`.
+
 ---
 
 ## Architecture
 
 ```mermaid
-graph TD
-    Client[CLI / REST / MCP] --> Runtime
-    Runtime --> Memory
-    Runtime --> Services
-    Services --> AI
+flowchart TD
+    I[CLI / REST / MCP] --> R[Runtime]
+    R --> M[Memory]
+    R --> S[Services]
+    S --> A[AI Providers]
+    S --> X[External Integrations]
 ```
 
-The Runtime owns all execution, state, and permissions. Memory is fully decoupled from AI. Services handle domain logic and call AI providers only when explicitly needed.
+The runtime owns execution and state. Memory is independent of AI, while domain services use providers and external integrations only when required.
 
-See [docs/architecture/ARCHITECTURE.md](docs/architecture/ARCHITECTURE.md) for layer definitions and [docs/architecture/ARCHITECTURE_INVARIANTS.md](docs/architecture/ARCHITECTURE_INVARIANTS.md) for enforced invariants.
+For the canonical design, see [`docs/architecture/`](docs/architecture/), the [architecture overview](docs/architecture/ARCHITECTURE.md), and [architecture invariants](docs/architecture/ARCHITECTURE_INVARIANTS.md).
 
 ---
 
 ## Configuration
 
-Jervis reads configuration from environment variables.
+Core local operation requires no environment variables.
 
-| Variable | Used by | Required | Description |
-|----------|---------|----------|-------------|
-| `NOTION_TOKEN` | Notion sync | No | Notion integration token |
-| `MASTER_CONTEXT_ID` | Notion sync | No | Notion page ID for context sync |
-| `TASKS_DB` | Notion sync | No | Notion database ID for task sync |
-| `PACKAGES_DB` | Notion sync | No | Notion database ID for project sync |
-| `MILESTONES_DB` | Notion sync | No | Notion database ID for milestone sync |
-| `ADRS_DB` | Notion sync | No | Notion database ID for ADR sync |
-| `SPECIFICATIONS_DB` | Notion sync | No | Notion database ID for spec sync |
-| `OPENAI_API_KEY` | OpenAI provider | No | Enables OpenAI chat |
-| `ANTHROPIC_API_KEY` | Anthropic provider | No | Enables Anthropic Claude chat |
-| `GOOGLE_API_KEY` | Gemini provider | No | Enables Google Gemini chat |
-| `OLLAMA_BASE_URL` | Ollama provider | No | Base URL for local Ollama instance |
+| Variable | Required | Default | Description |
+|-----------|----------|---------|-------------|
+| `NOTION_TOKEN` | No | — | Enables Notion integration |
+| `MASTER_CONTEXT_ID` | No | — | Notion page used for context synchronization |
+| `TASKS_DB` | No | — | Notion tasks database |
+| `PACKAGES_DB` | No | — | Notion projects/packages database |
+| `MILESTONES_DB` | No | — | Notion milestones database |
+| `ADRS_DB` | No | — | Notion ADR database |
+| `SPECIFICATIONS_DB` | No | — | Notion specifications database |
+| `OPENAI_API_KEY` | No | — | Enables the OpenAI provider |
+| `ANTHROPIC_API_KEY` | No | — | Enables the Anthropic provider |
+| `GOOGLE_API_KEY` | No | — | Enables the Google provider |
+| `OLLAMA_BASE_URL` | No | — | Ollama server base URL |
 
-No AI provider key is required to run the daemon or use core services.
+Jervis loads environment variables from the process environment and attempts to load a local `.env` file at startup.
 
 ---
 
@@ -146,77 +133,77 @@ No AI provider key is required to run the daemon or use core services.
 
 | Command | Purpose |
 |---------|---------|
-| daemon | Start runtime |
-| version | Version info |
-| planner | Planner operations |
-| sync | Synchronization |
-| calendar | Calendar integration |
-| chat | AI interaction |
-| mcp | MCP server |
-| api | REST API |
-| automation | Automation |
-
-Run `jervis <command> --help` for full flag documentation.
+| `version` | Show build and version information |
+| `planner` | Create and list planned tasks |
+| `sync` | Synchronize supported state with Notion |
+| `calendar` | Import or export iCalendar data |
+| `chat` | Send prompts through a configured AI provider |
+| `mcp` | Start the MCP server |
+| `api` | Start the REST API |
+| `automation` | Manage automation workflows |
+| `daemon` | Start the runtime daemon |
+
+Use command-specific help for available flags:
+
+```bash
+./bin/jervis planner --help
+./bin/jervis sync --help
+./bin/jervis calendar --help
+./bin/jervis chat --help
+./bin/jervis api --help
+./bin/jervis automation --help
+```
 
 ---
 
 ## Documentation
 
-### Start Here
-
 | Topic | Location |
 |-------|----------|
-| Architecture overview | [docs/architecture/ARCHITECTURE.md](docs/architecture/ARCHITECTURE.md) |
-| Documentation index | [docs/README.md](docs/README.md) |
-| Contributing | [CONTRIBUTING.md](CONTRIBUTING.md) |
-| Roadmap | [docs/adr/ROADMAP.md](docs/adr/ROADMAP.md) |
-
-<details>
-<summary>Additional documentation</summary>
-
-| Topic | Location |
-|--------|----------|
-| Architecture invariants | [docs/architecture/ARCHITECTURE_INVARIANTS.md](docs/architecture/ARCHITECTURE_INVARIANTS.md) |
-| Engineering principles | [docs/principles/engineering.md](docs/principles/engineering.md) |
-| Testing strategy | [docs/architecture/05_TESTING_STRATEGY.md](docs/architecture/05_TESTING_STRATEGY.md) |
-| Security model | [docs/architecture/06_SECURITY_MODEL.md](docs/architecture/06_SECURITY_MODEL.md) |
-| Release strategy | [docs/architecture/07_RELEASE_STRATEGY.md](docs/architecture/07_RELEASE_STRATEGY.md) |
-| CI/CD specification | [docs/architecture/CI_CD_SPECIFICATION.md](docs/architecture/CI_CD_SPECIFICATION.md) |
-| ADR process | [docs/architecture/ADR_GUIDE.md](docs/architecture/ADR_GUIDE.md) |
-| ADR decisions | [docs/adr/DECISIONS.md](docs/adr/DECISIONS.md) |
-| Governance & constitution | [docs/adr/CONSTITUTION.md](docs/adr/CONSTITUTION.md) |
-| Quality gates | [docs/adr/QUALITY_GATES.md](docs/adr/QUALITY_GATES.md) |
-| Runtime specs | [docs/specs/runtime/](docs/specs/runtime/) |
-| Security specs | [docs/specs/security/](docs/specs/security/) |
-| Repository governance | [docs/architecture/GITHUB_SETUP.md](docs/architecture/GITHUB_SETUP.md) |
-
-</details>
+| Documentation index | [`docs/README.md`](docs/README.md) |
+| Architecture | [`docs/architecture/`](docs/architecture/) |
+| Architecture overview | [`docs/architecture/ARCHITECTURE.md`](docs/architecture/ARCHITECTURE.md) |
+| Architecture invariants | [`docs/architecture/ARCHITECTURE_INVARIANTS.md`](docs/architecture/ARCHITECTURE_INVARIANTS.md) |
+| Specifications | [`docs/specs/`](docs/specs/) |
+| ADRs | [`docs/adr/`](docs/adr/) |
+| Engineering principles | [`docs/principles/engineering.md`](docs/principles/engineering.md) |
+| Security model | [`docs/architecture/06_SECURITY_MODEL.md`](docs/architecture/06_SECURITY_MODEL.md) |
+| Testing strategy | [`docs/architecture/05_TESTING_STRATEGY.md`](docs/architecture/05_TESTING_STRATEGY.md) |
+| Roadmap | [`docs/adr/ROADMAP.md`](docs/adr/ROADMAP.md) |
+| Contributing | [`CONTRIBUTING.md`](CONTRIBUTING.md) |
 
 ---
 
 ## Development
 
 ```bash
-make build        # Build binary to bin/jervis
-make test         # Run tests with race detector
-make lint         # Run golangci-lint
+make test
+make lint
+make build
 ```
 
-See [docs/README.md](docs/README.md) for architecture, specifications, ADRs, and engineering principles.
-
-All pull requests must pass CI, lint, and test gates before merge. See [CONTRIBUTING.md](CONTRIBUTING.md) for branch conventions and commit format.
+`make test` runs the Go test suite with the race detector. `make lint` requires `golangci-lint`.
 
 ---
 
 ## Contributing
 
-Jervis follows [Conventional Commits](https://www.conventionalcommits.org/) and a strict architectural review process for changes to layer boundaries or invariants. Read [CONTRIBUTING.md](CONTRIBUTING.md) before opening a pull request.
+Contributions should preserve Jervis's architectural boundaries, deterministic runtime behavior, and test requirements.
+
+Read [`CONTRIBUTING.md`](CONTRIBUTING.md) before opening a pull request.
 
 ---
 
 ## Roadmap
 
-See [docs/adr/ROADMAP.md](docs/adr/ROADMAP.md) for the full roadmap.
+The project roadmap is organized around four high-level areas:
+
+- Core runtime and memory
+- Domain services and local automation
+- AI provider integrations
+- Interface and ecosystem expansion
+
+See [`docs/adr/ROADMAP.md`](docs/adr/ROADMAP.md) for the canonical roadmap and current scope.
 
 ---
 
