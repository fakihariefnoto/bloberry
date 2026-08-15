# Task group — 12 domain: job

**Depends on:** `02-migrations`, `03-platform`, `17-storage`, `08-domain-object` (writer), `07-domain-folder` (writer). The Redis-backed queue + in-process worker (ADR-8, TRD R12).

- [ ] **`repository.go` / `repository/repository.go`** — job CRUD, state transitions, `{tenant_id, state, created_at}` listing.
- [ ] **`usecase.go` / `usecase/usecase.go`** — enqueue (extract/bundle/subtree_delete), status-by-id, retry (fresh `attempts`), the worker loop.
- [ ] **`handler.go` / `handler/handler.go`** — `POST /v1/jobs` (enqueue), `GET /v1/jobs`, `GET /v1/jobs/:id`.
- [ ] **`worker/`** — the Redis `job:queue` consumer, bounded concurrency well below `GOMAXPROCS` (a 2 GB extraction must not degrade every tenant's request handling — TRD R12), streaming during work.
- [ ] **Kinds:**
  - **extract** (PRD M11/AP4): archive uploaded → queued → extracted into a target folder. **Atomic commit from a staging prefix** — a failed job leaves the target folder unchanged (PRD AP-E2). **Safety ceilings server-side** (TRD R6): decompressed-size + ratio ceilings, reject absolute paths/`..`/symlinks, quota checked per entry not once up front.
  - **bundle** (PRD M11/AP5): N objects → one archive, served from a signed link. Only objects the caller can read are included; an unreadable path fails the whole request (silent omissions in a download are data loss).
  - **subtree_delete** (PRD M21/TA-E1): large folder deletes as a tracked job with object-count progress.
- [ ] **Progress is real** — `progress_done`/`progress_total` drive honest bars; `failure_code` (machine-readable, stable) + `failure_message` (human) both always set.
- [ ] **`attempts`** — failed jobs retry up to the cap; terminal codes (`archive_rejected`) don't retry (retrying a bomb on loop is a self-inflicted DoS).
- [ ] **Interface naming + mocks.**

**tests:** extract atomicity (failed job → target unchanged); zip-bomb/`..`/symlink rejected; bundle refuses an unreadable path; subtree_delete reports object counts; worker concurrency bound respected; terminal code doesn't retry.
