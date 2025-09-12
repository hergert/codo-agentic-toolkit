---
description: Map current behavior/UX precisely with evidence; no code edits
argument-hint: [feature-key]
allowed-tools: Read, Bash(git ls-files*), Bash(git grep -n*), Bash(find *)
---
# Input: $1 = feature key (e.g., payments-retry)

1) READ FIRST. Skim @README.md and any @docs/* mentioning "$1".
2) Create `docs/feature-maps/$1.md` (≤ 300 tokens) with:
   - Behavior & contracts: endpoints/CLI/handlers; copy signatures and routes.
   - UX/UI flow: user steps & states; validation/errors/timeouts.
   - Files & entry points: exact paths + line spans (**cap 25**; use `git grep -n` / `find -maxdepth 3`).
   - Data & side-effects: DB tables, external APIs, queues, flags.
   - Existing tests: list specs covering these contracts (`*_test.go`, `*.spec.ts`, `test_*.py`, `*Test.java`) and what they assert (**cap 10**).
   - Invariants & risky edges we must preserve.
   - Open questions as TODO(human).
3) No source edits. Output a one-screen TL;DR. Keep working set lean: avoid listing binary/lock files.