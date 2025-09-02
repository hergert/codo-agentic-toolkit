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

