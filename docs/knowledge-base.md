# Agentic Coding Knowledge Base (Living Doc)

Purpose: A concise, up‑to‑date guide to practical agentic coding. This is stack‑agnostic and complements any repo. Emphasis: read first, tight plans, small diffs, contract‑level tests, and human control over commits/deploys.

## Golden Rules

1. Read first: collect only the context that matters. Ask if unsure.
2. YAGNI: implement only what the current ticket/spec needs.
3. KISS: simplest working solution wins; avoid premature abstractions.
4. No speculation: don’t author docs/scripts unless explicitly requested.
5. Contract tests only: cover public APIs/handlers/CLIs, critical edge cases, and invariants. Avoid testing language guarantees or private internals.
6. Minimise the diff: prefer small, cohesive changes you can explain in one breath.
7. Fix root causes: prefer fundamentals over patches/symptom masking.
8. Separate concerns: isolate IO/HTTP from business logic; keep code testable.
9. Accuracy over speed: pause and ask when requirements are unclear.
10. Guardrails matter: block sensitive paths, installs, commits, and prod deploys by default.
11. Tests are the exit gate: changes are done when contract tests pass.
12. Human owns the final actions: commits, PRs, deploys only when explicitly requested.

## Core Loop

- Think → Plan (pseudo‑code) → Show & confirm → Implement → Test → Review.
- Always state expected file touches before edits and verify with commands after.

## Context Engineering

- Provide only relevant files/snippets; avoid full dumps. Cite exact paths and line anchors when possible.
- Summarize unknowns as TODO(human) with specific questions to resolve.
- Keep a short TL;DR for the session so reviewers can follow quickly.

## Planning (tight, reviewable)

Include in a short plan doc for each task:

- Goal: one paragraph of the problem and outcome.
- Acceptance checks: 3-7 bullets a reviewer can run manually.
- Touches: explicit files to edit/create/delete; intent per file.
- Diff outline: pseudo‑code for key changes (functions, structs, routes).
- Risk notes: migrations, perf concerns, breakage detection.
- Test strategy: what contract tests to add/modify and where.
- Out of scope: list non‑goals to enforce YAGNI.

## Mapping Before Changing (read‑first)

- Build an evidence‑backed feature map when unfamiliar with an area:
  - User behavior: flows, endpoints, CLI, or URLs; note errors/timeouts.
  - Contracts: public handlers/functions and their responsibilities.
  - Files & entry points: precise paths and relevant line ranges.
  - Data/side‑effects: DB tables/queries, external APIs, queues, flags.
  - Existing tests: which `*_test.*` cover this and what they assert.
  - Invariants/edges: what must never break; list open questions.

## Pin the Working Context

- Snapshot the working set for a task: links to the feature map and plan, current git status and `git diff --stat`, top N relevant paths, and the Golden Rules excerpt.
- Use `/context-slim` to drop or summarize low‑value context; keep a short TL;DR.

## Testing (contract‑level)

- Focus on public behavior and safety invariants; keep tests fast and valuable.
- Aim for: one happy path + 2-3 failure/edge cases per important behavior.
- Don’ts: avoid testing private helpers, trivial getters/setters, or language/runtime guarantees.

Guidance by ecosystem:

- Go: table‑driven tests; prefer `httptest` for HTTP boundaries; `go test ./...` as the gate; use `go vet` and `go fmt`.
- JS/TS: prefer Vitest/Jest; isolate IO; use Prettier/ESLint; mock network/DB.
- Python: pytest; keep fixtures small and explicit; use ruff/black for hygiene.

## Small, Surgical Diffs

- Edit only the planned files, keeping changes minimal and cohesive.
- Avoid cross‑cutting refactors unless planned and approved. If discovered mid‑change, pause and propose a separate refactor plan.
- Stage only touched files; review diffs before asking for commit/PR.

## Guardrails & Permissions

- Auto‑plan: require `.claude/session/ALLOW_EDITS` marker before source edits.
- Deny writes to secrets/lockfiles/VCS internals/deploy manifests by default.
- Treat installation/network access as high‑risk; require explicit approval.
- Block commits/PR merges/prod deploys by default; let humans drive these.

## Review & Change Management

- Review against acceptance checks and Golden Rules; flag scope creep.
- Ensure tests cover contracts and critical edges; avoid over‑testing internals.
- Communicate with a concise patch summary and verification steps.
- Commit message template (when approved to commit):
  - Header: `feat(scope): one‑line intent`
  - Body: Why / How / Tests

## Anti‑Patterns to Avoid

- Big‑bang changes and speculative design.
- Leaking infrastructure concerns into business logic.
- Silent behavior changes without tests.
- Unbounded context dumps.

## Quick Checklists

Plan

- State goal, acceptance checks, touches, risks, tests, non‑goals.

Implement

- Touch only planned files; keep diffs small; run formatter/linter.

Test

- Add/adjust contract tests; run the suite; explain failures succinctly.

Review

- Compare diff to acceptance checks; flag creep; confirm exit criteria.

Verify (examples)

- Go: `go test ./...` | `go vet ./...` | `gofmt -l .`
- JS/TS: `pnpm test -i` | `pnpm lint` | `pnpm format`
- Python: `pytest -q` | `ruff check` | `black --check .`
