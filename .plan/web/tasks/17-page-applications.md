# Task group — 17 page: applications

**Depends on:** `01-setup`, `02-design-tokens`, `03-routing`, `05-core-infra`, `06-shared-components` (DataTable, StatusPill, RelativeTime, EmptyState, FormField). **Blocks:** `18-page-application-detail`. **Mockup:** [`web/mockup/applications.md`](../mockup/applications.md).

- [ ] **Layout — desktop** per the mockup: PageHeader "Applications" with `+ New app`, subtitle "Machine accounts that hold access keys.", DataTable (name, keys, last used, created, `⋮`).
- [ ] **Layout — mobile**: stacked cards.
- [ ] **New app modal** — name + description; on save → **immediately `application-detail`** for the new app (issuing a key is its real purpose — don't drop the user back on the list).
- [ ] **Row click** → `application-detail`; Keys column links into the app's detail filtered to keys.
- [ ] **⚠ keyless state** — an application with zero active keys gets the warning pill ("0 keys · never" or "all revoked") — a keyless app is usually a broken pipeline, not a cleaned-up one. Deliberate deviation from neutral "0".
- [ ] **Delete application** — `ConfirmDestructive`; **refused when the app has active keys** (silently orphans the keys authorizing production CI — `ERD.md` access-key lifecycle), with a message telling the user to revoke keys first.
- [ ] **Empty state** — "No applications yet · Register an app to issue it scoped access keys" with the New button.
- [ ] **Loading / error** — DataTable skeleton / inline banner with retry.
- [ ] **Permission-aware** — `member`/`viewer` never see the route (sidebar hidden + server guard); key management is admin-wide, not owner-only.

**tests:** new-app lands on detail; keyless warning pill; delete refused with active keys.
