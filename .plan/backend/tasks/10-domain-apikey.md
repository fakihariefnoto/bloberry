# Task group — 10 domain: apikey

**Depends on:** `02-migrations`, `03-platform`, `16-authz` (invalidation). Applications (non-human principals) + their access keys (PRD M10/TA4).

- [ ] **`repository.go` / `repository/repository.go`** — application CRUD, access-key CRUD, secret-hash lookup (the `{secret_hash}` unique index — every machine request authenticates through it), last-used update, revocation.
- [ ] **`usecase.go` / `usecase/usecase.go`** —
  - **Key issue**: argon2id-hash the secret **server-side**, store `secret_hash` + `last_four` only; the full secret returns **exactly once** in the response and is unrecoverable (PRD D5/M10). Scope from `scope_folder_ids` + `permissions`; empty scope = whole tenant; expiry optional.
  - **Key lookup on request**: hash the presented token → lookup → check revoked/expired → load scope → build the `Principal` (the access-key half of the two-schemes-one-principal seam, `domains.md` §4.7).
  - **Revocation**: set `revoked_at` (history kept for the audit trail), **explicitly invalidate the Redis cache entry** (PRD G5 — takes effect on the next request, not TTL-out).
  - **Delete refused while active keys exist** (`ERD.md` access-key lifecycle — deleting an app with live keys silently orphans CI).
- [ ] **`handler.go` / `handler/handler.go`** — `POST /v1/applications`, `GET /v1/applications`, `DELETE /v1/applications/:id`, `POST /v1/applications/:id/keys`, `GET /v1/applications/:id/keys`, `DELETE /v1/keys/:id`.
- [ ] **Distinct terminal error codes** — `key_revoked`/`key_expired` (401, **terminal — do not retry**, PRD AP-E1), distinct from a generic 401 so retry logic doesn't loop against a dead credential.
- [ ] **Rate limiting keyed per access key** — not per tenant (one bad integration doesn't throttle its tenant's dashboard), 429 + `Retry-After` (`platform/httpx`).
- [ ] **Interface naming + mocks.**

**tests:** secret shown once then unrecoverable; lookup by hash; revocation effective on the next request (cache invalidated); keyless app delete refused; `key_revoked` returns the distinct code.
