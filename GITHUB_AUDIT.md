# GitHub Repository Audit — jervis

This document provides a machine-verifiable checklist for the `jervis` repository governance and security state. It serves as the source of truth for the `scripts/verify_github.sh` validation tool.

## 1. Repository Foundation
- [ ] **Visibility**: Public
- [ ] **Description**: Set to "Local-first runtime and context operating system for AI agents."
- [ ] **Topics**: `runtime`, `ai-agents`, `context-management`, `golang`, `local-first`
- [ ] **Wiki**: Enabled
- [ ] **Discussions**: Enabled
- [ ] **Projects**: Enabled

## 2. Security Configuration
- [ ] **CodeQL**: Enabled and reporting
- [ ] **Secret Scanning**: Enabled
- [ ] **Push Protection**: Enabled
- [ ] **Dependabot**: Enabled (Version & Security)
- [ ] **Dependency Graph**: Enabled
- [ ] **Private Vulnerability Reporting**: Enabled

## 3. Branch Protection: main
- [ ] **Require PR**: Enabled
- [ ] **Required Approvals**: 1
- [ ] **Dismiss Stale Reviews**: Enabled
- [ ] **Require Code Owner Reviews**: Enabled
- [ ] **Status Checks**: Required (`test`, `lint`, `security`)
- [ ] **Conversation Resolution**: Required
- [ ] **Linear History**: Required (No merge commits)
- [ ] **Force Push**: Restricted
- [ ] **Deletion**: Restricted

## 4. Branch Protection: develop
- [ ] **Require PR**: Enabled
- [ ] **Status Checks**: Required (`test`)
- [ ] **Linear History**: Required
- [ ] **Force Push**: Restricted (Limited to Owners)

## 5. Automation & Workflow
- [ ] **CI**: `ci.yml` passing
- [ ] **Release**: `release.yml` configured
- [ ] **CodeQL**: `codeql.yml` configured
- [ ] **Coverage**: `coverage.yml` configured
- [ ] **Benchmarks**: `benchmark.yml` configured
- [ ] **Docs**: GitHub Pages enabled on `gh-pages` branch

## 6. Community & Metadata
- [ ] **License**: Apache-2.0
- [ ] **CODEOWNERS**: Valid and present
- [ ] **Labels**: 20/20 taxonomy present
- [ ] **Milestones**: v0.1.0 through v1.0.0 present

---
**Last Audit Date**: 2026-07-29
**Audit Git SHA**: N/A (Pending first push)
**GitHub CLI Version**: `gh version`
