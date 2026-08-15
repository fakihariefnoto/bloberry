# Task group — 09 domain: share

**Depends on:** `02-migrations`, `03-platform`, `16-authz`, `17-storage` (redirect vs proxy). Two kinds, one collection (`ERD.md` share-links note).

- [ ] **`repository.go` / `repository/repository.go`** — share-link CRUD, slug lookup (globally unique per install, PRD D6), hits increment, revocation.
- [ ] **`usecase.go` / `usecase/usecase.go`** —
  - **Signed links** (PRD M7): create with caller-specified TTL (max from config), revocable before expiry. Download resolves through the object's backend — redirect for cloud drivers, proxy for disk (`architecture.md` §3.2, ADR-3). **TTL short (5 min for API-issued downloads)** because a presigned URL outlives its revocation (TRD R11) — stated in the UI, not buried.
  - **Short URLs** (PRD M8/D6): random unguessable slug by default (a short URL is a capability), optional requested slug, `slug` unique per install, permanent until revoked unless TTL given.
  - **Public visibility** (PRD M6): flips `objects.visibility`; stable public URL at `/o/<file_id>` (survives renames, PRD M4). Served via redirect with the un-publish caveat (ADR-3).
  - **Expired/revoked** → `link_expired` 410 rendered as the human HTML page (`web/mockup/link-expired.md`) — never a JSON envelope for `/s/` (the consumer is a person in a chat window).
- [ ] **`handler.go` / `handler/handler.go`** — `POST /v1/shares/signed`, `POST /v1/shares/short`, `GET /v1/shares`, `DELETE /v1/shares/:id`, `GET /s/:slug` (public), `POST /v1/objects/:id/visibility`.
- [ ] **Hit tracking** — `hit_count`/`last_accessed_at` updated on resolution (the "is this link still being used" answer before someone revokes it).
- [ ] **Authz** — `share` permission to create; the resolved download honors the object's current permissions at read time.
- [ ] **Interface naming + mocks.**

**tests:** signed link honors TTL + revocation (revoked link 410s); short-slug uniqueness across tenants; random slug unguessable; expired link renders HTML 410 not JSON; public URL stable across rename; hit count increments on each resolution.
