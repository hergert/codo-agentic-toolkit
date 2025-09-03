---
description: Compact solo flow — map → plan → prime (no code edits)
allowed-tools: Edit, Bash(git grep*), Bash(find *), Bash(git diff*), Bash(git status*)
---

Think deeply and enumerate trade‑offs. If the plan seems risky, keep thinking and propose safer alternatives before edits.

# Input: $1 = feature key

1. Produce `docs/feature-maps/$1.md` (behavior/contracts, UX flow, files, data/side‑effects, existing tests, invariants, TODOs).
2. Create `docs/specs/$1-plan.md` with Goal, 3-7 acceptance checks, explicit file touches, diff outline, risks, contract test strategy, non‑goals.
3. Write `docs/context/$1.md` (links to map/plan, `git status`, `git diff --stat`, top 20 relevant paths, Golden Rules excerpt, ≤200‑token TL;DR).
4. Stop. Ask for approval (`/approve "$1"`) or run `/fast-on` before edits.
