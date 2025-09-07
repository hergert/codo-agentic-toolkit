# Codo — agentic toolkit installer (CLI)

Install the Codo CLI and manage the **dotclaude** pack in any repo.

## Install

### One-liner (macOS/Linux)
```bash
curl -fsSL https://raw.githubusercontent.com/hergert/codo-agentic-toolkit/main/dist/codo-install.sh | bash
```
This downloads the latest release from GitHub, verifies SHA256 checksums, and installs `codo` to `~/.local/bin`.

### Homebrew (optional)
```bash
brew tap hergert/homebrew-tap
brew install hergert/homebrew-tap/codo
```

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