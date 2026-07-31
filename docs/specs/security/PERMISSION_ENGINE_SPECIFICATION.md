# Runtime Permission Engine Architecture Specification

## 1. Overview
This document specifies the canonical architecture for the **Runtime Permission Engine** (`internal/runtime/permissions`) in Project Jervis.

The Permission Engine is the central, non-bypassable security authority inside the Runtime layer. It validates capability requests and enforces capability-based access control (CBAC) for all operations prior to execution.

---

## 2. Core Principles & Decision Flow

### 2.1 High-Level Authorization Flow

```
PermissionEngine
       │
       ▼
Can this subject perform this action?
       │
       ├───► YES ───► Continue Execution
       │
       └───► NO  ───► Return PermissionDenied Error
```

### 2.2 Owned Responsibilities
1. **Capability Authorization**: Authorize requests (`Subject`, `Resource`, `Action`) against active security policies.
2. **Default Deny Policy**: Enforce a strict Default Deny posture; requests are denied unless explicitly permitted by an active policy rule.
3. **Explicit Deny Precedence**: Any matching `Deny` rule unconditionally overrides all matching `Allow` rules.
4. **Deterministic Evaluation**: Evaluate permissions synchronously using pure value semantics without side-effects or dynamic ambient context.

### 2.3 Explicit Non-Responsibilities & Architectural Boundaries
- **No AI Awareness**: Zero imports or references to `internal/aiprovider` (Invariant 2).
- **No Memory Engine Dependencies**: Operates independently of `internal/memory`.
- **No Service Layer Dependencies**: Operates independently of `internal/services`.
- **No Persistence or Filesystem Storage**: Policy storage is in-memory only for Phase 1.
- **No Background Concurrency**: Zero goroutines, channels, or background worker threads.

---

## 3. Evaluation Pipeline

```
   Access Request (Subject, Resource, Action)
                       │
                       ▼
   ┌──────────────────────────────────────────┐
   │ Stage 1: Request Structural Validation   │
   │ - Validate Subject, Resource, Action     │
   └──────────────────────────────────────────┘
                       │
                       ▼
   ┌──────────────────────────────────────────┐
   │ Stage 2: Policy Retrieval                │
   │ - Retrieve active Policy definitions     │
   └──────────────────────────────────────────┘
                       │
                       ▼
   ┌──────────────────────────────────────────┐
   │ Stage 3: Deny-First Rule Evaluation      │
   │ - If ANY rule returns EffectDeny -> DENY │
   └──────────────────────────────────────────┘
                       │
                       ▼
   ┌──────────────────────────────────────────┐
   │ Stage 4: Allow Rule Check                │
   │ - If ANY rule returns EffectAllow ->ALLOW│
   └──────────────────────────────────────────┘
                       │
                       ▼
   ┌──────────────────────────────────────────┐
   │ Stage 5: Default Deny Fallback           │
   │ - No match found -> Return DecisionDeny  │
   └──────────────────────────────────────────┘
```

---

## 4. Decision Resolution Rules

1. **Validation Failure**: If the request is malformed, return `DecisionDeny` with `errs.ErrValidationFailed`.
2. **Explicit Deny Match**: If any registered policy rule matches the capability and returns `EffectDeny`, immediately halt evaluation and return `DecisionDeny`.
3. **Explicit Allow Match**: If no `Deny` rule matched and at least one registered policy rule matches the capability and returns `EffectAllow`, return `DecisionAllow`.
4. **Default Deny Fallback**: If no rules match the request, return `DecisionDeny` with reason `"default deny policy enforced"`.

---

## 5. Performance & Complexity Guarantees

- **Time Complexity**: $O(P \cdot R)$ where $P$ is total registered policies and $R$ is average rules per policy.
- **Space Complexity**: $O(1)$ memory overhead per evaluation call.
- **Thread Safety & Isolation**: Single-threaded call stack evaluation guarantees 100% deterministic test execution.
