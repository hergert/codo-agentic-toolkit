#!/usr/bin/env python3
import json, sys, os, shlex

E = json.load(sys.stdin)
tool = E.get("tool_name") or ""
ti = E.get("tool_input") or {}
cmd = ti.get("command", "")
path = (ti.get("file_path") or "")

# Allow docs/.claude edits always
if tool in ("Edit","Write","MultiEdit") and (path.startswith("docs/") or path.startswith(".claude/")):
    sys.exit(0)

# Fast Mode or plan-approved gate for source edits
fast  = os.path.exists(".claude/session/FAST_MODE")
approved = os.path.exists(".claude/session/ALLOW_EDITS")
if tool in ("Edit","Write","MultiEdit") and not (fast or approved):
    print("✋ edits blocked: /approve <task> or /fast-on", file=sys.stderr)
    sys.exit(2)

# Secrets / sensitive paths guard
BLOCK = (
    ".env", "/.env", ".env.", "/.env.", "/.env.local", "/.env.production", "/.env.development",
    "/.git/", "package-lock.json", "pnpm-lock.yaml", "yarn.lock",
    "id_rsa", "id_ed25519", "known_hosts",
    "secrets.", "secret.", "credentials.", "config/secrets",
    "google-services.json", "GoogleService-Info.plist", "serviceAccountKey.json", "firebase"
)
if any(b in path for b in BLOCK):
    print(f"✋ blocked write to sensitive: {path}", file=sys.stderr)
    sys.exit(2)

# rm -rf guard (allow only inside trees/ for ALL targets)
if tool.startswith("Bash(") and "rm -rf" in cmd:
    parts = shlex.split(cmd)
    targets = [p for p in parts if not p.startswith("-") and p not in ("rm",)]
    if not targets or not all(t.startswith("trees/") and ".." not in t for t in targets):
        print("✋ rm -rf blocked (restrict to trees/ or use git worktree remove)", file=sys.stderr)
        sys.exit(2)

# Commit/PR/tag gate
if tool.startswith("Bash(git commit") or tool.startswith("Bash(gh pr ") or tool.startswith("Bash(git tag"):
    if not (fast or os.path.exists(".claude/session/ALLOW_COMMITS")):
        print("✋ commits/tags/PR merges are human-only (use /fast-on commits or set ALLOW_COMMITS)", file=sys.stderr)
        sys.exit(2)

# Cloudflare Workers production deploy gate
if "wrangler deploy --env production" in cmd:
    if not os.path.exists(".claude/session/ALLOW_PROD_DEPLOY"):
        print("✋ production deploy blocked (set ALLOW_PROD_DEPLOY)", file=sys.stderr)
        sys.exit(2)

# Supabase DB migrate/reset gate
if cmd.startswith("supabase db reset") or cmd.startswith("supabase db push"):
    if not os.path.exists(".claude/session/ALLOW_DB_MIGRATE"):
        print("✋ Supabase DB migration/reset blocked (set ALLOW_DB_MIGRATE)", file=sys.stderr)
        sys.exit(2)

# trigger.dev deploy gate
if cmd.startswith("npx trigger.dev deploy") or cmd.startswith("npx @trigger.dev/cli deploy"):
    if not os.path.exists(".claude/session/ALLOW_TRIGGER_DEPLOY"):
        print("✋ trigger.dev deploy blocked (set ALLOW_TRIGGER_DEPLOY)", file=sys.stderr)
        sys.exit(2)

# Flutter / Mobile release gate
if cmd.startswith("fastlane") or cmd.startswith("flutter build ipa") or cmd.startswith("flutter build appbundle"):
    if not os.path.exists(".claude/session/ALLOW_MOBILE_RELEASE"):
        print("✋ mobile release blocked (set ALLOW_MOBILE_RELEASE)", file=sys.stderr)
        sys.exit(2)

sys.exit(0)
