---
description: Plan a safe multi-step refactor with checkpoints and rollback
allowed-tools: Edit
---
Think deeply and enumerate trade‑offs. If the plan seems risky, keep thinking and propose safer alternatives before edits.
1) Produce `docs/refactors/$1-plan.md` with phases, dependency graph, and test strategy.
2) Mark clear checkpoints and rollback steps per phase.
3) Do not change code. Human approval required before /execute.
