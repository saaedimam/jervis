# Dependency Policy

## 1. Allowed Dependency Types
To maintain security, determinism, and instant startup times, Project Jervis strictly limits external dependencies:
- **Standard Library First**: Standard Go library packages (`net/http`, `os`, `sync`, `log/slog`, `context`, `encoding/json`, `time`) must be used for core functionality before considering external packages.
- **Approved Core Infrastructure Dependencies**:
  - CLI Parsing & Terminal UI: `github.com/spf13/cobra`, `github.com/charmbracelet/bubbletea`
  - Pure-Go SQLite Driver: `modernc.org/sqlite` (Cgo-free)
  - WASM Plugin Runtime: `github.com/tetratelabs/wazero`
  - Structured Logging: `log/slog` (Standard Library)
  - Vector Distance Math: Pure Go mathematical primitives or audited minimal SIMD wrappers.

## 2. Forbidden Dependencies
The following categories of dependencies are **strictly forbidden**:
- **Cgo-based C Libraries**: Dynamic links to C shared libraries requiring host gcc/clang toolchains at runtime.
- **Heavy Monolithic Frameworks**: Web frameworks (e.g., Gin, Echo) when `net/http` standard library suffices.
- **Unmaintained Single-Author Utilities**: Packages with single maintainers, no commit activity for >12 months, or lack of test suites.
- **AI SDK Monoliths**: Heavily coupled multi-provider wrappers. HTTP interactions with OpenAI, Claude, Gemini, and Ollama must use lightweight, direct REST/JSON clients built on standard `net/http`.

## 3. Dependency Approval Rules
Adding any new external dependency to `go.mod` requires:
1. **Architect Approval**: Formal verification that standard library cannot fulfill the requirement.
2. **Security & Vulnerability Audit**: Clean report from `govulncheck` with zero known vulnerabilities.
3. **License Audit**: Permissive open-source license verification (MIT, Apache 2.0, BSD-3-Clause).
4. **Maintenance Verification**: Minimum 2 active maintainers and recent repository updates within 6 months.

## 4. Update Policy
- **Minor & Patch Updates**: Evaluated monthly. Automated via dependabot/renovate PRs with full CI test suite verification.
- **Major Version Updates**: Requires explicit architectural review and contract testing.
- **Lockfile Enforcement**: `go.sum` is committed to version control and verified during CI builds via `go mod verify`.

## 5. Security Policy
- **Automated Vulnerability Scanning**: Continuous integration runs `govulncheck` on every commit.
- **Vulnerability SLA**: Critical CVEs must be patched or dependency replaced within 48 hours; High severity within 7 days.
- **Zero Supply-Chain Tampering**: Checksums strictly validated against Go checksum database (`sum.golang.org`).

## 6. Licensing Policy
- **Permissive Licenses Allowed**: MIT, Apache-2.0, BSD-2-Clause, BSD-3-Clause, MPL-2.0.
- **Forbidden Licenses**: Copyleft licenses (GPLv2, GPLv3, AGPLv3, SSPL) are explicitly prohibited to protect the project core.

## 7. Vendoring Policy
- **Deterministic Vendoring**: Production release builds enforce `go build -mod=vendor` to ensure 100% reproducible builds independent of remote proxy availability.
- **Vendor Directory**: The `vendor/` directory is committed to source control for release branches.
