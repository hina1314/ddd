#!/usr/bin/env bash

set -euo pipefail

release_tag="${1:-}"
if [[ ! "$release_tag" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "usage: $0 vMAJOR.MINOR.PATCH" >&2
  exit 1
fi

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$script_dir/.." && pwd)"
dist_dir="$repo_root/dist"

if [[ -e "$dist_dir" ]]; then
  echo "refusing to overwrite existing directory: $dist_dir" >&2
  exit 1
fi

mkdir -p "$dist_dir"

stage_root="$(mktemp -d)"
trap 'rm -rf "$stage_root"' EXIT

build_target() {
  local goos="$1"
  local goarch="$2"
  local extension="$3"
  local package="study-api_${release_tag}_${goos}_${goarch}"
  local stage="$stage_root/$package"

  mkdir -p "$stage/config"

  CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" \
    go build -trimpath -ldflags="-s -w" \
    -o "$stage/study-api${extension}" "$repo_root"

  cp -R "$repo_root/config/i18n" "$stage/config/i18n"
  cp "$repo_root/app.env.example" "$stage/app.env.example"
  cp "$repo_root/readme.md" "$stage/README.md"
  cp -R "$repo_root/ops" "$stage/ops"

  if [[ "$goos" == "windows" ]]; then
    if command -v zip >/dev/null 2>&1; then
      (cd "$stage_root" && zip -qr "$dist_dir/${package}.zip" "$package")
    elif command -v powershell.exe >/dev/null 2>&1 && command -v cygpath >/dev/null 2>&1; then
      STAGE_WINDOWS="$(cygpath -w "$stage")" \
        ARCHIVE_WINDOWS="$(cygpath -w "$dist_dir/${package}.zip")" \
        powershell.exe -NoProfile -NonInteractive -Command \
        'Compress-Archive -Path $env:STAGE_WINDOWS -DestinationPath $env:ARCHIVE_WINDOWS -Force'
    else
      echo "zip or PowerShell Compress-Archive is required" >&2
      exit 1
    fi
  else
    (cd "$stage_root" && tar -czf "$dist_dir/${package}.tar.gz" "$package")
  fi
}

build_target linux amd64 ""
build_target windows amd64 ".exe"

(cd "$dist_dir" && sha256sum *.tar.gz *.zip > checksums.txt)

echo "release artifacts created in: $dist_dir"
