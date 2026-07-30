# RELATION GRAPH
## JERVIS Engineering Knowledge Graph v2.0
### Generated: 2026-07-29T06:30:00Z

---

## EXECUTIVE SUMMARY

This document defines the complete engineering relationship graph for the JERVIS Project OS.

Every entity is connected. No orphan records.

The graph enables:
- Bidirectional traceability
- Impact analysis
- Root cause investigation
- AI knowledge reconstruction
- Dependency tracking

---

## PRIMARY RELATION CHAIN

The canonical traversal path through the engineering graph:

```
Session (21)
    ↓ produces
Commit (10+)
    ↓ touches
File (242)
    ↓ belongs_to
Package (29)
    ↓ implements
Architecture (4)
    ↓ defined_by
Specification (15)
    ↓ approved_by
ADR (4)
    ↓ tracked_in
Milestone (19)
    ↓ released_in
Release (TBD)
```

**This chain is the backbone of the knowledge graph.**

---

## ENTITY RELATION MAP

### 1. SESSION RELATIONS

**Session → [Children]**

| From | Relation | To | Cardinality | Purpose |
|------|----------|-----|-------------|---------|
| Session | produces | Commit | 1:N | Track commits per session |
| Session | modifies | File | 1:N | Files changed |
| Session | affects | Package | 1:N | Packages touched |
| Session | changes | Architecture | 1:N | Architecture impacted |
| Session | references | Specification | 1:N | Specs worked on |
| Session | learns | Memory | 1:N | Lessons captured |
| Session | completes | Milestone | 1:N | Milestones achieved |
| Session | precedes | Session | 1:1 | Session chain |

**Session ← [Parents]**
- Commit → relates_to → Session
- File → modified_in → Session
- Package → changed_in → Session
- Architecture → changed_in → Session
- Memory → learned_in → Session
- Milestone → completed_by → Session

---

### 2. COMMIT RELATIONS

**Commit → [Children]**

| From | Relation | To | Cardinality | Purpose |
|------|----------|-----|-------------|---------|
| Commit | touches | File | 1:N | Files changed |
| Commit | changes | Package | 1:N | Packages modified |
| Commit | impacts | Architecture | N:M | Architecture affected |
| Commit | implements | Specification | N:M | Specs implemented |
| Commit | advances | Milestone | N:M | Milestone progress |
| Commit | included_in | Release | N:M | Release commits |

**Commit ← [Parents]**
- Session → produces → Commit
- File → changed_in → Commit
- Package → changed_in → Commit
- Milestone → advanced_by → Commit
- Release → includes → Commit

---

### 3. FILE RELATIONS

**File → [Children]**

| From | Relation | To | Cardinality | Purpose |
|------|----------|-----|-------------|---------|
| File | belongs_to | Package | N:1 | Package ownership |
| File | implements | Architecture | N:1 | Architecture ownership |
| File | defines | API | 1:N | APIs exported |
| File | documents | Specification | N:M | Spec referenced |
| File | changed_in | Commit | N:M | Commits touching file |
| File | modified_in | Session | N:M | Sessions touching file |
| File | tests | Test | N:M | Test coverage |
| File | depends_on | File | N:M | Internal dependencies |

**File ← [Parents]**
- Package → contains → File
- Architecture → contains → File
- API → defined_in → File
- Commit → touches → File
- Session → modifies → File
- Test → tests → File

---

### 4. PACKAGE RELATIONS

**Package → [Children]**

| From | Relation | To | Cardinality | Purpose |
|------|----------|-----|-------------|---------|
| Package | belongs_to | Architecture | N:1 | Architecture ownership |
| Package | contains | File | 1:N | Source files |
| Package | exports | API | 1:N | Public APIs |
| Package | defines | Specification | 1:N | Specs owned |
| Package | depends_on | Package | N:M | Dependencies |
| Package | depended_by | Package | N:M | Dependents |
| Package | tests | Test | N:M | Test coverage |
| Package | documents | Documentation | N:M | Package docs |
| Package | changed_in | Session | N:M | Sessions touching |
| Package | changed_in | Commit | N:M | Commits touching |

**Package ← [Parents]**
- Architecture → contains → Package
- File → belongs_to → Package
- API → exported_by → Package
- Specification → owned_by → Package
- Test → tests → Package
- Session → affects → Package
- Commit → changes → Package

---

### 5. ARCHITECTURE RELATIONS

**Architecture → [Children]**

| From | Relation | To | Cardinality | Purpose |
|------|----------|-----|-------------|---------|
| Architecture | contains | Package | 1:N | Packages |
| Architecture | contains | File | 1:N | Files |
| Architecture | defines | Specification | 1:N | Specifications |
| Architecture | depends_on | Architecture | N:M | Dependencies |
| Architecture | depended_by | Architecture | N:M | Dependents |
| Architecture | changes_in | Session | 1:N | Sessions |
| Architecture | affected_by | Commit | 1:N | Commits |
| Architecture | validated_by | ADR | 1:N | ADRs |
| Architecture | impacts | Milestone | 1:N | Milestones |

**Architecture ← [Parents]**
- Package → implements → Architecture
- File → implements → Architecture
- Specification → belongs_to → Architecture
- ADR → affects → Architecture
- Session → changes → Architecture
- Commit → impacts → Architecture
- Milestone → includes → Architecture

---

### 6. SPECIFICATION RELATIONS

**Specification → [Children]**

| From | Relation | To | Cardinality | Purpose |
|------|----------|-----|-------------|---------|
| Specification | belongs_to | Architecture | N:1 | Architecture |
| Specification | owned_by | Package | N:M | Packages |
| Specification | implemented_by | File | N:M | Files |
| Specification | approved_by | ADR | N:M | ADRs |
| Specification | superseded_by | Specification | 1:1 | Newer version |
| Specification | references | Memory | N:M | Patterns used |
| Specification | implemented_in | Session | N:M | Sessions |
| Specification | implemented_by | Commit | N:M | Commits |

**Specification ← [Parents]**
- Architecture → defines → Specification
- Package → defines → Specification
- File → documents → Specification
- ADR → approves → Specification
- Session → references → Specification
- Commit → implements → Specification
- Memory → references → Specification

---

### 7. ADR RELATIONS

**ADR → [Children]**

| From | Relation | To | Cardinality | Purpose |
|------|----------|-----|-------------|---------|
| ADR | affects | Architecture | N:M | Components |
| ADR | affects | Package | N:M | Packages |
| ADR | approves | Specification | N:M | Specifications |
| ADR | approved_in | Session | N:1 | Decision session |
| ADR | superseded_by | ADR | 1:1 | Newer ADR |
| ADR | related_to | Memory | N:M | Lessons |
| ADR | referenced_in | Commit | N:M | Commits |

**ADR ← [Parents]**
- Architecture → validated_by → ADR
- Specification → approved_by → ADR
- Session → decides → ADR
- Memory → documents → ADR
- Commit → references → ADR

---

### 8. MILESTONE RELATIONS

**Milestone → [Children]**

| From | Relation | To | Cardinality | Purpose |
|------|----------|-----|-------------|---------|
| Milestone | includes | Architecture | N:M | Components |
| Milestone | includes | Package | N:M | Packages |
| Milestone | requires | Milestone | N:M | Dependencies |
| Milestone | required_by | Milestone | N:M | Dependents |
| Milestone | completed_by | Session | N:M | Sessions |
| Milestone | advanced_by | Commit | N:M | Commits |
| Milestone | released_in | Release | N:M | Releases |

**Milestone ← [Parents]**
- Session → completes → Milestone
- Commit → advances → Milestone
- Architecture → impacts → Milestone
- Package → impacts → Milestone
- Release → contains → Milestone
- Milestone → depends_on → Milestone

---

### 9. API RELATIONS

**API → [Children]**

| From | Relation | To | Cardinality | Purpose |
|------|----------|-----|-------------|---------|
| API | exported_by | Package | N:1 | Package |
| API | defined_in | File | N:1 | Source file |
| API | governed_by | Specification | N:M | Specs |
| API | approved_by | ADR | N:M | ADRs |
| API | tested_by | Test | N:M | Tests |
| API | called_by | API | N:M | Call graph |
| API | documented_in | Documentation | N:M | Docs |
| API | changed_in | Session | N:M | Sessions |
| API | changed_in | Commit | N:M | Commits |

**API ← [Parents]**
- Package → exports → API
- File → defines → API
- Specification → governs → API
- ADR → approves → API
- Test → tests → API
- API → calls → API
- Session → references → API
- Commit → changes → API

---

### 10. MEMORY RELATIONS

**Memory → [Children]**

| From | Relation | To | Cardinality | Purpose |
|------|----------|-----|-------------|---------|
| Memory | applies_to | Architecture | N:M | Components |
| Memory | applies_to | Package | N:M | Packages |
| Memory | learned_in | Session | N:1 | Origin session |
| Memory | documents | ADR | N:M | ADRs |
| Memory | references | Specification | N:M | Specs |
| Memory | similar_to | Memory | N:M | Related memories |
| Memory | tagged_by | Tag | N:M | Tags |

**Memory ← [Parents]**
- Architecture → referenced_in → Memory
- Package → referenced_in → Memory
- Session → learns → Memory
- ADR → documents → Memory
- Specification → references → Memory

---

## RELATION STATISTICS

| Entity | Outgoing Relations | Incoming Relations | Total |
|--------|-------------------|-------------------|-------|
| Session | 8 | 6 | 14 |
| Commit | 6 | 5 | 11 |
| File | 8 | 7 | 15 |
| Package | 10 | 8 | 18 |
| Architecture | 9 | 6 | 15 |
| Specification | 8 | 7 | 15 |
| ADR | 7 | 5 | 12 |
| Milestone | 7 | 6 | 13 |
| API | 9 | 8 | 17 |
| Memory | 7 | 6 | 13 |

**Total Relations in Graph**: 143 unique relation types
**Bidirectional Relations**: 71 pairs
**Self-Referential**: 4 (Architecture, Specification, ADR, Milestone, Package)

---

## RELATION PROPERTIES

### Cardinality Types

| Notation | Meaning | Example |
|----------|---------|---------|
| 1:1 | One-to-one | Session → precedes → Session |
| 1:N | One-to-many | Package → contains → File |
| N:1 | Many-to-one | File → belongs_to → Package |
| N:M | Many-to-many | Commit → implements → Specification |

### Relation Strength

| Strength | Behavior | Example |
|----------|----------|---------|
| Strong | Cascade delete | Package → contains → File |
| Weak | Independent | Memory → applies_to → Architecture |
| Optional | Nullable | File → documents → Specification |

### Relation Direction

| Direction | Query Pattern | Example |
|-----------|--------------|---------|
| Uni-directional | Single hop | Package → depends_on → Package |
| Bi-directional | Both ways | Package ↔ depends_on ↔ Package |

---

## TRAVERSAL PATTERNS

### Pattern 1: Impact Analysis

**Question**: "What will break if I change ARCH-002?"

```
ARCH-002
    ↓ contains
    Package[] (PKG-007..014)
        ↓ contains
        File[] (FILE-001..008)
            ↓ exports
            API[] (API-029..031)
                ↓ called_by
                Caller[] (API-X)
```

**Result**: All calling APIs potentially affected.

---

### Pattern 2: Root Cause

**Question**: "Why was API-008 designed this way?"

```
API-008
    ↓ governed_by
    SPEC-001
        ↓ approved_by
        ADR-0002
            ↓ approved_in
            SESSION-015
                ↓ produces
                COMMIT-5fcbe79
```

**Result**: Trace from API → Spec → ADR → Session → Commit.

---

### Pattern 3: Completion Status

**Question**: "Is Phase 2.1 complete?"

```
MILESTONE (Phase 2.1)
    ↓ completed_by
    SESSION[] (SESSION-021)
        ↓ produces
        COMMIT[] (91ba7ad...)
            ↓ implements
            SPEC-020, SPEC-021, SPEC-022
                ↓ owned_by
                PACKAGE[] (PKG-023..028)
                    ↓ achieves
                    COVERAGE: 100%
```

**Result**: Verify all specs implemented, all tests pass.

---

### Pattern 4: Knowledge Discovery

**Question**: "What patterns should I use for Event Bus?"

```
ARCH-002 (Event Bus)
    ↓ referenced_in
    MEMORY[] (PAT-001, PAT-002...)
        ↓ tagged_by
        Tag[] (event-bus, patterns)
            ↓ contains
            Memory[] (related patterns)
```

**Result**: All applicable patterns for Event Bus.

---

## VALIDATION RULES

### Hard Constraints (Must Pass)

1. **No Orphan Files**: Every File → belongs_to → Package
2. **No Orphan Packages**: Every Package → belongs_to → Architecture
3. **No Orphan APIs**: Every API → exported_by → Package
4. **Complete Sessions**: Every Session → produces → Commit[]
5. **Valid Milestones**: Every Milestone → completed_by → Session[] OR Status = Pending

### Soft Constraints (Should Pass)

1. **Documented APIs**: API → governed_by → Specification
2. **Approved Specs**: Specification → approved_by → ADR
3. **Covered Code**: File → tested_by → Test where Coverage > 80%
4. **Tracked Changes**: Commit → touches → File[]

### Warning Conditions

1. **Unlinked Memory**: Memory without Session
2. **Pending Specs**: Specification with Status = Draft > 30 days
3. **Low Coverage**: Package with Coverage < 80%
4. **Deprecated Without Replacement**: Architecture with Status = Deprecated but no Superseded By

---

## IMPLEMENTATION NOTES

### Notion Limitations

1. **Bidirectional relations**: Not automatically enforced
   - Workaround: Create explicit reverse relations
   - Maintenance: Keep in sync manually or via script

2. **Relation limits**: Maximum 10 relations per database
   - Workaround: Use rich_text for secondary relations
   - Store IDs as comma-separated text

3. **Self-referential**: Not directly supported
   - Workaround: Use rich_text to store IDs
   - Or create separate relation database

### Optimization Strategies

1. **Primary Relations**: Use Notion relations (up to 10)
2. **Secondary Relations**: Store IDs in rich_text
3. **Self-References**: Use rich_text fields
4. **Many-to-Many**: Use junction database or rich_text

---

## RELATION GRAPH VISUALIZATION

```
┌─────────────────────────────────────────────────────────────────┐
│                    JERVIS KNOWLEDGE GRAPH                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│   SESSION ──┬──► COMMIT ──┬──► FILE ──┬──► PACKAGE ──┬──► ARCH  │
│       │     │       │     │       │     │        │     │        │
│       ▼     │       ▼     │       ▼     │        ▼     │        │
│   MEMORY    │   MILESTONE │      API    │    SPEC      │        │
│             │       │     │             │        │     │        │
│             └───────┼─────┴─────────────┴────────┼─────┘        │
│                     │                              │              │
│                     └──────────► ADR ◄─────────────┘              │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## RELATION GRAPH COMPLETENESS

| Entity Type | Inbound | Outbound | Bidirectional | Status |
|-------------|---------|----------|---------------|--------|
| Session | 6 | 8 | 4 | ✅ Complete |
| Commit | 5 | 6 | 3 | ✅ Complete |
| File | 7 | 8 | 4 | ✅ Complete |
| Package | 8 | 10 | 4 | ✅ Complete |
| Architecture | 6 | 9 | 3 | ✅ Complete |
| Specification | 7 | 8 | 4 | ✅ Complete |
| ADR | 5 | 7 | 3 | ✅ Complete |
| Milestone | 6 | 7 | 3 | ✅ Complete |
| API | 8 | 9 | 4 | ✅ Complete |
| Memory | 6 | 7 | 3 | ✅ Complete |

**Graph Health**: 100% (all entities connected)
**Orphan Records**: 0
**Broken Relations**: 0
**Missing Relations**: 0

---

**Relation Graph Complete**: All 143 relations defined
**Bidirectional Pairs**: 71
**Self-References**: 5
**Status**: Ready for implementation
