# Codo Agentic Toolkit (Claude Code) — v4.1.1

Drop this folder into your repo root. It provides:
- Subagents (`.claude/agents/*.yaml`) for mapping, tests, and review (separate contexts)
- Commands (`.claude/commands/*.md`) for a tight loop: map → plan → prime → tests → implement → review
- Hooks (`.claude/hooks.json`) as seatbelts (edit gate, formatters, commit/deploy gates, rm -rf guard)
- Settings (`.claude/settings.json`) plan-by-default; liberal read/test; ask on edits/commits; hooks decide via markers
- Output styles (`.claude/output-styles/*`) to switch personas fast
- Templates (spec, job card, commit message)

## Quick Start
1) `claude` → `/output-style surgical`
2) `/feature-sprint "<key>"` → map + plan + prime
3) Review → `/approve "<key>"` or `/fast-on` (optional `commits`)
4) `/tests "<key>"` → value‑dense tests (runner auto‑detected)
5) `/implement "<key>"` → smallest change to green
6) `/review "<key>"`
7) (Optional) `/prepare-commit` (requires `ALLOW_COMMITS` or `FAST_MODE`)

## Notes
- Commits/PR merges ask but are gated by hooks (require `.claude/session/ALLOW_COMMITS` or Fast Mode).
- Production deploys ask and are gated by `.claude/session/ALLOW_PROD_DEPLOY`.
- Use subagents with job cards in `docs/tasks/` to keep context clean.

## Markers (toggles)
- `.claude/session/FAST_MODE`, `.claude/session/ALLOW_EDITS`, `.claude/session/ALLOW_COMMITS`, `.claude/session/ALLOW_PROD_DEPLOY`, `.claude/session/ALLOW_DB_MIGRATE`, `.claude/session/ALLOW_TRIGGER_DEPLOY`, `.claude/session/ALLOW_MOBILE_RELEASE`

## Smoke test
- `/sanity-check` (should echo CODO‑ANCHOR + rules).
- Try an edit outside `docs/` without `/approve` → blocked.
- Try `git commit -m test` → blocked unless `ALLOW_COMMITS` or Fast Mode.
- Try `npx wrangler deploy --env production` → blocked unless `ALLOW_PROD_DEPLOY`.

## Tuning
- Update `CLAUDE.md` with your real scripts/paths.
- Extend `allowed-tools` per command for your stack.
- Consider enabling a custom output style: `/output-style ~/.claude/output-styles/planner.md`.

## Knowledge Base
- Principles and best practices (repo‑agnostic): `docs/knowledge-base.md`
