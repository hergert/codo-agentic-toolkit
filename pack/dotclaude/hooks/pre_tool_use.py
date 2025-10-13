#!/usr/bin/env python3
import json, sys, os, shlex, re
from datetime import datetime, timedelta

E = json.load(sys.stdin)
tool = E.get("tool_name") or ""
ti = E.get("tool_input") or {}
cmd = ti.get("command", "")
path = (ti.get("file_path") or "")

SENSITIVE_NAMES = {
    ".env",
    ".env.local",
    ".env.production",
    ".env.development",
    "id_rsa",
    "id_ed25519",
    "known_hosts",
    "serviceAccountKey.json",
    "GoogleService-Info.plist",
    "google-services.json",
}
SENSITIVE_DIR_SUBSTR = ("/.git/", "/config/secrets/")

if tool == "Read":
    base = os.path.basename(path)
    if base in SENSITIVE_NAMES or any(sub in path for sub in SENSITIVE_DIR_SUBSTR):
        out = {
            "hookSpecificOutput": {
                "hookEventName": "PreToolUse",
                "permissionDecision": "ask",
                "permissionDecisionReason": f"Read of sensitive path: {path}"
            }
        }
        print(json.dumps(out))
        sys.exit(0)

    # ask on huge or minified assets
    try:
        if path and os.path.exists(path):
            size = os.path.getsize(path)
            if size > 300_000:  # ~300 KB
                print(json.dumps({"hookSpecificOutput": {
                    "hookEventName": "PreToolUse",
                    "permissionDecision": "ask",
                    "permissionDecisionReason": f"Large file ({size//1024}KB). Summarize or open smaller slice?"
                }}))
                sys.exit(0)
        base_lower = os.path.basename(path).lower()
        if base_lower.endswith((".min.js", ".min.css", ".map", ".bundle.js", ".lock", "package-lock.json", "pnpm-lock.yaml", "yarn.lock")):
            print(json.dumps({"hookSpecificOutput": {
                "hookEventName": "PreToolUse",
                "permissionDecision": "ask",
                "permissionDecisionReason": f"Likely generated/minified asset: {base_lower}. Proceed?"
            }}))
            sys.exit(0)
    except Exception:
        pass

# Allow docs/.claude edits always
if tool in ("Edit","Write","MultiEdit") and (path.startswith("docs/") or path.startswith(".claude/")):
    sys.exit(0)

def has_fresh_plan(hours: int = 48) -> bool:
    specs_dir = os.path.join("docs", "specs")
    if not os.path.isdir(specs_dir):
        return False
    cutoff = datetime.utcnow() - timedelta(hours=hours)
    for entry in os.listdir(specs_dir):
        if not entry.endswith("-plan.md"):
            continue
        candidate = os.path.join(specs_dir, entry)
        try:
            mtime = os.path.getmtime(candidate)
        except OSError:
            continue
        if datetime.utcfromtimestamp(mtime) >= cutoff:
            return True
    return False

if tool in ("Edit","Write","MultiEdit") and not (path.startswith("docs/") or path.startswith(".claude/")):
    if not has_fresh_plan():
        out = {
            "hookSpecificOutput": {
                "hookEventName": "PreToolUse",
                "permissionDecision": "ask",
                "permissionDecisionReason": "No recent plan in docs/specs/*-plan.md (≤48h). Proceed only if this is a trivial fix."
            }
        }
        print(json.dumps(out))
        sys.exit(0)

# Secrets / sensitive paths guard
bn = os.path.basename(path)
BLOCK_BASENAMES = {".env", ".env.local", ".env.production", ".env.development",
                   "package-lock.json","pnpm-lock.yaml","yarn.lock",
                   "id_rsa","id_ed25519","known_hosts",
                   "serviceAccountKey.json","GoogleService-Info.plist","google-services.json"}
BLOCK_DIR_SUBSTR = ("/.git/", "/config/secrets/")
if bn in BLOCK_BASENAMES or any(s in path for s in BLOCK_DIR_SUBSTR):
    print(f"✋ blocked write to sensitive: {path}", file=sys.stderr)
    sys.exit(2)

# rm recursive force guard — block any rm with -r and -f (in any order) unless confined to trees/
if tool.startswith("Bash(") and re.search(r'\brm\b', cmd) and re.search(r'\-(?:[^\s]*r[^\s]*f|[^\s]*f[^\s]*r)', cmd):
    parts = shlex.split(cmd)
    targets = [p for p in parts if not p.startswith("-") and p != "rm"]
    if not targets or not all(t.startswith("trees/") and ".." not in t for t in targets):
        print("✋ destructive rm blocked (restrict to trees/ or use git worktree remove)", file=sys.stderr)
        sys.exit(2)

# Commit/PR/tag/push gate (anywhere in command)
lower = cmd.lower()
if any(x in lower for x in ["git commit", "git tag", "git push", "gh pr create", "gh pr merge"]):
    print("✋ commits/tags/pushes/PR merges are human-only. Use /prepare-commit to stage & draft.", file=sys.stderr)
    sys.exit(2)

# Cloudflare Workers production deploy gate
if "wrangler deploy --env production" in cmd:
    if not os.path.exists(".claude/session/ALLOW_PROD_DEPLOY"):
        print("✋ production deploy blocked (set ALLOW_PROD_DEPLOY)", file=sys.stderr)
        print("Create the gate file: `mkdir -p .claude/session && touch .claude/session/ALLOW_PROD_DEPLOY`", file=sys.stderr)
        sys.exit(2)

# Supabase DB migrate/reset gate
if cmd.startswith("supabase db reset") or cmd.startswith("supabase db push"):
    if not os.path.exists(".claude/session/ALLOW_DB_MIGRATE"):
        print("✋ Supabase DB migration/reset blocked (set ALLOW_DB_MIGRATE)", file=sys.stderr)
        print("Create the gate file: `mkdir -p .claude/session && touch .claude/session/ALLOW_DB_MIGRATE`", file=sys.stderr)
        sys.exit(2)

# trigger.dev deploy gate
if cmd.startswith("npx trigger.dev deploy") or cmd.startswith("npx @trigger.dev/cli deploy"):
    if not os.path.exists(".claude/session/ALLOW_TRIGGER_DEPLOY"):
        print("✋ trigger.dev deploy blocked (set ALLOW_TRIGGER_DEPLOY)", file=sys.stderr)
        print("Create the gate file: `mkdir -p .claude/session && touch .claude/session/ALLOW_TRIGGER_DEPLOY`", file=sys.stderr)
        sys.exit(2)

# Flutter / Mobile release gate
if cmd.startswith("fastlane") or cmd.startswith("flutter build ipa") or cmd.startswith("flutter build appbundle"):
    if not os.path.exists(".claude/session/ALLOW_MOBILE_RELEASE"):
        print("✋ mobile release blocked (set ALLOW_MOBILE_RELEASE)", file=sys.stderr)
        print("Create the gate file: `mkdir -p .claude/session && touch .claude/session/ALLOW_MOBILE_RELEASE`", file=sys.stderr)
        sys.exit(2)

sys.exit(0)
