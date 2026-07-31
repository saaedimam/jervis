# jervis

[![Go Version](https://img.shields.io/github/go-mod/go-version/saaedimam/jervis)](https://go.dev/)
[![CI](https://github.com/saaedimam/jervis/actions/workflows/ci.yml/badge.svg)](https://github.com/saaedimam/jervis/actions/workflows/ci.yml)
[![CodeQL](https://github.com/saaedimam/jervis/actions/workflows/codeql.yml/badge.svg)](https://github.com/saaedimam/jervis/actions/workflows/codeql.yml)
[![Codecov](https://codecov.io/gh/saaedimam/jervis/branch/main/graph/badge.svg)](https://codecov.io/gh/saaedimam/jervis)
[![License](https://img.shields.io/github/license/saaedimam/jervis)](LICENSE)

> **Status**
>
> Active development. Interfaces may evolve before v1.0. Architectural invariants are documented under `docs/architecture/`.

Jervis is a local-first runtime and context operating system for deterministic automation, persistent memory, and AI-assisted workflows. The runtime does not require an AI provider; AI integrations are enabled only when the corresponding provider is configured.

---

## Why Jervis?

- **Determinism over probability.** The runtime executes deterministically. AI is isolated downstream and never controls system state.
- **Local sovereignty.** All data stays on your machine by default. Cloud access is opt-in and explicitly authorized per service.
- **Single-direction architecture.** Dependencies flow strictly downward. No upward calls. No cycles.
- **Extensible interfaces.** CLI, REST API, MCP server, and daemon mode over the same runtime core.

> **Design Principle**
>
> Runtime owns execution. AI consumes context. AI never owns state.

---

## Who is Jervis for?

Jervis is intended for developers building:

- Local AI agents
- MCP tools
- Personal automation
- Knowledge systems
- Deterministic workflows

It is not intended to be a cloud orchestration platform or a hosted AI service.

---

## Features

### Core Features (No External Dependencies)

- Local-first runtime for long-lived agent workflows
- Task and project management
- Persistent memory store
- Automation services

### Optional Integrations

Requires configuration to enable:

- **Notion sync** — `NOTION_TOKEN` + database IDs
- **Calendar** — iCal import/export
- **AI chat** — OpenAI, Anthropic, Gemini, or Ollama API keys

### Interfaces

- CLI
- REST API
- MCP server

---

## Quick Start

**Prerequisites**

- Go (see `go.mod`)
- macOS (Linux and Windows support planned)

**Build from source**

```bash
git clone https://github.com/saaedimam/jervis.git
cd jervis
make build
```

**Run**

```bash
./bin/jervis daemon
```

**Verify**

```bash
./bin/jervis version
```

---

## Minimal Example

```bash
./bin/jervis planner --create --id task-001 --title "Review architecture invariants"
./bin/jervis planner --list
```

---

## Architecture

```mermaid
graph TD
    Client[CLI / REST / MCP] --> Runtime
    Runtime --> Memory
    Runtime --> Services
    Services --> AI
```

The Runtime owns all execution, state, and permissions. Memory is fully decoupled from AI. Services handle domain logic and call AI providers only when explicitly needed.

See [docs/architecture/ARCHITECTURE.md](docs/architecture/ARCHITECTURE.md) for layer definitions and [docs/architecture/ARCHITECTURE_INVARIANTS.md](docs/architecture/ARCHITECTURE_INVARIANTS.md) for enforced invariants.

---

## Configuration

Jervis reads configuration from environment variables.

| Variable | Used by | Required | Description |
|----------|---------|----------|-------------|
| `NOTION_TOKEN` | Notion sync | No | Notion integration token |
| `MASTER_CONTEXT_ID` | Notion sync | No | Notion page ID for context sync |
| `TASKS_DB` | Notion sync | No | Notion database ID for task sync |
| `PACKAGES_DB` | Notion sync | No | Notion database ID for project sync |
| `MILESTONES_DB` | Notion sync | No | Notion database ID for milestone sync |
| `ADRS_DB` | Notion sync | No | Notion database ID for ADR sync |
| `SPECIFICATIONS_DB` | Notion sync | No | Notion database ID for spec sync |
| `OPENAI_API_KEY` | OpenAI provider | No | Enables OpenAI chat |
| `ANTHROPIC_API_KEY` | Anthropic provider | No | Enables Anthropic Claude chat |
| `GOOGLE_API_KEY` | Gemini provider | No | Enables Google Gemini chat |
| `OLLAMA_BASE_URL` | Ollama provider | No | Base URL for local Ollama instance |

No AI provider key is required to run the daemon or use core services.

---

## CLI Reference

| Command | Purpose |
|---------|---------|
| daemon | Start runtime |
| version | Version info |
| planner | Planner operations |
| sync | Synchronization |
| calendar | Calendar integration |
| chat | AI interaction |
| mcp | MCP server |
| api | REST API |
| automation | Automation |

Run `jervis <command> --help` for full flag documentation.

---

## Documentation

### Start Here

| Topic | Location |
|-------|----------|
| Architecture overview | [docs/architecture/ARCHITECTURE.md](docs/architecture/ARCHITECTURE.md) |
| Documentation index | [docs/README.md](docs/README.md) |
| Contributing | [CONTRIBUTING.md](CONTRIBUTING.md) |
| Roadmap | [docs/adr/ROADMAP.md](docs/adr/ROADMAP.md) |

<details>
<summary>Additional documentation</summary>

| Topic | Location |
|--------|----------|
| Architecture invariants | [docs/architecture/ARCHITECTURE_INVARIANTS.md](docs/architecture/ARCHITECTURE_INVARIANTS.md) |
| Engineering principles | [docs/principles/engineering.md](docs/principles/engineering.md) |
| Testing strategy | [docs/architecture/05_TESTING_STRATEGY.md](docs/architecture/05_TESTING_STRATEGY.md) |
| Security model | [docs/architecture/06_SECURITY_MODEL.md](docs/architecture/06_SECURITY_MODEL.md) |
| Release strategy | [docs/architecture/07_RELEASE_STRATEGY.md](docs/architecture/07_RELEASE_STRATEGY.md) |
| CI/CD specification | [docs/architecture/CI_CD_SPECIFICATION.md](docs/architecture/CI_CD_SPECIFICATION.md) |
| ADR process | [docs/architecture/ADR_GUIDE.md](docs/architecture/ADR_GUIDE.md) |
| ADR decisions | [docs/adr/DECISIONS.md](docs/adr/DECISIONS.md) |
| Governance & constitution | [docs/adr/CONSTITUTION.md](docs/adr/CONSTITUTION.md) |
| Quality gates | [docs/adr/QUALITY_GATES.md](docs/adr/QUALITY_GATES.md) |
| Runtime specs | [docs/specs/runtime/](docs/specs/runtime/) |
| Security specs | [docs/specs/security/](docs/specs/security/) |
| Repository governance | [docs/architecture/GITHUB_SETUP.md](docs/architecture/GITHUB_SETUP.md) |

</details>

---

## Development

```bash
make build        # Build binary to bin/jervis
make test         # Run tests with race detector
make lint         # Run golangci-lint
```

See [docs/README.md](docs/README.md) for architecture, specifications, ADRs, and engineering principles.

All pull requests must pass CI, lint, and test gates before merge. See [CONTRIBUTING.md](CONTRIBUTING.md) for branch conventions and commit format.

---

## Contributing

Jervis follows [Conventional Commits](https://www.conventionalcommits.org/) and a strict architectural review process for changes to layer boundaries or invariants. Read [CONTRIBUTING.md](CONTRIBUTING.md) before opening a pull request.

---

## Roadmap

See [docs/adr/ROADMAP.md](docs/adr/ROADMAP.md) for the full roadmap.

---

## License

Jervis is licensed under the [Apache License 2.0](LICENSE).
