# Task group — 26 page: admin-tenant-detail

**Depends on:** `01-setup`, `02-design-tokens`, `03-routing`, `05-core-infra`, `06-shared-components` (PageHeader, ByteSize, StatusPill, ConfirmDestructive). **Blocks:** none. **Mockup:** [`web/mockup/admin-tenant-detail.md`](../mockup/admin-tenant-detail.md).

- [ ] **Layout — desktop** per the mockup: back link, PageHeader (tenant name + slug/status/since + Suspend + Reassign backend), the backend-unreachable banner (real provider error), three stat cards, CONFIG panel (backend, health, quota, used), PEOPLE & ACCESS panel (summary counts), USAGE breakdown (this month).
- [ ] **Layout — mobile**: stacked.
- [ ] **Backend-unreachable banner** — shows the **raw provider error** (PRD PA-E1 — the one documented exception to never passing provider errors through, `backend/domains.md` §8). The fix (re-enter credentials) lives in `admin-backend-detail`.
- [ ] **Suspend / Reactivate** — toggle with a plain confirm stating reads still work / writes blocked (PRD PA-E2); suspended renders the banner + disabled controls.
- [ ] **Reassign backend** — backend picker; the confirm carries the ADR-4 contract ("new objects go to the new backend; existing objects keep resolving") — the platform-admin equivalent of `tenant-settings`'s change.
- [ ] **Stat cards** — skeleton on load; sparklines + direction.
- [ ] **PEOPLE & ACCESS is summary-only** with counts; "Full management lives in the tenant's own surfaces" stated (a platform admin who tries to edit a member's role from here must be told where to go).
- [ ] **Usage breakdown** — same rate-card rule as `usage`: missing card renders "unknown", never $0.
- [ ] **Delete tenant NOT here** — it lives in `admin-tenants` (row action) and `tenant-settings` (owner's DANGER ZONE); deleting from a summary page invites fat-finger catastrophe.
- [ ] **A11y** — the unreachable-banner error announced on arrival; Suspend returns focus to the banner when confirmed.

**tests:** raw-provider-error banner; suspend/reactivate confirm wording; reassign-backend ADR-4 contract; no delete affordance on this page.
