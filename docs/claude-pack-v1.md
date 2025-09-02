# Claude Code Agentic Tools Pack (v1)

A ready‑to‑drop **`.claude/`** bundle that bakes in context‑engineering and agentic coding best practices. Copy these files into your repo, then tweak paths/models/tool patterns as needed.

> Works well with: monorepos, TDD flows, PR‑driven releases, and human‑in‑the‑loop approvals.

---

## 0) One‑time setup

```bash
mkdir -p .claude/commands .claude/output-styles .claude/templates
```

* Add these files to your repo (commit them).
* Start Claude REPL in your repo: `claude`
* Run `/permissions` to review/adjust rules and `/config` to verify settings.

---

## 1) Project memory: `./CLAUDE.md`

> Encodes team conventions that become the model’s default context for this repo.

```md
# CLAUDE.md (project memory)

## Mission
Ship reliable, well‑tested code via **plan → implement → test → review → ship** loops. Prefer **minimal diffs**, explicit tests, and deterministic tooling.

## Coding Standards
- Language/toolchain: Node 20 + pnpm; Python 3.11 + uv; Go 1.22
- Lint/format: eslint + prettier; ruff + black; gofmt + golangci-lint
- Tests: Vitest/Jest; Pytest; Go test. Always add acceptance tests before code.

## Workflow Conventions
- **Spec first** using `/plan-feature` or `/design-spec`.
- Never edit lockfiles or secrets. Prefer small, surgical diffs.
- Commit message template: `feat(scope): one-line intent` + body with **Why/How/Tests**.
- Require green tests before opening a PR (`/ship`).

## Useful Context Slots
@README.md
@CONTRIBUTING.md
@.claude/templates/spec.md
@.claude/templates/commit_message_template.txt

## Risk Boundaries
- Deny writes to `.env`, `.git/`, lockfiles, CI config, deployment manifests unless explicitly approved.
- For migrations/refactors: create a plan with `/refactor-plan` and get human approval.
```

---

## 2) Settings: `.claude/settings.json`

> Opinionated defaults. Adjust the patterns for your stack. Use `/permissions` to inspect/apply.

```json
{
  "defaultMode": "plan",
  "additionalDirectories": ["docs", "scripts", ".github"],
  "rules": [
    { "tool": "Edit|Write|MultiEdit", "scope": "src/**|tests/**", "policy": "ask" },
    { "tool": "Edit|Write|MultiEdit", "scope": "package.json|pyproject.toml|go.mod", "policy": "ask" },
    { "tool": "Edit|Write|MultiEdit", "scope": "**/*.env|**/.env|**/.git/**|**/package-lock.json|**/pnpm-lock.yaml|**/yarn.lock|**/.github/**|**/deploy/**", "policy": "deny" },

    { "tool": "Bash(npm test*)|Bash(pnpm test*)|Bash(uv run pytest*)|Bash(pytest*)|Bash(go test*)", "policy": "allow" },
    { "tool": "Bash(npm run build*)|Bash(pnpm build*)|Bash(uv run *)|Bash(pip install*)|Bash(pnpm install*)", "policy": "ask" },

    { "tool": "Bash(git status*)|Bash(git diff*)|Bash(git add*)|Bash(git restore --staged*)|Bash(git checkout*)|Bash(git switch*)", "policy": "allow" },
    { "tool": "Bash(git commit -m*)|Bash(git tag*)", "policy": "ask" },

    { "tool": "Bash(gh pr create*)|Bash(gh pr view*)|Bash(gh pr status*)|Bash(gh pr merge*)", "policy": "ask" }
  ],
  "notes": "Plan mode by default. Escalate to acceptEdits only after /plan-feature or /design-spec is approved. Use /compact before long runs; run /cost to watch budget."
}
```

> **Why:** Defaults to **analysis‑only** (plan) to prevent overeager edits, then you escalate intentionally. Rules enforce safe Bash patterns and file scopes.

---

## 3) Hooks: `.claude/hooks.json`

> Deterministic guardrails + automation. Import/edit via `/hooks`.

```json
{
  "PreToolUse": [
    {
      "matcher": "Edit|MultiEdit|Write",
      "hooks": [
        {
          "type": "command",
          "command": "python3 - <<'PY'\nimport json,sys,os,fnmatch\ne=json.load(sys.stdin)\np=e.get('tool_input',{}).get('file_path','') or ''\nBLOCKERS=['.env','/package-lock.json','/pnpm-lock.yaml','/yarn.lock','/.git/']\nif any(b in p for b in BLOCKERS):\n  print('Blocked write to sensitive file:', p)\n  sys.exit(2)\nsys.exit(0)\nPY"
        }
      ]
    },
    {
      "matcher": "Bash(pip install*)|Bash(pnpm install*)|Bash(npm install*)",
      "hooks": [
        {
          "type": "command",
          "command": "echo 'Install blocked by default. Use human approval or sandbox.' ; exit 2"
        }
      ]
    }
  ],
  "PostToolUse": [
    {
      "matcher": "Edit|MultiEdit|Write",
      "hooks": [
        { "type": "command", "command": "jq -r '.tool_input.file_path' | { read f; case \\"$f\\" in *.ts|*.tsx|*.js|*.jsx) npx prettier --write \\"$f\\" ;; *.py) uv run ruff check --fix \\"$f\\"; uv run black \\"$f\\" ;; *.go) go fmt \\"$f\\" ;; esac; }" }
      ]
    }
  ],
  "Notification": [
    {
      "matcher": "SubagentStop|Stop",
      "hooks": [
        { "type": "command", "command": "echo '[Claude] Session ended. Ensure tests are green before /ship.'" }
      ]
    }
  ]
}
```

> **Why:** Blocks high‑risk writes, prevents ad‑hoc installs, and auto‑formats after edits—tightening the feedback loop and keeping diffs surgical.

---

## 4) Output styles

### `.claude/output-styles/surgical.md`

```md
# Surgical Style

- Prefer **small, localized diffs** with clear rationale.
- Always suggest **unit/acceptance tests first** and reference them in code.
- Leave `TODO(human)` where external knowledge/credentials are required.
- Before editing: summarize the plan and expected diff footprint.
- After editing: show a **patch summary** and how to verify (commands).
```

### `.claude/output-styles/learning-onboarding.md`

```md
# Learning Onboarding Style

- Narrate decisions and trade‑offs succinctly.
- Insert short comments explaining non‑obvious code paths.
- Emit a final "What to read next" list with repo files and external docs.
```

---

## 5) Templates

### `.claude/templates/spec.md`

```md
# Feature Spec: $1

## Problem & Goal
(What user story are we solving? Success criteria?)

## Acceptance Criteria
- [ ] …
- [ ] …

## Risks / Constraints
- Perf / Security / Migration notes

## Test Plan (write these tests first)
- [ ] Unit: …
- [ ] Integration: …
- [ ] E2E: …

## Rollout & Monitoring
- Flag: … | Metrics: … | Rollback: …
```

### `.claude/templates/commit_message_template.txt`

```
feat(scope): one‑line intent

Why
-
How
-
Tests
-
```

---

## 6) Commands (slash‑style markdown)

Place these in `.claude/commands/`. They appear as `/plan-feature`, `/implement-diff`, etc., inside the REPL.

> **Frontmatter keys:** `description`, `allowed-tools`, optional `model`, optional `notes`. Use `$1`, `$2` for args. `@path` pulls file contents into context.

### `.claude/commands/plan-feature.md`

```md
---
description: Draft a tight plan + tests before any code changes
allowed-tools: Edit
model: claude-3.7-code
---
1) Create a spec from template for "$1" at `docs/specs/$1.md` using @.claude/templates/spec.md.
2) Extract minimal context (files/dirs) needed; list them explicitly.
3) Propose a **diff plan** (files to touch, function signatures, risks).
4) DO NOT modify code. Place all edits only in the new spec file.
5) End with a checklist: tests to add, commands to run, and a token budget note.
```

### `.claude/commands/design-spec.md`

```md
---
description: Turn an issue or request into an acceptance-test-first spec
allowed-tools: Edit
---
1) Read @README.md and any referenced design docs.
2) Generate `docs/specs/$1.md` with acceptance criteria and a failing test list.
3) Summarize unknowns as TODO(human) and propose the smallest viable scope.
```

### `.claude/commands/implement-diff.md`

```md
---
description: Implement the approved plan with surgical diffs and tests
allowed-tools: Edit, MultiEdit, Bash(npm test*), Bash(pnpm test*), Bash(uv run pytest*), Bash(pytest*), Bash(go test*), Bash(git add*), Bash(git restore --staged*)
model: claude-3.7-code
---
1) Read `docs/specs/$1.md`. Summarize plan and confirm scope.
2) Add tests FIRST. Run tests: !`pnpm test -i` (or project equivalent). Expect failures.
3) Implement minimal code to satisfy tests with **small diffs**. Avoid touching lockfiles.
4) Re-run tests until green. Show a patch summary.
5) Stage only touched files: !`git add -p` (or explicit paths).
6) Prepare commit message using @.claude/templates/commit_message_template.txt (do not commit yet).
```

### `.claude/commands/fix-tests.md`

```md
---
description: Diagnose failing tests and propose minimal fixes
allowed-tools: Edit, Bash(npm test*), Bash(pnpm test*), Bash(uv run pytest*), Bash(pytest*)
---
1) Run tests quietly and capture failures.
2) Propose the smallest changes to fix; prefer test updates only when spec dictates.
3) Re-run tests until green; produce a concise rationale for each change.
```

### `.claude/commands/review-diff.md`

```md
---
description: Review staged changes for correctness, risk, and completeness
allowed-tools: Bash(git diff --staged*), Bash(git status*), Bash(git restore --staged*)
---
1) Show staged diff. Check against `docs/specs/$1.md` acceptance criteria.
2) Flag risky edits (security, migrations). Suggest test gaps.
3) If issues: unstage problematic hunks and propose safer alternatives.
```

### `.claude/commands/ship.md`

```md
---
description: Open a PR after tests are green with a clean summary
allowed-tools: Bash(git commit -m*), Bash(gh pr create*), Bash(gh pr view*)
---
1) Verify tests are passing; if not, exit with instructions to run /fix-tests.
2) Commit with the prepared message. Avoid committing large regenerated files.
3) Open a PR with title from the commit header and body summarizing the spec, risks, and test evidence.
```

### `.claude/commands/refactor-plan.md`

```md
---
description: Plan a safe multi-step refactor with checkpoints and rollback
allowed-tools: Edit
---
1) Produce `docs/refactors/$1-plan.md` with phases, dependency graph, and test strategy.
2) Mark clear checkpoints and rollback steps per phase.
3) Do not change code. Human approval required before /implement-diff.
```

### `.claude/commands/code-map.md`

```md
---
description: Build a lightweight map of key modules, entrypoints, and tests
allowed-tools: Edit
---
1) Create or update `docs/code-map.md` listing modules, owners, and primary tests.
2) Link to critical paths and known footguns.
```

### `.claude/commands/triage-errors.md`

```md
---
description: Triage a log file or error paste, propose fixes, and add tests
allowed-tools: Edit, Bash(npm test*), Bash(pnpm test*), Bash(uv run pytest*), Bash(pytest*)
---
1) Summarize recurring errors (paste or path as $1). Identify likely root causes.
2) Propose tests to reproduce, then minimal code changes.
3) Implement only if spec exists; otherwise generate a spec first.
```

---

## 7) Daily control loop (how to use)

1. `/output-style surgical` (or onboarding style for new team members).
2. `/plan-feature Add payment retries` → Review spec & plan.
3. `/permissions` → switch to `acceptEdits` for the session if approved.
4. `/implement-diff Add payment retries` → Tests → Minimal code → Green.
5. `/review-diff Add payment retries` → Address findings.
6. `/ship` → PR created with clear summary and evidence.

> Sprinkle `/compact` before long sessions and `/cost` to watch token budget. Use `/status` to see tool approvals in effect.

---

## 8) Notes & tailoring

* Adjust **Bash patterns** to your ecosystem (bun, rye, tox, maven/gradle, cargo).
* If you use MCP tools (e.g., GitHub, Jira, internal search), add explicit allow rules per tool (no wildcards).
* For high‑risk repos, keep `defaultMode: "plan"` and only enable edits in narrow scopes.
* Pair this with CI gates (lint, tests, security) so `/ship` mirrors prod reality.

