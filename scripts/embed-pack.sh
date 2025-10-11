#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PACK_DIR="$ROOT/pack"
TARGET="$ROOT/cli/internal/pack/embedded_base.zip"

if ! command -v python3 >/dev/null 2>&1; then
  echo "python3 is required to refresh embedded pack" >&2
  exit 1
fi

python3 - "$PACK_DIR" "$TARGET" <<'PY'
import os, sys, zipfile

pack_dir = sys.argv[1]
target = sys.argv[2]
os.makedirs(os.path.dirname(target), exist_ok=True)
with zipfile.ZipFile(target, 'w', zipfile.ZIP_DEFLATED) as z:
    for base, dirs, files in os.walk(pack_dir):
        rel_dir = os.path.relpath(base, pack_dir)
        if rel_dir != '.':
            zinfo = zipfile.ZipInfo(rel_dir.rstrip('/') + '/')
            zinfo.external_attr = 0o755 << 16
            z.writestr(zinfo, '')
        for name in files:
            full = os.path.join(base, name)
            rel = os.path.relpath(full, pack_dir)
            z.write(full, arcname=rel)
PY

echo "Embedded pack refreshed at $TARGET"
