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
3) Daily usage (tight loop):
   - `/map-feature "<key>"` → evidence-backed map
   - `/plan-tight "<key>"` → lean plan (files, pseudo‑code, tests, risks)
   - `/prime-context "<key>"` → pin context; then `/sanity-check`
   - Review → `/approve-plan "<key>"` to unlock source edits (auto‑plan gate)
   - `/write-tests-contract "<key>"` → contract-level tests
   - `/implement-diff-min "<key>"` → smallest change to green
   - `/review-diff "<key>"` → crisp decision
   - Optional: `/parallel-worktrees "<key>" 3` to explore variants; use `/context-slim` to keep tokens lean
   - Note: commit/PR/deploy blocked by policy; you run them explicitly.

## Safety
- Auto‑plan: source edits require `.claude/session/ALLOW_EDITS` marker; docs and `.claude/` remain writable.
- Hooks block sensitive writes and deny commits/PR merges/prod deploys; adjust in `.claude/hooks.json`.
- Settings default to plan‑mode; see `.claude/settings.json`.

## Tuning
- Update `CLAUDE.md` with your real scripts/paths.
- Extend `allowed-tools` per command for your stack.
- Consider enabling a custom output style: `/output-style ~/.claude/output-styles/planner.md`.

## Knowledge Base
- Principles and best practices (repo‑agnostic): `docs/knowledge-base.md`
