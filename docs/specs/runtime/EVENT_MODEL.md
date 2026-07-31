# Event Model Specification

## 1. Overview
This document defines the canonical Event Envelope and Data Structure for Project Jervis. Every message transmitted across the Jervis Runtime Event Bus MUST strictly conform to this specification.

---

## 2. Event Envelope Structure

An Event is an immutable data structure containing standardized headers and an abstract payload payload:

| Field Name | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `EventID` | `types.EventID` | MUST | Unique, non-empty identifier for the individual event instance. |
| `EventType` | `string` | MUST | Namespaced dot-separated event classification (e.g. `runtime.lifecycle.booted`). |
| `Source` | `string` | MUST | Identifier of the originating component or subsystem. |
| `Timestamp` | `types.Timestamp` | MUST | Canonical UTC creation timestamp. |
| `CorrelationID` | `string` | MUST | Identifier tracing a root request or multi-step workflow. |
| `CausationID` | `string` | MUST | Identifier of the immediate parent event that triggered this event. |
| `Priority` | `int` | MUST | Handler execution priority order (Higher numeric values execute first). |
| `Payload` | `any` | SHOULD | Domain-specific data object or struct payload. |
| `Metadata` | `map[string]string` | OPTIONAL | Arbitrary key-value metadata headers (context, tracing, attributes). |
| `Version` | `string` | MUST | Semantic version of the event schema (e.g. `1.0.0`). |

---

## 3. Field Constraints & Rules

### 3.1 `EventID`
- MUST be initialized using `types.NewEventID()`.
- MUST NOT be empty or zero-valued.
- MUST be unique per generated event.

### 3.2 `EventType`
- MUST follow lowercase dot-separated namespace format: `<layer>.<component>.<verb>`.
- Examples:
  - `runtime.lifecycle.starting`
  - `runtime.session.created`
  - `memory.timeline.appended`
- MUST NOT contain whitespace, special characters, or uppercase letters.

### 3.3 `Source`
- MUST identify the originating subsystem package path or logical component name.
- Example: `internal/runtime/lifecycle`.

### 3.4 `Timestamp`
- MUST be populated using `types.Now()` at the exact moment of event instantiation.
- MUST be in UTC time zone.

### 3.5 `CorrelationID` & `CausationID`
- `CorrelationID` MUST match the root event ID of the initiating flow. If the event is the root, `CorrelationID` MUST equal `EventID`.
- `CausationID` MUST equal the `EventID` of the direct cause. If the event has no parent, `CausationID` MUST equal `EventID`.

### 3.6 `Priority`
- Integer value ranging from `-100` (lowest priority / background) to `+100` (highest priority / critical system events).
- Standard default priority MUST be `0`.

### 3.7 `Payload`
- MUST be read-only and treated as immutable once attached to the event envelope.
- Handlers MUST NOT mutate payload state directly.

### 3.8 `Metadata`
- Map of string key-value pairs for metadata propagation (e.g., `trace_id`, `user_id`, `session_id`).
- MUST NOT contain nil or unserializable types.

### 3.9 `Version`
- MUST specify the event schema semantic version. Default version for initial events MUST be `1.0.0`.

---

## 4. Immutability & Safety Rules

1. **Envelope Freeze**: Once an event envelope is constructed, its headers and metadata MUST NOT be modified.
2. **No Mutating Handlers**: Handlers receiving an event MUST NOT alter payload properties.
3. **Defensive Cloning**: If a handler requires mutating payload state, it MUST create a deep copy before performing modifications.
