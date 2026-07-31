# Jervis Documentation Index

> Last Updated: 2026-07-31
> Owner: @saaedimam

This directory contains all project documentation. Root-level files are reserved for:
- README.md — project overview
- LICENSE — legal terms
- CONTRIBUTING.md — contribution guidelines
- CHANGELOG.md — release history

## Navigation

### Principles & Vision
| Document | Purpose |
|----------|---------|
| [principles/engineering.md](principles/engineering.md) | Core philosophy, design rules, coding standards |

### Architecture
| Document | Purpose |
|----------|---------|
| [architecture/ARCHITECTURE_INVARIANTS.md](architecture/ARCHITECTURE_INVARIANTS.md) | 15 immutable rules — **single source of truth** |
| [architecture/ARCHITECTURE.md](architecture/ARCHITECTURE.md) | Layer definitions, system flows, diagrams |
| [architecture/ARCHITECTURE_REVIEW_REPORT.md](architecture/ARCHITECTURE_REVIEW_REPORT.md) | Review outcomes |
| [architecture/ADR_GUIDE.md](architecture/ADR_GUIDE.md) | ADR format and process |
| [architecture/CI_CD_SPECIFICATION.md](architecture/CI_CD_SPECIFICATION.md) | CI/CD pipeline design |

### Specifications — Runtime
| Document | Purpose |
|----------|---------|
| [specs/runtime/BUS_SPECIFICATION.md](specs/runtime/BUS_SPECIFICATION.md) | Event bus contract |
| [specs/runtime/BUS_CONTRACTS.md](specs/runtime/BUS_CONTRACTS.md) | Bus message contracts |
| [specs/runtime/BUS_PIPELINE.md](specs/runtime/BUS_PIPELINE.md) | Bus middleware pipeline |
| [specs/runtime/BUS_SEQUENCE.md](specs/runtime/BUS_SEQUENCE.md) | Bus lifecycle sequence |
| [specs/runtime/EVENT_BUS_SPECIFICATION.md](specs/runtime/EVENT_BUS_SPECIFICATION.md) | Event bus spec |
| [specs/runtime/EVENT_CONTRACTS.md](specs/runtime/EVENT_CONTRACTS.md) | Event contracts |
| [specs/runtime/EVENT_MODEL.md](specs/runtime/EVENT_MODEL.md) | Event data model |
| [specs/runtime/EVENT_IMPLEMENTATION_PLAN.md](specs/runtime/EVENT_IMPLEMENTATION_PLAN.md) | Event bus implementation |
| [specs/runtime/MIDDLEWARE_SPECIFICATION.md](specs/runtime/MIDDLEWARE_SPECIFICATION.md) | Middleware contract |
| [specs/runtime/MIDDLEWARE_CONTRACTS.md](specs/runtime/MIDDLEWARE_CONTRACTS.md) | Middleware interface contracts |
| [specs/runtime/MIDDLEWARE_PIPELINE.md](specs/runtime/MIDDLEWARE_PIPELINE.md) | Middleware execution pipeline |
| [specs/runtime/MIDDLEWARE_ORDERING.md](specs/runtime/MIDDLEWARE_ORDERING.md) | Middleware ordering rules |
| [specs/runtime/DISPATCH_PIPELINE.md](specs/runtime/DISPATCH_PIPELINE.md) | Dispatcher pipeline |
| [specs/runtime/DISPATCH_SEQUENCE.md](specs/runtime/DISPATCH_SEQUENCE.md) | Dispatcher sequence |

### Specifications — Security
| Document | Purpose |
|----------|---------|
| [specs/security/PERMISSION_MODEL.md](specs/security/PERMISSION_MODEL.md) | Permission model |
| [specs/security/PERMISSION_CONTRACTS.md](specs/security/PERMISSION_CONTRACTS.md) | Permission contracts |
| [specs/security/PERMISSION_ENGINE_SPECIFICATION.md](specs/security/PERMISSION_ENGINE_SPECIFICATION.md) | Permission engine spec |
| [specs/security/PERMISSION_IMPLEMENTATION_PLAN.md](specs/security/PERMISSION_IMPLEMENTATION_PLAN.md) | Permission implementation |

### ADRs & Planning
| Document | Purpose |
|----------|---------|
| [adr/DECISIONS.md](adr/DECISIONS.md) | Key architectural decisions |
| [adr/ROADMAP.md](adr/ROADMAP.md) | Implementation roadmap |
| [adr/PROJECT_METRICS.md](adr/PROJECT_METRICS.md) | Project health metrics |
| [adr/PROJECT_STRUCTURE.md](adr/PROJECT_STRUCTURE.md) | Directory layout |
| [adr/QUALITY_GATES.md](adr/QUALITY_GATES.md) | Quality gate definitions |

### Releases
| Document | Purpose |
|----------|---------|
| [releases/v1.0.0_report.md](releases/v1.0.0_report.md) | Final release report |
| [releases/v1.0.0_checklist.md](releases/v1.0.0_checklist.md) | Release checklist |
