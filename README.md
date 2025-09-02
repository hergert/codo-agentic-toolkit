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
   - `/map-feature "<feature>"` → evidence-backed map
   - `/plan-tight "<feature>"` → minimal, file-explicit plan
   - `/prime-context "<feature>"` → snapshot and pin context
   - Get approval
   - `/write-tests-contract "<feature>"` → add contract-level tests
   - `/implement-diff-min "<feature>"` → smallest change to green
   - `/review-diff "<feature>"` → reviewer checklist
   - Note: commit/PR/deploy are blocked by policy; you run them explicitly.

## Safety
- Hooks block sensitive writes and deny commits/PR merges/prod deploys; adjust in `.claude/hooks.json`.
- Settings default to plan-mode; source edits ask; commits/PRs are denied; see `.claude/settings.json`.

## Tuning
- Update `CLAUDE.md` with your real scripts/paths.
- Extend `allowed-tools` per command for your stack.
- Consider enabling a custom output style: `/output-style ~/.claude/output-styles/planner.md`.

## Knowledge Base
- Principles and best practices (repo‑agnostic): `docs/knowledge-base.md`
