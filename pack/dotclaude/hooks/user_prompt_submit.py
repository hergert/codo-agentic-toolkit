#!/usr/bin/env python3
import json, sys, os, datetime

E = json.load(sys.stdin)
prompt = (E.get("prompt") or "")

# 1) Early block for obviously dangerous patterns in the prompt itself
if "rm -rf /" in prompt or ("curl " in prompt and "| sh" in prompt):
    print("✋ blocked: dangerous pattern in prompt", file=sys.stderr)
    sys.exit(2)

# 2) Prepend a tiny Golden Rules TL;DR once per session
sentinel = ".claude/session/RULES_INJECTED"
if not os.path.exists(sentinel):
    rules = "KISS · YAGNI · Small diffs · Contract tests only · TSH: 3-7 high-signal tests; avoid trivial platform checks."
    ts = datetime.datetime.now().isoformat(timespec='seconds')
    print(f"[{ts}] Rules: {rules}")
    os.makedirs(os.path.dirname(sentinel), exist_ok=True)
    open(sentinel, "w").close()

sys.exit(0)
