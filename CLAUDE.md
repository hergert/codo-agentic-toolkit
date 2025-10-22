# CLAUDE.md (project memory)

Project Tag: CODO-ANCHOR

## GOLDEN RULES - Read first, enforce always

1. YAGNI - Write only what the current ticket/spec demands.
2. KISS - Prefer the simplest working solution; avoid needless abstractions.
3. No speculative output - Skip READMEs/docs/scripts unless explicitly requested.
4. Tests via `go test ./...` - Green tests are the exit gate ✅.
5. Find the root cause - Fix fundamentals, never patch symptoms.
6. Minimise the diff - The smaller, clearer change wins.
7. Accuracy over speed - If unsure, ask.
8. KISS architecture - Separate HTTP handling from business logic; code must be testable/readable.
9. STOP adding complexity - No extra features unless explicitly requested.
10. Contract-level testing only - Test public APIs/critical edge-cases/safety invariants; skip language guarantees & private internals.
11. NEVER commit changes - Do NOT commit unless the user explicitly asks. The user commits.
12. NEVER deploy to production - deployments require explicit human approval.

## Default loop

Think → Plan (pseudo‑code) → **Show & confirm** → Implement → Tests green → (User handles commit/PR).

### Agent heuristics

- `reader <key>` (scout) after `/plan` when unknown codepaths remain or “Files to touch” > ~8; review the brief, update the plan, then continue.
- `reviewer <key>` (audit) after `/execute` succeeds and before `/prepare-commit` on large diffs (>8 files) or risky areas; skim the one-page output.

## Conventions

- Default model: **Claude 3.5 Sonnet**. For thorny planning spikes run `/model opusplan`, then switch back to Sonnet before writing code.
- Prefer small, surgical diffs with clear rationale.
- Document only what’s necessary to make reviews effective (spec, plan, tests, risk notes).
- Contract tests only; keep the suite fast and lean.

### Test Selection Heuristic (TSH)

- Write **3-7** high-signal tests per change.
- Test **contracts & risky logic** (branching/state/concurrency/parsing/security/cross-module), not the framework.
- Avoid trivial HTTP status/JSON-shape tests unless behavior is custom.
- Each test defends an invariant users care about; if a small mutation would slip through, add one targeted test.
