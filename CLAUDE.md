# Codo Agent Project Memory (CLAUDE.md)

> Purpose: Spec-first, plan-before-code, small focused loops, deterministic wrappers, human-in-the-loop on risk.

## Project Norms
- Coding style: TypeScript strict, ESLint + Prettier
- Tests: Vitest/Jest. Unit-first; integration where relevant.
- Branching: trunk-based; feature branches `feat/<slug>`
- CI gates: lint, typecheck, unit tests must pass

## Agentic Workflow
1) SPEC: Write user story + acceptance criteria
2) PLAN: Produce a task list/diff plan (no edits yet)
3) IMPLEMENT: Apply smallest viable change sets
4) VERIFY: Run tests; iterate until green
5) COMMIT: One well-described commit per scope

## Useful Paths
- Source: ./src
- Tests: ./tests
- Docs: ./docs

## Tools & Scripts
- Test: `npm test --silent`
- Lint: `npm run lint`
- Typecheck: `npm run typecheck`
- Format: `npm run format`

## Human-in-the-loop
- Require human approval for config changes, dependency bumps, and any edit touching security-sensitive files.

## Context Tips
- Include only files needed for current task
- Prefer summaries/notes over raw dumps
- Keep a running "research pack" in ./docs/notes.md (task-scoped)
