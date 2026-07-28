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
- `feat`: A new feature (MINOR).
- `fix`: A bug fix (PATCH).
- `docs`: Documentation only changes.
- `refactor`: A code change that neither fixes a bug nor adds a feature.
- `perf`: A code change that improves performance.
- `test`: Adding missing tests or correcting existing tests.
- `build`: Changes to build system or deps.
- `ci`: Changes to CI configuration.
- `chore`: Maintenance tasks.

## 4. Git Tagging Strategy
Tags must be used for all releases and significant milestones:
- `v0.1.0-alpha.1` (Initial exploration)
- `v0.1.0-beta.1` (Feature complete for milestone)
- `v0.1.0-rc.1` (Release candidate)
- `v0.1.0` (Stable release)

## 5. Pull Request (PR) Rules
1. **Mandatory Checklist**: Every PR must satisfy the Architectural, Security, and Performance checklists in the PR template.
2. **Linear History**: Rebase feature branches before merging.
3. **Mandatory Review**: At least one approval from a CODEOWNER.
4. **CI Compliance**: All status checks must pass (Lint, Test, Security, CodeQL, Coverage).

## 6. Branch Protection Rules
### `main` branch
- Require PR before merging.
- Require 1 Approval.
- Require Status Checks (Test, Lint, CodeQL, Security).
- Require Conversation Resolution.
- Require Linear History (Squash/Rebase).
- Restrict Force Push & Deletion.

### `develop` branch
- Require PR.
- Require Status Checks.
- Require Linear History.

## 7. Security Policy
Enabled features:
- Dependabot (Weekly updates)
- CodeQL Analysis (Every push/PR to main/develop)
- Secret Scanning
- Dependency Graph
- Private Vulnerability Reporting

## 8. Release Automation
Handled via **GoReleaser**:
1. Tag a release: `git tag -a v0.1.0 -m "Release v0.1.0"`
2. Push tag: `git push origin v0.1.0`
3. CI (`release.yml`) triggers GoReleaser to build binaries, generate checksums, and create a GitHub Release.

## 9. Labels Taxonomy
- `kind/bug`: Defects.
- `kind/feature`: New capabilities.
- `kind/refactor`: Internal cleanup.
- `kind/docs`: Documentation.
- `kind/security`: Vulnerabilities.
- `kind/performance`: Optimizations.
- `kind/architecture`: Structural changes.
- `priority/high`, `priority/medium`, `priority/low`.
- `status/blocked`, `status/in-review`, `status/needs-info`.
- `phase/runtime`, `phase/memory`, `phase/services`, `phase/providers`, `phase/interfaces`.

## 10. Automated Setup (GitHub CLI)
Use `gh` to configure the environment:
```bash
gh repo edit --enable-discussions --enable-wiki --enable-projects
gh api -X PATCH repos/:owner/:repo/import/settings -f private_vulnerability_reporting=enabled
```
