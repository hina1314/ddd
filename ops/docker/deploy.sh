#!/usr/bin/env bash
set -euo pipefail

# Install this file as root-owned /usr/local/sbin/study-deploy. Do not execute
# the copy in the checkout with sudo: the checkout is writable by the SSH user.
repo=/opt/apps/ddd
compose_dir="$repo/ops/docker"
image_repo=ghcr.io/hina1314/ddd

if (( $# != 1 )) || [[ ! $1 =~ ^[0-9a-f]{40}$ ]]; then
  echo "usage: study-deploy <full-main-commit-sha>" >&2
  exit 2
fi
if (( EUID != 0 )); then
  echo "study-deploy must run as root" >&2
  exit 2
fi

sha=$1
image="$image_repo:sha-$sha"
cd "$compose_dir"

exec 9>/var/lock/study-deploy.lock
flock -n 9 || { echo "another deployment is running" >&2; exit 1; }

git_safe=(git -c "safe.directory=$repo" -C "$repo")
checkout_sha=$("${git_safe[@]}" rev-parse HEAD)
checkout_changes=$("${git_safe[@]}" status --porcelain --untracked-files=all)
if [[ $checkout_sha != "$sha" || -n $checkout_changes ]]; then
  echo "checkout is not the requested clean commit $sha" >&2
  exit 1
fi

env_file="$compose_dir/.env"
if [[ ! -f $env_file ]] || [[ $(grep -c '^STUDY_IMAGE=' "$env_file") != 1 ]]; then
  echo "expected exactly one STUDY_IMAGE in $env_file" >&2
  exit 1
fi
old_image=$(sed -n 's/^STUDY_IMAGE=//p' "$env_file")
if [[ ! $old_image =~ ^[a-zA-Z0-9._/@:-]+$ ]]; then
  echo "current image is invalid" >&2
  exit 1
fi

set_image() {
  local next_image=$1 temp
  temp=$(mktemp "$compose_dir/.env.deploy.XXXXXXXX")
  if ! sed "s|^STUDY_IMAGE=.*|STUDY_IMAGE=$next_image|" "$env_file" > "$temp"; then
    rm -f -- "$temp"
    return 1
  fi
  chown --reference="$env_file" "$temp"
  chmod --reference="$env_file" "$temp"
  mv -f -- "$temp" "$env_file"
}

port=$(sed -n 's/^API_HOST_PORT=//p' "$env_file" | tail -n 1)
port=${port:-3000}
if [[ ! $port =~ ^[0-9]{1,5}$ ]] || (( 10#$port < 1 || 10#$port > 65535 )); then
  echo "invalid API_HOST_PORT: $port" >&2
  exit 1
fi

echo "pulling $image"
docker pull "$image"

# Run the SQL files from the same checkout as the image. An already-exited
# Compose migrate service would otherwise be reused for later deployments.
echo "running migrations for $sha"
docker compose -f "$compose_dir/compose.yml" run --rm -T migrate

set_image "$image"
if ! docker compose -f "$compose_dir/compose.yml" up -d --no-deps api; then
  deploy_failed=1
else
  deploy_failed=0
  for attempt in $(seq 1 30); do
    headers=$(curl --noproxy '*' -sS -D - -o /dev/null --max-time 2 \
      "http://127.0.0.1:$port/readyz" 2>/dev/null | tr -d '\r' || true)
    if grep -q '^HTTP/1.1 200 ' <<< "$headers" &&
       grep -qi "^X-App-Version: sha-$sha$" <<< "$headers"; then
      echo "deployed $image; ready on 127.0.0.1:$port"
      exit 0
    fi
    sleep 1
  done
  deploy_failed=1
fi

if (( deploy_failed )); then
  echo "new API failed readiness; restoring $old_image" >&2
  set_image "$old_image"
  docker compose -f "$compose_dir/compose.yml" up -d --no-deps api || true
  exit 1
fi
