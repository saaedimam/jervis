# Architecture Decision Records (ADR)

## ADR-0001: Initial Architecture Baseline

**Date:** 2026-07-29  
**Status:** Superseded by ADR-0002

---

## ADR-0002: Architecture Reconciliation - Runtime Ownership & AI Decoupling

**Date:** 2026-07-29  
**Status:** Accepted

### Context
Initial architecture drafts placed AI inside the core runtime loop as a "Cognitive Layer," implicitly allowing AI logic to drive orchestration and event dispatching. This created a coupling where the runtime depended on AI availability and vendor interfaces. 

### Decision
1. **Ownership Hierarchy**: Redefine system architecture into 5 distinct, single-direction layers:
   `OS -> Runtime -> Memory Engine -> Service Layer -> AI Provider Layer -> Interfaces`.
2. **Runtime Ownership**: The Runtime layer (Observer, Event Bus, Scheduler, Lifecycle, Session Manager, Permissions, Configuration) owns system state, execution control, and security permissions.
3. **AI Provider Decoupling**: AI is strictly a consumer of context payloads and a producer of text/tool responses. AI providers reside *below* the Service Layer and have *zero* ownership of runtime lifecycle or event flow.
4. **Mandatory 15 Design Rules**: Enforce 15 architecture rules preventing runtime/memory/bus dependencies on AI and prohibiting cyclic coupling.
5. **AI Independence**: The runtime, memory engine, and service layer must remain 100% operational even if all AI providers are removed or offline.

### Consequences
- **Positive**: Complete vendor independence, rock-solid security (permissions checked by Runtime before execution), deterministic local workflows, zero runtime downtime if AI providers fail.
- **Negative**: Requires strict interface boundary discipline; services cannot directly pass execution tokens to LLMs without going through memory and provider abstraction interfaces.
