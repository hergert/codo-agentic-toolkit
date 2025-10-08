---
description: Remove completed/dirty worktrees safely
allowed-tools: Bash(git worktree list:*), Bash(git worktree remove:*), Bash(rm -rf:*)
---
# Input: $1 = feature/task key

1) List current worktrees: run `!git worktree list`.
2) For each matching `trees/$1-v*`:
   - Try to remove via Git: run `!git worktree remove -f <path> 2>/dev/null || true`.
   - If a directory remains without a registered worktree, remove it: run `!rm -rf <path>`.
3) Print remaining worktrees: run `!git worktree list`.
