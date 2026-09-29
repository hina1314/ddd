#!/usr/bin/env bash
set -euo pipefail

release_tag="${1:-}"
if [[ ! "$release_tag" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "usage: $0 vMAJOR.MINOR.PATCH" >&2
  exit 1
fi

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$script_dir/.." && pwd)"
extract_root="$(mktemp -d)"
trap 'rm -rf "$extract_root"' EXIT

verify_package() {
  local package="$1" extension="$2" stage="$extract_root/$1"
  test -s "$stage/app-api$extension"
  test -f "$stage/app.env.example"
  test -f "$stage/README.md"
  test -f "$stage/ops/docker/compose.yml"
  test -f "$stage/config/i18n/en.yml"

  # Resolve the migration mount exactly as the bundled Compose file does.
  local migration_dir
  migration_dir="$(cd "$stage/ops/docker/../../db/migration" && pwd)"
  local source
  for source in "$repo_root"/db/migration/*.sql; do
    test -f "$source"
    cmp "$source" "$migration_dir/$(basename "$source")"
  done
  echo "Verified release package: $package"
}

linux_package="app-api_${release_tag}_linux_amd64"
windows_package="app-api_${release_tag}_windows_amd64"
tar -xzf "$repo_root/dist/$linux_package.tar.gz" -C "$extract_root"
unzip -q "$repo_root/dist/$windows_package.zip" -d "$extract_root"
verify_package "$linux_package" ""
verify_package "$windows_package" ".exe"
(cd "$repo_root/dist" && sha256sum -c checksums.txt)
