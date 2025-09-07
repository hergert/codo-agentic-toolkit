#!/usr/bin/env bash
set -euo pipefail

# --- configure your repo (owner/org) and binary repo name ---
OWNER="${CODO_OWNER:-hergert}"
REPO="${CODO_REPO:-codo-agentic-toolkit}"
BIN_DIR="${CODO_BIN:-$HOME/.local/bin}"

# --- detect platform ---
os() {
  case "$(uname -s)" in
    Linux)  echo "linux" ;;
    Darwin) echo "darwin" ;;
    *) echo "unsupported OS: $(uname -s)" >&2; exit 1 ;;
  esac
}
arch() {
  case "$(uname -m)" in
    x86_64|amd64)  echo "amd64" ;;
    arm64|aarch64) echo "arm64" ;;
    *) echo "unsupported arch: $(uname -m)" >&2; exit 1 ;;
  esac
}

OS="$(os)"; ARCH="$(arch)"
CHECKSUMS="checksums.txt"
BASE="https://github.com/${OWNER}/${REPO}/releases/latest/download"

mkdir -p "$BIN_DIR"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT
cd "$TMP"

echo "→ fetching checksums..."
curl -fsSLO "$BASE/$CHECKSUMS"

# resolve the exact asset name for this OS/arch from checksums
ASSET="$(grep -Eo "codo_.*_${OS}_${ARCH}\.tar\.gz" "$CHECKSUMS" | head -n1)"
if [ -z "${ASSET:-}" ]; then
  echo "could not find a codo tarball for ${OS}/${ARCH} in $CHECKSUMS" >&2
  exit 1
fi

echo "→ downloading $ASSET ..."
curl -fsSLO "$BASE/$ASSET"

# verify checksum
echo "→ verifying checksum..."
if command -v sha256sum >/dev/null 2>&1; then
  grep "  $ASSET" "$CHECKSUMS" | sha256sum -c -
elif command -v shasum >/dev/null 2>&1; then
  grep "  $ASSET" "$CHECKSUMS" | shasum -a 256 -c -
else
  echo "(!) sha256sum/shasum not found; skipping checksum verification"
fi

# extract and install
echo "→ installing to $BIN_DIR/codo"
tar -xzf "$ASSET" codo
install -m 0755 codo "$BIN_DIR/codo"

# PATH hint
if ! command -v codo >/dev/null 2>&1; then
  case ":$PATH:" in
    *:"$BIN_DIR":*) : ;;
    *) echo "ℹ️  $BIN_DIR is not on PATH. Add:  export PATH=\"$BIN_DIR:\$PATH\"" ;;
  esac
fi

echo "✅ done"
"$BIN_DIR/codo" version || true