# Single-server Docker lab

The Compose stack runs PostgreSQL, a one-shot database migration, and the API.
PostgreSQL has no host port; the API listens only on host loopback port 3001
by default (`API_HOST_PORT` controls the host port).
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
sudo docker compose config --quiet
sudo docker compose up -d
sudo docker compose ps
curl --fail --include http://127.0.0.1:3001/readyz
```

`migrate` must complete successfully before `api` starts. For later schema
changes, the deployment procedure must run migrations again before replacing
the API. Keep `postgres.env` and `app.env` on the server only. Back up the
PostgreSQL data volume separately; do not use `docker compose down --volumes`
when stopping this stack.

## Automatic deployment from GitHub Actions

On a push to `main` (including a merged pull request), CI verifies the code,
builds and publishes the image, then calls `Deploy to lab server`. PRs and
feature-branch pushes do not deploy. CI/image failures prevent deployment.
The existing manual trigger is retained for retries; check that the image for
that commit has been published before using it. Superseded commits are skipped
if `main` has advanced before the server checkout is changed.

The workflow uses SSH to update `/opt/apps/app` to the exact commit, then calls a
root-owned helper. The helper pulls the matching immutable image, runs that
commit's migrations, updates the API and checks both `/readyz` and
`X-App-Version`. If API startup fails, it attempts to restore the previous image; it does
**not** reverse database migrations. Keep schema changes backward-compatible.
The single API instance can have a brief interruption during replacement.

### Replace the old lab with a fresh deployment

For a full reset of the old `/opt/apps/ddd` deployment using the `liao`
account, follow [the fresh deployment steps](../../docs/fresh-lab-deployment.md).
Those steps delete the old Compose project's PostgreSQL volume and generate
new credentials; use them only when the old data can be discarded.

### First-time server setup

One-time server setup, after this workflow and helper have reached `main`.
The workflow requires an existing checkout; it does not create the directory,
clone the repository, generate credentials or install the privileged helper.
Run these commands on the deployment server, not on your local development PC.

First create the checkout as the same account configured in `DEPLOY_USER`
(the following commands assume the `deploy` account already exists):

```bash
sudo install -d -o deploy -g "$(id -gn deploy)" -m 0755 /opt/apps/app
sudo -u deploy git clone git@github.com:hina1314/ddd.git /opt/apps/app
sudo -iu deploy
cd /opt/apps/app
git fetch origin main
git checkout main
git pull --ff-only
exit
```

The deploy account needs repository read access for the clone and subsequent
fetches. If your checkout already exists elsewhere, put it at the expected path;
creating an empty `/opt/apps/app` directory is insufficient.

Then, as the server administrator, initialize the configuration once and install
the helper:

```bash
cd /opt/apps/app
sudo bash ops/docker/prepare.sh
sudoedit ops/docker/.env
sudo install -o root -g root -m 0755 ops/docker/deploy.sh /usr/local/sbin/app-deploy
sudo visudo -f /etc/sudoers.d/app-deploy
```

Set `APP_IMAGE` to `ghcr.io/hina1314/ddd:sha-<full-commit-sha>` from the
successful image publication step. Set `API_HOST_PORT` to the intended port.
`prepare.sh` refuses to overwrite existing configuration; skip it if all three
configuration files already exist. Docker Engine, the Compose plugin, Git,
OpenSSL and curl must be available on the server.

Put this single line in the sudoers file, then save it:

```text
deploy ALL=(root) NOPASSWD: /usr/local/sbin/app-deploy
```

The helper validates its SHA argument and is copied outside the writable
checkout. Do not allow passwordless `sudo docker` or add the deploy SSH key's
account to the `docker` group: both grant broad host control. For a real
production server, also protect the checkout from modification by the SSH
account; this is a single-server learning setup.
Reinstall the root-owned helper whenever `ops/docker/deploy.sh` changes.

Create a dedicated SSH key for the workflow. Add only its **public** key to the
deploy account's `~/.ssh/authorized_keys`; put its **private** key in GitHub Actions
repository secret `DEPLOY_SSH_KEY` (never commit it). Configure:

- Repository variables: `DEPLOY_HOST` (server address), `DEPLOY_USER` (`deploy`).
- Repository secret: `DEPLOY_KNOWN_HOSTS` (the verified SSH host-key line for
  that address; check its fingerprint against the server before saving it).

The server's deploy account must also be able to run `git fetch origin main`
without an interactive password. Verify the SSH connection and `git fetch`
before running the workflow. For a manual retry, start the workflow from the
`main` branch in the GitHub Actions tab. If the image is private, the server's root Docker client
must already be logged in to GHCR with pull access.
