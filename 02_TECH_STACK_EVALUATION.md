# Tech Stack Evaluation

## 1. Evaluation Methodology
This document presents a rigorous evaluation of candidate programming languages for implementing the Jervis runtime core, memory engine, services, AI provider adapters, and interfaces. 

Only three candidate languages are evaluated: **Go (Golang)**, **Rust**, and **Python**.

Each language is scored across 15 weighted criteria on a scale of 1 to 10 (where 10 is optimal). The weighted total determines the single canonical language for Project Jervis.

---

## 2. Weighted Criteria & Scoring Matrix

| Criterion | Weight | Go (Score / Weighted) | Rust (Score / Weighted) | Python (Score / Weighted) |
| :--- | :---: | :---: | :---: | :---: |
| 1. Startup Time | 10% | 9 (0.90) | 10 (1.00) | 3 (0.30) |
| 2. Memory Footprint | 8% | 8 (0.64) | 10 (0.80) | 3 (0.24) |
| 3. Cross-Platform Support | 8% | 10 (0.80) | 9 (0.72) | 7 (0.56) |
| 4. Static Single-Binary | 10% | 10 (1.00) | 10 (1.00) | 2 (0.20) |
| 5. AI Ecosystem & APIs | 6% | 7 (0.42) | 6 (0.36) | 10 (0.60) |
| 6. Plugin Architecture | 7% | 9 (0.63) | 8 (0.56) | 7 (0.49) |
| 7. SQLite Integration | 6% | 9 (0.54) | 9 (0.54) | 9 (0.54) |
| 8. Vector / Embedding Support | 5% | 7 (0.35) | 7 (0.35) | 10 (0.50) |
| 9. MCP Compatibility | 7% | 9 (0.63) | 8 (0.56) | 9 (0.63) |
| 10. CLI Ecosystem | 7% | 10 (0.70) | 9 (0.63) | 7 (0.49) |
| 11. Concurrency Model | 8% | 10 (0.80) | 9 (0.72) | 4 (0.32) |
| 12. Maintainability | 6% | 9 (0.54) | 7 (0.42) | 5 (0.30) |
| 13. Learning Curve | 4% | 9 (0.36) | 4 (0.16) | 10 (0.40) |
| 14. Debugging & Tooling | 4% | 9 (0.36) | 8 (0.32) | 8 (0.32) |
| 15. Packaging & Distribution | 4% | 10 (0.40) | 9 (0.36) | 3 (0.12) |
| **TOTAL WEIGHTED SCORE** | **100%** | **8.97 / 10.0 (89.7%)** | **8.42 / 10.0 (84.2%)** | **5.71 / 10.0 (57.1%)** |

---

## 3. Analysis of Candidates

### Go (Golang) — Selected Winner (Score: 89.7%)
- **Strengths**: Near-instant CLI startup (<15ms), low memory overhead (~15-30MB), native cross-compilation (`GOOS`/`GOARCH`), single static self-contained binary distribution without external runtimes, built-in CSP concurrency (goroutines/channels) ideally suited for the Event Bus and Scheduler, outstanding CLI ecosystem (`cobra`, `bubbletea`, `viper`), pure Go Cgo-free SQLite drivers (`modernc.org/sqlite`), and zero-friction WASM/RPC plugin sandboxing (`wazero`, `hashicorp/go-plugin`).
- **Weaknesses**: Slightly less raw native performance than Rust; AI/ML client libraries are HTTP-REST based rather than native C bindings (which is acceptable since Jervis uses remote/local HTTP LLM providers).

### Rust — Rejected (Score: 84.2%)
- **Reasons for Rejection**: 
  - Excessive compilation times and complex memory lifetime management slow down developer velocity for domain service development.
  - Higher maintenance friction and steep learning curve for community plugin authors.
  - Complex FFI bindings for dynamic runtime plugin reloading compared to Go's RPC/WASM sandboxing.

### Python — Rejected (Score: 57.1%)
- **Reasons for Rejection**:
  - High startup latency (200ms+ interpreter boot time), violating local-first CLI performance standards.
  - High baseline memory overhead (60MB-150MB+), unsuitable for a lightweight background desktop daemon/menu bar widget.
  - Global Interpreter Lock (GIL) creates async concurrency bottlenecks for multi-threaded Event Bus execution.
  - Extremely fragile distribution model; lacks static single-binary compilation without bloated packagers (PyInstaller/cx_Freeze).

---

## 4. Final Language Decision

**Selected Language: Go (Golang 1.22+)**

Go is officially selected as the canonical implementation language for Project Jervis. All core packages (`internal/runtime`, `internal/memory`, `internal/services`, `internal/aiprovider`, `internal/interfaces`, `cmd/`) will be written in standard Go.
