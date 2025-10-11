---
description: Stage touched files and draft a commit message (no push)
allowed-tools: Bash(git add:*), Bash(git status:*), Bash(git diff --staged:*), Bash(git commit -F*), Bash(cat*), Bash(mkdir -p*)
---
1) Stage explicit paths from the plan or last diff (avoid lockfiles/generated assets).

2) Draft the commit message from @.claude/templates/commit_message_template.txt.
   - Keep it concise and action-oriented (no generated footers).
   - Save it to `.claude/.session-commit-message.txt` and show it:

   !`mkdir -p .claude >/dev/null 2>&1 || true`
   !`MSG_PATH=".claude/.session-commit-message.txt" && echo "$MSG_PATH"`
   (Claude: write the message into `$MSG_PATH`, then)
   !`cat .claude/.session-commit-message.txt`

3) If commits are allowed (marker present), run:
   - `git commit -F .claude/.session-commit-message.txt`
   Otherwise:
   - Print instructions for manual commit:  
     `git commit -F .claude/.session-commit-message.txt`
