# Task group — 07 domain: folder

**Depends on:** `02-migrations`, `03-platform`, `16-authz`, `12-domain-job` (`job.Enqueuer` for large subtree deletes).

- [ ] **`repository.go` / `repository/repository.go`** — list children (one indexed query via `{tenant_id, parent_id, name}`), get-by-path, subtree query (via `ancestors` multikey), insert/update/delete, descendants rewrite for moves.
- [ ] **`usecase.go` / `usecase/usecase.go`** — create (with missing parents), rename, move, delete-subtree, list. **Cycle prevention** (`target.ancestors` contains the moved folder's id → refuse, PRD TA-E2, `folder_cycle`). **Move rewrites `path` + `ancestors` for every descendant**; bounded inline, above threshold → `job.Enqueuer` subtree_delete (PRD M21/TA-E1). **Zero storage-backend copies** (ADR-7).
- [ ] **`handler.go` / `handler/handler.go`** — `POST /v1/folders`, `GET /v1/folders/:id/children`, `PATCH /v1/folders/:id`, `POST /v1/folders/:id/move`, `DELETE /v1/folders/:id` (subtree), `GET /v1/folders/tree`.
- [ ] **Name-collision** — `{tenant_id, parent_id, name}` unique (no two siblings share a name) → `name_conflict` 409.
- [ ] **Authz on every operation** — resolver gate on the folder's `ancestors` per the `16-authz` model (read on list/stat, write on create/rename/move, delete on subtree-delete).
- [ ] **Interface naming + mocks.**

**tests:** list children latency shape (single indexed query); move rewrites descendants' `ancestors`; moving into own descendant refused (`folder_cycle`); large-subtree delete becomes a job; rename leaves object `file_id`s untouched.
