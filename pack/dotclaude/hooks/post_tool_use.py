#!/usr/bin/env python3
import json, sys, subprocess, shlex, os

E = json.load(sys.stdin)
f = (E.get("tool_input") or {}).get("file_path", "") or ""

if f and os.path.exists(f):
    # Format only the changed file when possible
    if f.endswith(".go"):
        subprocess.call([
            "bash",
            "-lc",
            f"command -v goimports >/dev/null 2>&1 && goimports -w {shlex.quote(f)} || go fmt {shlex.quote(f)}"
        ])
    elif f.endswith((".ts", ".tsx", ".js", ".jsx")):
        subprocess.call(["bash","-lc", f"command -v npx >/dev/null 2>&1 && npx -y prettier --write {shlex.quote(f)} || true"])
        subprocess.call(["bash","-lc", f"command -v npx >/dev/null 2>&1 && npx -y eslint --fix {shlex.quote(f)} || true"])
    elif f.endswith(".py"):
        subprocess.call(["bash","-lc", f"command -v ruff >/dev/null 2>&1 && ruff check --fix {shlex.quote(f)} || true"])
        subprocess.call(["bash","-lc", f"command -v black >/dev/null 2>&1 && black -q {shlex.quote(f)} || true"])
    elif f.endswith(".dart"):
        subprocess.call(["bash","-lc", f"command -v dart >/dev/null 2>&1 && dart format {shlex.quote(f)} || true"])

# Show concise diff stat for situational awareness
subprocess.call(["bash","-lc", "git add -N . >/dev/null 2>&1 || true; git diff --shortstat || true"])
print("[format] done:", f)
