# Jervis Vision

## Vision
Jervis is a local-first, extensible, runtime-centric personal OS and service platform designed to augment developer productivity and automate personal workflows. The Jervis Runtime owns system lifecycle, event dispatching, state, memory, and security permissions. AI models (local or cloud) act strictly as pluggable processing providers, enabling Jervis to operate independently and deterministically even when no AI provider is configured or available.

## Scope
- **Core Runtime**: A deterministic, local-first runtime managing process lifecycles, event bus, background scheduling, session states, and permission enforcement.
- **Memory Engine**: Unified memory architecture managing working context, timeline logging, episodic storage, semantic knowledge, and memory compression without dependency on AI providers.
- **Service Layer**: Core business services managing projects, meetings, habits, planning, integrations (e.g., Notion, Calendar), and local automation.
- **AI Provider Abstraction**: Modular provider layer supporting replaceable AI engines (OpenAI, Anthropic Claude, Google Gemini, Ollama, local LLMs).
- **Multi-Interface Architecture**: Exposing Jervis capabilities through CLI, MCP Server, REST API, Desktop UI, and Menu Bar interfaces.

## Non-goals
- Placing AI in control of runtime lifecycle, event distribution, or permission authorization.
- Cloud-only multi-tenant SaaS hosting (focus is single-user, local-first environments).
- Re-training or fine-tuning foundation models within the runtime process space.
