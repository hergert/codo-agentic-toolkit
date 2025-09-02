---
description: Mark plan as approved for source edits in this session
allowed-tools: Edit
---
# Input: $1 = feature/task key
1) Write a marker file: `.claude/session/ALLOW_EDITS` with content `$1 approved by human`. Keep it local to this repo/session.
2) Confirm: future source edits are now permitted by hooks.
