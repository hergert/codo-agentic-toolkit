---
description: Snapshot the exact context we will use; pin diffs/specs/rules
allowed-tools: Edit, Bash(git status*), Bash(git diff*), Bash(git grep*)
---
# Input: $1 = feature/task name

1) Create `docs/context/$1.md` containing:
   - Link to `docs/feature-maps/$1.md` and `docs/specs/$1-plan.md` if present.
   - Current git status (untracked/modified files) and a short `git diff --stat`.
   - A context include list: up to 20 most relevant paths for this task.
   - The GOLDEN RULES excerpt to keep in scope.
2) Provide a one‑screen TL;DR (≤200 tokens) to prepend in this session.
3) Suggest `/output-style surgical` and `/compact` if the context is large.

