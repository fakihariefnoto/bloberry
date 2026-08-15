# Task group — 20 page: audit

**Depends on:** `01-setup`, `02-design-tokens`, `03-routing`, `05-core-infra`, `06-shared-components` (DataTable, DateRangePicker, RelativeTime, CopyableCode, EmptyState). **Blocks:** none (investigation surface). **Mockup:** [`web/mockup/audit.md`](../mockup/audit.md).

- [ ] **Layout — desktop** per the mockup: PageHeader "Audit log", filter bar (DateRangePicker + Action + Principal + Clear + Export), DataTable (time, action, principal, target, `⋮`), the redirect-download footer caption.
- [ ] **Layout — mobile**: stacked cards with a details expand.
- [ ] **Filter → URL state** — `?from&to&action&principal` so an investigation can be shared/bookmarked; Clear resets.
- [ ] **Export CSV** — exports the *current filter* with a caption of the range.
- [ ] **Row expand** — shows the full event: action-specific `metadata`, IP, user-agent (upload expands to size/hash/backend; share.create to link + TTL).
- [ ] **Target link** — object targets → `file-detail`; a deleted target opens the "file no longer exists" state, keeping the event context.
- [ ] **Redirect-download limitation stated** — footer caption on every width: "Downloads on the redirect path record link issuance, not each byte read (default is 5-min signed URLs)" (ADR-3 honesty, `ERD.md` share-links note). Not buried in a tooltip.
- [ ] **Retention note** — "Events are retained for 365 days" (monthly retention job, `ERD.md` Q2) so a long look-back isn't silently empty.
- [ ] **Empty states** — no events at all → "Nothing has happened yet · Mutations of files, keys and grants appear here"; no events in range → "No events in this range · Widen the date window or clear a filter" with a Clear action. Distinct.
- [ ] **Permission-aware** — `member`/`viewer` never see the route; `tenant_admin` sees the full tenant log; install-wide audit is platform-admin (`admin-usage`/`admin-tenant-detail`).

**tests:** filter URL round-trip; CSV of the filtered set; redirect-limitation caption present; two distinct empties; target link for a deleted file.
