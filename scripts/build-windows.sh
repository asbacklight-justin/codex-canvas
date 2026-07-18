#!/usr/bin/env bash
set -euo pipefail

project_dir="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
output_exe="$project_dir/build/bin/Codex-Canvas.exe"
output_zip="$project_dir/build/bin/Codex-Canvas-1.1.0-Windows-x64.zip"
release_readme="$project_dir/build/windows/README-Windows.txt"

if command -v wails >/dev/null 2>&1; then
  wails_command="$(command -v wails)"
else
  wails_command="$(go env GOPATH)/bin/wails"
fi

if [[ ! -x "$wails_command" ]]; then
  echo "Wails CLI is required: go install github.com/wailsapp/wails/v2/cmd/wails@v2.11.0" >&2
  exit 1
fi

cd "$project_dir/frontend"
npm run build

cd "$project_dir"
env GOCACHE="${TMPDIR:-/tmp}/codex-canvas-go-cache-windows" \
  "$wails_command" build \
    -platform windows/amd64 \
    -o Codex-Canvas.exe \
    -nopackage \
    -m \
    -nosyncgomod \
    -skipbindings \
    -s \
    -trimpath

rm -f "$output_zip"
zip -j "$output_zip" "$output_exe" "$release_readme" >/dev/null
echo "$output_zip"
