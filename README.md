# iotplatform

Multi-tenant B2B example built with [go-zero](https://go-zero.dev/). It demonstrates IoT-oriented services.

## Clone

This repository uses [Git submodules](https://git-scm.com/book/en/v2/Git-Tools-Submodules) for shared Cursor rules and AI context (`.cursorrules`, `.github/ai-context`, `.ai-context/zero-skills`).

```bash
git clone --recurse-submodules <repository-url>
```

If you already cloned without submodules:

```bash
git submodule update --init --recursive
```

## Development

### IDE

**Cursor** or **Visual Studio Code**

### Unit tests

Run unit tests in a clean Go container:

```bash
bash scripts/test/test.sh
```

### Run locally

**Start** the stack:

```bash
bash scripts/dev/up.sh
```

**View logs:**

```bash
bash scripts/dev/logs.sh
```

**Stop** and remove volumes:

```bash
bash scripts/dev/down.sh
```

### Debugging

Start the stack in debug mode (same app ports + extra Delve ports):

```bash
bash scripts/dev/debug.sh
```

Then attach your Go debugger to:

- `platform-api`: `localhost:2345`
- `transform-rpc`: `localhost:2346`

## Deploy

For small and medium projects, Docker Compose is usually enough. For large systems, Kubernetes is the usual choice for availability, scaling, and orchestration. go-zero fits well with both: Compose for simpler setups and cloud-native patterns when you move to Kubernetes.

### Option 1: Docker Compose

#### GitHub

In GitHub: **Settings → Environments**

- Create an environment named `staging`
- Create an environment named `production`

**Configure GitHub secrets**

- Environment secrets (per environment: `staging` / `production`):
  - `DEPLOY_HOST` (server hostname/IP)
  - `DEPLOY_PORT` (optional, defaults to 22)
  - `DEPLOY_USER` (SSH username)
  - `DEPLOY_SSH_KEY` (private key for SSH)
  - `DEPLOY_PATH` (working directory on server)
- Repository secrets:
  - `DOCKERHUB_USERNAME`
  - `DOCKERHUB_TOKEN`

In the `production` environment, enable **Required reviewers** so production deploys require manual approval.

Branch mapping:

- Pushes to `develop` deploy to **staging**
- Pushes to `main` deploy to **production**

#### Server

Upload everything from the `scripts/deploy/with-docker-compose/` directory to the `$DEPLOY_PATH/` directory on the server.


In the `$DEPLOY_PATH/` directory, run `mv .env.example .env`, then edit `.env`.

After changing `.env`, update `config/*.yaml` on the server so those files match the hostnames, ports, credentials, and other values in `.env`.

### Option 2: Kubernetes

This repository does not include Kubernetes manifests. See the [go-zero documentation](https://go-zero.dev/) for cloud-native deployment guidance.
