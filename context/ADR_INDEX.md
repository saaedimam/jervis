# ADR Index

---

### ADR-0001
- **Title**: Initial Architecture Baseline
- **Status**: Superseded (by ADR-0002)
- **Date**: 2026-07-28

---

### ADR-0002
- **Title**: Architecture Reconciliation - Runtime Ownership & AI Decoupling
- **Status**: Accepted
- **Date**: 2026-07-28
- **Summary**: Replaced AI-centric runtime with canonical 5-tier architecture (`OS -> Runtime -> Memory Engine -> Service Layer -> AI Provider Layer -> Interfaces`). Selected Go (Golang 1.22+) as canonical implementation language. Enforced 15 mandatory design rules.
