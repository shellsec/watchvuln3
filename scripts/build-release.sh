#!/usr/bin/env bash
# 交叉编译 Windows / Linux / macOS 发布包
# 用法: ./scripts/build-release.sh [version]
set -euo pipefail

VERSION="${1:-v3.1.0}"
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
DIST="$ROOT/dist"
LDFLAGS="-s -w -X main.Version=${VERSION}"

mkdir -p "$DIST"
cd "$ROOT"

build() {
  local goos="$1" goarch="$2" out="$3"
  echo "==> ${goos}/${goarch} -> dist/${out}"
  CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" \
    go build -trimpath -ldflags "$LDFLAGS" -o "$DIST/$out" .
}

build windows amd64 watchvuln-windows-amd64.exe
build linux   amd64 watchvuln-linux-amd64
build linux   arm64 watchvuln-linux-arm64
build darwin  amd64 watchvuln-darwin-amd64
build darwin  arm64 watchvuln-darwin-arm64

echo ""
echo "Done. Artifacts in dist/:"
ls -lh "$DIST"
