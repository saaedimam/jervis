# Architectural Decision Record (ADR) Governance Guide

## 1. Purpose & Scope
This document specifies the mandatory governance process for proposing, reviewing, approving, and deprecating Architectural Decision Records (ADRs) within Project Jervis.

---

## 2. When an ADR is Required
An ADR MUST be created whenever a proposed change:
- Alters component boundaries, layer relationships, or runtime ownership.
- Modifies any of the 15 Architectural Invariants.
- Introduces, replaces, or removes external dependencies.
- Changes persistent storage schemas or context formats.
- Modifies security, capability authorization, or plugin sandboxing models.

---

## 3. Mandatory ADR Structure
Every ADR MUST include the following standardized sections:
- **Title**: Sequential identifier and descriptive title (`ADR-XXXX: Title`).
- **Date**: Date of submission in `YYYY-MM-DD` format.
- **Status**: Current lifecycle state (`Proposed`, `Accepted`, `Rejected`, `Superseded`).
- **Context**: Problem statement and architectural background driving the proposal.
- **Decision**: Precise normative description of the change being adopted.
- **Consequences**: Explicit enumeration of positive outcomes, negative tradeoffs, and operational risks.
- **Compliance & Invariants**: Verification that the decision complies with all 15 Architectural Invariants.

---

## 4. ADR Approval Workflow
1. **Authoring**: The proposer creates a draft ADR document under `docs/adrs/ADR-XXXX.md`.
2. **Review Submission**: The proposer submits a pull request containing exclusively the draft ADR.
3. **Architectural Evaluation**: The Architecture Board reviews the proposal against system principles and invariants.
4. **Voting Threshold**: Approval SHALL require unanimous consent from all members of the Architecture Board.
5. **Merging & Status Update**: Upon approval, the status MUST be updated to `Accepted`, merged into main, and recorded in `DECISIONS.md`.

---

## 5. Deprecation & Superseding Workflow
- An accepted ADR SHALL NOT be edited retroactively except to update its status.
- If a new architectural decision replaces an existing ADR, a new ADR MUST be created.
- The new ADR MUST explicitly mark its status as `Accepted` and list the superseded ADR.
- The original ADR status MUST be updated to `Superseded by ADR-XXXX`.

---

## 6. Architecture Review Process
- The Architecture Board SHALL conduct quarterly architecture reviews to audit codebase compliance with accepted ADRs.
- Non-compliant code detected during reviews MUST be logged as high-priority architectural defect issues and resolved within 14 days.
