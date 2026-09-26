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

## Manually triggered deployment from GitHub Actions

The `Deploy to lab server` workflow deploys the current `main` commit only.
First confirm that the CI `publish-image` job for that commit is green. The
workflow uses SSH to update `/opt/apps/ddd` to the exact commit, then calls a
root-owned helper. The helper pulls the matching immutable image, runs that
commit's migrations, updates the API and checks both `/readyz` and
`X-App-Version`. If API startup fails, it attempts to restore the previous image; it does
**not** reverse database migrations. Keep schema changes backward-compatible.
The single API instance can have a brief interruption during replacement.

One-time server setup, after this workflow and helper have reached `main`:

```bash
cd /opt/apps/ddd
git fetch origin main
git checkout main
git pull --ff-only
sudo install -o root -g root -m 0755 ops/docker/deploy.sh /usr/local/sbin/study-deploy
sudo visudo -f /etc/sudoers.d/study-deploy
```

Put this single line in the sudoers file, then save it:

```text
liao ALL=(root) NOPASSWD: /usr/local/sbin/study-deploy
```

The helper validates its SHA argument and is copied outside the writable
checkout. Do not allow passwordless `sudo docker` or add the deploy SSH key's
account to the `docker` group: both grant broad host control. For a real
production server, also protect the checkout from modification by the SSH
account; this is a single-server learning setup.
Reinstall the root-owned helper whenever `ops/docker/deploy.sh` changes.

Create a dedicated SSH key for the workflow. Add only its **public** key to
`liao`'s `~/.ssh/authorized_keys`; put its **private** key in the GitHub Actions
repository secret `DEPLOY_SSH_KEY` (never commit it). Configure:

- Repository variables: `DEPLOY_HOST` (server address), `DEPLOY_USER` (`liao`).
- Repository secret: `DEPLOY_KNOWN_HOSTS` (the verified SSH host-key line for
  that address; check its fingerprint against the server before saving it).

The server's `liao` account must also be able to run `git fetch origin main`
without an interactive password. Verify the SSH connection and `git fetch`
before running the workflow. Start the workflow from the `main` branch in the
GitHub Actions tab. If the image is private, the server's root Docker client
must already be logged in to GHCR with pull access.
