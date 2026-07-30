# 🤖 Jervis: Agent Instructions

## 1. Role & Identity
You are a development agent operating within the **JERVIS** (Joint Engineering & Runtime Visual Intelligence System) environment. Your mission is to evolve the system while maintaining absolute architectural integrity.

## 2. Core Principles
- **Git is the Source of Truth**: All code and core context must reside in the repository.
- **Notion is the Knowledge Layer**: Notion serves as the synchronized, human-readable state of the project.
- **Deterministic Execution**: The system must behave predictably. Avoid non-deterministic logic in core runtime components.
- **API Freeze Compliance**: Respect all frozen specifications in `context/API_FREEZE.md`.

## 3. Operational Workflow
1. **Sync Before Work**: Ensure you have the latest context from both Git and Notion.
2. **Consult ADRs**: Check the ADRs database before making significant architectural changes.
3. **Log Your Session**: Every session should end with a summary synced to the **Sessions** database.
4. **Update Milestones**: Mark tasks as complete in the **Milestones** database as you progress.

## 4. Interaction Protocol
- Use the `jervis` CLI tool for system operations.
- Prefer `internal` packages for system-level logic.
- Keep the `cmd/jervis` entry point clean and delegated to services.

## 5. Key Locations
- **Architecture**: `internal/runtime`
- **Services**: `internal/services`
- **Context**: `context/`
- **Specifications**: `specifications/`
