---
description: Open a PR after tests are green with a clean summary
allowed-tools: Bash(git commit -m*), Bash(gh pr create*), Bash(gh pr view*)
---
1) Verify tests are passing; if not, exit with instructions to run /fix-tests.
2) Commit with the prepared message. Avoid committing large regenerated files.
3) Open a PR with title from the commit header and body summarizing the spec, risks, and test evidence.

