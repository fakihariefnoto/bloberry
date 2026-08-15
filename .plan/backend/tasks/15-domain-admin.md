# Task group — 15 domain: admin

**Depends on:** `02-migrations`, `03-platform`, `17-storage` (registry, health), `13-domain-usage` (reader), `06-domain-tenant` (reader). Platform-admin only. The credentials surface.

- [ ] **`repository.go` / `repository/repository.go`** — storage-backend CRUD, rate-card update, install stats.
- [ ] **`usecase.go` / `usecase/usecase.go`** —
  - **Backend registration** (PRD M1/M20/PA1): driver config + **credentials envelope-encrypted at rest** (`platform/crypto`), write-only — never returned, never echoed (TRD R7). Null `tenant_id` = install-level pool (PRD D4); many backends per driver type, `name` unique per install (`ERD.md` storage-backends note).
  - **Health check** (PRD M19/PA-E1): periodic ticker + on-demand test; the **raw provider error is visible to platform admins only** — the one documented exception to never passing provider errors through (`domains.md` §8).
  - **Rate card** (PRD M18/PA4): storage $/GB-mo, egress $/GB, per-1k-requests; the input side of `usage`'s estimates.
  - **Delete refused while tenants are assigned** (`ERD.md` storage-backends note — deleting a backend with live tenants must be refused).
- [ ] **`handler.go` / `handler/handler.go`** — `POST/GET /v1/admin/backends`, `GET/PATCH /v1/admin/backends/:id`, `POST /v1/admin/backends/:id/test`, `PUT /v1/admin/backends/:id/rate-card`, `GET /v1/admin/tenants`, `PATCH /v1/admin/tenants/:id/quota`, `GET /v1/admin/stats`.
- [ ] **Platform-admin gating** — every admin route guarded by the `platform_admin` role, distinct from tenant-admin routes.
- [ ] **Interface naming + mocks.**

**tests:** credentials never returned by any endpoint; raw provider error only on the health view; delete refused with assigned tenants; rate-card drives the usage estimate; platform-admin gate on every route.
