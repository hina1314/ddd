#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"

if [[ "$EUID" -ne 0 ]]; then
  echo "Run this script with sudo so app.env can be owned by container UID 65532." >&2
  exit 1
fi
if ! command -v openssl >/dev/null 2>&1; then
  echo "openssl is required to generate local lab credentials." >&2
  exit 1
fi

for file in .env postgres.env app.env; do
  if [[ -e "$file" ]]; then
    echo "Refusing to overwrite existing $file" >&2
    exit 1
  fi
done

umask 077
cp .env.example .env
cp postgres.env.example postgres.env
cp app.env.example app.env

db_password="$(openssl rand -hex 24)"
token_key="$(openssl rand -hex 16)"

sed -i "s|^POSTGRES_PASSWORD=.*$|POSTGRES_PASSWORD=${db_password}|" postgres.env
sed -i "s|^DB_SOURCE=.*$|DB_SOURCE=postgres://study:${db_password}@db:5432/study?sslmode=disable|" app.env
sed -i "s|^TOKEN_SYMMETRIC_KEY=.*$|TOKEN_SYMMETRIC_KEY=${token_key}|" app.env

chown 65532:65532 app.env
chmod 600 .env postgres.env app.env

echo "Local lab config created. Set STUDY_IMAGE in .env to the CI-published SHA tag before starting the API."
echo "Credentials were not printed. Keep these files on the server and out of Git."
