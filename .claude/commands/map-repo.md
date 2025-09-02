---
description: Build a lightweight repo map and symbol index to aid navigation
allowed-tools:
  - Bash(rg*)
  - Bash(ctags*)
  - Bash(npm run typecheck*)
  - Edit
  - Write
---
1) If available, !`rg -n "^export|function|class" src > docs/repo-map.txt`
2) If universal-ctags installed, !`ctags -R -f docs/tags .`
3) Run !`npm run typecheck --silent` to gather type errors into `docs/typecheck.txt`
4) Summarize key modules and their dependencies in `docs/repo-overview.md`
