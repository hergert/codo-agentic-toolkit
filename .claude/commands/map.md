---
description: Map current behavior/UX precisely with evidence; no code edits
allowed-tools: Edit, Bash(git grep*), Bash(git diff*), Bash(find *), Bash(go test -list*)
---
# Input: $1 = feature key (e.g., payments-retry)

1) READ FIRST. Skim @README.md and any @docs/* mentioning "$1".
2) Create `docs/feature-maps/$1.md` with:
   - Behavior & contracts: endpoints/CLI/handlers; copy signatures and routes.
   - UX/UI flow: user steps & states; validation/errors/timeouts.
   - Files & entry points: exact paths with line spans (use `git grep -n`/`find`).
   - Data & side‑effects: DB tables, external APIs, queues, flags.
   - Existing tests: list the files/specs covering these contracts (e.g., `*_test.go`, `*.spec.ts`, `test_*.py`, `*Test.java`) and what they assert.
   - Invariants & risky edges we must preserve.
   - Open questions as TODO(human).
3) No source edits. Output a one‑screen TL;DR.
