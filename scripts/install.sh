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
CHECKSUMS="checksums.txt"
VERSION="${CODO_VERSION:-latest}"

mkdir -p "$BIN_DIR"

install_from_release() {
  local tmp base asset
  tmp="$(mktemp -d)"
  if [ "$VERSION" = "latest" ]; then
    base="https://github.com/${OWNER}/${REPO}/releases/latest/download"
  else
    base="https://github.com/${OWNER}/${REPO}/releases/download/${VERSION}"
  fi
  (
    set -euo pipefail
    trap 'rm -rf "$tmp"' EXIT
    cd "$tmp"
    curl -fsSLO "$base/$CHECKSUMS"
    pattern="codo_[^[:space:]]*_${OS}_${ARCH}\.tar\.gz"
    asset="$(grep -Eo "$pattern" "$CHECKSUMS" | head -n1 || true)"
    if [ -z "$asset" ]; then
      echo "Release asset for ${OS}/${ARCH} not found in $CHECKSUMS"
      exit 1
    fi
    curl -fsSLO "$base/$asset"
    if command -v shasum >/dev/null 2>&1; then
      grep "  $asset" "$CHECKSUMS" | shasum -a 256 -c -
    elif command -v sha256sum >/dev/null 2>&1; then
      grep "  $asset" "$CHECKSUMS" | sha256sum -c -
    else
      echo "(!) sha256 tool not found; skipping checksum verification"
    fi
    tar -xzf "$asset" codo
    install -m 0755 codo "$BIN_DIR/codo"
  )
}

install_from_source() {
  command -v go >/dev/null 2>&1 || {
    echo "Go toolchain not found; cannot build from source" >&2
    return 1
  }

  local tmp tarball
  tmp="$(mktemp -d)"
  tarball="https://github.com/${OWNER}/${REPO}/archive/refs/heads/main.tar.gz"
  (
    set -euo pipefail
    trap 'rm -rf "$tmp"' EXIT
    cd "$tmp"
    curl -fsSL "$tarball" | tar -xz --strip-components=1
    cd cli
    go build -o "$BIN_DIR/codo" .
  )
}

if install_from_release; then
  echo "✅ codo installed to $BIN_DIR/codo"
else
  echo "Falling back to source build..."
  install_from_source
  echo "✅ codo built from source at $BIN_DIR/codo"
fi

"$BIN_DIR/codo" version || true
if ! command -v codo >/dev/null 2>&1; then
  echo ""
  echo "Add to PATH (bash/zsh):"
  echo "  export PATH=\"\$PATH:$BIN_DIR\""
fi
