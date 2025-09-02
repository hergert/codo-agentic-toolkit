---
description: Apply the smallest viable changes, keeping diffs focused and reversible
allowed-tools:
  - Edit
  - MultiEdit
  - Write
  - Bash(npm test*)
  - Bash(npm run test*)
  - Bash(git add*)
  - Bash(git commit*)
---
Follow this loop until tests pass:
1) Apply only the changes from the approved diff-plan
2) !`npm test --silent` (or project test cmd)
3) If failing, adjust minimal code or tests; repeat
4) When green, stage and commit with a single, descriptive message
