# Task group — 13 domain: usage

**Depends on:** `02-migrations`, `03-platform`, `06-domain-tenant` (reader). Hourly metering + cost estimation (PRD G7/M18, ADR-3's egress estimate).

- [ ] **`repository.go` / `repository/repository.go`** — snapshot upsert on `{tenant_id, period}` (idempotent), query-by-tenant, install-wide aggregation.
- [ ] **`usecase.go` / `usecase/usecase.go`** —
  - **Hourly bucket** per tenant: bytes stored, object count, egress, request count, `estimated_cost` from the then-current rate card (`ERD.md` usage-snapshots — historical figures don't silently change when a rate card is edited).
  - **Egress is estimated** (±10%, PRD G7): the default download path is a redirect and Bloberry never sees the transfer (ADR-3) — the estimate is labeled, never presented as a bill.
  - **Quota reconciliation** — the hourly metering job reconciles the denormalized `used_bytes`/`used_objects` counters against the truth (quota checks read the counter because they sit on the write path; drift is corrected here).
  - **Rate-card rule**: absent card ⇒ cost "unknown", never $0 (a zero reads as free — PRD M18).
- [ ] **`handler.go` / `handler/handler.go`** — `GET /v1/usage` (per-tenant, `?period`), `GET /v1/admin/usage` (install-wide, platform admin).
- [ ] **Scheduled ticker** — the hourly metering job runs as an **in-process ticker** (not a systemd timer, `detail-infra` decision) with the single-instance assumption recorded.
- [ ] **Interface naming + mocks.**

**tests:** idempotent upsert (re-run doesn't double-count); estimated-cost computed at snapshot time from the then-current card; unknown-not-zero for a card-less backend; counter reconciliation corrects drift; egress estimate flag set.
