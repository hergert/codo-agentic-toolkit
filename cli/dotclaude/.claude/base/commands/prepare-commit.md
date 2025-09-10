---
description: Stage touched files and draft a commit message (no push)
allowed-tools: Bash(git add*), Bash(git status*), Bash(git diff --staged*), Bash(git commit -m*), Edit
---
1) Stage explicit paths from the plan or last diff; avoid lockfiles and large generated assets.
2) Draft the commit message using @.claude/templates/commit_message_template.txt and print for review.
3) If commits are allowed (marker present), run `git commit -m <msg>`; otherwise, print the message and instructions for manual commit.

