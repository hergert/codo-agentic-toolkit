---
description: Prime the session context WITHOUT writing files; keep ≥70% headroom
argument-hint: [feature-key]
allowed-tools: Read, Bash(git status:*), Bash(git diff --stat:*), Bash(git ls-files:*), Bash(git grep -n:*)
---

# Input: $1 = feature/task key (optional)

# Prime (no writes)

1. Show `git status -s` and a short `git diff --stat`.
2. Read (only):
   - @README.md
   - @docs/knowledge-base.md
   - @docs/feature-maps/$1.md (if present)
   - The 5 most relevant experience shards from `docs/experience/ledger.md`
     (rank by tag/path overlap with "$1" and recency).
3. Build a one-screen TL;DR (≤200 tokens) with:
   - What we're solving now, the top 3 files to open next, and the risks to watch.
   - A TSH-driven “candidate tests” sketch (3-7), each tied to a real invariant.
   - A proposed "include list" (max 20 files) for subsequent steps.

# Slim tokens

4. Run `/context` and propose exactly what to unload/summarize to keep ≥70% free
   (unused agents, long memory sections, large messages). DO NOT write any files.
5. If a scope plan is missing/outdated, suggest `/plan "$1"`—do not create it yourself.
