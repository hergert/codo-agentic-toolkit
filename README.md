# Codo — agentic toolkit installer (CLI)

Install the Codo CLI and manage the **dotclaude** pack in any repo.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/hergert/codo-agentic-toolkit/main/scripts/install.sh | bash
```
This downloads the latest release from GitHub, verifies SHA256 checksums, and installs `codo` to `~/.local/bin`.

### Shell completion
```bash
codo completion bash|zsh|fish|powershell > <your-shell-completions-dir>
```

## Use

```bash
# init the toolkit in a repo (TUI wizard)
codo init

# headless init (skip TUI)
codo init --stacks "cloudflare-workers,go"

# update (safe: only overwrites clean files; diverged files → *.codo.new)
codo update

# uninstall (backs up to .codo-backup/<timestamp>/)
codo remove

# self-update the CLI
codo self-update

# status & doctor
codo status
codo doctor

Run `codo doctor` after `codo init` to confirm hooks are executable and Python 3 is available.
```

## What gets installed

- `dotclaude/CLAUDE.md` → `CLAUDE.md`
- `dotclaude/.claude/base/**` → `.claude/**`
- Stack overlays from `dotclaude/.claude/stacks/<stack>/**` (if selected)

The toolkit provides:
- Subagents for mapping, tests, and review
- Commands for tight development loops
- Hooks as safety gates
- Settings for plan-by-default mode
- Output styles to switch personas

## Codo agentic loop (minimal)
1) `/scope "<key>"`      # map + plan in one file; stop for human review
2) `/prime "<key>"`      # context preflight (status, include list, keep ≥70% headroom)
3) `/build "<key>"`      # contract tests → smallest diff → green
4) `/review "<key>"`     # split-role review → APPROVE / REQUEST-CHANGES
(5) `/prepare-commit`    # stage & draft message (never commits; you commit & push)

Golden Rules still apply: YAGNI, KISS, smallest diff, contract tests only, and no commits or deploys by Claude. Claude will ask before every source edit—confirm in chat to proceed.
