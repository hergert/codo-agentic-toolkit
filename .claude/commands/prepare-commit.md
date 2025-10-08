---
description: Stage touched files and draft a commit message (no push)
allowed-tools: Bash(git add:*), Bash(git status:*), Bash(git diff --staged:*), Edit
---
1) Stage explicit paths from the plan or last diff; avoid lockfiles and large generated assets.
2) Draft the commit message using @.claude/templates/commit_message_template.txt and print for review.
3) Print a **Run manually** banner that includes the commit message and the exact commands for the human to run (e.g. `git commit -F /tmp/commit-msg.txt && git push`).
