#!/usr/bin/env python3
import json, sys, os, shlex, re
from datetime import datetime, timedelta

E = json.load(sys.stdin)
tool = E.get("tool_name") or ""
ti = E.get("tool_input") or {}
cmd = ti.get("command", "")
path = ti.get("file_path") or ""

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


def has_fresh_plan(hours: int = 48) -> bool:
    specs_dir = os.path.join("docs", "specs")
    if not os.path.isdir(specs_dir):
        return False
    cutoff = datetime.utcnow() - timedelta(hours=hours)
    for name in os.listdir(specs_dir):
        if not name.endswith("-plan.md"):
            continue
        try:
            mtime = os.path.getmtime(os.path.join(specs_dir, name))
        except OSError:
            continue
        if datetime.utcfromtimestamp(mtime) >= cutoff:
            return True
    return False


def ask(reason: str) -> None:
    print(
        json.dumps(
            {
                "hookSpecificOutput": {
                    "hookEventName": "PreToolUse",
                    "permissionDecision": "ask",
                    "permissionDecisionReason": reason,
                }
            }
        )
    )
    sys.exit(0)


def deny(message: str) -> None:
    print(message, file=sys.stderr)
    sys.exit(2)


if tool == "Read":
    base = os.path.basename(path)
    if base in SENSITIVE_NAMES or any(sub in path for sub in SENSITIVE_DIR_SUBSTR):
        ask(f"Read of sensitive path: {path}")

# Always allow docs/.claude edits
if tool in ("Edit", "Write", "MultiEdit") and (
    path.startswith("docs/") or path.startswith(".claude/")
):
    sys.exit(0)

if tool in ("Edit", "Write", "MultiEdit") and not (
    path.startswith("docs/") or path.startswith(".claude/")
):
    if not has_fresh_plan():
        ask("No recent plan in docs/specs/*-plan.md (≤48h). Proceed only if this is a trivial fix.")

if tool in ("Edit", "Write", "MultiEdit"):
    base = os.path.basename(path)
    if base in SENSITIVE_NAMES or any(sub in path for sub in SENSITIVE_DIR_SUBSTR):
        deny(f"✋ blocked write to sensitive: {path}")

if tool.startswith("Bash(") and re.search(r"\brm\b", cmd) and re.search(
    r"\-(?:[^\s]*r[^\s]*f|[^\s]*f[^\s]*r)", cmd
):
    parts = shlex.split(cmd)
    targets = [p for p in parts if not p.startswith("-") and p != "rm"]
    if not targets or not all(t.startswith("trees/") and ".." not in t for t in targets):
        deny("✋ destructive rm blocked (restrict to trees/ or use git worktree remove)")

lower = cmd.lower()
if any(x in lower for x in ["git commit", "git tag", "git push", "gh pr create", "gh pr merge"]):
    deny("✋ commits/tags/pushes/PR merges are human-only. Use /prepare-commit.")

if (
    cmd.startswith("fastlane")
    or cmd.startswith("flutter build ipa")
    or cmd.startswith("flutter build appbundle")
) and not os.path.exists(".claude/session/ALLOW_MOBILE_RELEASE"):
    deny("✋ mobile release blocked (set ALLOW_MOBILE_RELEASE)")

sys.exit(0)
