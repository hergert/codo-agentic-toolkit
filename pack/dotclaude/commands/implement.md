---
description: Implement the approved plan with the smallest viable change
argument-hint: [feature-key]
allowed-tools: Read, Edit, MultiEdit,
  Bash(go test ./...),
  Bash(pnpm test*), Bash(npm test*), Bash(yarn test*), Bash(bun test*), Bash(vitest*),
  Bash(uv run pytest*), Bash(pytest*),
  Bash(mvn test*), Bash(gradle test*),
  Bash(flutter test*)
---
# Input: $1 = feature/task key

1) Re‑read `docs/specs/$1-plan.md`; restate the exact files to touch.
2) Edit only those files; minimise diff.
3) Detect the test runner (go/JS‑TS/python/JVM/flutter) and run the full suite until green.
   - If unrelated failures appear (not touched by plan), stop and record TODO(human) with file + failing test name.
4) Show a concise patch summary.
5) Stop. Do not commit or open PRs.

