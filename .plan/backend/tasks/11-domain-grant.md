# Task group — 11 domain: grant

**Depends on:** `02-migrations`, `03-platform`, `16-authz` (invalidation), `07-domain-folder` (reader). The folder-level RBAC layer on top of the role floor (PRD M9/D7).

- [ ] **`repository.go` / `repository/repository.go`** — grant CRUD, load-by-principal (`{tenant_id, principal_type, principal_id}` partial `revoked_at:null` — the resolver's grant load), load-by-folder (`{tenant_id, folder_id}` for the UI).
- [ ] **`usecase.go` / `usecase/usecase.go`** —
  - **Allow-only, most-specific-wins, no deny** (PRD D7): create validates the permission set (read/write/delete/share — `admin` is not a grant-level permission), folder subtree, optional expiry.
  - **Revocation** sets `revoked_at` (kept for audit, `ERD.md` grants note), **explicitly invalidates the principal's cache entry** (next-request effect, PRD G5).
  - **Principal must be in the tenant** (a grant to an outside identity is inert and confusing; refuse at create).
- [ ] **`handler.go` / `handler/handler.go`** — `POST /v1/grants`, `GET /v1/grants`, `GET /v1/grants?folder=…`, `DELETE /v1/grants/:id`.
- [ ] **Interface naming + mocks.**

**tests:** allow-only resolution (a deny-shaped request is refused at create, not evaluated); revocation invalidates the cache on the next request; principal-not-in-tenant refused; most-specific-grant-wins is the resolver's job (`16-authz`), this file's tests cover the store + invalidation.
