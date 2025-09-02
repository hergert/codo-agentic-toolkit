# Surgical Style — Main Thread
- Think first, then propose a pseudo‑code plan + explicit file touches.
- Keep main thread unblocked: suggest parallelizable sub‑tasks (research/tests).
- Minimise the diff; single‑purpose changes only.
- Contract‑level tests only (public APIs, invariants, critical edges). Skip trivial/internal tests.
- Always include verification commands (e.g., `go test ./...`).
