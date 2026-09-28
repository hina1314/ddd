#!/usr/bin/env bash
set -euo pipefail

module="${1:-}"
app_name="${2:-}"
if [[ ! "$module" =~ ^[A-Za-z0-9._~-]+(/[A-Za-z0-9._~-]+)+$ ]] ||
   [[ ! "$app_name" =~ ^[a-z0-9][a-z0-9-]*$ ]]; then
  echo "usage: $0 module/path app-name" >&2
  exit 1
fi

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$repo_root"
current_module="$(sed -n 's/^module[[:space:]]\+//p' go.mod)"

while IFS= read -r -d '' file; do
  sed -i "s|\"${current_module}/|\"${module}/|g" "$file"
done < <(find . -path ./.git -prune -o -name '*.go' -type f -print0)

go mod edit -module="$module"
for file in app.env.example ops/docker/app.env.example; do
  sed -i "s/^APP_NAME=.*/APP_NAME=${app_name}/" "$file"
  sed -i "s/^TOKEN_ISSUER=.*/TOKEN_ISSUER=${app_name}/" "$file"
done
go mod tidy
go run github.com/google/wire/cmd/wire@v0.6.0 ./internal/di
go test ./...

echo "Initialized ${app_name} with module ${module}"
