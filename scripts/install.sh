#!/usr/bin/env bash
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

OS="$(os)"
ARCH="$(arch)"
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
if ! command -v codo >/dev/null 2>&1; then
  echo ""
  echo "Add to PATH (bash/zsh):"
  echo "  export PATH=\"\$PATH:$BIN_DIR\""
fi
