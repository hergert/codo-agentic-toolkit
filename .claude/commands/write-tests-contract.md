---
description: Add or refine only contract-level tests that matter
allowed-tools: Edit, Bash(go test ./*), Bash(go test -run*)
---
# Input: $1 = feature/task name

1) From `docs/specs/$1-plan.md`, derive a minimal test set that:
   - Exercises public APIs/handlers/CLIs and critical invariants.
   - Covers at least one happy path + 2–3 edge/failure cases with business value.
   - Avoids trivial assertions or tests of private helpers.
2) For Go repos (detected via go.mod):
   - Use table‑driven tests; prefer `httptest` for HTTP boundaries.
   - Keep tests fast (no network unless mocked); avoid global state.
3) Create/modify only the necessary `*_test.go` files listed in the plan.
4) Run `!go test ./...`; summarize failures clearly. If non-plan files fail, propose smallest fixes OR add TODO(human) if the fix is uncertain.

