# Agentic Coding Knowledge Base

## Core Loop
- Plan → Implement → Test → Review → Ship.
- Prefer minimal, localized diffs with explicit tests.
- Keep humans in the loop for risky edits and scope changes.

## Context Engineering
- Provide only the files and snippets needed for a task.
- Use templates and commands to standardize prompts and outputs.
- Summarize unknowns as TODO(human) to avoid guesswork.

## Guardrails & Permissions
- Block high‑risk paths: secrets, lockfiles, VCS internals, deployment manifests.
- Default to analysis/plan; escalate to edits intentionally and in narrow scopes.
- Auto‑format after edits to keep diffs clean.

## Test‑First Practice
- Add acceptance criteria and failing tests before code.
- Start with the smallest test that proves the requirement.
- Iterate until green; avoid speculative generalization.

## Diff Strategy
- Make the smallest change that passes the tests.
- Avoid cross‑cutting refactors unless planned and approved.
- Stage only touched files; review diffs before commit.

## Change Management
- Commit messages: `feat(scope): one‑line intent` followed by Why/How/Tests.
- Open PRs with evidence: spec link, passing tests, risk notes.
- Pair with CI gates: lint, typecheck, tests.

## Ecosystem Notes
- JS/TS: pnpm, eslint, prettier; prefer Vitest/Jest.
- Python: uv, ruff, black, pytest.
- Go: gofmt, golangci‑lint, `go test`.

## Prompt Patterns
- plan‑feature: generate spec and diff plan before edits.
- implement‑diff: add tests first, then minimal code, then stage.
- review‑diff: sanity‑check staged hunks vs acceptance criteria.

## Risk Boundaries
- No ad‑hoc installs without approval.
- No edits to lockfiles, CI, or secrets unless explicitly approved.
- For refactors: write a refactor‑plan with checkpoints and rollback.

