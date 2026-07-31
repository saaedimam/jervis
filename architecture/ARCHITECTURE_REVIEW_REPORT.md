# Architecture Review Report: Project Jervis Reconciliation

**Author:** Principal Software Architect  
**Date:** 2026-07-29  
**Status:** Completed & Approved Baseline  

---

## Executive Summary

An architecture audit of Project Jervis revealed a fundamental structural inconsistency in the initial baseline. The previous design placed AI inside the core execution loop as a controlling "Cognitive Layer." This created an improper ownership model where the system runtime was dependent upon external AI providers and vendor APIs.

This reconciliation report details the identified flaws, the structural corrections applied, and the canonical 5-tier architecture now frozen for all future implementation.

---

## Inconsistencies Found & Proposed Corrections

### Inconsistency 1: AI Placed Inside Runtime Loop ("Cognitive Layer")
- **Location**: Previous `ARCHITECTURE.md` (Runtime Layers 2 & 4), `DECISIONS.md` (ADR-0001).
- **Why it was wrong**: Placing AI inside the runtime loop made AI the driver of system state, event dispatching, and orchestration. If an AI provider failed, went offline, or was removed, the runtime would collapse. It also violated system security by allowing AI inference to drive execution before permission checks.
- **Proposed Correction**: Removed the "Cognitive Layer" from the runtime. Replaced with a top-down **Runtime Ownership Model** where the Runtime owns the system, state, and permissions. AI is moved down to Layer 4 (`AI Provider Layer`), acting strictly as a pluggable context consumer.

### Inconsistency 2: Presentation Layer at the Top Driving Orchestration Directly
- **Location**: Previous `ARCHITECTURE.md` (Presentation Layer -> Orchestration Layer).
- **Why it was wrong**: UIs (CLI/API) were invoking orchestration logic directly without going through Runtime Permission enforcement or Event Bus distribution, mixing UI handling with business logic.
- **Proposed Correction**: Placed `Interfaces` at Layer 5. Interfaces contain zero business logic and must communicate with the system strictly by issuing requests to the `Runtime Permissions` module and `Event Bus`.

### Inconsistency 3: Memory Dependent on AI Summarization/Embeddings
- **Location**: Previous `ARCHITECTURE.md` (Data Layer / Memory Model).
- **Why it was wrong**: Memory operations implicitly relied on LLM summarization and embedding models, coupling core data storage to AI availability.
- **Proposed Correction**: Isolated the `Memory Engine` at Layer 2. Working Memory, Timeline logging, Retrieval indexing, and Heuristic Compression operate 100% deterministically without depending on AI providers.

### Inconsistency 4: Lack of Explicit Subcomponent Hierarchy & Security Boundaries
- **Location**: Previous `PROJECT_STRUCTURE.md` & `PRINCIPLES.md`.
- **Why it was wrong**: Components like Planner, Notion, Calendar, and Observer were unorganized or floating, making dependency cycles likely. Plugins could potentially bypass security controls.
- **Proposed Correction**: Enforced canonical subcomponent mappings across 5 distinct tiers, backed by 15 mandatory design rules. Plugins are explicitly prohibited from bypassing Runtime Permission checks.

---

## Canonical Architecture Hierarchy

```
OS
↓
Runtime (Observer, Event Bus, Scheduler, Lifecycle, Session Manager, Permissions, Configuration)
↓
Memory Engine (Working Memory, Episodic Memory, Semantic Memory, Timeline, Retrieval, Compression, Knowledge Store)
↓
Service Layer (Planner, Projects, Meetings, Habits, Notion, Calendar, Automation)
↓
AI Provider Layer (OpenAI, Claude, Gemini, Ollama, Local Models, Future Providers)
↓
Interfaces (CLI, MCP Server, REST API, Desktop, Menu Bar, Future Interfaces)
```

---

## Summary of 15 Mandatory Design Rules Enforced

1. Observer never calls AI.
2. Event Bus never calls AI.
3. Memory never depends on AI.
4. Services never depend on a specific provider.
5. Runtime knows nothing about OpenAI, Claude, Gemini or any vendor.
6. Everything communicates through interfaces.
7. No cyclic dependency.
8. The runtime must continue working if every AI provider is removed.
9. Interfaces never contain business logic.
10. Business logic never depends on UI.
11. Storage is implementation detail.
12. Plugins cannot bypass Runtime.
13. Permissions are enforced before execution.
14. Runtime owns state.
15. AI only consumes context and produces responses.
