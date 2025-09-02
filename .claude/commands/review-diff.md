---
description: Review staged changes for correctness, risk, and completeness
allowed-tools: Bash(git diff --staged*), Bash(git status*), Bash(git restore --staged*)
---
1) Show staged diff. Check against `docs/specs/$1.md` acceptance criteria.
2) Flag risky edits (security, migrations). Suggest test gaps.
3) If issues: unstage problematic hunks and propose safer alternatives.

