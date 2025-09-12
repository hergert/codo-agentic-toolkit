# Active Task Context

## Completed: Restructure dotclaude → pack

### Directory Structure Changes
- Renamed top-level `dotclaude/` to `pack/` with cleaner organization:
  - `pack/dotclaude/` - base content directly (no nested base folder)
  - `pack/stacks/` - stack overlays
  - `pack/docs/` - documentation
  - `pack/CLAUDE.md` - project instructions

### Code Updates
- **cli/cmd/init.go:58** - Updated local dev check from `os.Stat("dotclaude")` to `os.Stat("pack")`
- **cli/cmd/update.go:30** - Updated local dev check from `os.Stat("dotclaude")` to `os.Stat("pack")`
- **cli/internal/pack/fsloader.go:14-15** - Path constants: `baseRoot = "dotclaude"`, `stacksRoot = "stacks"`
- **cli/internal/pack/resolve.go:103** - Cache extraction: `extractDir := filepath.Join(cacheDir, "pack")`
- **.github/workflows/release.yml:67-71** - Package from `pack` directory in releases
- **.github/workflows/ci.yml:65** - Package `pack` directory for latest builds

### Build System
- **cli/justfile** - Consolidates build commands, handles temporary `.embedded_pack` copying
- **cli/embed.go:11** - Embeds pack via `//go:embed all:.embedded_pack`
- Build process: copy `../pack` → `.embedded_pack` → build → cleanup

### Testing Verified
- ✅ `codo init --offline` - Uses embedded pack successfully
- ✅ `codo init` with local `pack/` directory - Detects and uses local pack
- ✅ `codo update` - Falls back to embedded pack when download fails
- ✅ `codo remove --dry-run` - Correctly identifies all managed files

### Known Issue (Harmless)
- Compile-time warning: `pattern all:.embedded_pack: no matching files found`
- Expected since `.embedded_pack` only exists during build
- Does not affect runtime functionality