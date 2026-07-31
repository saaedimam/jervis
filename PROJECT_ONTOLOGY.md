# Jervis Engineering Knowledge Graph - Project Ontology

## Entity Types

### Core Engineering Artifacts

| Type | Description | Example |
|------|-------------|---------|
| Architecture | High-level system component | ARCH-001 (Runtime) |
| Package | Go module package | PKG-014 (eventbus) |
| File | Source code file | FILE-0014 (eventbus.go) |
| API | Exported interface/function | API-042 |
| Test | Test suite/case | TEST-103 |
| Specification | Engineering specification | SPEC-001 |
| ADR | Architecture Decision Record | ADR-0002 |

### Process Artifacts

| Type | Description | Example |
|------|-------------|---------|
| Session | Development session | SESSION-017 |
| Commit | Git commit | COMMIT-91ba7ad |
| Milestone | Project milestone | Phase 1.4.1 |
| Release | Version release | REL-0.1.0 |
| Task | Development task | TASK-087 |
| Bug | Issue/bug | BUG-029 |

### Knowledge Artifacts

| Type | Description | Example |
|------|-------------|---------|
| Pattern | Design pattern | Pattern-001 |
| Anti-Pattern | Anti-pattern to avoid | AntiPattern-003 |
| Lesson | Lesson learned | Lesson-007 |
| Rule | Coding rule | Rule-015 |
| Research | Research spike | Research-002 |
| Prompt | Reusable prompt | Prompt-012 |
| Risk | Identified risk | Risk-005 |
| TechnicalDebt | Technical debt item | Debt-011 |

---

## Relationship Types

### Containment (hierarchical)

| Relationship | Description | Example |
|--------------|-------------|---------|
| `contains` | Parent-child containment | Architecture → Package → File |
| `belongs_to` | Membership | File → Package, Package → Architecture |
| `part_of` | Component membership | API → Package |

### Implementation

| Relationship | Description | Example |
|--------------|-------------|---------|
| `implements` | Realization | Package → Specification |
| `exports` | Public API exposure | File → API |
| `validates` | Test coverage | Test → API |
| `documents` | Documentation | Specification → Architecture |

### Approval & Governance

| Relationship | Description | Example |
|--------------|-------------|---------|
| `approved_by` | Decision approval | Specification → ADR |
| `supersedes` | Version replacement | ADR-0002 → ADR-0001 |
| `defined_by` | Specification definition | Architecture → Specification |

### Process Flow

| Relationship | Description | Example |
|--------------|-------------|---------|
| `introduced_by` | Creation | File → Commit |
| `affects` | Impact | Commit → Architecture |
| `produces` | Generation | Session → Commit |
| `tracked_by` | Progress tracking | ADR → Milestone |

---

## Constraints (Validation Rules)

### Mandatory Relationships

1. **File → Package**: Every file MUST belong to exactly one Package
2. **Package → Architecture**: Every package MUST belong to exactly one Architecture
3. **Package → Specification**: Every frozen Package MUST implement at least one Specification
4. **Specification → ADR**: Every Specification MUST be approved by an ADR
5. **Session → Commit**: Every Session MUST reference exactly one Commit
6. **Commit → Files**: Every Commit MUST touch at least one File

### Cardinality Rules

| Source | Relation | Target | Cardinality |
|--------|----------|--------|-------------|
| File | belongs_to | Package | 1:1 |
| Package | belongs_to | Architecture | 1:1 |
| API | exported_by | File | 1:1 |
| Test | validates | API | 1:N |
| Specification | approved_by | ADR | N:1 |
| Session | produces | Commit | 1:1 |
| Commit | touches | File | 1:N |

---

## Validation Queries

### Integrity Checks

1. **Orphan Files**: Find files not linked to packages
2. **Orphan Packages**: Find packages not linked to architecture
3. **Unapproved Specs**: Find specifications without ADR approval
4. **Missing Tests**: Find APIs without test coverage
5. **Dangling Commits**: Find commits not linked to sessions

---

## Naming Conventions

### ID Format

- Architecture: `ARCH-###` (3 digits)
- Package: `PKG-###` (3 digits)
- File: `FILE-####` (4 digits)
- API: `API-###` (3 digits)
- Test: `TEST-###` (3 digits)
- Specification: `SPEC-###` (3 digits)
- ADR: `ADR-####` (4 digits)
- Session: `SESSION-###` (3 digits)
- Commit: `COMMIT-<hash>` (short hash)
- Milestone: `MILE-<name>`
- Pattern: `PATTERN-###`
- Lesson: `LESSON-###`
- Rule: `RULE-###`

### Title Conventions

- Architecture: `<Layer>: <Name>` (e.g., "Layer 1: Runtime")
- Package: `<path>` (e.g., "internal/runtime/eventbus")
- File: `<filename>` (e.g., "eventbus.go")
- Specification: `<Component> <Type>` (e.g., "Event Bus Specification")

---

## AI Query Patterns

### Context Reconstruction

```
Given FILE-0014:
1. Find Package: FILE-0014 → PKG-014
2. Find Architecture: PKG-014 → ARCH-002
3. Find Specifications: ARCH-002 → SPEC-001, SPEC-002
4. Find ADR: SPEC-001 → ADR-0002
5. Find Session: Latest session touching FILE-0014
6. Find Commit: Session → COMMIT-xxx

Result: Complete lineage from file to decision
```

### Impact Analysis

```
Given ADR-0002:
1. Find affected Specifications: ADR-0002 ← SPEC-001, SPEC-010
2. Find affected Architectures: SPEC-001 → ARCH-002
3. Find affected Packages: ARCH-002 → PKG-014, PKG-015
4. Find affected Files: PKG-014 → FILE-0014, FILE-0015

Result: All files affected by an architectural decision
```

### Next Task Determination

```
Given HANDOFF-001:
1. Read Completed: What was done
2. Read Blocked: What's blocking
3. Read Next Task: Explicit next step
4. Query Architecture: ARCH-004 status
5. Query Packages: PKG-026..028 status
6. Query Specs: SPEC-020..022 frozen status

Result: Validated next action with full context
```

---

## Multi-Agent Collaboration

### Agent Responsibilities

| Agent | Databases | Writes | Reads |
|-------|-----------|--------|-------|
| Implementation Agent | Package, File, API | PKG-xxx, FILE-xxx, API-xxx | ARCH-xxx, SPEC-xxx |
| Architecture Agent | Architecture, Spec, ADR | ARCH-xxx, SPEC-xxx, ADR-xxxx | All |
| Knowledge Agent | Memory, Pattern, Lesson | Pattern-xxx, Lesson-xxx | All |
| Sync Agent | All | Commit, Session, Handoff | Git |
| Review Agent | All | Validation status | All |

### Handoff Protocol

1. **Before session**: Query HANDOFF-xxx for context
2. **During session**: Update current Session, Files
3. **After commit**: Create Commit entry, link to Session
4. **After session**: Create new Handoff with explicit Next Task
5. **Validation**: Ensure all constraints satisfied

---

## Versioning

This ontology is versioned as: **v1.0.0**

Breaking changes require ADR approval.
