---
description: Summarize errors/logs and propose likely root causes and next actions
allowed-tools:
  - Bash(cat*)
  - Bash(grep*)
  - Bash(rg*)
---
Read logs from `logs/*.log` (or provided path) and produce:
- Top recurring errors
- Probable root causes (ranked)
- Next actions (tests, instrumentation, code changes)
