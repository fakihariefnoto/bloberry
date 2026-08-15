# Task group — 17 screen: shares

**Depends on:** `02-design-tokens`, `03-routing`, `04-core-infra`. **Blocks:** `25-flows` (share flow). **Mockup:** [`mobile/mockup/shares.md`](../mockup/shares.md).

- [ ] **Layout — populated** per the mockup: ACTIVE links first (type icon, name, kind + expiry line, hit count, status pill, `⋮`), then EXPIRED / REVOKED grouped under a muted header (no actions, Copy disabled for dead links).
- [ ] **Layout — empty** — "Nothing shared yet · Share any file to create a link — signed links expire, short URLs stick until revoked." with a Share a file action.
- [ ] **Loading** — skeleton rows.
- [ ] **Error** — inline banner + retry.
- [ ] **Revoke** — `⋮` → Revoke with a plain confirm stating the consequence ("12 people have opened this this week — revoking kills the link now"). No Undo toast.
- [ ] **Hits-first** — active links sort by hits desc (the revoke decision); `0 hits · never` is safe to revoke.
- [ ] **Copy URL** — copies the full link with confirmation.
- [ ] **Tap the linked file** → `file-detail`.
- [ ] **`+` Share** → file picker → `share-sheet` (sharing a *new* file from this tab).
- [ ] **Footer caption** — "Public objects live on the file, not here" (prevents the "made it public, why isn't it here" hunt).
- [ ] **A11y** — status pills carry text labels (not color alone); rows ≥48dp.

**tests:** active/expired grouping; revoke confirm with hit count; hits-desc sort; dead-link copy disabled; footer caption present.
