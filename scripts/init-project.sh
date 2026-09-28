#!/usr/bin/env bash
set -euo pipefail

module="${1:-}"
app_name="${2:-}"
kit_version="${3:-}"
kit_module="github.com/hina1314/kit"
if [[ ! "$module" =~ ^[A-Za-z0-9._~-]+(/[A-Za-z0-9._~-]+)+$ ]] ||
   [[ ! "$app_name" =~ ^[a-z0-9][a-z0-9-]*$ ]]; then
  echo "usage: $0 module/path app-name [kit-version]" >&2
  exit 1
fi
if [[ -n "$kit_version" && ! "$kit_version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?$ ]]; then
  echo "kit-version must be a semantic version or a Go pseudo-version" >&2
  exit 1
fi

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$repo_root"
current_module="$(sed -n 's/^module[[:space:]]\+//p' go.mod)"
escaped_module="$(printf '%s' "$current_module" | sed 's/[.[\*^$\\]/\\&/g')"

while IFS= read -r -d '' file; do
  sed -i "s|\"${escaped_module}/|\"${module}/|g" "$file"
done < <(find . \( -path ./.git -o -path ./.local -o -path ./build -o -path ./dist \) -prune -o -name '*.go' -type f -print0)

go mod edit -module="$module"
go mod edit -dropreplace="$kit_module"
if [[ -n "$kit_version" ]]; then
  go get "$kit_module@$kit_version"
fi
for file in app.env.example ops/docker/app.env.example; do
  sed -i "s/^APP_NAME=.*/APP_NAME=${app_name}/" "$file"
  sed -i "s/^TOKEN_ISSUER=.*/TOKEN_ISSUER=${app_name}/" "$file"
done
go mod tidy
go run github.com/google/wire/cmd/wire ./internal/di
go test ./...

echo "Initialized ${app_name} with module ${module}"
echo "Using github.com/hina1314/kit at the version pinned in go.mod."
