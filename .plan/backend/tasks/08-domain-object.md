# Task group — 08 domain: object

**Depends on:** `02-migrations`, `03-platform`, `16-authz`, `17-storage`, `06-domain-tenant` (quota checker via narrow interface). The `file_id` lifecycle — the busiest domain.

- [ ] **`repository.go` / `repository/repository.go`** — object CRUD, stat-by-id, list-by-folder (indexed), `{tenant_id, content_hash}` lookup (dedup), state transitions, soft-delete.
- [ ] **`usecase.go` / `usecase/usecase.go`** —
  - **Presigned upload path** (default): `PresignPut` against the tenant's backend, **quota checked before presign** (`tenant.QuotaChecker`), `object_pending` state written first.
  - **Direct upload path**: streaming `io.Copy` to the driver (never `io.ReadAll`), hard body ceiling.
  - **Multipart/resumable**: init → part presign → complete; `multipart_uploads` record with `parts_received` for resume (PRD MB2/G10); abort cleans provider-side (dangling multiparts accrue invisible cost).
  - **Two-phase commit** (TRD R8): metadata `pending` → byte write confirms → promote `active`. A crash between leaves either a pending row (swept) or an unreferenced blob (reported).
  - **Stable identity**: `id` is the `file_id`, never changes across move/rename/visibility/backend (PRD M4/G4).
  - **Move/rename/visibility/delete**: rename + move update `path`/`ancestors` (no provider copy); visibility flips `public`/`private`; delete is soft (S5) with the partial unique index.
- [ ] **`handler.go` / `handler/handler.go`** — `POST /v1/objects` (direct), `POST /v1/objects/presign-put`, multipart `POST /v1/objects/:id/multipart/init` + `PUT /v1/objects/:id/multipart/parts/:n` + `POST /v1/objects/:id/multipart/complete`, `GET /v1/objects/:id`, `GET /v1/objects/:id/download`, `GET /v1/objects/:id/raw` (disk-driver proxy), `PATCH /v1/objects/:id` (rename/visibility), `POST /v1/objects/:id/move`, `DELETE /v1/objects/:id`.
- [ ] **Name collision handling** — `name_conflict` 409 with the replace/keep-both/cancel choice left to the client (PRD MV-E2).
- [ ] **Authz** — resolver on `objects.ancestors` (the denormalized array is what makes the hot path one indexed query, `ERD.md`).
- [ ] **Interface naming + mocks.**

**tests:** quota checked before presign; pending→active two-phase (crash leaves pending, sweep finds it); multipart resume re-sends only missing parts; `file_id` stable across move+rename+visibility+backend; download via redirect vs raw proxy per driver capability; name-collision 409.
