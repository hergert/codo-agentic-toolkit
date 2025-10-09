---
description: Split-role review (security, readability, tests) → single decision
argument-hint: [feature-key]
allowed-tools: Read, Bash(git diff:*), Bash(git status:*)
---
1) Compare current diff to `docs/specs/$1-plan.md` acceptance checks and **Golden Rules**.
2) Run three short passes and consolidate:
   - 🔒 **Security**: secrets/config, injection surfaces, risky I/O.
   - 📖 **Readability/Perf**: clarity, needless abstraction, hot paths.
   - ✅ **Tests**: contract coverage per plan; missing edges; **reject trivial platform tests** (status=200/json-shape) unless contract-specific; prefer invariants/assertions that catch realistic mutations.
3) Output a markdown **table** with columns: `file | line | issue | rule/acceptance | suggested fix`. Keep the `issue` cell ≤120 chars, the `suggested fix` cell ≤160 chars, and include acceptance-check IDs when possible.
4) Then add a single line: `DECISION: APPROVE` or `DECISION: REQUEST-CHANGES`.
