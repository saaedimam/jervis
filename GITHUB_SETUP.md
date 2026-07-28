# GitHub Repository Setup & Governance — jervis

This document defines the canonical repository standards and the declarative governance model for the `jervis` project.

## 1. Declarative Governance Model
Governance is managed as **Infrastructure-as-Code (IaC)**. The single source of truth for all repository settings, protections, and labels is:

👉 **[.github/governance.yaml](file:///Users/ioriimasu/dev/jervis/.github/governance.yaml)**

### Layers of Governance
1.  **Desired State**: Defined in `governance.yaml`.
2.  **Provisioning**: Executed via `scripts/setup_github.sh`.
3.  **Verification**: Continuously audited via `scripts/verify_github.sh` and GitHub Actions.

## 2. Branching Strategy: GitHub Flow + Release Branches
| Branch Pattern | Purpose | Base Branch | Target Branch |
| :--- | :--- | :--- | :--- |
| `main` | Production-ready state. Always stable. | - | - |
| `develop` | Integration branch for current cycle. | `main` | `main` (via release) |
| `feature/*` | New features or improvements. | `develop` | `develop` |

## 3. Commit & Tagging Conventions
- **Commits**: Must follow [Conventional Commits](https://www.conventionalcommits.org/).
- **Tags**: Required for all releases (e.g., `v0.1.0-alpha.1`, `v0.1.0`).

## 4. Continuous Compliance
The `Governance Audit` workflow runs on every push and PR to ensure the live repository does not drift from the specification.

### Audit Artifacts
On every successful audit, two files are generated and uploaded as CI artifacts:
1.  `GITHUB_AUDIT.md`: A human-readable summary of the current compliance state.
2.  `governance-report.json`: A machine-readable report for external monitoring.

## 5. Security & Supply Chain
- **Action Pinning**: All GitHub Actions are pinned by full commit SHA.
- **Linting**: Shell scripts are validated with `shellcheck`; workflows with `actionlint`.
- **Secrets**: Secret scanning and push protection are enabled globally.

---
**Governance Version**: v1.0.0
**Owner**: @ioriimasu
