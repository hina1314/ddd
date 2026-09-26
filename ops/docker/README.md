# Single-server Docker lab

The Compose stack runs PostgreSQL, a one-shot database migration, and the API.
PostgreSQL has no host port; the API listens only on host loopback port 3000.
The database uses the named `postgres_data` volume. The volume is **not a backup**.

On the server, use a checkout of this repository and run once:

```bash
sudo bash ops/docker/prepare.sh
```

Edit `ops/docker/.env` and replace `replace-me` with the full SHA image tag
published by the green CI run on `main`. If the GHCR package is private, log in
to `ghcr.io` with a token that has `read:packages` before pulling.

Then start the stack from the Compose directory:

```bash
cd ops/docker
docker compose config --quiet
docker compose up -d
docker compose ps
curl --fail --include http://127.0.0.1:3000/readyz
```

`migrate` must complete successfully before `api` starts. For later schema
changes, the deployment procedure must run migrations again before replacing
the API. Keep `postgres.env` and `app.env` on the server only. Back up the
PostgreSQL data volume separately; do not use `docker compose down --volumes`
when stopping this stack.
