---
description: Prepare a release with changelog and tags (dry-run by default)
allowed-tools:
  - Bash(npm run build*)
  - Bash(git tag*)
  - Bash(git push*)
---
1) Generate a CHANGELOG fragment from recent commits
2) Build artifacts with !`npm run build` (ask if not allowed)
3) Propose a semver bump and tag (no push unless approved)
