#!/usr/bin/env python3
import json, sys
focus = (
  "When compacting, keep ONLY: the current plan (docs/specs/*-plan.md), "
  "active files, and the latest failing test output (if any). "
  "Summarize or drop long chat, dormant agents, and old diffs. Aim ≥70% headroom."
)
print(json.dumps({ "hookSpecificOutput": { "hookEventName": "PreCompact", "additionalContext": focus }}))
