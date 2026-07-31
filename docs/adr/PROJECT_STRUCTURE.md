# Repository Structure

The `jervis` project follows a strict 5-tier layer directory layout enforcing single-direction dependency flow:

```
jervis/
├── cmd/                        # Entrypoints for Interfaces
│   ├── jervis/                 # CLI interface binary
│   ├── mcp/                    # MCP Server interface binary
│   └── daemon/                 # Desktop/Menu Bar REST daemon
├── internal/                   # Private core implementation
│   ├── runtime/                # LAYER 1: Runtime (System Owner)
│   │   ├── observer/           # System health & metric observation (no AI calls)
│   │   ├── eventbus/           # Asynchronous event broker (no AI calls)
│   │   ├── scheduler/          # Cron & background task scheduler
│   │   ├── lifecycle/          # Boot & shutdown process control
│   │   ├── session/            # Session context & state manager
│   │   ├── permissions/        # Security capability authorization
│   │   └── config/             # Configuration loading & validation
│   ├── memory/                 # LAYER 2: Memory Engine (Independent of AI)
│   │   ├── working/            # Active context sliding window
│   │   ├── episodic/           # Interaction history episodes
│   │   ├── semantic/           # Relational knowledge graph
│   │   ├── timeline/           # Immutable event timeline
│   │   ├── retrieval/          # Search & index retrieval pipeline
│   │   ├── compression/        # Heuristic context compression
│   │   └── store/              # Document & knowledge storage drivers
│   ├── services/               # LAYER 3: Service Layer (Domain Business Logic)
│   │   ├── planner/            # Task planning & execution logic
│   │   ├── projects/           # Local repository & project management
│   │   ├── meetings/           # Calendar sync & meeting agendas
│   │   ├── habits/             # Routine & habit tracking logic
│   │   ├── notion/             # Notion API integration service
│   │   ├── calendar/           # Calendar API service
│   │   └── automation/         # Local script & workflow automation
│   ├── aiprovider/             # LAYER 4: AI Provider Layer (Context Consumer)
│   │   ├── provider.go         # Standard Provider & Model interfaces
│   │   ├── openai/             # OpenAI client adapter
│   │   ├── claude/             # Anthropic Claude client adapter
│   │   ├── gemini/             # Google Gemini client adapter
│   │   ├── ollama/             # Local Ollama client adapter
│   │   └── local/              # llama.cpp / GGUF local model bindings
│   └── interfaces/             # LAYER 5: Client Interfaces (Zero Business Logic)
│       ├── cli/                # Terminal interface handlers
│       ├── mcp/                # MCP protocol server implementation
│       ├── rest/               # HTTP REST endpoint handlers
│       ├── desktop/            # Desktop GUI bindings
│       └── menubar/            # System tray / menu bar widget handlers
├── pkg/                        # Publicly exportable interfaces & SDKs
│   └── plugin/                 # Dynamic plugin interface contracts
├── docs/                       # Architectural specs, ADRs, & guides
├── scripts/                    # Build, cross-compile, & test scripts
└── tests/                      # Integration & end-to-end test suites
```

## Versioning
- Semantic Versioning 2.0.0 (`MAJOR.MINOR.PATCH`).
- Public contracts reside under `pkg/` and interface commands.
