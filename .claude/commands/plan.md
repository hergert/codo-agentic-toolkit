---
description: Tight, reviewable plan with explicit file touches; stop for approval
allowed-tools: Edit
---
Think deeply and enumerate trade‑offs. If the plan seems risky, keep thinking and propose safer alternatives before edits.
# Input: $1 = feature/task key

1) Create `docs/specs/$1-plan.md` with only essentials:
   - Goal: one paragraph (problem → desired behavior).
   - Acceptance checks: 3–7 bullets a reviewer can run.
   - Files to touch: exact list; per‑file intent.
   - Diff outline: pseudo‑code for functions/structs/routes.
   - Risks: migrations/perf/breakages + detection method.
   - Test strategy: contract tests only (files/functions to add).
   - Out of scope: enforce YAGNI.
2) Stop here. Wait for human approval.
