---
description: Add only high‑value contract tests (public APIs/invariants)
allowed-tools: Edit, Bash(go test ./...), Bash(go test -run*)
---
# Input: $1 = feature/task key

1) From `docs/specs/$1-plan.md`, propose a minimal test set that:
   - Hits public API/handlers/CLI; 1 happy path + 2–3 valuable edge/failure cases.
   - Uses table‑driven tests; `httptest` for HTTP boundaries.
   - Stays fast: no network unless mocked; isolate global state.
2) Create/modify only the planned `*_test.go` files.
3) Run `!go test ./...`; summarize failures; avoid fixing unrelated code unless the plan requires it (YAGNI).
