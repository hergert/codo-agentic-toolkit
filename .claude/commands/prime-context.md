---
description: Snapshot working set (specs, diffs, rules) and pin a lean context
allowed-tools: Edit, Bash(git status*), Bash(git diff*), Bash(git grep*)
---
# Input: $1 = feature/task key

1) Create `docs/context/$1.md` with:
   - Links to `docs/feature-maps/$1.md` and `docs/specs/$1-plan.md`.
   - `git status` and a short `git diff --stat`.
   - Context include list: up to 20 high‑leverage files for this task.
   - Golden Rules excerpt.
   - A one‑screen TL;DR (≤200 tokens) to prepend in this session.
2) Suggest running `/context` and propose drops/summaries to keep tokens low.
