# Event Bus Middleware Ordering Specification

## 1. Canonical Ordering Policy: FIFO (First In, First Out)

The Jervis Event Bus adopts **FIFO (Registration Order)** as its single canonical ordering model for middleware.

```
Registration:   bus.Use(M1) -> bus.Use(M2) -> bus.Use(M3)

Pre-Dispatch:   M1 Pre  --->  M2 Pre  --->  M3 Pre
                                               │
                                               ▼
                                      Dispatcher.Dispatch
                                               │
Post-Dispatch:  M1 Post <---  M2 Post <---  M3 Post
```

---

## 2. Rationale for FIFO Selection

1. **Intuitive & Predictable**: Developers expect middleware registered first to intercept events first (e.g. logging/auditing registered first sees the raw event before authentication or rate limiting).
2. **Standard Onion Model**: Matches standard HTTP middleware design (e.g. Go `net/http` middleware, Express.js, standard inter-process wrappers).
3. **Deterministic Testing**: Registration order directly defines execution order, guaranteeing 100% deterministic test execution.

---

## 3. Immutability of Order

- Middleware order is fixed at registration time (`Use()`).
- Middleware order MUST NOT change dynamically during an active dispatch loop.
- Dynamic registration during event dispatch does not affect the currently executing dispatch chain.
