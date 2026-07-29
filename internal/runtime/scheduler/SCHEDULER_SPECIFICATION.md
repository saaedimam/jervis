# Runtime Scheduler Specification (Phase 1.5.0 - DRAFT)

## 1. Overview & Responsibilities

The Runtime Scheduler is the canonical background task execution engine of Project Jervis. It manages recurring (Cron), interval-based, and deferred (Once) jobs within the deterministic Runtime layer.

### Core Responsibilities
1. **Multi-Strategy Scheduling**: Support Cron expressions, fixed intervals, and one-time deferred executions.
2. **Deterministic Ticking**: Provide a `Tick(time.Time)` method for manual engine driving and testing.
3. **Panic Isolation**: (TODO) Wrap job executions in `recover()` blocks to prevent a single failing job from crashing the OS.
4. **Thread-Safe Registry**: Allow dynamic job registration and removal at runtime.
5. **FIFO Execution Order**: If multiple jobs are due at the same tick, execute them in order of registration.
6. **Zero Side Effects (on system state)**: Jobs themselves may have side effects, but the scheduler engine is a pure execution coordinator.

---

## 2. Architecture & Subsystem Boundaries

```
internal/runtime/scheduler/
├── contracts/        # Interfaces (Job, Schedule, Registry, Scheduler)
├── model/            # Concrete implementations of Job and Schedules
├── errors/           # Canonical scheduler errors
├── registry/         # Thread-safe job storage
├── engine/           # Ticking and execution logic
└── scheduler.go      # Facade orchestrating background workers and registry
```

---

## 3. Supported Schedules

### Interval Schedule
- Triggers at a fixed `time.Duration`.
- Calculation: `NextRun = LastRun + Interval`.

### Once Schedule (Deferred)
- Triggers exactly once at or after a specific `time.Time`.
- Calculation: `NextRun = FixedTime` (returns zero time after first execution).

### Cron Schedule (Planned)
- Triggers based on standard Cron expressions (`min hour dom month dow`).
- Supports `*` (any) and discrete numeric values.

---

## 4. Implementation Invariants

- **Deterministic Core**: The `engine.Tick(now)` logic must be pure and deterministic based on the provided `now` time and the internal state of `lastRunAt` timestamps.
- **Async Execution Engine**: While the core logic is deterministic, the `Start()` method uses a background goroutine to drive the ticker at a defined resolution (e.g., 1 second).
- **NoAI Policy**: The scheduler must not depend on any AI Provider services.
- **Layer 1 Ownership**: The scheduler must not depend on Memory, Service, or AI layers.

---

## 5. Exit Criteria (Phase 1.5)

- [x] Implementation of core contracts and registry.
- [x] Implementation of Interval and Once schedules.
- [x] 100% statement coverage on registry and engine logic.
- [ ] Panic isolation for job handlers.
- [ ] Basic Cron expression support.
- [x] Facade with background worker loop.
