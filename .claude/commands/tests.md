---
description: Add only high‑value contract tests (public APIs/invariants)
allowed-tools: Edit, Bash(go test ./...), Bash(go test -run*), Bash(pnpm test*), Bash(npm test*), Bash(yarn test*), Bash(bun test*), Bash(vitest*), Bash(uv run pytest*), Bash(pytest*), Bash(mvn test*), Bash(gradle test*), Bash(flutter test*)
---
# Input: $1 = feature/task key

1) Detect runner from repo (prefer go/pytest/vitest if present). From `docs/specs/$1-plan.md`, propose a minimal test set that:
   - Hits public API/handlers/CLI; 1 happy path + 2–3 valuable edge/failure cases.
   - Uses table‑driven tests; `httptest` for HTTP boundaries.
   - Stays fast: no network unless mocked; isolate global state.
2) Create/modify only the planned test files.
3) Run the detected suite and summarize failures (e.g., Go: `!go test ./...`; JS/TS: `!pnpm test -i` or `!npm test`; Python: `!pytest -q`). Avoid fixing unrelated code unless the plan requires it (YAGNI).
