# Task group — 05 first-deploy dry run

**Depends on:** `01-vps-runner`, `02-process-model`, `03-reverse-proxy-dns`, `04-environments-secrets`. **Blocks:** trusting the pipeline for production.

Per `infra/README.md` — the manual dry-run that validates the whole chain against staging before production relies on it.

- [ ] **`ci.yml` manually dispatched and green** against the current code — generate/lint/security/test/build all pass on the self-hosted runner.
- [ ] **`deploy.yml` manually dispatched against `staging`** — the full sequence runs: pull → generate → build → **index migrations** → install to `/opt/bloberry` → `systemctl restart` → health wait.
- [ ] **Staging serves** — `https://staging.bloberry.example.com` (or the staging host) answers `/healthz` and `/readyz`; the dashboard loads; a real login works against the staging Mongo/Redis.
- [ ] **A write path verified on staging** — a small file uploads through the whole chain (presigned or direct) and downloads back, confirming the disk driver + volume + Caddy streaming all actually work end-to-end.
- [ ] **Rollback path exercised** — a second dispatch of `deploy.yml` with a previous binary confirms a redeploy over a bad one is just a rerun (systemd + the install step are idempotent).
- [ ] **The single-instance assumption confirmed** — exactly one `bloberry` process is running after the deploy (two would double-run every ticker).

**verification:** the staging deployment is reachable, healthy, and round-trips a file; the whole sequence took a single manual dispatch each.
