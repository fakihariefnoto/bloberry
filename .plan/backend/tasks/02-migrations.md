# Task group — 02 migrations

**Depends on:** `01-setup.md` (migration tool wired). **Blocks:** every domain file. One migration per entity group from `ERD.md`, in dependency order (referenced collections before their foreign keys). Each = the collection + its compound indexes from `ERD.md`'s index table + any Mongo validators.

- [ ] **Baseline: `users`** (with embedded `oauth_identities: []` and `settings: {}` per `ERD.md` collection shapes — no separate collections). Indexes: `email` unique.
- [ ] **`tenants`** — status, `used_bytes`/`used_objects` counters, `billing` (dormant, PRD NG5), `default_backend_id` (referenced later).
- [ ] **`memberships`** — role enum on `memberships` (the per-tenant role home, `domains.md` §1.2). Index: `{user_id, tenant_id}` unique.
- [ ] **`invitations`** — `token_hash` unique, `expires_at` TTL.
- [ ] **`applications`** — non-human principals.
- [ ] **`access_keys`** — `secret_hash` unique (argon2id), `{tenant_id, application_id}`; `scope_folder_ids` + `permissions` arrays; `application_id`/`user_id` mutually exclusive (a key is one or the other).
- [ ] **`storage_backends`** — nullable `tenant_id` (NULL = install-level, PRD D4); `credentials_encrypted` binary; `capabilities` object; `rate_card` object; `health_*` fields.
- [ ] **`folders`** — `ancestors[]` (ordered, excludes self), `path`, `depth`; indexes: `{tenant_id, parent_id, name}` unique, `{tenant_id, ancestors}` multikey, `{tenant_id, path}` unique.
- [ ] **`objects`** — the `file_id` home (stable forever, PRD M4); `backend_id` per-object (ADR-4); `ancestors[]` denormalized; `state` pending/active/deleting; `content_hash` (dedup key); `deleted_at` nullable (soft delete, S5). Indexes: `{tenant_id, folder_id, name}` unique partial `deleted_at:null`, `{tenant_id, ancestors}` multikey, `{tenant_id, content_hash}`, `{state, created_at}` partial `state:"pending"`, `{backend_id}`.
- [ ] **`multipart_uploads`** — `parts_received` array; `expires_at` TTL (abandoned uploads self-expire).
- [ ] **`grants`** — `principal_type`/`principal_id`, `permissions`, `revoked_at` nullable; indexes: `{tenant_id, principal_type, principal_id}` partial `revoked_at:null`, `{tenant_id, folder_id}`.
- [ ] **`share_links`** — kind signed/short, `slug` unique per install (PRD D6), `token_hash`, `hit_count`; index `{tenant_id, object_id}`.
- [ ] **`jobs`** — kind/state/payload/result/progress/failure_*; index `{tenant_id, state, created_at}`.
- [ ] **`usage_snapshots`** — `{tenant_id, period}` unique (idempotent hourly upsert); decimal cost.
- [ ] **`audit_events`** — append-only; indexes `{tenant_id, created_at:-1}`, `{tenant_id, target_type, target_id}`. **Standard collection + retention job** (ERD-Q2 resolution), not capped/time-series.

**tests:** migrations run clean on a fresh Mongo; every unique index rejects the duplicate case it's named for (two siblings with the same name, duplicate slug, duplicate secret_hash); the partial `deleted_at:null` index allows re-uploading a soft-deleted name.
