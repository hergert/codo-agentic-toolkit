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
   - `/plan-feature "<feature>"` → spec + diff plan
   - `/design-spec "<issue>"` → acceptance-test-first spec
   - `/implement-diff "<feature>"` → tests first + minimal code
   - `/fix-tests` → minimal fixes to get green
   - `/review-diff "<feature>"` → review staged changes
   - `/ship` → open PR once tests are green
   - `/compact` → summarize & trim context

## Safety
- Hooks block writes to `.env`, `.git/`, and `package-lock.json`; adjust in `.claude/hooks.json`.
- Permissions ask before builds or pushes; see `.claude/settings.json`.

## Tuning
- Update `CLAUDE.md` with your real scripts/paths.
- Extend `allowed-tools` per command for your stack.
- Consider enabling a custom output style: `/output-style ~/.claude/output-styles/planner.md`.

## Pack V1 Docs
- Setup and usage: `docs/claude-pack-v1.md`
- Knowledge base (general, repo‑agnostic): `docs/knowledge-base.md`
