# Release Strategy

## 1. Semantic Versioning
Project Jervis strictly adheres to **Semantic Versioning 2.0.0** (`MAJOR.MINOR.PATCH`):
- `MAJOR`: Breaking changes to public interfaces (`pkg/plugin/`), CLI command contracts, or core storage schemas requiring manual user migration.
- `MINOR`: Backward-compatible new capabilities, new AI provider adapters, new domain services, or new client interfaces.
- `PATCH`: Backward-compatible bug fixes, security patches, performance optimizations, or documentation updates.

---

## 2. Release Channels
- **Stable Channel**: Production releases tagged as `vX.Y.Z`. Thoroughly validated against full unit, integration, contract, performance, and security test suites.
- **Beta / RC Channel**: Release candidates tagged as `vX.Y.Z-rc.N`. Used for community testing prior to major or minor releases.
- **Nightly / Dev Channel**: Automated build artifacts generated from the main branch tagged as `nightly-YYYYMMDD`.

---

## 3. Binary Distribution & Target Architectures
Jervis compiles to single static binaries using Go cross-compilation (`CGO_ENABLED=0`). Every release publishes binaries for the following target OS and Architecture pairs:

| Target Platform | Architecture | Binary File Name |
| :--- | :--- | :--- |
| **macOS (Apple Silicon)** | `darwin/arm64` | `jervis-darwin-arm64` |
| **macOS (Intel)** | `darwin/amd64` | `jervis-darwin-amd64` |
| **Linux (64-bit)** | `linux/amd64` | `jervis-linux-amd64` |
| **Linux (ARM64)** | `linux/arm64` | `jervis-linux-arm64` |
| **Windows (64-bit)** | `windows/amd64` | `jervis-windows-amd64.exe` |

---

## 4. Packaging & Publishing Pipelines

### GitHub Releases & GoReleaser
- **Automated Pipeline**: Triggered automatically on git tag push (`v*`).
- **Build Tooling**: `goreleaser` compiles static binaries, generates SHA256 checksum manifests (`checksums.txt`), generates release notes, and signs artifacts with GPG/Cosign.

### Homebrew Tap Integration
- **Official Tap**: Maintain official Homebrew formula `jervis.rb` in `saaedimam/homebrew-jervis`.
- **Installation Command**: `brew install saaedimam/jervis/jervis`.
- **Automated Formula Updates**: GoReleaser automatically updates Homebrew formula URL and SHA256 hashes upon stable tag release.

---

## 5. Upgrade Strategy
- **Self-Update Command**: `jervis update` checks GitHub Releases API for newer stable versions.
- **Atomic Replacement**: Downloads the platform-specific compressed archive, verifies SHA256 checksum against signed manifest, extracts binary to temporary path, and atomically replaces the existing binary.
- **Schema Migration**: On first boot following binary upgrade, the `Runtime Lifecycle` module runs automatic backward-compatible schema migrations on local SQLite/Timeline databases.

---

## 6. Rollback Strategy
- **Backup Before Migration**: On binary upgrade, the existing binary is preserved as `jervis.bak` and database state is snapshotted to `~/.jervis/backups/vX.Y.Z-pre-migration.db`.
- **Self-Rollback Command**: `jervis update rollback` restores the backed-up binary and database snapshot if post-upgrade health checks fail.
