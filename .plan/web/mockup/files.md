# Screen — files

## Purpose & context

- **User goal**: browse the tenant's folder tree and manage objects — the default landing after login and the reason the app exists.
- **Entry points**: default after login/tenant-switch; sidebar Files; any deep link `/files/:folderId`; back from `file-detail`.
- **Exit points**: click a file → `file-detail`; click a folder row or breadcrumb → `files` at that folder; Share action → `share-dialog`; Permissions → `grant-dialog`; Move → `move-picker`; Delete folder → `confirm-destructive`; Extract/bundle started → toast linking `jobs`; drop files → `upload-queue` docks (no navigation).
- **Data needed**: folder path (breadcrumbs), child folders + objects (name, type, size, visibility, modified, uploader), the tenant's quota bar. Real objects from `ERD.md`: `folders`, `objects` (with `ancestors`, `visibility`, `content_hash`, `size_bytes`), `grants` (permission awareness).

## States

- [x] Loading (skeleton rows at `size.row-height`)
- [x] Empty folder
- [x] Populated (happy path)
- [x] Empty search/results
- [x] Error (inline banner above table, table blank)
- [x] Domain-specific — viewer/permission-limited (write actions disabled with reason)
- [x] Domain-specific — near/exceeded quota (warning/error banner, uploads blocked)
- [x] Domain-specific — selection mode with bulk-action bar

## Style reference

- **Components used**: `AppShell`, `TenantSwitcher`, `Breadcrumbs` (primary navigation — no tree pane), `DataTable` (virtualized, server-side cursor pagination, sort, filter, selection), `Dropzone` (wraps the whole table area), `UploadQueue` (docked), `StatusPill` (visibility), `FileTypeIcon`, `ByteSize`, `RelativeTime`, `EmptyState`, `ConfirmDestructive`.
- Web + desktop both get the sidebar shell; mobile-width shows the drawer.
- No token deltas.

## Wireframe — desktop (populated)

```
┌──────────┬─────────────────────────────────────────────────────────────┐
│ ⬡        │  Acme Inc ▾        Jane Doe ▾         (search) (upload)     │
│ Bloberry │                                                             │
│          │  Files                           [Sort ▾] [Filter ▾]        │
│ ──────── │  Root › Projects › Web › 2026                               │
│ WORKSPACE│  ───────────────────────────────────────────────────────    │
│ ▸ Files  │  Drop files anywhere to upload · Max 5 GB per file          │
│  Shares  │     ┌──────────────────────────────────────────────────────┐│
│   Jobs   │  │ Name                Size   Visibility  Modified  ⋮  │    │
│ ──────── │  │ 📁 _assets                  ─          Mar 12    ⋮  │    │
│ ADMIN    │  │ 📁 public           ─        public    Mar 12    ⋮  │    │
│   Apps   │  │ 📁 scripts         ─        ─          Mar 11    ⋮  │    │
│   Members│  │ 🖼 hero.png         2.4 MB  public     Mar 12    ⋮   │    │
│   Audit  │  │ 📄 index.html      18.2 KB  ─         Mar 12    ⋮   │    │
│   Usage  │  │ 📄 main.ts         9.8 KB   ─         Mar 11    ⋮   │    │
│   Settings│ │ └──────────────────────────────────────────────────────┘ │
│ ──────── │  Showing 6 of 248          ‹ 1 2 3 … 42 ›                   │
│ PLATFORM │  ───────────────────────────────────────────────────────    │
│   Tenants│  Storage  312 GB of 500 GB ████████████░░░░  62%            │
│  Backends│                                                             │
│   Usage  │                                                             │
│ ──────── │                                                             │
│ [👤] Jane│                                                             │
└──────────┴─────────────────────────────────────────────────────────────┘
```

## Wireframe — mobile (populated)

```
┌────────────────────────────┐
│ ☰ ⬡    Acme Inc ▾   👤     │
├────────────────────────────┤
│ Files          [Sort][Flt] │
│ Root › Projects › 2026     │
│ ─────────────────────────  │
│ Drop files anywhere        │
│ ┌────────────────────────┐ │
│ │ 📁 _assets        ⋮  │   │
│ │ 📁 public   public ⋮ │   │
│ │ 📁 scripts        ⋮  │   │
│ │ 🖼 hero.png  2.4 MB ⋮ │   │
│ │ 📄 index.html   ⋮    │   │
│ │ 📄 main.ts      ⋮    │   │
│ └────────────────────────┘ │
│ Storage 312/500 GB   62%   │
├────────────────────────────┤
│  Files   Uploads  Shares  ▸│
└────────────────────────────┘
```

## Wireframe — selection mode with bulk actions (desktop)

```
┌──────────┬───────────────────────────────────────────────────────────┐
│ (sidebar)│  Acme Inc ▾        Jane Doe ▾         (search) (upload)   │
│          │  2 selected                                   [Download]  │
│          │  Root › Projects › Web › 2026                  [Move]     │
│          │  ┌──────────────────────────────────────────────────────┐ │
│          │  │ ☐ Name               Size   Visibility  Modified  ⋮  │ │
│          │  │ ☑ 🖼 hero.png         2.4 MB  public     Mar 12    ⋮  │ │
│          │  ☑ 📄 index.html      18.2 KB  ─         Mar 12    ⋮  │   │
│          │  │ ☐ 📄 main.ts         9.8 KB   ─         Mar 11    ⋮  │ │
│          │  └──────────────────────────────────────────────────────┘ │
│          │                          [Delete (2)]                     │
└──────────┴───────────────────────────────────────────────────────────┘
```

## Interactions

- **Breadcrumbs are the navigation** — each segment links up the hierarchy; beyond 4 levels it collapses to `Root / … / Parent / Current` with the ellipsis opening a dropdown of elided segments. Every segment is a drop target for drag-to-move.
- **Dropzone** wraps the whole table area — dropping anywhere uploads into the current folder. Drag-over → primary border/fill. Rejected (size/type) → error border, message reverts after 3s. Name collision → inline replace/keep-both/cancel per file, queue stays open (PRD MV-E2).
- **Row actions** (hover on pointer, overflow `⋮` on touch): Share, Permissions (folders), Move, Rename, Extract (archive), Download, Delete. Delete of a folder states the **real object count** and requires typed-name confirmation (PRD TA-E1); above a threshold it becomes a tracked job → toast links to `jobs`.
- **Selection mode**: checkbox column + bulk bar replaces the filter bar (Download / Move / Delete). Bulk delete of 2+ items confirms once, stating the count.
- **Sort/filter** are server-side (cursor pagination, no offset — `web/components.md`). Sort by name/size/modified; filter by visibility/type/uploader.
- **Quota**: under the table at all times. ≥80% → `color.warning`; exceeded → `color.error`, uploads blocked with "Quota exceeded — reads still work, new uploads are paused" and a link to `usage`.
- **Permission awareness**: a viewer sees write actions **disabled with a reason tooltip** ("Requires `write` on this folder"), never hidden, never an error wall (PRD MV4).
- **Loading**: skeleton rows at `size.row-height`; never a spinner. **Error**: inline banner with retry, table blank (no stale rows). **Empty folder**: distinct copy — "This folder is empty · Drop files in or press Upload" vs "No matches for this filter".
- **Virtualized rows**: a 10,000-object folder renders a windowed list with sticky header — the page never scrolls horizontally (PRD G3).
- **Breakpoints**: sidebar collapses to 72px nav rail below ~1200px, to a drawer under 768px (no bottom tabs on web — it's an admin surface). The table becomes stacked cards under 640px.
