#!/usr/bin/env bash
set -euo pipefail
BIN="${1:?pass compiled macOS binary path}"
VERSION="${VERSION:-0.1.0}"
ROOT="$(mktemp -d)"
trap 'rm -rf "$ROOT"' EXIT
APP="$ROOT/Applications/DevOrchestra.app/Contents"
mkdir -p "$APP/MacOS" dist
install -m 755 "$BIN" "$APP/MacOS/DevOrchestra"
cat > "$APP/Info.plist" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"><dict>
<key>CFBundleName</key><string>DevOrchestra</string>
<key>CFBundleExecutable</key><string>DevOrchestra</string>
<key>CFBundleIdentifier</key><string>io.github.jahangard.devorchestra</string>
<key>CFBundlePackageType</key><string>APPL</string>
<key>CFBundleVersion</key><string>$VERSION</string>
<key>CFBundleShortVersionString</key><string>$VERSION</string>
</dict></plist>
EOF
pkgbuild --root "$ROOT" --identifier io.github.jahangard.devorchestra --version "$VERSION" --install-location / "dist/DevOrchestra-${VERSION}.pkg"
