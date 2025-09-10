---
description: Review changes against acceptance checks and Golden Rules
allowed-tools: Bash(git diff*), Bash(git status*)
---
1) Compare diff to `docs/specs/$1-plan.md` acceptance checks.
2) Flag any scope creep (YAGNI/KISS violations) and risky changes.
3) List missing contract tests or over‑testing of internals.
4) Output a short decision: approve / request‑changes with exact bullets.

