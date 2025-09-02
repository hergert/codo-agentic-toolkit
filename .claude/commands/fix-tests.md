---
description: Diagnose failing tests and make the smallest change to make them pass
allowed-tools:
  - Edit
  - MultiEdit
  - Write
  - Bash(npm test*)
  - Bash(npm run test*)
  - Bash(git add*)
  - Bash(git commit*)
---
1) !`npm test --silent`
2) Identify the smallest change to pass failing tests
3) Edit only relevant files; re-run !`npm test --silent` until all green
4) Commit with message "tests: fix <brief>"
