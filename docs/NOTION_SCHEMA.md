# NOTION CANONICAL SCHEMA
## JERVIS Engineering Knowledge Graph v2.0
### Generated: 2026-07-29T06:20:00Z

---

## OVERVIEW

This document defines the complete property schema for every database in the JERVIS Notion workspace.

Each database is designed to support:
- Unique canonical IDs
- Rich metadata
- Bidirectional relations
- Queryable properties
- AI-friendly structure

---

## 1. 🏛️ ARCHITECTURE REGISTRY

**Database ID**: `d3dcb133-f96e-4e8e-944f-5825c2d1eee0`

**Purpose**: High-level system components and their relationships

### Properties

| Property | Type | Required | Options/Notes |
|----------|------|----------|---------------|
| **Architecture ID** | title | ✅ | ARCH-001, ARCH-002, etc. |
| **Name** | rich_text | ✅ | Human-readable name |
| **Layer** | select | ✅ | Layer 1, Layer 2, Layer 3, Layer 4, Layer 5 |
| **Status** | select | ✅ | In Progress, Complete, Deprecated |
| **Version** | rich_text | ✅ | Semantic version |
| **Owner** | people | ❌ | Assigned engineer |
| **Purpose** | rich_text | ✅ | Why this exists |
| **Responsibilities** | rich_text | ✅ | Bullet list of duties |
| **Dependencies** | relation | ❌ | → Architecture (self-reference) |
| **Dependents** | relation | ❌ | → Architecture (self-reference) |
| **Coverage** | number | ✅ | Test coverage % (0-100) |
| **Risk** | select | ✅ | Low, Medium, High |
| **Frozen** | checkbox | ✅ | API freeze status |
| **Related Packages** | relation | ✅ | → Package Registry |
| **Related Specifications** | relation | ✅ | → Specification Registry |
| **Related ADRs** | relation | ❌ | → ADR Database |
| **Related Sessions** | relation | ❌ | → Session Database |
| **Last Updated** | date | ✅ | Auto-updated timestamp |

### Example Entry

```
Architecture ID: ARCH-002
Name: Event Bus Engine
Layer: Layer 1
Status: Complete
Version: 1.0.0
Purpose: Synchronous in-process event routing
Responsibilities:
  • Priority-based dispatch
  • Panic isolation
  • Middleware chains
  • Aggregate errors
Coverage: 100
Risk: Low
Frozen: ☑️
Related Packages: PKG-007..014
Related Specifications: SPEC-001..006
Related ADRs: ADR-0002
```

---

## 2. 📦 PACKAGE REGISTRY

**Database ID**: `9c8bb7d5-5675-4cc5-b1b7-6a9c1ac3fe2f`

**Purpose**: Go packages and their metadata

### Properties

| Property | Type | Required | Options/Notes |
|----------|------|----------|---------------|
| **Package ID** | title | ✅ | PKG-001, PKG-002, etc. |
| **Name** | rich_text | ✅ | Package name (e.g., eventbus) |
| **Module** | rich_text | ✅ | Full import path |
| **Purpose** | rich_text | ✅ | One-line description |
| **Architecture** | relation | ✅ | → Architecture Registry |
| **Specifications** | relation | ❌ | → Specification Registry |
| **Source Path** | url | ✅ | GitHub link |
| **Coverage** | number | ✅ | Test coverage % |
| **Complexity** | number | ❌ | Cyclomatic complexity |
| **Exported APIs** | number | ❌ | Count of exported symbols |
| **Tests** | number | ❌ | Number of test files |
| **Files** | number | ❌ | Number of source files |
| **Frozen** | checkbox | ✅ | API freeze status |
| **Status** | select | ✅ | Active, Deprecated, Planning |
| **Owner** | people | ❌ | Maintainer |
| **Dependencies** | rich_text | ❌ | Comma-separated PKG IDs |
| **Dependents** | rich_text | ❌ | Comma-separated PKG IDs |
| **Related Sessions** | relation | ❌ | → Session Database |
| **Last Updated** | date | ✅ | Auto-updated |

### Example Entry

```
Package ID: PKG-014
Name: eventbus
Module: internal/runtime/eventbus
Purpose: Event Bus facade - canonical entry point
Architecture: ARCH-002
Source Path: https://github.com/...
Coverage: 100
Complexity: 3
Exported APIs: 8
Tests: 2
Files: 1
Frozen: ☑️
Status: Active
Dependencies: PKG-007..013
```

---

## 3. 📋 SPECIFICATION REGISTRY

**Database ID**: `f30e0d51-a787-421a-ad6b-77935f7d2e53`

**Purpose**: Engineering specifications and requirements

### Properties

| Property | Type | Required | Options/Notes |
|----------|------|----------|---------------|
| **Spec ID** | title | ✅ | SPEC-001, SPEC-002, etc. |
| **Name** | rich_text | ✅ | Specification title |
| **Version** | rich_text | ✅ | Spec version |
| **Status** | select | ✅ | Draft, Frozen, Superseded |
| **Architecture** | relation | ✅ | → Architecture Registry |
| **Packages** | relation | ❌ | → Package Registry |
| **Frozen** | checkbox | ✅ | Spec freeze status |
| **Superseded By** | relation | ❌ | → Specification (self) |
| **Related ADR** | relation | ❌ | → ADR Database |
| **Markdown File** | url | ✅ | Link to spec document |
| **Complexity** | select | ❌ | Simple, Moderate, Complex |
| **Acceptance Criteria** | rich_text | ✅ | Bullet list |
| **Implementation Notes** | rich_text | ❌ | Technical details |
| **Change History** | rich_text | ❌ | Version history |
| **Related Sessions** | relation | ❌ | → Session Database |
| **Last Updated** | date | ✅ | Auto-updated |

### Example Entry

```
Spec ID: SPEC-001
Name: Event Bus Specification
Version: 1.0.0
Status: Frozen
Architecture: ARCH-002
Packages: PKG-007..014
Frozen: ☑️
Markdown File: /docs/EVENT_BUS_SPECIFICATION.md
Complexity: Moderate
Acceptance Criteria:
  • Synchronous dispatch
  • Priority ordering
  • Panic isolation
  • 100% coverage
Related ADR: ADR-0002
```

---

## 4. 📄 FILE REGISTRY

**Database ID**: `d5b8d71a-c568-4288-9443-f3deb8b316bc`

**Purpose**: Source code files and their metadata

### Properties

| Property | Type | Required | Options/Notes |
|----------|------|----------|---------------|
| **File ID** | title | ✅ | FILE-0001, FILE-0002, etc. |
| **Path** | rich_text | ✅ | Relative path from repo root |
| **Package** | relation | ✅ | → Package Registry |
| **Architecture** | relation | ✅ | → Architecture Registry |
| **Specification** | relation | ❌ | → Specification Registry |
| **Language** | select | ✅ | Go, Markdown, YAML, JSON, Shell, Python, SQL |
| **Exports** | rich_text | ❌ | Exported symbols (comma list) |
| **Imports** | rich_text | ❌ | Imported packages |
| **Coverage** | rich_text | ❌ | Coverage % or N/A |
| **Frozen** | checkbox | ✅ | API freeze status |
| **Owner** | people | ❌ | File maintainer |
| **Status** | select | ✅ | Active, Deprecated, Planning |
| **Last Commit** | relation | ❌ | → Commit Database |
| **Last Session** | relation | ❌ | → Session Database |
| **API Count** | number | ❌ | Exported API count |
| **Lines** | number | ❌ | Lines of code |
| **Hash** | rich_text | ✅ | MD5 for change detection |
| **Created** | date | ❌ | File creation date |
| **Updated** | date | ✅ | Last modification |

### Example Entry

```
File ID: FILE-0008
Path: internal/runtime/eventbus/eventbus.go
Package: PKG-014
Architecture: ARCH-002
Language: Go
Exports: EventBus, New, Publish, Subscribe, Use, Count
Coverage: 100%
Frozen: ☑️
Status: Active
Lines: 110
Hash: a1b2c3d4e5f6...
Updated: 2026-07-29
```

---

## 5. ⚡ API REGISTRY

**Database ID**: `5e2dad61-5186-46f7-be6b-e7e5c3715f04`

**Purpose**: Exported APIs and their contracts

### Properties

| Property | Type | Required | Options/Notes |
|----------|------|----------|---------------|
| **API ID** | title | ✅ | API-001, API-002, etc. |
| **Name** | rich_text | ✅ | Function/method name |
| **Package** | relation | ✅ | → Package Registry |
| **File** | relation | ✅ | → File Registry |
| **Signature** | rich_text | ✅ | Full Go signature |
| **Inputs** | rich_text | ✅ | Input parameters |
| **Outputs** | rich_text | ✅ | Return values |
| **Errors** | rich_text | ❌ | Error types |
| **Version** | rich_text | ✅ | API version |
| **Breaking** | checkbox | ✅ | Breaking change flag |
| **Coverage** | number | ❌ | Test coverage % |
| **Tests** | relation | ❌ | → Test Database |
| **Specification** | relation | ❌ | → Specification Registry |
| **ADR** | relation | ❌ | → ADR Database |
| **Status** | select | ✅ | Stable, Experimental, Deprecated |
| **Related Sessions** | relation | ❌ | → Session Database |
| **Last Updated** | date | ✅ | Auto-updated |

### Example Entry

```
API ID: API-008
Name: EventBus.Publish
Package: PKG-014
File: FILE-0008
Signature: func (eb *EventBus) Publish(event eventcontracts.Event) error
Inputs: event - The event to publish
Outputs: error - AggregateError on handler failures
Version: 1.0.0
Breaking: ☐
Coverage: 100
Status: Stable
```

---

## 6. 📅 ADR DATABASE (Engineering Timeline)

**Database ID**: `abc5d892-1299-4813-b8bf-a143d6c8c73c`

**Purpose**: Architecture Decision Records

### Properties

| Property | Type | Required | Options/Notes |
|----------|------|----------|---------------|
| **ADR ID** | title | ✅ | ADR-0001, ADR-0002, etc. |
| **Title** | rich_text | ✅ | Decision title |
| **Status** | select | ✅ | Proposed, Accepted, Superseded, Rejected |
| **Date** | date | ✅ | Decision date |
| **Author** | people | ✅ | Decision author |
| **Context** | rich_text | ✅ | Background |
| **Problem** | rich_text | ✅ | What needed solving |
| **Decision** | rich_text | ✅ | What was decided |
| **Alternatives** | rich_text | ❌ | Options considered |
| **Consequences** | rich_text | ❌ | Positive impacts |
| **Trade-offs** | rich_text | ❌ | Negative impacts |
| **Affected Components** | relation | ❌ | → Architecture Registry |
| **Affected Packages** | relation | ❌ | → Package Registry |
| **Related Specifications** | relation | ❌ | → Specification Registry |
| **Superseded By** | relation | ❌ | → ADR (self) |
| **Review History** | rich_text | ❌ | Review log |
| **Timeline** | rich_text | ❌ | Key dates |
| **Related Sessions** | relation | ❌ | → Session Database |
| **Last Updated** | date | ✅ | Auto-updated |

### Example Entry

```
ADR ID: ADR-0002
Title: Event Bus Synchronous Dispatch
Status: Accepted
Date: 2026-07-20
Author: Engineering Team
Context: Runtime needs reliable event routing
Problem: Async dispatch introduces complexity
Decision: 100% synchronous with panic isolation
Alternatives: Async with channels, goroutine pools
Consequences: Deterministic, testable, simple
Trade-offs: No parallel dispatch
Affected Components: ARCH-002
Related Specs: SPEC-001..006
```

---

## 7. ✅ QUALITY GATES (Milestones)

**Database ID**: `39ae6e23-2bc1-4e34-a7b0-a1da9410b081`

**Purpose**: Project milestones and phases

### Properties

| Property | Type | Required | Options/Notes |
|----------|------|----------|---------------|
| **Milestone ID** | title | ✅ | M-001, PHASE-1.1, etc. |
| **Name** | rich_text | ✅ | Milestone name |
| **Phase** | select | ✅ | Phase 1.1, Phase 1.2, ..., Phase 5 |
| **Status** | select | ✅ | Done, In Progress, Pending, Blocked |
| **Coverage** | rich_text | ❌ | Test coverage achieved |
| **Start Date** | date | ❌ | When work began |
| **Target Date** | date | ❌ | Planned completion |
| **Completion Date** | date | ❌ | Actual completion |
| **Dependencies** | relation | ❌ | → Milestone (self) |
| **Dependents** | relation | ❌ | → Milestone (self) |
| **Related Sessions** | relation | ❌ | → Session Database |
| **Related Commits** | relation | ❌ | → Commit Database |
| **Deliverables** | rich_text | ✅ | What was delivered |
| **Risks** | rich_text | ❌ | Known risks |
| **Blockers** | rich_text | ❌ | Current blockers |
| **Last Updated** | date | ✅ | Auto-updated |

### Example Entry

```
Milestone ID: PHASE-1.2
Name: Event Bus Complete
Phase: Phase 1.2
Status: Done
Coverage: 100%
Completion Date: 2026-07-20
Deliverables:
  • Event Bus facade
  • Registry, Dispatcher, Middleware
  • 8 packages, 100% coverage
Related Sessions: SESSION-015
```

---

## 8. 🔄 COMMIT INTELLIGENCE

**Database ID**: `69c5145a-b84c-43e5-83b2-05d746a80e26`

**Purpose**: Git commit tracking and impact analysis

### Properties

| Property | Type | Required | Options/Notes |
|----------|------|----------|---------------|
| **Commit ID** | title | ✅ | Short SHA (e.g., 91ba7ad) |
| **SHA** | rich_text | ✅ | Full SHA |
| **Message** | rich_text | ✅ | Commit message |
| **Author** | people | ✅ | Commit author |
| **Branch** | rich_text | ✅ | Git branch |
| **Date** | date | ✅ | Commit timestamp |
| **Changed Files** | relation | ✅ | → File Registry |
| **Changed Packages** | relation | ✅ | → Package Registry |
| **Architecture Impact** | relation | ❌ | → Architecture Registry |
| **Specification Impact** | relation | ❌ | → Specification Registry |
| **Related Session** | relation | ✅ | → Session Database |
| **Related Milestone** | relation | ❌ | → Milestone Database |
| **Related Release** | relation | ❌ | → Release Database |
| **Review Status** | select | ❌ | Approved, Pending, Rejected |
| **Last Updated** | date | ✅ | Auto-updated |

### Example Entry

```
Commit ID: 91ba7ad
SHA: 91ba7ad7f5...
Message: chore: finalize governance verification
Author: Engineering
Branch: main
Date: 2026-07-29
Changed Files: FILE-XXXX
Changed Packages: PKG-XXX
Related Session: SESSION-021
Review Status: Approved
```

---

## 9. 🤖 AI HANDOFF (Session Database)

**Database ID**: `c1e36ebb-a3fc-4aea-a3d2-ac8214e1e40a`

**Purpose**: Development session logs and context handoff

### Properties

| Property | Type | Required | Options/Notes |
|----------|------|----------|---------------|
| **Session ID** | title | ✅ | SESSION-021, etc. |
| **Date** | date | ✅ | Session date |
| **Objective** | rich_text | ✅ | Session goal |
| **Summary** | rich_text | ✅ | What was accomplished |
| **Work Completed** | rich_text | ✅ | Detailed list |
| **Files Modified** | relation | ❌ | → File Registry |
| **Commits** | relation | ✅ | → Commit Database |
| **Architecture Changes** | relation | ❌ | → Architecture Registry |
| **Problems** | rich_text | ❌ | Issues encountered |
| **Solutions** | rich_text | ❌ | How issues were resolved |
| **AI Discussion** | rich_text | ❌ | Key AI interactions |
| **Lessons Learned** | relation | ❌ | → Lessons Database |
| **Remaining Work** | rich_text | ❌ | What's left |
| **Next Session** | rich_text | ❌ | Planned next steps |
| **Related Components** | relation | ❌ | Multi-select |
| **Related Specifications** | relation | ❌ | → Specification Registry |
| **Status** | select | ✅ | Complete, In Progress, Failed |
| **Last Updated** | date | ✅ | Auto-updated |

### Example Entry

```
Session ID: SESSION-021
Date: 2026-07-29
Objective: Complete Phase 2.1 Working Memory
Summary: Implemented Working Memory and Timeline
Work Completed:
  • Working Memory FIFO store
  • Timeline immutable ledger
  • 100% coverage achieved
Commits: 91ba7ad
Status: Complete
```

---

## 10. 🧠 ENGINEERING MEMORY

**Database ID**: `38a76b5b-b20e-498e-b6e9-e643c2ae7d8b`

**Purpose**: Patterns, lessons, and institutional knowledge

### Properties

| Property | Type | Required | Options/Notes |
|----------|------|----------|---------------|
| **Memory ID** | title | ✅ | MEM-001, PAT-001, etc. |
| **Type** | select | ✅ | Pattern, Lesson, Anti-Pattern, Research, Rule |
| **Name** | rich_text | ✅ | Memory title |
| **Category** | select | ❌ | Architecture, Testing, Performance, Security |
| **Summary** | rich_text | ✅ | One-line summary |
| **Context** | rich_text | ✅ | When this applies |
| **Content** | rich_text | ✅ | Full explanation |
| **Code Examples** | rich_text | ❌ | Code snippets |
| **Benefits** | rich_text | ❌ | Why use this |
| **Risks** | rich_text | ❌ | Potential downsides |
| **Related Components** | relation | ❌ | → Architecture Registry |
| **Related Packages** | relation | ❌ | → Package Registry |
| **Related Sessions** | relation | ❌ | → Session Database |
| **References** | rich_text | ❌ | Links, docs |
| **Tags** | multi_select | ❌ | Searchable tags |
| **Last Updated** | date | ✅ | Auto-updated |

### Example Entry

```
Memory ID: PAT-001
Type: Pattern
Name: Synchronous Event Dispatch
Category: Architecture
Summary: Use synchronous dispatch for determinism
Context: Event Bus, Observer
Content: Always prefer synchronous execution...
Benefits: Deterministic, testable, simple
Related Components: ARCH-002, ARCH-004
Tags: event-bus, patterns, sync
```

---

## 11. 🔗 DEPENDENCY GRAPH

**Database ID**: `1de04b92-6fe3-4756-b85d-c9370f838a3b`

**Purpose**: Package and component dependencies

### Properties

| Property | Type | Required | Options/Notes |
|----------|------|----------|---------------|
| **Relation ID** | title | ✅ | DEP-001, etc. |
| **Source** | relation | ✅ | → Package Registry |
| **Target** | relation | ✅ | → Package Registry |
| **Type** | select | ✅ | Imports, Uses, Implements, DependsOn |
| **Strength** | select | ✅ | Strong, Weak, Optional |
| **Direction** | select | ✅ | Uni, Bi |
| **Created** | date | ❌ | When relation established |
| **Verified** | checkbox | ❌ | Verified by CI |
| **Notes** | rich_text | ❌ | Additional context |

### Example Entry

```
Relation ID: DEP-001
Source: PKG-014 (eventbus)
Target: PKG-007 (contracts)
Type: Imports
Strength: Strong
Direction: Uni
Verified: ☑️
```

---

## PROPERTY TYPE REFERENCE

### Title
- Unique identifier
- Immutable
- Used for canonical IDs

### Rich Text
- Multi-line text
- Supports formatting
- URLs, mentions

### Select
- Single choice
- Predefined options
- Color coding

### Multi-Select
- Multiple choices
- Predefined options
- Tag-style display

### Relation
- Links to other database
- Bidirectional when supported
- Shows related pages

### Number
- Numeric values
- Supports formatting
- Formulas possible

### Checkbox
- Boolean
- Checked/Unchecked
- Quick filtering

### Date
- Date and time
- Date ranges
- Filtering by date

### URL
- Web links
- GitHub, docs
- File paths

### People
- Team members
- Assignees
- Notifications

---

## SCHEMA VERSIONING

**Version**: 2.0.0
**Last Updated**: 2026-07-29
**Schema Change Policy**: 
- Minor additions: Add properties freely
- Major changes: Requires migration plan
- Breaking changes: Requires ADR

**Next Review**: After Phase 2.2 completion

---

**Schema Complete**: All 11 databases defined
**Total Properties Defined**: 200+
**Relations Defined**: 50+
**Status**: Ready for implementation
