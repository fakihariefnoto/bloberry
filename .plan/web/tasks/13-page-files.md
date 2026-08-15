# Task group — 13 page: files

**Depends on:** `01-setup`, `02-design-tokens`, `03-routing`, `05-core-infra`, `06-shared-components` (DataTable, Breadcrumbs, Dropzone, UploadQueue, StatusPill, FileTypeIcon, EmptyState, ConfirmDestructive, Toast, PermissionDenied). **Blocks:** `33-flows` (upload, share, extract, delete flows all route through here). **Mockup:** [`web/mockup/files.md`](../mockup/files.md).

- [ ] **Layout — desktop populated** per the mockup: `AppShell` + tenant switcher, PageHeader "Files" with sort/filter controls, Breadcrumbs (`Root › Projects › Web › 2026`), the full-table-area Dropzone with the constraint caption ("Drop files anywhere to upload · Max 5 GB per file"), the DataTable (name/size/visibility/modified + `⋮`), the cursor-pager footer, the quota bar under the table.
- [ ] **Layout — mobile populated**: drawer shell, breadcrumbs (collapsed), table becomes **stacked cards** under 640px (rows become labeled lines, per `web-screen/patterns.md`).
- [ ] **Breadcrumbs are the navigation** — each segment links up the hierarchy; collapse to `Root / … / Parent / Current` beyond 4 levels with the ellipsis opening a dropdown of elided segments; every segment is a drop target for drag-to-move.
- [ ] **Loading** — skeleton rows at `size.row-height` (never a spinner).
- [ ] **Empty folder** — distinct copy "This folder is empty · Drop files in or press Upload" with an action.
- [ ] **Filtered-to-empty** — "No matches for this filter", distinct from the empty-folder state, with a Clear-filter action.
- [ ] **Error** — inline banner above the table with Retry, **table area blank** (no stale rows).
- [ ] **Sort/filter server-side** — sort by name/size/modified; filter by visibility/type/uploader; **cursor pagination, never offset** (a collection being written to with offset skips and repeats rows). Virtualized rows for 10,000-object folders; sticky header; page never scrolls horizontally.
- [ ] **Row actions** — hover on pointer, always in the overflow `⋮` on touch (never hover-only): Share, Permissions (folders), Move, Rename, Extract (archive), Download, Delete.
- [ ] **Permission-aware rendering** — a viewer sees write actions **disabled with a reason tooltip** ("Requires `write` on this folder"), never hidden, never an error wall (PRD MV4) — via `PermissionDenied`.
- [ ] **Selection mode** — checkbox column + bulk-action bar replacing the filter bar when anything is selected: Download / Move / Delete (bulk delete confirms once with the count).
- [ ] **Dropzone behaviors** — drag-over → primary border/fill; rejected (size/type) → error border, message reverts after 3s; drop → upload flow. Name collision → inline replace/keep-both/cancel per file, queue stays open (PRD MV-E2).
- [ ] **Quota** — under the table always. ≥80% → warning (`color.warning`); exceeded → error (`color.error`), uploads blocked with "Quota exceeded — reads still work, new uploads are paused" + a link to `usage`.
- [ ] **Public visibility** renders as a warning pill, not a success (the semantic mapping in `02-design-tokens`).
- [ ] **Delete folder** — `ConfirmDestructive` stating the **real object count** ("10,342 objects" not "this folder"), typed-name for irreversible; above a threshold becomes a tracked job → toast links to `jobs` (PRD TA-E1, M21).
- [ ] **Extract/bundle started** → toast linking to `jobs` (feeds `33-flows`).
- [ ] **Breakpoints** — sidebar → 72px nav rail below ~1200px → drawer under 768px; table → stacked cards under 640px.

**tests:** DataTable four states on this page (loading/empty/filtered-empty/error); cursor pagination; bulk-selection bar; permission-disabled rows; quota banner threshold; dropzone reject-and-revert.
