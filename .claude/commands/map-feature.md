---
description: Map the current behavior of a feature with zero code edits
allowed-tools: Edit, Bash(git grep*), Bash(git diff*), Bash(find *), Bash(go test -list*)
---
# Input: $1 = short feature name (e.g., payments-retry)

1) READ FIRST. Skim @README.md, @docs/*, and existing specs for "$1".
2) Build an evidence-backed map at `docs/feature-maps/$1.md` containing:
   - User-visible behavior: flows, URLs, commands, or endpoints; include screenshots/route names if present.
   - UX/UI flow: enumerate steps/states; note validation, errors, timeouts.
   - Contracts: public handlers/functions/CLI commands this feature exposes.
   - Files & entry points: precise paths + line ranges (use `git grep -n`, `find`); list owners if available.
   - Data & side-effects: DB tables/queries, external APIs, queues, feature flags.
   - Tests covering it: list existing `*_test.go` (or stack equivalents) and what they assert.
   - Invariants & edge cases we must not break.
   - Open questions (TODO(human)) where uncertainty remains.
3) DO NOT modify source code. Produce a short "next-reads" list for speed.

