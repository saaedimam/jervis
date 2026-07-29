# SCHEMA MIGRATION PLAN
## JERVIS Notion Workspace Enhancement
### Generated: 2026-07-29T06:25:00Z

---

## EXECUTIVE SUMMARY

**Approach**: Manual schema enhancement via Notion UI
**Rationale**: Preserves existing database IDs and relations
**Timeline**: 4-6 hours for complete enhancement
**Risk**: Low - Databases exist, only adding properties

---

## MIGRATION STRATEGY

### Phase 1: Foundation Properties (Priority: CRITICAL)
Add essential properties first. These enable basic functionality.

### Phase 2: Relation Properties (Priority: HIGH)
Add relation properties to enable graph traversal.

### Phase 3: Metadata Properties (Priority: MEDIUM)
Add richness properties for complete documentation.

### Phase 4: Validation (Priority: CRITICAL)
Verify all properties exist and are correctly typed.

---

## DETAILED MIGRATION MAP

### DATABASE 1: Architecture Registry

**ID**: d3dcb133-f96e-4e8e-944f-5825c2d1eee0

| Step | Current | New | Type | Required | Default | Notes |
|------|---------|-----|------|----------|---------|-------|
| 1 | Name (title) | Architecture ID | title | ✅ | - | Rename existing or add new |
| 2 | - | Name | rich_text | ✅ | - | Human-readable |
| 3 | - | Layer | select | ✅ | Layer 1 | Options: Layer 1-5 |
| 4 | - | Status | select | ✅ | In Progress | Options: In Progress, Complete, Deprecated |
| 5 | - | Version | rich_text | ✅ | 1.0.0 | Semantic version |
| 6 | - | Purpose | rich_text | ✅ | - | One-line description |
| 7 | - | Responsibilities | rich_text | ✅ | - | Bullet list |
| 8 | - | Coverage | number | ✅ | 0 | 0-100 |
| 9 | - | Risk | select | ✅ | Low | Options: Low, Medium, High |
| 10 | - | Frozen | checkbox | ✅ | ☐ | Boolean |
| 11 | - | Interfaces | rich_text | ❌ | - | Public APIs |
| 12 | - | Dependencies | relation | ❌ | - | → Architecture self-ref |
| 13 | - | Dependents | relation | ❌ | - | → Architecture self-ref |
| 14 | - | Related Packages | relation | ❌ | - | → Package Registry |
| 15 | - | Related Specifications | relation | ❌ | - | → Specification Registry |
| 16 | - | Related ADRs | relation | ❌ | - | → ADR Database |
| 17 | - | Related Sessions | relation | ❌ | - | → Session Database |
| 18 | - | Owner | people | ❌ | - | Assigned engineer |
| 19 | - | Last Updated | date | ✅ | now() | Auto timestamp |

**Migration Order**: 1 → 2-5 → 6-10 → 11-19
**Estimated Time**: 30 minutes
**Verification Query**: Check all 19 properties exist

---

### DATABASE 2: Package Registry

**ID**: 9c8bb7d5-5675-4cc5-b1b7-6a9c1ac3fe2f

| Step | Current | New | Type | Required | Default | Notes |
|------|---------|-----|------|----------|---------|-------|
| 1 | Name (title) | Package ID | title | ✅ | - | PKG-001 format |
| 2 | - | Name | rich_text | ✅ | - | Package name |
| 3 | - | Module | rich_text | ✅ | - | Import path |
| 4 | - | Purpose | rich_text | ✅ | - | One-line description |
| 5 | - | Status | select | ✅ | Active | Options: Active, Deprecated, Planning |
| 6 | - | Frozen | checkbox | ✅ | ☐ | Boolean |
| 7 | - | Coverage | number | ✅ | 0 | Percentage |
| 8 | - | Architecture | relation | ✅ | - | → Architecture Registry |
| 9 | - | Complexity | number | ❌ | 1 | Cyclomatic |
| 10 | - | Exported APIs | number | ❌ | 0 | Count |
| 11 | - | Tests | number | ❌ | 0 | Count |
| 12 | - | Files | number | ❌ | 0 | Count |
| 13 | - | Source Path | url | ❌ | - | GitHub link |
| 14 | - | Dependencies | rich_text | ❌ | - | PKG IDs |
| 15 | - | Dependents | rich_text | ❌ | - | PKG IDs |
| 16 | - | Specifications | relation | ❌ | - | → Specification Registry |
| 17 | - | Related Sessions | relation | ❌ | - | → Session Database |
| 18 | - | Owner | people | ❌ | - | Maintainer |
| 19 | - | Last Updated | date | ✅ | now() | Auto timestamp |

**Migration Order**: 1 → 2-8 → 9-13 → 14-19
**Estimated Time**: 30 minutes

---

### DATABASE 3: Specification Registry

**ID**: f30e0d51-a787-421a-ad6b-77935f7d2e53

| Step | Current | New | Type | Required | Default | Notes |
|------|---------|-----|------|----------|---------|-------|
| 1 | Name (title) | Spec ID | title | ✅ | - | SPEC-001 format |
| 2 | - | Name | rich_text | ✅ | - | Spec title |
| 3 | - | Version | rich_text | ✅ | 1.0.0 | Semver |
| 4 | - | Status | select | ✅ | Draft | Options: Draft, Frozen, Superseded |
| 5 | - | Frozen | checkbox | ✅ | ☐ | Boolean |
| 6 | - | Architecture | relation | ✅ | - | → Architecture |
| 7 | - | Complexity | select | ❌ | Moderate | Options: Simple, Moderate, Complex |
| 8 | - | Acceptance Criteria | rich_text | ✅ | - | Bullet list |
| 9 | - | Markdown File | url | ❌ | - | Doc link |
| 10 | - | Implementation Notes | rich_text | ❌ | - | Technical details |
| 11 | - | Change History | rich_text | ❌ | - | Versions |
| 12 | - | Packages | relation | ❌ | - | → Package Registry |
| 13 | - | Related ADR | relation | ❌ | - | → ADR Database |
| 14 | - | Superseded By | relation | ❌ | - | → Specification self |
| 15 | - | Related Sessions | relation | ❌ | - | → Session Database |
| 16 | - | Last Updated | date | ✅ | now() | Auto timestamp |

**Migration Order**: 1 → 2-8 → 9-16
**Estimated Time**: 25 minutes

---

### DATABASE 4: File Registry

**ID**: d5b8d71a-c568-4288-9443-f3deb8b316bc

| Step | Current | New | Type | Required | Default | Notes |
|------|---------|-----|------|----------|---------|-------|
| 1 | Name (title) | File ID | title | ✅ | - | FILE-0001 format |
| 2 | - | Path | rich_text | ✅ | - | Relative path |
| 3 | - | Language | select | ✅ | Go | Options: Go, Markdown, YAML, JSON, Shell |
| 4 | - | Status | select | ✅ | Active | Options: Active, Deprecated, Planning |
| 5 | - | Frozen | checkbox | ✅ | ☐ | Boolean |
| 6 | - | Package | relation | ✅ | - | → Package Registry |
| 7 | - | Architecture | relation | ✅ | - | → Architecture Registry |
| 8 | - | Lines | number | ❌ | 0 | LOC |
| 9 | - | API Count | number | ❌ | 0 | Exported APIs |
| 10 | - | Coverage | rich_text | ❌ | N/A | Percentage or N/A |
| 11 | - | Exports | rich_text | ❌ | - | Symbols |
| 12 | - | Imports | rich_text | ❌ | - | Packages |
| 13 | - | Hash | rich_text | ✅ | - | MD5 for sync |
| 14 | - | Specification | relation | ❌ | - | → Specification |
| 15 | - | Last Commit | relation | ❌ | - | → Commit |
| 16 | - | Last Session | relation | ❌ | - | → Session |
| 17 | - | Owner | people | ❌ | - | Maintainer |
| 18 | - | Created | date | ❌ | - | File creation |
| 19 | - | Updated | date | ✅ | now() | Last modification |

**Migration Order**: 1 → 2-7 → 8-13 → 14-19
**Estimated Time**: 35 minutes

---

### DATABASE 5: API Registry

**ID**: 5e2dad61-5186-46f7-be6b-e7e5c3715f04

| Step | Current | New | Type | Required | Default | Notes |
|------|---------|-----|------|----------|---------|-------|
| 1 | Name (title) | API ID | title | ✅ | - | API-001 format |
| 2 | - | Name | rich_text | ✅ | - | Function name |
| 3 | - | Version | rich_text | ✅ | 1.0.0 | Semver |
| 4 | - | Status | select | ✅ | Stable | Options: Stable, Experimental, Deprecated |
| 5 | - | Breaking | checkbox | ✅ | ☐ | Breaking change |
| 6 | - | Package | relation | ✅ | - | → Package Registry |
| 7 | - | File | relation | ✅ | - | → File Registry |
| 8 | - | Signature | rich_text | ✅ | - | Full signature |
| 9 | - | Inputs | rich_text | ✅ | - | Parameters |
| 10 | - | Outputs | rich_text | ✅ | - | Returns |
| 11 | - | Errors | rich_text | ❌ | - | Error types |
| 12 | - | Coverage | number | ❌ | 0 | Test coverage |
| 13 | - | Specification | relation | ❌ | - | → Specification |
| 14 | - | ADR | relation | ❌ | - | → ADR |
| 15 | - | Tests | relation | ❌ | - | → Tests |
| 16 | - | Related Sessions | relation | ❌ | - | → Session |
| 17 | - | Last Updated | date | ✅ | now() | Auto timestamp |

**Migration Order**: 1 → 2-7 → 8-12 → 13-17
**Estimated Time**: 30 minutes

---

### DATABASE 6: ADR Database

**ID**: abc5d892-1299-4813-b8bf-a143d6c8c73c

| Step | Current | New | Type | Required | Default | Notes |
|------|---------|-----|------|----------|---------|-------|
| 1 | Name (title) | ADR ID | title | ✅ | - | ADR-0001 format |
| 2 | - | Title | rich_text | ✅ | - | Decision title |
| 3 | - | Status | select | ✅ | Proposed | Options: Proposed, Accepted, Superseded, Rejected |
| 4 | - | Date | date | ✅ | today() | Decision date |
| 5 | - | Author | people | ✅ | - | Decision author |
| 6 | - | Context | rich_text | ✅ | - | Background |
| 7 | - | Problem | rich_text | ✅ | - | What needed solving |
| 8 | - | Decision | rich_text | ✅ | - | What was decided |
| 9 | - | Alternatives | rich_text | ❌ | - | Options considered |
| 10 | - | Consequences | rich_text | ❌ | - | Positive impacts |
| 11 | - | Trade-offs | rich_text | ❌ | - | Negative impacts |
| 12 | - | Affected Components | relation | ❌ | - | → Architecture |
| 13 | - | Affected Packages | relation | ❌ | - | → Package |
| 14 | - | Related Specifications | relation | ❌ | - | → Specification |
| 15 | - | Superseded By | relation | ❌ | - | → ADR self |
| 16 | - | Review History | rich_text | ❌ | - | Review log |
| 17 | - | Timeline | rich_text | ❌ | - | Key dates |
| 18 | - | Related Sessions | relation | ❌ | - | → Session |
| 19 | - | Last Updated | date | ✅ | now() | Auto timestamp |

**Migration Order**: 1 → 2-8 → 9-11 → 12-19
**Estimated Time**: 35 minutes

---

### DATABASE 7: Milestones Database

**ID**: 39ae6e23-2bc1-4e34-a7b0-a1da9410b081

| Step | Current | New | Type | Required | Default | Notes |
|------|---------|-----|------|----------|---------|-------|
| 1 | Name (title) | Milestone ID | title | ✅ | - | PHASE-1.1 format |
| 2 | - | Name | rich_text | ✅ | - | Milestone name |
| 3 | - | Phase | select | ✅ | Phase 1.1 | Options: Phase 1.1, 1.2, ..., 5 |
| 4 | - | Status | select | ✅ | Pending | Options: Done, In Progress, Pending, Blocked |
| 5 | - | Coverage | rich_text | ❌ | TBD | Test coverage |
| 6 | - | Start Date | date | ❌ | - | When began |
| 7 | - | Target Date | date | ❌ | - | Planned completion |
| 8 | - | Completion Date | date | ❌ | - | Actual completion |
| 9 | - | Deliverables | rich_text | ✅ | - | What delivered |
| 10 | - | Risks | rich_text | ❌ | - | Known risks |
| 11 | - | Blockers | rich_text | ❌ | - | Current blockers |
| 12 | - | Dependencies | relation | ❌ | - | → Milestone self |
| 13 | - | Dependents | relation | ❌ | - | → Milestone self |
| 14 | - | Related Sessions | relation | ❌ | - | → Session |
| 15 | - | Related Commits | relation | ❌ | - | → Commit |
| 16 | - | Last Updated | date | ✅ | now() | Auto timestamp |

**Migration Order**: 1 → 2-9 → 10-16
**Estimated Time**: 25 minutes

---

### DATABASE 8: Commit Intelligence

**ID**: 69c5145a-b84c-43e5-83b2-05d746a80e26

| Step | Current | New | Type | Required | Default | Notes |
|------|---------|-----|------|----------|---------|-------|
| 1 | Name (title) | Commit ID | title | ✅ | - | Short SHA |
| 2 | - | SHA | rich_text | ✅ | - | Full SHA |
| 3 | - | Message | rich_text | ✅ | - | Commit message |
| 4 | - | Author | people | ✅ | - | Commit author |
| 5 | - | Branch | rich_text | ✅ | main | Git branch |
| 6 | - | Date | date | ✅ | now() | Commit timestamp |
| 7 | - | Changed Files | relation | ✅ | - | → File Registry |
| 8 | - | Changed Packages | relation | ✅ | - | → Package Registry |
| 9 | - | Architecture Impact | relation | ❌ | - | → Architecture |
| 10 | - | Specification Impact | relation | ❌ | - | → Specification |
| 11 | - | Related Session | relation | ✅ | - | → Session |
| 12 | - | Related Milestone | relation | ❌ | - | → Milestone |
| 13 | - | Related Release | relation | ❌ | - | → Release |
| 14 | - | Review Status | select | ❌ | Pending | Options: Approved, Pending, Rejected |
| 15 | - | Last Updated | date | ✅ | now() | Auto timestamp |

**Migration Order**: 1 → 2-8 → 9-15
**Estimated Time**: 25 minutes

---

### DATABASE 9: AI Handoff (Session Database)

**ID**: c1e36ebb-a3fc-4aea-a3d2-ac8214e1e40a

| Step | Current | New | Type | Required | Default | Notes |
|------|---------|-----|------|----------|---------|-------|
| 1 | Name (title) | Session ID | title | ✅ | - | SESSION-021 format |
| 2 | - | Date | date | ✅ | today() | Session date |
| 3 | - | Objective | rich_text | ✅ | - | Session goal |
| 4 | - | Summary | rich_text | ✅ | - | Accomplished |
| 5 | - | Work Completed | rich_text | ✅ | - | Detailed list |
| 6 | - | Status | select | ✅ | In Progress | Options: Complete, In Progress, Failed |
| 7 | - | Files Modified | relation | ❌ | - | → File |
| 8 | - | Commits | relation | ✅ | - | → Commit |
| 9 | - | Architecture Changes | relation | ❌ | - | → Architecture |
| 10 | - | Problems | rich_text | ❌ | - | Issues encountered |
| 11 | - | Solutions | rich_text | ❌ | - | Resolutions |
| 12 | - | AI Discussion | rich_text | ❌ | - | Key interactions |
| 13 | - | Remaining Work | rich_text | ❌ | - | What's left |
| 14 | - | Next Session | rich_text | ❌ | - | Next steps |
| 15 | - | Related Components | relation | ❌ | - | Multi-select |
| 16 | - | Related Specifications | relation | ❌ | - | → Specification |
| 17 | - | Lessons Learned | relation | ❌ | - | → Memory |
| 18 | - | Last Updated | date | ✅ | now() | Auto timestamp |

**Migration Order**: 1 → 2-6 → 7-9 → 10-18
**Estimated Time**: 30 minutes

---

### DATABASE 10: Engineering Memory

**ID**: 38a76b5b-b20e-498e-b6e9-e643c2ae7d8b

| Step | Current | New | Type | Required | Default | Notes |
|------|---------|-----|------|----------|---------|-------|
| 1 | Name (title) | Memory ID | title | ✅ | - | MEM-001 or PAT-001 |
| 2 | - | Type | select | ✅ | Pattern | Options: Pattern, Lesson, Anti-Pattern, Research, Rule |
| 3 | - | Name | rich_text | ✅ | - | Memory title |
| 4 | - | Category | select | ❌ | Architecture | Options: Architecture, Testing, Performance, Security |
| 5 | - | Summary | rich_text | ✅ | - | One-liner |
| 6 | - | Context | rich_text | ✅ | - | When applies |
| 7 | - | Content | rich_text | ✅ | - | Full explanation |
| 8 | - | Code Examples | rich_text | ❌ | - | Snippets |
| 9 | - | Benefits | rich_text | ❌ | - | Why use |
| 10 | - | Risks | rich_text | ❌ | - | Downsides |
| 11 | - | Related Components | relation | ❌ | - | → Architecture |
| 12 | - | Related Packages | relation | ❌ | - | → Package |
| 13 | - | Related Sessions | relation | ❌ | - | → Session |
| 14 | - | References | rich_text | ❌ | - | Links |
| 15 | - | Tags | multi_select | ❌ | - | Searchable |
| 16 | - | Last Updated | date | ✅ | now() | Auto timestamp |

**Migration Order**: 1 → 2-7 → 8-10 → 11-16
**Estimated Time**: 25 minutes

---

### DATABASE 11: Dependency Graph

**ID**: 1de04b92-6fe3-4756-b85d-c9370f838a3b

| Step | Current | New | Type | Required | Default | Notes |
|------|---------|-----|------|----------|---------|-------|
| 1 | Name (title) | Relation ID | title | ✅ | - | DEP-001 |
| 2 | - | Source | relation | ✅ | - | → Package |
| 3 | - | Target | relation | ✅ | - | → Package |
| 4 | - | Type | select | ✅ | Imports | Options: Imports, Uses, Implements, DependsOn |
| 5 | - | Strength | select | ✅ | Strong | Options: Strong, Weak, Optional |
| 6 | - | Direction | select | ✅ | Uni | Options: Uni, Bi |
| 7 | - | Created | date | ❌ | - | When established |
| 8 | - | Verified | checkbox | ❌ | ☐ | CI verified |
| 9 | - | Notes | rich_text | ❌ | - | Context |
| 10 | - | Last Updated | date | ✅ | now() | Auto timestamp |

**Migration Order**: 1 → 2-6 → 7-10
**Estimated Time**: 15 minutes

---

## EXECUTION PLAN

### Day 1: Foundation (2 hours)

**Hour 1: Core Databases**
- [ ] Architecture Registry (30 min)
- [ ] Package Registry (30 min)

**Hour 2: Supporting Databases**
- [ ] Specification Registry (25 min)
- [ ] File Registry (35 min)

### Day 2: Relations (2 hours)

**Hour 1: API & ADR**
- [ ] API Registry (30 min)
- [ ] ADR Database (30 min)

**Hour 2: Tracking Databases**
- [ ] Milestones (25 min)
- [ ] Commit Intelligence (25 min)
- [ ] Quick validation (10 min)

### Day 3: Memory (2 hours)

**Hour 1: Session & Memory**
- [ ] AI Handoff (30 min)
- [ ] Engineering Memory (25 min)
- [ ] Dependency Graph (15 min)

**Hour 2: Validation & Cleanup**
- [ ] Verify all properties exist
- [ ] Run population scripts
- [ ] Test relations
- [ ] Generate validation report

---

## ROLLBACK PLAN

If migration fails:

1. **Export all data** before starting (Notion CSV export)
2. **Document current state** (screenshots)
3. **Add properties incrementally** (one database at a time)
4. **Test each database** before moving to next
5. **Keep old data** until verified

---

## VERIFICATION CHECKLIST

After each database:

- [ ] All properties visible in Notion UI
- [ ] Property types match specification
- [ ] Select options correctly defined
- [ ] Relations point to correct databases
- [ ] Required properties enforced
- [ ] Default values set
- [ ] Can create new entry with all properties

---

## POST-MIGRATION TASKS

1. **Populate canonical IDs** for all existing entries
2. **Create relations** between existing pages
3. **Populate metadata** fields
4. **Test bidirectional relations**
5. **Generate validation report**

---

**Migration Plan Complete**: All 11 databases mapped
**Total Properties to Add**: ~200
**Estimated Duration**: 6 hours
**Risk Level**: Low
**Recommended**: Execute over 3 days
