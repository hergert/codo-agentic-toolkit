---
description: Final checklist to merge and monitor
allowed-tools:
  - Bash(gh pr view*)
  - Bash(gh pr merge*)
---
1) Ensure CI green and approvals present
2) Merge the PR using !`gh pr merge --squash --auto`
3) Add a post-merge monitoring note to `docs/rollout/<date>-<slug>.md`
