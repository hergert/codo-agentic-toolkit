#!/usr/bin/env python3
import json, sys, re, datetime

E = json.load(sys.stdin)
prompt = (E.get("prompt") or "")

# 1) Early block for obviously dangerous patterns in the prompt itself
BAD = [r"rm\s+-rf\s+/", r"curl\s+[^|]+\|\s*sh"]
if any(re.search(pat, prompt, re.I) for pat in BAD):
    print("✋ blocked: dangerous pattern in prompt", file=sys.stderr)
    sys.exit(2)

# 2) Prepend a tiny Golden Rules TL;DR (cheap; improves reliability)
rules = "KISS · YAGNI · Small diffs · Contract tests only · Ask if unsure."
timestamp = datetime.datetime.now().isoformat(timespec='seconds')
print(f"[{timestamp}] Rules: {rules}")

sys.exit(0)

