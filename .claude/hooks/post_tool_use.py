#!/usr/bin/env python3
import json
import shutil
import subprocess
import sys
from pathlib import Path

E = json.load(sys.stdin)
TOOL_INPUT = E.get("tool_input") or {}
TARGET = TOOL_INPUT.get("file_path") or ""

STATE_DIR = Path(".claude") / "session"
HINT_CACHE = STATE_DIR / "format_hints.json"
_hint_state = None


def hint_once(key: str, message: str) -> None:
    global _hint_state
    STATE_DIR.mkdir(parents=True, exist_ok=True)
    if _hint_state is None:
        if HINT_CACHE.exists():
            try:
                _hint_state = json.loads(HINT_CACHE.read_text(encoding="utf-8"))
            except Exception:
                _hint_state = {}
        else:
            _hint_state = {}
    if _hint_state.get(key):
        return
    _hint_state[key] = True
    try:
        HINT_CACHE.write_text(json.dumps(_hint_state), encoding="utf-8")
    except Exception:
        pass
    print(json.dumps({"systemMessage": message}))


def run(cmd, cwd=None):
    result = subprocess.run(
        cmd,
        cwd=cwd,
        text=True,
        capture_output=True,
    )
    if result.stdout:
        sys.stdout.write(result.stdout)
    if result.stderr:
        sys.stderr.write(result.stderr)
    return result


def find_up(start: Path, target: str):
    cur = start
    for _ in range(10_000):
        candidate = cur / target
        if candidate.exists():
            return candidate
        if cur.parent == cur:
            break
        cur = cur.parent
    return None


def handle_go(path: Path):
    go = shutil.which("go")
    if not go:
        hint_once("go", "Install Go to enable go fmt/go build checks.")
        return
    goimports = shutil.which("goimports")
    if goimports:
        result = run([goimports, "-w", str(path)])
        if result.returncode != 0:
            hint_once("goimports-error", "goimports failed; inspect output above.")
    else:
        hint_once(
            "goimports",
            "Optional: install goimports (`go install golang.org/x/tools/cmd/goimports@latest`).",
        )
        run([go, "fmt", str(path)])
    pkg_dir = path.parent if path.parent != Path("") else Path(".")
    run([go, "build"], cwd=str(pkg_dir))


def handle_ts(path: Path):
    npx = shutil.which("npx")
    if not npx:
        hint_once("npx", "Install Node.js (provides `npx`) to format TS/JS files.")
        return
    prettier = run([npx, "--no-install", "prettier", "--write", str(path)])
    if (
        prettier.returncode != 0
        and "Cannot find module" in (prettier.stderr or "")
        and "prettier" in (prettier.stderr or "")
    ):
        hint_once(
            "prettier",
            "Prettier is missing; add it to devDependencies and run `pnpm|npm install`.",
        )
    eslint = run([npx, "--no-install", "eslint", "--fix", str(path)])
    if (
        eslint.returncode != 0
        and "Cannot find module" in (eslint.stderr or "")
        and "eslint" in (eslint.stderr or "")
    ):
        hint_once(
            "eslint",
            "ESLint is missing; add it to devDependencies and run `pnpm|npm install`.",
        )
    if path.suffix in {".ts", ".tsx"}:
        tsconfig = find_up(path.parent, "tsconfig.json")
        if tsconfig:
            tsc = run(
                [npx, "--no-install", "tsc", "--noEmit", "-p", str(tsconfig.parent)],
                cwd=str(tsconfig.parent),
            )
            if (
                tsc.returncode != 0
                and "Cannot find module" in (tsc.stderr or "")
                and "typescript" in (tsc.stderr or "")
            ):
                hint_once(
                    "tsc",
                    "TypeScript is missing; install `typescript` in devDependencies.",
                )


def handle_py(path: Path):
    ruff = shutil.which("ruff")
    if ruff:
        run([ruff, "check", "--fix", str(path)])
    else:
        hint_once("ruff", "Optional: install `ruff` for fast linting (`pip install ruff`).")
    black = shutil.which("black")
    if black:
        run([black, "-q", str(path)])
    else:
        hint_once("black", "Optional: install `black` for formatting (`pip install black`).")
    py3 = shutil.which("python3")
    if py3:
        run([py3, "-m", "py_compile", str(path)])
    else:
        hint_once("python3", "Install Python 3 to enable syntax checks (`python3 -m py_compile`).")


def handle_dart(path: Path):
    dart = shutil.which("dart")
    if not dart:
        hint_once("dart", "Install Dart SDK to enable format/analyze.")
        return
    run([dart, "format", str(path)])
    pubspec = find_up(path.parent, "pubspec.yaml")
    if pubspec:
        run([dart, "analyze"], cwd=str(pubspec.parent))


def main():
    target = Path(TARGET)
    if not target.exists():
        return
    suffix = target.suffix.lower()
    if suffix == ".go":
        handle_go(target)
    elif suffix in {".ts", ".tsx", ".js", ".jsx"}:
        handle_ts(target)
    elif suffix == ".py":
        handle_py(target)
    elif suffix == ".dart":
        handle_dart(target)

    subprocess.call(
        ["bash", "-lc", "git add -N . >/dev/null 2>&1 || true; git diff --shortstat || true"]
    )
    print("[format] done:", target)


if __name__ == "__main__":
    main()
