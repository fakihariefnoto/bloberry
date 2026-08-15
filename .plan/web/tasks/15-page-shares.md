# Task group — 15 page: shares

**Depends on:** `01-setup`, `02-design-tokens`, `03-routing`, `05-core-infra`, `06-shared-components` (DataTable, StatusPill, CopyableCode, EmptyState, ConfirmDestructive). **Blocks:** `33-flows` (share flow). **Mockup:** [`web/mockup/shares.md`](../mockup/shares.md).

- [ ] **Layout — desktop** per the mockup: PageHeader "Shares" with `+ Share`, status filter tabs (All / Active / Expired / Revoked), DataTable (file, kind, created, hits, status, `⋮`), the footer caption about public objects.
- [ ] **Layout — mobile**: stacked cards (each row's columns become labeled lines).
- [ ] **`?status=` URL state** — the filter is real URL state, so it survives refresh.
- [ ] **Hits-first default sort** — active links sort by hit count desc (the "is this link load-bearing" answer); sortable.
- [ ] **Row actions** — Copy URL, Open (→ `file-detail` of the target), Revoke. Revoke = **plain confirm** (not typed-name) stating what dies: "13 people have opened this in the last week — revoking kills the link now".
- [ ] **Selection + bulk revoke** — checkbox column on the active tab; bulk revoke confirms once with the count.
- [ ] **Target deleted** — a share whose object is soft-deleted shows a muted "File in trash" tag and Copy URL disabled (copying a dead URL is worse than not offering it).
- [ ] **Loading / error** — DataTable skeleton rows / inline error banner with retry.
- [ ] **Empty states differ** — no shares → "Nothing shared yet · Press Share on any file to create a link"; expired filter → "No expired links · Nothing here has outlived its TTL". The second must not read like the first.
- [ ] **Public-objects footer caption** — "Public objects are not listed here — visibility lives on the file itself" (prevents the "I made it public, why isn't it here" hunt).

**tests:** status filter → URL state; hits-default sort; revoke confirm states hit count; target-deleted disables copy; the two distinct empties.
