# Codo Agentic Toolkit (Claude Code)

Drop this folder into the root of your repo. It provides:
- `.claude/commands/*.md` — reusable slash commands for spec → plan → implement → verify loops
- `.claude/hooks.json` — deterministic guardrails (block sensitive edits, auto-format after edits)
- `.claude/settings.json` — baseline permission patterns (tune to your stack)
- `.claude/output-styles/*.md` — optional Planner/Reviewer system styles
- `CLAUDE.md` — project memory (norms, paths, workflow)

## Quick Start
1) Open a terminal in your repo.
2) Run `claude` to start interactive mode.
3) Use these commands:
   - `/spec` → write acceptance criteria first
   - `/plan` → propose steps & diff strategy
   - `/diff-plan` → preview patches (no edits)
   - `/implement` → apply smallest diffs + run tests
   - `/fix-tests` → minimal fixes to get green
   - `/pr-open "<branch>" "<title>"` → open PR
   - `/ship` → merge once CI passes
   - `/compact` → summarize & trim context

## Safety
- Hooks block writes to `.env`, `.git/`, and `package-lock.json`; adjust in `.claude/hooks.json`.
- Permissions ask before builds or pushes; see `.claude/settings.json`.

## Tuning
- Update `CLAUDE.md` with your real scripts/paths.
- Extend `allowed-tools` per command for your stack.
- Consider enabling a custom output style: `/output-style ~/.claude/output-styles/planner.md`.
