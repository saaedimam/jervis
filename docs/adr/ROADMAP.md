# Jervis Roadmap

## Phase 1: Core Runtime & Memory Engine Foundations
- Implement **Runtime Layer**: Observer, Event Bus, Scheduler, Lifecycle, Session Manager, Permissions, Configuration.
- Build **Memory Engine**: Working Memory, Timeline logging, Episodic & Semantic Memory stores, Retrieval pipeline, Compression.
- Ensure 100% deterministic operation without any AI Provider configured.

## Phase 2: Domain Services & Local Automation
- Implement **Service Layer**: Planner, Projects, Meetings, Habits, Notion integration, Calendar sync, Automation engine.
- Enforce Runtime Permission checks prior to all service executions.
- Build local workflow automation and CLI interface contracts.

## Phase 3: AI Provider Layer Integration
- Implement **AI Provider Abstraction Layer**: Standard interfaces for context consumption and response generation.
- Build adapters for OpenAI, Anthropic Claude, Google Gemini, local Ollama, and llama.cpp bindings.
- Integrate Services with AI Provider Layer for opt-in intelligence augmentation.

## Phase 4: Multi-Interface & Ecosystem Expansion
- Expose **Interfaces**: Model Context Protocol (MCP) Server, REST API daemon, Desktop application, Menu Bar widget.
- Build dynamic plugin loader bound to Runtime Permission checks.
- Establish multi-interface client synchronization.
