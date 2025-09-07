# Codo Installation System Refactor Plan

## Goal
Create a one-liner installation method (`curl | bash`) that downloads the latest codo binary from GitHub releases with checksum verification, fix the deprecated Homebrew configuration, and remove unnecessary embedded pack templates.

## Acceptance Checks
- [ ] `curl -fsSL .../dist/codo-install.sh | bash` installs codo to `~/.local/bin`
- [ ] Installer verifies SHA256 checksums before installation
- [ ] `brew tap youruser/homebrew-tap && brew install youruser/homebrew-tap/codo` works
- [ ] All embedded pack templates removed from `cli/internal/pack/`
- [ ] `codo init` uses only `dotclaude/` directory as source
- [ ] README shows clean, minimal install instructions
- [ ] GoReleaser generates `checksums.txt` with each release

## Files to Touch

1. **CREATE `dist/codo-install.sh`** - curl|bash installer script with checksum verification
2. **MODIFY `cli/.goreleaser.yaml`** - Switch from `brews` to `homebrew_casks`
3. **DELETE `cli/internal/pack/templates/`** - Remove all embedded pack files
4. **MODIFY `cli/internal/pack/pack.go`** - Remove `//go:embed` directives
5. **REPLACE `README.md`** - Minimal install docs with curl|bash primary method

## Diff Outline

### `dist/codo-install.sh` (NEW)
```bash
#!/usr/bin/env bash
set -euo pipefail

OWNER="${CODO_OWNER:-<youruser>}"
REPO="${CODO_REPO:-codo}"
BIN_DIR="${CODO_BIN:-$HOME/.local/bin}"

# Platform detection
os() { case "$(uname -s)" in Linux) echo linux;; Darwin) echo darwin;; esac }
arch() { case "$(uname -m)" in x86_64) echo amd64;; arm64) echo arm64;; esac }

# Main flow:
# 1. Download checksums.txt from /releases/latest/download/
# 2. Parse to find correct asset name for OS/arch
# 3. Download tarball
# 4. Verify with sha256sum -c or shasum -a 256 -c
# 5. Extract and install to $BIN_DIR
# 6. Check PATH and print hint if needed
```

### `cli/.goreleaser.yaml` (MODIFY)
```yaml
# REMOVE:
brews:
  - name: codo
    # ...

# ADD:
homebrew_casks:
  - name: codo
    repository:
      owner: youruser
      name: homebrew-tap
    directory: Casks
    install: |
      bin.install "codo"
```

### `cli/internal/pack/pack.go` (MODIFY)
```go
// REMOVE all //go:embed directives
// REMOVE embedded FS references
// Keep only FilesFromDotclaudeFS() using os.DirFS("dotclaude")
```

### `README.md` (REPLACE)
```markdown
# Codo — agentic toolkit installer (CLI)

Install the Codo CLI and manage the **dotclaude** pack in any repo.

## Install

### One-liner (macOS/Linux)
curl -fsSL https://raw.githubusercontent.com/<youruser>/codo-agentic-toolkit/main/dist/codo-install.sh | bash

### Homebrew (optional)
brew tap <youruser>/homebrew-tap
brew install <youruser>/homebrew-tap/codo

## Use
codo init                    # TUI wizard
codo init --stacks "go,cf"   # headless
codo update                  # safe update (conflicts → *.codo.new)
codo remove                  # backup to .codo-backup/<timestamp>/
codo self-update            # update CLI binary
```

## Risks

### Migration Risks
- **Embedded pack removal**: Could break existing installations if dotclaude/ missing
  - Detection: CI will fail on `codo init`
  - Mitigation: Keep dotclaude/ in repo root

### Performance Risks
- **Checksum verification**: Adds ~1s to install time
  - Detection: Time install script
  - Mitigation: Make verification optional with flag

### Breaking Changes
- **Homebrew formula → cask**: Existing brew users must reinstall
  - Detection: Old formula will 404
  - Mitigation: Document migration in release notes

## Test Strategy

### Contract Tests Only

1. **`test/install_test.sh`** - Test installer script
   ```bash
   # Test checksum verification works
   # Test PATH detection
   # Test error handling for bad checksums
   ```

2. **`cli/cmd/init_test.go`** - Verify dotclaude loading
   ```go
   func TestInitFromDotclaude(t *testing.T) {
       // Ensure FilesFromDotclaudeFS works
       // No embedded pack fallback
   }
   ```

3. **`cli/cmd/selfupdate_test.go`** - Add checksum verification
   ```go
   func TestSelfUpdateChecksum(t *testing.T) {
       // Download checksums.txt
       // Verify binary matches
   }
   ```

## Out of Scope (YAGNI)

- GPG signing (checksums sufficient)
- Remote pack fetching (use local dotclaude/)
- Auto-update mechanism (user controls updates)
- Windows installer (WSL works fine)
- Docker image (unnecessary complexity)
- Package managers beyond Homebrew