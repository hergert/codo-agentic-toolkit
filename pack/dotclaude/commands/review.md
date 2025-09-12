---
description: Review changes against acceptance checks and Golden Rules
argument-hint: [feature-key]
allowed-tools: Read, Bash(git diff*), Bash(git status*)
---

1. Compare diff to `docs/specs/$1-plan.md` acceptance checks.
2. Flag any scope creep (YAGNI/KISS violations) and risky changes.
3. List missing contract tests or over‑testing of internals; ignore style nits unless the linter fails.
4. Output **one line** decision: `APPROVE` or `REQUEST-CHANGES`, then bullets with file:line and the acceptance check they violate.
