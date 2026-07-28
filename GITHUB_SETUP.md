# GitHub Repository Setup & Governance — jervis

This document defines the canonical repository standards, governance policies, and operational workflows for the `jervis` project.

## 1. Repository Standards
- **Visibility**: Public
- **Primary Branch**: `main`
- **Development Branch**: `develop`
- **Merge Strategy**: Squash and Merge (preferred for features), Merge Commit (for releases/hotfixes).
- **Branch Protection**: Required for `main` and `develop`.

## 2. Branching Strategy: GitHub Flow + Release Branches
| Branch Pattern | Purpose | Base Branch | Target Branch |
| :--- | :--- | :--- | :--- |
| `main` | Production-ready state. Always stable. | - | - |
| `develop` | Integration branch for current cycle. | `main` | `main` (via release) |
| `feature/*` | New features or improvements. | `develop` | `develop` |
| `fix/*` | Non-critical bug fixes. | `develop` | `develop` |
| `release/vX.Y.Z` | Preparation for a new production release. | `develop` | `main` & `develop` |
| `hotfix/vX.Y.Z` | Critical production fixes. | `main` | `main` & `develop` |

## 3. Commit Conventions
All commits must follow the **Conventional Commits** specification:

```text
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Types:**
- `feat`: A new feature (correlates with MINOR in SemVer).
- `fix`: A bug fix (correlates with PATCH in SemVer).
- `docs`: Documentation only changes.
- `style`: Changes that do not affect the meaning of the code.
- `refactor`: A code change that neither fixes a bug nor adds a feature.
- `perf`: A code change that improves performance.
- `test`: Adding missing tests or correcting existing tests.
- `build`: Changes that affect the build system or external dependencies.
- `ci`: Changes to CI configuration files and scripts.
- `chore`: Other changes that don't modify src or test files.

## 4. Pull Request (PR) Rules
1. **Linear History**: Rebase feature branches before merging.
2. **Squash and Merge**: Use for all `feature/*` and `fix/*` PRs into `develop`.
3. **Atomic PRs**: One PR per logical change.
4. **Mandatory Review**: At least one approval from a CODEOWNER.
5. **CI Compliance**: All status checks must pass (Lint, Test, Security).

## 5. Issue Templates
Standardized forms are available in `.github/ISSUE_TEMPLATE/`:
- `Bug Report`: For defects and regressions.
- `Feature Request`: For new capabilities.
- `Architecture Proposal`: For structural changes.
- `ADR Proposal`: For formal decision records.
- `Performance`: For optimization tasks.
- `Security`: For vulnerability reporting.

## 6. Labels & Milestones
### Standard Labels
- `bug`: Something isn't working.
- `feature`: New functionality.
- `architecture`: Structural design changes.
- `adr`: Architecture Decision Records.
- `performance`: Speed or resource optimizations.
- `security`: Security-related issues.
- `triage`: Needs classification.
- `critical`: High-priority / blocking.
- `good first issue`: Suitable for new contributors.

### Milestones
- `Phase 1.0 — Runtime Foundation`
- `Phase 2.0 — Context Engine`
- `Phase 3.0 — OS Integration`

## 7. CODEOWNERS Policy
- The `@ioriimasu` team/user is the default owner of all code.
- Changes to `ARCHITECTURE.md` and `internal/runtime/` require explicit approval from the primary owner.

## 8. Release Workflow (Semantic Versioning)
1. Branch from `develop` to `release/vX.Y.Z`.
2. Update `CHANGELOG.md` and version files.
3. Open PR into `main`.
4. Merge into `main` and tag `vX.Y.Z`.
5. Back-merge `main` into `develop`.

## 9. Branch Protection Rules
**`main` branch:**
- Require pull request reviews before merging.
- Required approvals: 1.
- Dismiss stale pull request approvals when new commits are pushed.
- Require status checks to pass before merging.
- Require conversation resolution before merging.
- Enforce restrictions for administrators.

**`develop` branch:**
- Same as `main`, but allows force pushes for rebase cleanup by owners.
