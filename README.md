# Codo Agentic Toolkit (Claude Code) — v4 Complete

Drop this folder into your repo root. It provides:
- Subagents (`.claude/agents/*.yaml`) for mapping, tests, and review (separate contexts)
- Commands (`.claude/commands/*.md`) for a tight loop: map → plan → prime → tests → implement → review
- Hooks (`.claude/hooks.json`) as seatbelts (edit gating, formatters, commit/prod deploy blocks, rm -rf block)
- Settings (`.claude/settings.json`) plan-by-default; liberal read/test; ask on edits/commits
- Output styles (`.claude/output-styles/*`) to switch personas fast
- Templates (spec, job card, commit message)

## Quick Start
1) Open a terminal in your repo.
2) Run `claude` to start interactive mode.
3) Quick Start
   - `claude` → `/output-style surgical`
   - `/feature-sprint "<key>"` → map + plan + prime
   - Review → `/approve "<key>"` or `/fast-on` (optional `commits`)
   - `/tests "<key>"` → value-dense tests
   - `/implement "<key>"` → smallest change to green
   - `/review "<key>"`
   - (Optional) `/prepare-commit` (requires `ALLOW_COMMITS` or `FAST_MODE`)

## Notes
- Commits/PR merges are gated by hooks (require `.claude/session/ALLOW_COMMITS` or Fast Mode).
- Production deploys require explicit marker `.claude/session/ALLOW_PROD_DEPLOY`.
- Use subagents with job cards in `docs/tasks/` to keep context clean.

## Tuning
- Update `CLAUDE.md` with your real scripts/paths.
- Extend `allowed-tools` per command for your stack.
- Consider enabling a custom output style: `/output-style ~/.claude/output-styles/planner.md`.

## Knowledge Base
- Principles and best practices (repo‑agnostic): `docs/knowledge-base.md`
