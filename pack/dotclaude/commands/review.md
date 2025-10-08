---
description: Split-role review (security, readability, tests) → single decision
argument-hint: [feature-key]
allowed-tools: Read, Bash(git diff:*), Bash(git status:*)
---
1) Compare current diff to `docs/specs/$1-plan.md` acceptance checks and **Golden Rules**.
2) Run three short passes and consolidate:
   - 🔒 **Security**: secrets/config, injection surfaces, risky I/O.
   - 📖 **Readability/Perf**: clarity, needless abstraction, hot paths.
   - ✅ **Tests**: contract coverage per plan; missing edges; over-testing internals.
3) Output **one line**: `APPROVE` or `REQUEST-CHANGES`.
4) Then bullets: `file:line → issue → which acceptance check / rule it violates`.
