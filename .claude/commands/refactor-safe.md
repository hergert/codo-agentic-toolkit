---
description: Perform a small, behavior-preserving refactor with tests green at each step
allowed-tools:
  - Edit
  - MultiEdit
  - Bash(npm test*)
  - Bash(git add*)
  - Bash(git commit*)
---
1) Identify a tiny refactor (rename, extract function, etc.)
2) Apply the change; run !`npm test --silent`
3) If green, commit "refactor: <desc>"
4) If red, revert to last green diff and adjust
