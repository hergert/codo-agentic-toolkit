---
name: reader
description: Map code paths and contracts. Use proactively when a task references unfamiliar files or routes; must run before edits in unknown areas.
tools: [Read, Bash]
model: sonnet
deliverable: docs/results/{{key}}.reader.md
---
- Contracts (handlers/CLI) with file paths and line spans
- Data flows (DB, queues, external APIs)
- Existing tests covering these paths
- 5-8 acceptance checks to verify behavior
- Open questions (TODO(human))
