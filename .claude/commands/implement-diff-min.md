---
description: Implement the approved plan with the smallest possible change
allowed-tools: Edit, MultiEdit, Bash(go test ./*), Bash(gofmt -w*), Bash(go vet*)
---
# Input: $1 = feature/task name

1) Re-read `docs/specs/$1-plan.md` and restate the exact files to touch.
2) Edit only those files; keep changes minimal and cohesive.
3) Re-run `!go test ./...` until green; show a concise patch summary.
4) Stop. Do not commit or open PRs. Ask the user to review the diff.

