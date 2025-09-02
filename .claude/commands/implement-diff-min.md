---
description: Implement the approved plan with the smallest viable change
allowed-tools: Edit, MultiEdit, Bash(go test ./...), Bash(gofmt -w*), Bash(go vet*)
---
# Input: $1 = feature/task key

1) Re‑read `docs/specs/$1-plan.md`; restate the exact files to touch.
2) Edit only those files; minimise diff.
3) Run `!go test ./...` until green; show a concise patch summary.
4) Stop. Do not commit or open PRs.
