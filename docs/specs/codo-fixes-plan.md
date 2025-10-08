# CODO Fixes & Improvements Plan

## Goal
Fix critical module path mismatches, self-update configuration errors, and operational inconsistencies. Add remote pack fetching with checksums, improve installer with SHA256 verification, and streamline documentation. Ensure release coordinates, self-update mechanism, and hook configurations work correctly without false positives or misleading documentation.

## Acceptance Checks
- [ ] `go build ./...` succeeds with unified module path
- [ ] `codo self-update` correctly detects and updates from GitHub releases  
- [ ] `curl | bash` installer exists at `dist/codo-install.sh` with SHA256 verification
- [ ] Remote pack fetch works: `codo init --pack-url <url>` downloads and verifies
- [ ] No duplicate pack directories (only `dotclaude/` remains)
- [ ] Settings descriptions match actual hook behaviors (commits/deploys)
- [ ] Pre-tool hooks don't block legitimate `firebase` subdirectories
- [ ] `codo status --strict` exits non-zero when drift exists
- [ ] README has concise install/use section with proper examples
- [ ] CI builds and releases work with correct module paths

## Files to Touch

### Critical Fixes
1. `cli/go.mod` - Update module path to `github.com/hergert/codo-agentic-toolkit/cli`
2. `cli/main.go` - Update import to new module path
3. `cli/.goreleaser.yaml` - Fix ldflags version path  
4. `cli/cmd/selfupdate.go` - Fix repo string and UpdateTo usage
5. `cli/internal/pack/` - Remove old embedded templates directory
6. `.claude/settings.json` - Fix note about commit/deploy policies
7. `dotclaude/.claude/base/settings.json` - Fix note about commit/deploy policies  
8. `.claude/hooks/pre_tool_use.py` - Tighten secrets guard patterns
9. `dotclaude/.claude/base/hooks/pre_tool_use.py` - Tighten secrets guard patterns

### New Features
10. `dist/codo-install.sh` - Create installer script with SHA256 verification (new file)
11. `cli/internal/pack/resolve.go` - Add remote pack fetching with checksums (new file)
12. `cli/cmd/status.go` - Add --strict flag for CI use
13. `cli/cmd/doctor.go` - Add Python version check
14. `docs/PACKS.md` - Document pack overlay system (new file)
15. `README.md` - Replace install/use section with concise version

### Cleanup
16. `.claude/settings.local.json` - Remove invalid Bash patterns
17. `dotclaude/.claude/base/settings.local.json` - Remove invalid Bash patterns

## Diff Outline

### Module path unification
```go
// cli/go.mod
- module github.com/hergert/codo
+ module github.com/hergert/codo-agentic-toolkit/cli

// cli/main.go  
- import "github.com/hergert/codo/cmd"
+ import "github.com/hergert/codo-agentic-toolkit/cli/cmd"

// cli/.goreleaser.yaml
ldflags:
-  - -X github.com/hergert/codo/cmd.version={{.Version}}
+  - -X github.com/hergert/codo-agentic-toolkit/cli/cmd.version={{.Version}}
```

### Self-update fix
```go
// cli/cmd/selfupdate.go
- latest, found, err := selfupdate.DetectLatest("hergert/codo")
+ latest, found, err := selfupdate.DetectLatest("hergert/codo-agentic-toolkit")

- if err := selfupdate.UpdateTo(latest.AssetURL, os.Args[0], ...); err != nil {
+ if err := selfupdate.UpdateTo(latest, os.Args[0]); err != nil {
```

### Remote pack resolver
```go
// cli/internal/pack/resolve.go (new)
func Resolve(tag string) (string, error) {
    baseURL := "https://github.com/hergert/codo-agentic-toolkit/releases"
    if tag == "latest" {
        tag = "latest/download"
    } else {
        tag = "download/" + tag
    }
    
    packURL := fmt.Sprintf("%s/%s/dotclaude-pack.zip", baseURL, tag)
    checksumURL := fmt.Sprintf("%s/%s/dotclaude-pack.sha256", baseURL, tag)
    
    // Download to ~/.codo/packs/<tag>/
    cacheDir := filepath.Join(os.Getenv("HOME"), ".codo", "packs", tag)
    // ... download, verify checksum, extract
    return filepath.Join(cacheDir, "dotclaude"), nil
}
```

### Status --strict
```go
// cli/cmd/status.go
var strictFlag bool

func init() {
    statusCmd.Flags().BoolVar(&strictFlag, "strict", false, "Exit non-zero if drift exists")
}

func runStatus(cmd *cobra.Command, args []string) error {
    // ... existing status logic
    if strictFlag && hasDrift {
        return fmt.Errorf("drift detected")
    }
    return nil
}
```

### Enhanced installer script
```bash
#!/usr/bin/env bash
# dist/codo-install.sh
set -euo pipefail

OWNER="hergert"
REPO="codo-agentic-toolkit"
BIN_DIR="${CODO_BIN:-$HOME/.local/bin}"

os() {
  case "$(uname -s)" in
    Linux) echo "linux" ;;
    Darwin) echo "darwin" ;;
    *) echo "unsupported OS"; exit 1 ;;
  esac
}
arch() {
  case "$(uname -m)" in
    x86_64|amd64) echo "amd64" ;;
    arm64|aarch64) echo "arm64" ;;
    *) echo "unsupported arch"; exit 1 ;;
  esac
}

OS="$(os)"; ARCH="$(arch)"
ASSET="codo_${OS}_${ARCH}.tar.gz"
CHECKSUMS="checksums.txt"

mkdir -p "$BIN_DIR"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT
cd "$TMP"

BASE="https://github.com/${OWNER}/${REPO}/releases/latest/download"
curl -fsSLO "$BASE/$ASSET"
curl -fsSLO "$BASE/$CHECKSUMS"

if command -v shasum >/dev/null 2>&1; then
  grep "  $ASSET" "$CHECKSUMS" | shasum -a 256 -c -
elif command -v sha256sum >/dev/null 2>&1; then
  grep "  $ASSET" "$CHECKSUMS" | sha256sum -c -
else
  echo "(!) sha256 tool not found; skipping checksum verification"
fi

tar -xzf "$ASSET" codo
install -m 0755 codo "$BIN_DIR/codo"
echo "✅ codo installed to $BIN_DIR/codo"
"$BIN_DIR/codo" version || true
```

### Settings fixes
```json
// .claude/settings.json & dotclaude/.claude/base/settings.json
- "note": "Commits/PR merges/prod deploys blocked. Exit gate: green tests"
+ "note": "Source edits require explicit confirmation (ask). Commits/tags/PR merges are denied (human only). High-risk ops gated individually (ALLOW_PROD_DEPLOY / ALLOW_DB_MIGRATE / ALLOW_MOBILE_RELEASE / ALLOW_TRIGGER_DEPLOY). Exit gate: tests green."
```

### Secrets guard improvements
```python
# pre_tool_use.py (both locations)
bn = os.path.basename(path)
BLOCK_BASENAMES = {".env", ".env.local", ".env.production", ".env.development",
                   "package-lock.json","pnpm-lock.yaml","yarn.lock",
                   "id_rsa","id_ed25519","known_hosts",
                   "serviceAccountKey.json","GoogleService-Info.plist","google-services.json"}
BLOCK_DIR_SUBSTR = ("/.git/", "/config/secrets/")
if bn in BLOCK_BASENAMES or any(s in path for s in BLOCK_DIR_SUBSTR):
    # block logic
```

### Doctor Python check
```go
// cli/cmd/doctor.go
func checkPython() error {
    cmd := exec.Command("python3", "--version")
    output, err := cmd.Output()
    if err != nil {
        return fmt.Errorf("Python 3 not found (required for hooks)")
    }
    // Parse version, warn if < 3.6
    fmt.Printf("✓ Python: %s", output)
    return nil
}
```

### README replacement
```markdown
## Install

### One-liner (macOS/Linux)
\`\`\`bash
curl -fsSL https://raw.githubusercontent.com/hergert/codo-agentic-toolkit/main/dist/codo-install.sh | bash
# installs codo to ~/.local/bin; set CODO_BIN to override
\`\`\`

### Homebrew (optional)
\`\`\`bash
brew tap hergert/homebrew-tap
brew install hergert/homebrew-tap/codo
\`\`\`

### Shell completion
\`\`\`bash
codo completion bash|zsh|fish|powershell > <your-shell-completions-dir>
\`\`\`

## Use
\`\`\`bash
# TUI wizard
codo init

# headless (skip TUI)
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
\`\`\`

### What gets installed
- `dotclaude/CLAUDE.md` → `CLAUDE.md`
- `dotclaude/.claude/base/**` → `.claude/**`
- overlays from `dotclaude/.claude/stacks/<stack>/**` (overlay wins on conflicts)
```

## Risks
- **Import breakage**: Changing module paths requires updating all internal imports. Detection: `go build ./...` will fail.
- **Release asset mismatch**: GoReleaser config must match installer expectations. Detection: Test installer after first release.
- **Self-update compatibility**: Existing installations may fail to update. Detection: Test update from current version.
- **Hook behavior change**: Tightened patterns may change blocking behavior. Detection: Manual test with firebase paths.
- **Remote pack security**: Downloading packs from internet requires checksums. Detection: Verify SHA256 matches.
- **Python version**: Hooks may fail on Python 2. Detection: `codo doctor` will warn.

## Test Strategy
- `cli/cmd/selfupdate_test.go` - Test DetectLatest with correct repo
- `cli/internal/pack/resolve_test.go` - Test remote pack download and checksum verification
- `cli/cmd/status_test.go` - Test --strict flag behavior with drift
- `cli/internal/fsops/fsops_test.go` - Verify CopySafe conflict handling
- Manual tests:
  - Build with new module path: `cd cli && go build`
  - Test installer: `bash dist/codo-install.sh`  
  - Test self-update flow with mock release
  - Test remote pack fetch with checksum
  - Verify hooks don't block legitimate firebase directories
  - Test `codo status --strict` exit code in CI

## Out of Scope
- Adding telemetry or analytics
- Changing pack content structure beyond necessary fixes
- Creating extensive documentation beyond PACKS.md
- Adding update preview TUI (defer to next iteration)
- Switching to different self-update library (current fix is sufficient)
- Implementing Homebrew tap (mentioned in README but not implemented here)
- Adding features unrelated to the identified fixes
