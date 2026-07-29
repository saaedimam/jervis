# Memory Engine Specification (Phase 2.1 - DRAFT)

## 1. Overview & Responsibilities

The Memory Engine (Layer 2) is responsible for the persistent and transient storage of all system state, context, and history. It serves as the "brain" of Project Jervis, enabling agents to recall past events, access long-term knowledge, and maintain active situational awareness.

### Phase 2.1 Responsibilities
1. **Working Memory**: Manage an in-memory, sliding-window store of recent context.
2. **Timeline Ledger**: Record a chronological, immutable ledger of all system events.

---

## 2. Architecture & Subsystem Boundaries

```
internal/memory/
├── contracts/        # Shared memory interfaces and types
├── working/          # Sliding window context store
├── timeline/         # Append-only event ledger
├── store/            # Persistence drivers (SQLite)
├── semantic/         # Vector/distance search
├── episodic/         # Episode-based indexing
├── retrieval/        # Context lookup pipeline
└── compression/      # Context pruning and summarization
```

---

## 3. Component Design

### 3.1 Working Memory (`internal/memory/working`)
- **Nature**: In-memory, non-persistent.
- **Data Model**: `Entry` (ID, Content, Metadata, Timestamp).
- **Policy**: Sliding window (FIFO) with a configurable capacity (e.g., 50 entries).
- **Operations**: `Add`, `Get`, `List`, `Prune`.

### 3.2 Timeline Ledger (`internal/memory/timeline`)
- **Nature**: Persistent (via `internal/memory/store`), append-only.
- **Data Model**: `Event` (Wraps `runtime.Event` with memory-specific metadata).
- **Invariants**: Immutable once written. Chronological ordering enforced.
- **Operations**: `Append`, `Query` (by time range or event type).

---

## 4. Implementation Invariants

- **Determinism**: Memory access and retrieval MUST be deterministic.
- **Thread-Safety**: All concurrent access MUST use RWMutex.
- **No AI Awareness**: Memory primitives are "dumb" storage; AI logic resides in Layer 4.
- **Zero Cycles**: Memory must not depend on higher layers (Service, AI).

---

## 5. Exit Criteria (Phase 2.1)

- [ ] Implementation of Working Memory with sliding window policy.
- [ ] Implementation of Timeline Ledger with append-only semantics.
- [ ] 100% statement coverage across Phase 2.1 packages.
- [ ] Integration tests verifying event persistence into Timeline.
