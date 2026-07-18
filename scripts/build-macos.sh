#!/usr/bin/env bash
set -euo pipefail

project_dir="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
app_bundle="$project_dir/build/bin/Codex Canvas.app"
contents_dir="$app_bundle/Contents"
macos_dir="$contents_dir/MacOS"
resources_dir="$contents_dir/Resources"
universal_work_dir="$(mktemp -d "${TMPDIR:-/tmp}/codex-canvas-universal.XXXXXX")"

cleanup_universal_work_dir() {
  rm -rf "$universal_work_dir"
}
trap cleanup_universal_work_dir EXIT

mkdir -p "$macos_dir" "$resources_dir"

cd "$project_dir/frontend"
npm run build

cd "$project_dir"
build_architecture() {
  local go_arch="$1"
  local clang_arch="$2"
  local output_path="$universal_work_dir/Codex Canvas-$clang_arch"

  echo "Building macOS $clang_arch..."
  env \
    GOOS=darwin \
    GOARCH="$go_arch" \
    CGO_ENABLED=1 \
    CC="clang -arch $clang_arch" \
    CGO_CFLAGS="-arch $clang_arch -mmacosx-version-min=12.0" \
    CGO_LDFLAGS="-arch $clang_arch -framework UniformTypeIdentifiers -mmacosx-version-min=12.0" \
    MACOSX_DEPLOYMENT_TARGET=12.0 \
    GOCACHE="${TMPDIR:-/tmp}/codex-canvas-go-cache-$go_arch" \
    go build \
      -buildvcs=false \
      -tags "desktop,wv2runtime.download,production" \
      -ldflags "-w -s" \
      -o "$output_path" .
}

build_architecture amd64 x86_64
build_architecture arm64 arm64

lipo -create \
  "$universal_work_dir/Codex Canvas-x86_64" \
  "$universal_work_dir/Codex Canvas-arm64" \
  -output "$macos_dir/Codex Canvas"

cp "$project_dir/build/darwin/CodexCanvas.plist" "$contents_dir/Info.plist"

env GOCACHE="${TMPDIR:-/tmp}/codex-canvas-go-cache-tools" \
  go run ./cmd/iconpack "$project_dir/build/appicon.png" "$resources_dir/iconfile.icns"

codesign --force --deep --sign - "$app_bundle" >/dev/null
lipo -info "$macos_dir/Codex Canvas"
echo "$app_bundle"
