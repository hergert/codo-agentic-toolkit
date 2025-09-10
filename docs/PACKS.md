# PACKS.md - Pack Overlay System

## Overview

The codo pack system uses a **base + overlay** architecture to compose Claude configurations for different technology stacks.

## How It Works

1. **Base Layer** (`dotclaude/.claude/base/**`)
   - Core Claude configuration files
   - Universal settings, hooks, and commands
   - Installed to `.claude/**` in your project

2. **Stack Overlays** (`dotclaude/.claude/stacks/<stack>/**`)
   - Technology-specific configurations
   - Extend or override base settings
   - Available stacks: `cloudflare-workers`, `go`, `python`, `nextjs`, etc.

3. **Overlay Resolution**
   - Base files are installed first
   - Stack files overlay on top (overlay wins on conflicts)
   - Multiple stacks can be combined

## Example

When you run:
```bash
codo init --stacks "cloudflare-workers,go"
```

The resolution order is:
1. Install all files from `dotclaude/.claude/base/**`
2. Overlay files from `dotclaude/.claude/stacks/cloudflare-workers/**`
3. Overlay files from `dotclaude/.claude/stacks/go/**`

If both base and a stack contain `.claude/settings.json`, the stack version wins.

## Conflict Resolution

During updates (`codo update`):
- **Clean files** (matching checksums) are overwritten
- **Modified files** are preserved, new version written to `*.codo.new`
- Review `.codo.new` files manually to merge changes

## Custom Packs

Future versions will support:
- Remote pack fetching from GitHub releases
- Custom pack URLs with SHA256 verification
- Local pack directories