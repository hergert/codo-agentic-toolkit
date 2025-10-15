---
description: Plan a change end-to-end (map → plan) in one pass; stop for human review
argument-hint: [feature-key]
allowed-tools: Read, Edit, Bash(git ls-files:*), Bash(git grep -n:*), Bash(find :*)
---

# Input: $1 = feature key

1. READ FIRST. Skim @README.md and any @docs/\* mentioning "$1".
2. Produce **one file**: `docs/specs/$1-plan.md` (≤ 400 tokens) containing:
   - **Map TL;DR (≤120)**: current behavior/contracts, primary files (≤10) with line spans.
   - **Goal**: problem → desired behavior (1 paragraph).
   - **Acceptance checks (3-7)**: verifiable commands, not prose.
   - **Files to touch**: exact list; intent per file (keep it short). Only the plan file is written during /plan.
   - **Diff outline**: pseudo-code of changes (functions/structs/routes).
   - **Risks**: perf/migrations/breakage + detection method.
   - **Test strategy**: contract tests only to prove behavior (list files/functions).
   - **Out of scope**: enforce YAGNI.
3. Stop. Wait for human review in chat before any edits.
