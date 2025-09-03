---
description: Generate a one‑shot bash script that applies many small edits from JSON
allowed-tools: Edit
---
1) Create `scripts/batch_edit.sh` that reads `changes.json` with entries {"path","before","after"} or {"path","insert","at"}.
2) The script should validate paths, back up files, apply edits, and print a shortstat.
3) Warn that running the script still requires /approve-plan.

