---
description: Triage a log file or error paste, propose fixes, and add tests
allowed-tools: Edit, Bash(npm test*), Bash(pnpm test*), Bash(uv run pytest*), Bash(pytest*)
---
1) Summarize recurring errors (paste or path as $1). Identify likely root causes.
2) Propose tests to reproduce, then minimal code changes.
3) Implement only if spec exists; otherwise generate a spec first.

