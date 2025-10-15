---
description: Execute the minimal change needed to make contract tests green (tests → code → rerun)
argument-hint: [feature-key]
allowed-tools: Read, Edit, MultiEdit,
  Bash(go test ./...), Bash(pnpm test*), Bash(npm test*), Bash(yarn test*), Bash(bun test*), Bash(vitest*),
  Bash(uv run pytest*), Bash(pytest*), Bash(mvn test*), Bash(gradle test*), Bash(flutter test*)
---

# Input: $1 = feature key

1. Re-read `docs/specs/$1-plan.md`; restate exact files to touch.
2. **Tests first (contract-level)**: draft a short **Test Selection** list using the TSH — propose 3-7 cases that defend real invariants; briefly note which “obvious/trivial” tests you are **not** adding and why. Keep fast; mock network/DB.
3. Implement **smallest viable diff** that satisfies acceptance checks. Touch only planned files.
4. Run the suite until **green**. If unrelated fails appear, **stop** and add TODO(human) with file + test name.
5. Emit a **concise patch summary** (changed files + rationale).
6. Stop. Do not commit or open PRs.
