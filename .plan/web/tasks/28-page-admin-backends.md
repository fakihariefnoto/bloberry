# Task group — 27 page: admin-backends

**Depends on:** `01-setup`, `02-design-tokens`, `03-routing`, `05-core-infra`, `06-shared-components` (DataTable, StatusPill, EmptyState, FormField, ConfirmDestructive). **Blocks:** `28-page-admin-backend-detail`. **Mockup:** [`web/mockup/admin-backends.md`](../mockup/admin-backends.md).

- [ ] **Layout — desktop** per the mockup: PageHeader "Storage backends" with `+ Register`, DataTable **grouped by driver** (S3 / R2 / OSS / GCS / Disk headers — one install has many distinct accounts per driver, `ERD.md` storage-backends note), columns driver/name/bucket-prefix/health/tenants/`⋮`, the health-check footer.
- [ ] **Layout — mobile**: stacked cards grouped under driver headers.
- [ ] **Register modal** — driver selector, name, config (endpoint/bucket/prefix per driver), **write-only credentials fields** (never echoed on edit, PRD M20/R7), rate card. Notes "credentials are envelope-encrypted at rest".
- [ ] **Health** — polled every 5 minutes (in-process ticker); an unreachable backend expands to show the **raw provider error** (PA-E1), the one place it's legal.
- [ ] **Tenants column** answers the deletion question before it's asked: a backend with tenants assigned gets **delete refused** ("3 tenants use this — reassign them first"); only a 0-tenant backend can be deleted, and that still confirms.
- [ ] **Rate-card gap** — a backend with no rate card shows "no rate card" (`color.warning`) — `usage` shows "unknown" for its tenants, so the admin sees the gap where filling it matters.
- [ ] **Empty state** — "No storage backends · Register one to start assigning tenants" (the true first-run gate — nothing works until a backend exists).
- [ ] **Permission-aware** — platform_admin only.

**tests:** driver grouping; delete refused with assigned tenants; unreachable expands to raw error; register writes credentials once.
