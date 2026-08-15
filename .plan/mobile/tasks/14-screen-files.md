# Task group — 13 screen: files

**Depends on:** `02-design-tokens`, `03-routing`, `04-core-infra` (queue store, API). **Blocks:** `25-flows` (upload, share, revoke flows). **Mockup:** [`mobile/mockup/files.md`](../mockup/files.md).

- [ ] **Layout — populated** per the mockup: app bar (tenant switcher + search icon), breadcrumbs, folder rows (type icon, name, modified) then file rows (icon, name, size/visibility line), the quota bar, FAB (Camera/Files), 4-tab bar.
- [ ] **Layout — empty folder** — distinct copy "This folder is empty · Drop files in from the picker, or upload from this phone." with an Upload action + FAB.
- [ ] **Loading** — skeleton rows at `size.row-height` (guard rule 1, not a spinner).
- [ ] **Filtered-empty** — "No matches in this folder" (different from empty-folder).
- [ ] **Error** — inline banner + Retry, list blank (no stale rows); tab bar survives so navigation isn't dead.
- [ ] **Breadcrumbs** — tap any segment to jump up; long paths collapse to `Root › … › Parent › Current` with the elided segment opening a picker.
- [ ] **Row actions** — long-press → selection mode with a contextual action bar (Download / Share / Move / Delete), same action set as the sheet so no action is reachable only one way.
- [ ] **FAB** → `source-picker` sheet (Camera / Files). Camera/photos permission denied → inline explanation with "Open Settings" (not a dead end).
- [ ] **Upload flow** — pick → returns to `uploads` with items queued; this screen only changes the FAB's availability and quota banners.
- [ ] **Quota** — ≥80% warning banner; exceeded → error banner, FAB disabled with "Quota exceeded — reads still work" + a `usage` link.
- [ ] **Permission-aware** — a viewer's disabled action shows a caption ("Requires `write` on this folder"), not an error toast; disabled rows aren't hidden (PRD MV4).
- [ ] **Offline** — pull-to-refresh shows "Offline — showing last loaded"; queued uploads persist and sync on reconnect.
- [ ] **Overflow-safe** — long file names middle-truncate preserving the extension (guard rule 3); test with a 200-char filename.
- [ ] **Pull-to-refresh** on the list.

**tests:** skeleton-then-list; empty vs filtered-empty; quota banner thresholds; viewer caption on disabled actions; 200-char filename renders without overflow; pull-to-refresh.
