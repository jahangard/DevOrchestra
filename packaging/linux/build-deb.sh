#!/usr/bin/env bash
set -euo pipefail
BIN="${1:?pass compiled Linux binary path}"
VERSION="${VERSION:-0.1.0}"
ARCH="${ARCH:-amd64}"
ROOT="$(mktemp -d)"
trap 'rm -rf "$ROOT"' EXIT
mkdir -p "$ROOT/DEBIAN" "$ROOT/usr/bin" "$ROOT/usr/share/applications" dist
install -m 755 "$BIN" "$ROOT/usr/bin/devorchestra"
cat > "$ROOT/DEBIAN/control" <<EOF
Package: devorchestra
Version: $VERSION
Section: devel
Priority: optional
Architecture: $ARCH
Maintainer: DevOrchestra contributors
Description: Local-first GUI companion for Codex CLI
EOF
cat > "$ROOT/usr/share/applications/devorchestra.desktop" <<EOF
[Desktop Entry]
Name=DevOrchestra
Comment=Local-first Codex CLI workspace
Exec=/usr/bin/devorchestra
Terminal=false
Type=Application
Categories=Development;IDE;
EOF
dpkg-deb --build "$ROOT" "dist/devorchestra_${VERSION}_${ARCH}.deb"
