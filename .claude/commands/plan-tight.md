---
description: Produce a tight, reviewable plan with explicit file touches
allowed-tools: Edit
---
# Input: $1 = feature/task name

1) Create `docs/specs/$1-plan.md` with only the essentials:
   - Goal: one paragraph problem/desired outcome.
   - Acceptance checks: 3–7 bullet checks a reviewer can run.
   - Touches: exact files to edit/create/delete; per-file intent.
   - Diff outline: pseudo‑code of key changes (functions, structs, routes).
   - Risk notes: migrations, perf, breakages; how we’ll detect them.
   - Test strategy: contract tests only (list test files/functions to add).
   - Out of scope: explicitly list non-goals to enforce YAGNI.
2) Stop and ask for approval before any code changes.

