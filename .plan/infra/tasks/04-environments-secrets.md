# Task group — 04 environments + secrets + workflows

**Depends on:** `01-vps-runner.md` (runner exists), `02-process-model.md` (`.env.example` per stage). **Blocks:** `05-first-deploy.md`.

Per `infra/README.md` §Environments and §Workflows. Two GitHub Environments (`staging`, `production`), each with its own secrets; dev is local only.

- [ ] **GitHub Environment `staging` created.**
- [ ] **GitHub Environment `production` created.**
- [ ] **Every secret set in BOTH environments** from `backend/.env.example` (committed — the real values are not), one task per secret: `MONGO_URI`, `REDIS_ADDR`, `JWT_SECRET`, **`CREDENTIAL_ENCRYPTION_KEY`** (the special one — see `02`), `GOOGLE_CLIENT_ID`/`GOOGLE_CLIENT_SECRET`, `SMTP_HOST`/`_PORT`/`_USER`/`_PASSWORD`/`_FROM`, `PUBLIC_BASE_URL` (staging vs production value differs), `DISK_STORAGE_PATH`, `DISK_SIGNING_SECRET`, `MAX_OBJECT_SIZE`.
- [ ] **`ci.yml` generated** from `templates/github-actions/`, **`workflow_dispatch` only** (the fixed convention — CI checks are manual too), with the build order that is not optional (`architecture.md` §7): `make generate` (fails if the tree changes) → `cd web && bun install && bun run build` → `make lint` → `make security` → `make test` (incl. the 100%-branch gate on `internal/authz`) → `make build`. Steps 1–2 before any Go build (a backend build first embeds the previous frontend).
- [ ] **Conformance suite in CI** — the storage conformance suite runs against MinIO in a service container on every CI run; against real S3/R2/OSS/GCS on a **weekly schedule** (costs money + live credentials; not on every push).
- [ ] **`deploy.yml` generated** — `workflow_dispatch`, `runs-on: self-hosted`, `environment` input (`staging`|`production`): `git pull → make generate → bun run build → make build → run pending Mongo index migrations → install binary to /opt/bloberry → systemctl restart bloberry → wait for the health endpoint`. **Index migrations before the restart** (a missing migration is a silent latency cliff, not a loud failure — README's explicit ordering).

**verification:** `ci.yml` runs green on a manual dispatch; `deploy.yml` validates its inputs; every secret is set in the right environment (a secret referenced but unset fails loudly at first use, not silently).
