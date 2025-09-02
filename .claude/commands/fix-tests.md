---
description: Diagnose failing tests and propose minimal fixes
allowed-tools: Edit, Bash(npm test*), Bash(pnpm test*), Bash(uv run pytest*), Bash(pytest*)
---
1) Run tests quietly and capture failures.
2) Propose the smallest changes to fix; prefer test updates only when spec dictates.
3) Re-run tests until green; produce a concise rationale for each change.

