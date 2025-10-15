#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

if ! command -v go >/dev/null 2>&1; then
  echo "go toolchain is required to refresh embedded pack" >&2
  exit 1
fi

cd "$ROOT/cli/internal/pack"
go generate ./...

echo "Embedded pack refreshed at cli/internal/pack/embedded_base_gen.go"
