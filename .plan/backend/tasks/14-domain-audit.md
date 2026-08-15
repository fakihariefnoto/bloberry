# Task group — 14 domain: audit

**Depends on:** `02-migrations`, `03-platform`. Append-only event record (PRD M16/TA6) — the standard collection + retention job (ERD-Q2 resolution).

- [ ] **`repository.go` / `repository/repository.go`** — append, query per tenant (reverse-chronological via `{tenant_id, created_at:-1}`), query by target (`{tenant_id, target_type, target_id}`), hard-delete older than retention.
- [ ] **`usecase.go` / `usecase/usecase.go`** —
  - **Append-only**: nothing updates or deletes a single event (except the retention job's bulk delete).
  - **Event shape**: action (`object.upload`, `object.delete`, `grant.create`, `key.revoke`, `member.join`, …), principal type/id, target type/id, action-specific `metadata`, IP, user-agent, timestamp (`ERD.md` audit_events).
  - **The redirect-path limitation recorded honestly** (ADR-3): on the redirect download path this records **link issuance, not byte reads** — a documented limitation stated in the UI (`web/mockup/audit.md` footer), not hidden.
  - **Retention job** — monthly, hard-deletes past the window (default 365 days), in-process ticker.
- [ ] **`handler.go` / `handler/handler.go`** — `GET /v1/audit` with `?from,to,action,principal` filters (tenant_admin+).
- [ ] **Interface naming + mocks.**

**tests:** append-only (no update path); per-tenant isolation (tenant A never sees tenant B's events); retention job deletes only past-window rows; redirect-path events carry the issuance marker.
