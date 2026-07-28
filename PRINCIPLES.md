# Design Principles & Rules

## Core Design Principles

1. **Local-First**: Prioritize operating entirely on the user's local machine for privacy, security, and low latency.
2. **Runtime Ownership**: The Jervis Runtime owns all system state, event dispatching, and process lifecycles. AI never owns the runtime.
3. **Transparency & Auditability**: System events, permissions, and tool execution logs are recorded on an immutable timeline.
4. **User in Control**: Actions with side effects require authorization enforced by the Runtime Permissions module.
5. **Modularity & Decoupling**: Components interact strictly through defined interfaces. Storage and AI providers are implementation details.
6. **Graceful Degradation**: The system operates deterministically even when AI providers are offline or removed.

---

## Mandatory Design Rules

The following 15 rules must **NEVER** be violated under any circumstances:

1. **Observer never calls AI.**
2. **Event Bus never calls AI.**
3. **Memory never depends on AI.**
4. **Services never depend on a specific provider.**
5. **Runtime knows nothing about OpenAI, Claude, Gemini or any vendor.**
6. **Everything communicates through interfaces.**
7. **No cyclic dependency.**
8. **The runtime must continue working if every AI provider is removed.**
9. **Interfaces never contain business logic.**
10. **Business logic never depends on UI.**
11. **Storage is implementation detail.**
12. **Plugins cannot bypass Runtime.**
13. **Permissions are enforced before execution.**
14. **Runtime owns state.**
15. **AI only consumes context and produces responses.**

---

## Coding Standards

1. **Language & Style**: Strict adherence to the standard style guide of the chosen implementation language.
2. **Type Safety**: Enforce static typing and explicit interfaces across all layer boundaries.
3. **Documentation**: All public package contracts, layer interfaces, and services must be fully documented.
4. **Testing**: Unit tests are required for Runtime, Memory Engine, and Services. Integration tests for AI Provider adapters.
5. **Error Handling**: Use structured error types. Errors must be propagated through the Event Bus or returned cleanly; never swallowed silently.
