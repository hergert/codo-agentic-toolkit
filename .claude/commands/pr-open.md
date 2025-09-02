---
description: Create a GitHub PR with a structured template and checklist
allowed-tools:
  - Bash(git add*)
  - Bash(git commit*)
  - Bash(git checkout -b*)
  - Bash(gh pr create*)
---
1) Ensure branch exists; if not, !`git checkout -b feat/$1`
2) Stage and commit outstanding changes with a meaningful message if needed
3) !`gh pr create --fill --title "$2" --body "Implements: $2\nChecklist:\n- [ ] Spec approved\n- [ ] Tests added/updated\n- [ ] Lint/typecheck green\n- [ ] Risk summary attached"`
