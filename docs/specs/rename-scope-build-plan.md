# Rename scope→plan, build→execute

## Map TL;DR
**Current behavior**: Two slash commands exist: `/scope` (map+plan phase) and `/build` (implementation phase).

**Primary files** (4):
- `.claude/commands/scope.md:1-19` — command definition & prompt
- `.claude/commands/build.md:1-16` — command definition & prompt
- `pack/dotclaude/commands/scope.md` — source template (mirrors .claude/)
- `pack/dotclaude/commands/build.md` — source template (mirrors .claude/)

**References** (6 occurrences):
- `README.md:56,58` — agentic loop documentation
- `.claude/commands/prime.md:24` — suggests `/scope`
- `.claude/commands/refactor-plan.md:8` — mentions `/build`
- `pack/dotclaude/commands/prime.md:24`, `pack/dotclaude/commands/refactor-plan.md:8` — duplicates

## Goal
Rename `/scope` → `/plan` and `/build` → `/execute` for clearer semantics. "Plan" better conveys the scoping+planning phase; "execute" better conveys the implementation+test phase. Maintain all existing behavior and prompt content.

## Acceptance checks
1. `ls .claude/commands/` shows `plan.md` and `execute.md`, no `scope.md` or `build.md`
2. `ls pack/dotclaude/commands/` shows `plan.md` and `execute.md`, no `scope.md` or `build.md`
3. `grep -r "/scope\|/build" README.md .claude/ pack/` returns zero matches in working tree
4. `grep -r "/plan\|/execute" README.md .claude/commands/prime.md` shows updated references
5. File contents unchanged except for filename references and descriptions

## Files to touch
1. `.claude/commands/scope.md` → **rename** to `plan.md`, update description
2. `.claude/commands/build.md` → **rename** to `execute.md`, update description
3. `pack/dotclaude/commands/scope.md` → **rename** to `plan.md`, update description
4. `pack/dotclaude/commands/build.md` → **rename** to `execute.md`, update description
5. `README.md` — update lines 56, 58 with new command names
6. `.claude/commands/prime.md` — update line 24 reference
7. `.claude/commands/refactor-plan.md` — update line 8 reference
8. `pack/dotclaude/commands/prime.md` — update line 24 reference
9. `pack/dotclaude/commands/refactor-plan.md` — update line 8 reference

## Diff outline
```
# For each .claude/commands/{scope,build}.md:
git mv scope.md plan.md
git mv build.md execute.md

# In plan.md (was scope.md):
- description: "Scope a change end-to-end (map → plan) in one pass; stop for human review"
  → "Plan a change end-to-end (map → plan) in one pass; stop for human review"

# In execute.md (was build.md):
- description: "Write only the minimal change needed to make contract tests green (tests → code → rerun)"
  → "Execute the minimal change needed to make contract tests green (tests → code → rerun)"

# In README.md:
- "/scope" → "/plan"
- "/build" → "/execute"

# In prime.md, refactor-plan.md (both .claude/ and pack/):
- "/scope" → "/plan"
- "/build" → "/execute"
```

## Risks
- **Breakage**: Users with existing workflows referencing `/scope` or `/build` will get "command not found"
  - *Detection*: Manual testing — invoke old command names after rename
  - *Mitigation*: Document in release notes/changelog as breaking change
- **Pack sync**: If embedded pack and source pack drift, init/update may fail
  - *Detection*: Run `codo init --no-tui` in clean repo; verify both commands present
  - *Mitigation*: Keep `.claude/` and `pack/dotclaude/` identical

## Test strategy
**Manual contract tests** (no automated Go tests needed for markdown renames):
1. `codo init --no-tui` in temp repo → verify `.claude/commands/{plan,execute}.md` exist
2. Invoke `/plan test-feature` → verify command runs (reads plan.md prompt)
3. Invoke `/execute test-feature` → verify command runs (reads execute.md prompt)
4. Invoke `/scope test` → verify error "command not found"
5. Invoke `/build test` → verify error "command not found"
6. `grep -r "/scope\|/build" .claude/ pack/ README.md` → zero matches

## Out of scope
- No changes to command **behavior** or **prompt content** (only names/descriptions)
- No deprecation warnings or backwards-compatibility shims
- No changes to other commands besides prime.md and refactor-plan.md references
- No updates to external documentation beyond README.md
