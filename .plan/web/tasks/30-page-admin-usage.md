# Task group — 29 page: admin-usage

**Depends on:** `01-setup`, `02-design-tokens`, `03-routing`, `05-core-infra`, `06-shared-components` (PageHeader, DateRangePicker, ByteSize, DataTable). **Blocks:** none. **Mockup:** [`web/mockup/admin-usage.md`](../mockup/admin-usage.md).

Per the analytics-dashboard pattern — the install-wide cost view (PRD G7/PA3), the screen the admin opens when the storage invoice is higher than expected.

- [ ] **Layout — desktop** per the mockup: PageHeader "Install usage", global DateRangePicker, three stat cards (stored / tenants / est. total cost, with sparklines + direction), one primary area chart (total stored across tenants), "Cost by tenant" DataTable with Export CSV, the "unknown" footnote + egress estimate note.
- [ ] **Layout — mobile**: swipeable stat cards, chart full-width, tenant table stacked.
- [ ] **Period control global** — chart, stat cards, table, totals together. Presets 24h/7d/30d/This month/custom.
- [ ] **Cost by tenant is the load-bearing table** — default sort Est. cost desc so the top row answers "who's burning money".
- [ ] **The `unknown` row** — a tenant whose backend has no rate card renders "unknown" in `color.warning` with a "no rate card" hint; the footer states the rule once so the whole table's semantics are legible (same rule as `usage`; `admin-backend-detail` is where it's fixed).
- [ ] **Row click** → `admin-tenant-detail` (from there the tenant's audit/usage deepens). This page is attribution, not investigation.
- [ ] **Export CSV** — current filter/period, caption with range + "estimates" note.
- [ ] **Egress everywhere footnoted estimated** (±10%, `ERD.md` usage-snapshots note).
- [ ] **Loading** — skeleton stat cards + flat chart placeholder (never stale data or a blank axis).
- [ ] **Empty state** — "No usage data yet · Snapshots are taken hourly" (a fresh install, not an error).
- [ ] **A11y** — chart has a data-table fallback for screen readers; stat cards are `dl` structure.

**tests:** cost-desc table; unknown row never $0; period control drives everything; CSV carries the range + estimate note.
