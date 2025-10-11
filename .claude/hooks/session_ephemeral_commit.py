#!/usr/bin/env python3
import json, sys, os, shutil

E = json.load(sys.stdin)
event = E.get("hook_event_name") or ""
path = os.path.join(".claude", ".session-commit-message.txt")
os.makedirs(".claude", exist_ok=True)

def ok():
    # Tell Claude Code to hide our stdout (don’t add to context)
    print('{"suppressOutput": true}')
    sys.exit(0)

if event == "SessionStart":
    # Truncate/create the file; don’t prefill
    open(path, "w", encoding="utf-8").close()
    ok()

if event == "SessionEnd":
    try:
        if os.path.exists(path):
            os.remove(path)
    except Exception:
        pass
    ok()

ok()
