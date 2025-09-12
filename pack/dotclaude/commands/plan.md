---
description: Tight, reviewable plan with explicit file touches; stop for approval
argument-hint: [feature-key]
allowed-tools: Edit, Read
---
# Input: $1 = feature/task key

1. Create `docs/specs/$1-plan.md` with only essentials (**≤ 350 tokens**):
   - Goal: one paragraph (problem → desired behavior).
   - Acceptance checks: **3–7** bullets a reviewer can run (verifiable commands, not prose).
   - Files to touch: exact list; per-file intent.
   - Diff outline: pseudo-code for functions/structs/routes.
   - Risks: migrations/perf/breakages + detection method.
   - Test strategy: contract tests only (files/functions to add).
   - Out of scope: enforce YAGNI.
2. Stop here. Wait for human approval.