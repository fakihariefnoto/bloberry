# Task group — 08 monitoring + backups

**Depends on:** `02-process-model.md` (the running services), `01-vps-runner.md`. **Blocks:** nothing (operational readiness).

Per `infra/README.md` §Backups and §Monitoring — minimal and honest; knowing when the box is in trouble, not an observability platform.

- [ ] **`/healthz` liveness** served by the app (already in `backend/tasks/01-setup` — confirm it's wired through Caddy).
- [ ] **`/readyz` readiness** — checks Mongo and Redis; Caddy and the deploy workflow both wait on it.
- [ ] **Structured JSON logs to stdout → journald.** **Never logged:** access-key secrets, presigned URLs (they *are* credentials for their TTL), storage credentials.
- [ ] **Alert rules set** — disk > 80% on the objects volume, disk > 80% on root, any storage backend `unreachable` > 15 min, job queue depth > 100, `bloberry.service` restart-looping.
- [ ] **Nightly `mongodump`** to a *different* storage backend than any tenant uses, retained 30 days — the only irreplaceable state (lose it and every object becomes an unaddressable blob).
- [ ] **`CREDENTIAL_ENCRYPTION_KEY` backed up separately by hand** — in a password manager, deliberately NOT alongside the dump (backing it up with the database defeats the encryption design).
- [ ] **Local-disk objects: tenant's responsibility, stated** — the install docs say plainly those bytes have exactly the durability of that one volume (PRD NG4; Bloberry is not a backup product and doesn't silently become one).
- [ ] **Restore runbook written** — a Mongo restore **must be paired with a reconciliation run** (metadata and bytes are separately stored, `architecture.md` ADR-5): the sweep resolves orphaned blobs from today's uploads; a documented manual pass handles today's deleted objects resurrected as broken references.

**verification:** a test alert fires (a `curl` to a deliberately unhealthy endpoint); a nightly dump exists after one day; the restore runbook is a documented, rehearsed-on-staging procedure.
