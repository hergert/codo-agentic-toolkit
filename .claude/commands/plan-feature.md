---
description: Draft a tight plan + tests before any code changes
allowed-tools: Edit
model: claude-3.7-code
---
1) Create a spec from template for "$1" at `docs/specs/$1.md` using @.claude/templates/spec.md.
2) Extract minimal context (files/dirs) needed; list them explicitly.
3) Propose a **diff plan** (files to touch, function signatures, risks).
4) DO NOT modify code. Place all edits only in the new spec file.
5) End with a checklist: tests to add, commands to run, and a token budget note.

