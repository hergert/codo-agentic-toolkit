---
description: Set up N git worktrees so multiple Claude sessions can run in parallel
allowed-tools: Bash(git worktree add:*), Bash(git branch:*), Bash(mkdir -p:*), Edit
---
# Input: $1 = feature/task key, $2 = count (default 2)

1) Make a `trees/` folder; create branches `$1-v1..v$2` from current HEAD.
2) For each, run `git worktree add trees/$1-v{n} $1-v{n}`.
3) Create `trees/$1-v{n}/docs/context/READ_ME_FIRST.md` with quick start steps and the Golden Rules.
4) Print instructions to launch parallel Claude sessions in each worktree.
