# Runtime Permission Model Specification

## 1. Domain Entities & Value Objects

The Jervis Permission Model consists of seven immutable core concepts:

### 1.1 `Subject`
Represents the security entity requesting permission.
- **Examples**: `"runtime:system"`, `"user:admin"`, `"plugin:git"`.
- **Validation**: Must be a non-empty, lower-case, colon-separated namespace string.

### 1.2 `Resource`
Represents the target asset or system resource.
- **Examples**: `"fs:config.json"`, `"event:system.user.created"`, `"env:PATH"`.
- **Validation**: Must be a non-empty string. Supports exact match and prefix wildcard (`"fs:*"`).

### 1.3 `Action`
Represents the operation performed on the resource.
- **Examples**: `"read"`, `"write"`, `"execute"`, `"publish"`.
- **Validation**: Must be a non-empty lower-case string.

### 1.4 `Capability`
An immutable value struct binding `Subject`, `Resource`, and `Action`.

### 1.5 `Effect`
The evaluation output of a single rule:
- `EffectAllow`: Explicitly permits the capability.
- `EffectDeny`: Explicitly denies the capability.
- `EffectNeutral`: The rule does not match the capability.

### 1.6 `Rule`
An immutable evaluation statement containing target patterns and an intended `Effect` (`EffectAllow` or `EffectDeny`).

### 1.7 `Policy`
An immutable container grouping a set of related `Rule` statements identified by a unique `PolicyID`.

---

## 2. Evaluation Rules & Precedence Matrix

| Explicit Deny Match | Explicit Allow Match | Final Decision | Reason |
| :---: | :---: | :---: | :--- |
| **YES** | YES | **DENY** | Explicit Deny rule takes precedence |
| **YES** | NO | **DENY** | Explicit Deny rule matched |
| NO | **YES** | **ALLOW** | Explicit Allow rule matched |
| NO | NO | **DENY** | Default Deny policy enforced |
