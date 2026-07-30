# Contributing to Jervis

> Last Updated: 2026-07-31
> Owner: @saaedimam

Thank you for considering a contribution to Jervis! This document outlines the process and standards we follow.

## Getting Started

1. Fork the repository.
2. Clone your fork locally.
3. Create a feature branch: `git checkout -b feat/short-description`
4. Make your changes.
5. Run the full verification suite:
   ```bash
   go fmt ./...
   go vet ./...
   go build ./...
   go test -race ./...
   ```
6. Open a Pull Request against `main`.

## Branch Hygiene

| Branch | Purpose | Direct Push? |
|--------|---------|--------------|
| `main` | Production-ready, always green | No |
| `develop` | Integration branch for features | No |
| `feat/*` | Individual feature work | Yes (on your fork) |
| `fix/*` | Bug fixes | Yes (on your fork) |
| `docs/*` | Documentation updates | Yes (on your fork) |

- **No direct pushes** to `main` or `develop`.
- Open Pull Requests for all changes.
- Keep PRs focused: one logical change per PR.

## Commit Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>[(scope)]: <short description>

[optional body]

[optional footer(s)]
```

### Types

| Type | When to use |
|------|-------------|
| `feat` | New feature or capability |
| `fix` | Bug fix |
| `docs` | Documentation-only change |
| `style` | Formatting, semicolons, etc. (no code change) |
| `refactor` | Code restructuring without behavior change |
| `test` | Adding or updating tests |
| `chore` | Maintenance, dependency updates, CI changes |
| `build` | Build system or external dependency changes |
| `ci` | CI/CD configuration changes |

### Rules

- **Subject line**: max 50 characters, lowercase, no period at end.
- **Body**: Wrap at 72 characters. Explain *what* and *why*, not *how*.
- **References**: Include issue/PR numbers in footer when applicable.

**Example:**
```
feat(runtime): add session isolation for concurrent agents

Enforces workspace-scoped session boundaries so multiple
AI agents can run concurrently without state collision.

Closes #42
```

## Document Structure

- **Single source of truth**: If content lives in Notion, link to it — do not duplicate.
- **Ownership**: Every doc must have an `@owner` and `Last Updated: YYYY-MM-DD` in the header.
- **Size limit**: No doc > 200 lines. Split into sub-docs and link them.
- **Headers**: Sentence case for all headers.
- **Diagrams**: Wrap ASCII diagrams in triple-backtick blocks.

### File Locations

| Content | Location |
|---------|----------|
| Principles, vision | `docs/principles/` |
| Architecture specs | `docs/architecture/` |
| Runtime/service specs | `docs/specs/` |
| ADRs | `docs/adr/` |
| Release notes | `docs/releases/` |
| Contributing guide | Root `CONTRIBUTING.md` |

## Code Standards

- **Go**: Follow `gofumpt -extra` formatting. Run `golangci-lint` before submitting.
- **Error handling**: Never swallow errors. Propagate or log explicitly.
- **Testing**: Unit tests required for all new logic. Race detector must pass.
- **Interfaces**: Define interfaces at consumer boundaries, not producer.

## PR Template

When opening a Pull Request, include:

1. **What** changed and **why**.
2. **How** to verify (commands, expected output).
3. **Checklist**:
   - [ ] `go build ./...` passes
   - [ ] `go test -race ./...` passes
   - [ ] `go fmt ./...` clean
   - [ ] Documentation updated if interfaces changed
   - [ ] Commit messages follow Conventional Commits

## Questions?

Open an issue or reach out to `@saaedimam`.
