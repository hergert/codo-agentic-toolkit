# Job Cards

Create a job card per task to hand off to subagents.

Template: `docs/tasks/<key>.md`

- Objective: one sentence
- Inputs: up to 10 paths/logs/links
- Constraints: KISS, YAGNI, smallest diff, contract tests only, no commits
- Allowed tools: Read, Bash(go test ./...), (Edit? only if needed)
- Deliverable: docs/results/<key>.<role>.md
- Done when: 3-7 acceptance checks pass
