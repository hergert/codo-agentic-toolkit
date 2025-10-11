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

EPOCH = (1980, 1, 1, 0, 0, 0)
FILE_MODE = 0o100644 << 16
DIR_MODE = 0o040755 << 16

pack_dir = sys.argv[1]
target = sys.argv[2]
os.makedirs(os.path.dirname(target), exist_ok=True)
with zipfile.ZipFile(target, 'w', zipfile.ZIP_DEFLATED) as z:
    for base, dirs, files in os.walk(pack_dir):
        dirs.sort()
        files.sort()
        rel_dir = os.path.relpath(base, pack_dir)
        if rel_dir != '.':
            zinfo = zipfile.ZipInfo(rel_dir.rstrip('/') + '/')
            zinfo.date_time = EPOCH
            zinfo.external_attr = DIR_MODE
            z.writestr(zinfo, b'')
        for name in files:
            full = os.path.join(base, name)
            rel = os.path.relpath(full, pack_dir)
            zinfo = zipfile.ZipInfo(rel)
            zinfo.date_time = EPOCH
            zinfo.external_attr = FILE_MODE
            with open(full, 'rb') as fh:
                z.writestr(zinfo, fh.read())
PY

echo "Embedded pack refreshed at $TARGET"
