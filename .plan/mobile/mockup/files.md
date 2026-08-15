# Screen — files

## Purpose & context

- **User goal**: browse folders and files in the current tenant, upload from the phone, and act on items — the default landing after login, the reason the app exists.
- **Entry points**: default tab after login/tenant-switch; tab tap; deep link `/files/:folderId`.
- **Exit points**: tap a file → `file-detail`; tap a folder → `files` at that folder; breadcrumb segment → `files` at that level; search icon → `search`; FAB → `source-picker` sheet; long-press row → actions incl. Move → `folder-picker` sheet, Delete → `confirm-destructive`; tap a share → `file-detail` (from `shares`).
- **Data needed**: `folders` + `objects` for the current folder (name, type, size, visibility, modified), quota bar, `ancestors` for breadcrumbs.

## States

- [x] Loading (skeleton rows at `size.row-height`)
- [x] Empty folder
- [x] Populated (happy path)
- [x] Search-empty (from a filter)
- [x] Error (inline banner + retry; the tab bar survives so navigation isn't dead)
- [x] Domain-specific — permission-limited viewer (write actions disabled with reason)
- [x] Domain-specific — near/exceeded quota (warning/error banner, uploads blocked)
- [x] Domain-specific — offline (queue persists; reads show cached-or-blocked)
- [x] Domain-specific — long-press selection with action sheet

## Style reference

- **Components used**: 4-tab shell (`size.navbar-height` 64px), app bar with tenant switcher + search, breadcrumbs (collapsible), list rows, FAB (Camera/Files), bottom sheets (`source-picker`, `folder-picker`, `confirm-destructive`), `UploadQueue` access via the `uploads` tab.
- Mobile is a list, not a table — rows carry type icon, name, size/visibility line, and a trailing `⋮`. **No tree pane**: breadcrumbs + folder rows are the navigation (`web/navigation.md`'s breadcrumb decision applies to both platforms).
- No token deltas.

## Wireframe — mobile (populated)

```
┌────────────────────────────┐
│  ☰  Acme Inc ▾      🔍     │
├────────────────────────────┤
│  Root › Projects › Web›2026│
│ ────────────────────────── │
│ ┌────────────────────────┐ │
│ │ 📁 _assets        ⋮    │ │
│ │   Mar 12               │ │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │ 📁 public    public  ⋮ │ │
│ │   Mar 12               │ │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │ 📁 scripts        ⋮    │ │
│ │   Mar 11               │ │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │ 🖼 hero.png     2.4 MB ⋮│ │
│ │   public · Mar 12      │ │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │ 📄 index.html  18.2 KB⋮│ │
│ │   Mar 12               │ │
│ └────────────────────────┘ │
│ ────────────────────────── │
│ Storage 312/500 GB   62%   │
│                   ➕  FAB  │
├────────────────────────────┤
│ ●Files   Uploads  Shares  ▸│
└────────────────────────────┘
```

## Wireframe — mobile (empty folder)

```
┌───────────────────────────┐
│  ☰  Acme Inc ▾      🔍    │
├───────────────────────────┤
│  Root › Projects › 2026   │
│ ───────────────────────── │
│                           │
│      📁                   │
│                           │
│  This folder is empty     │
│                           │
│  Drop files in from the   │
│  picker, or upload from   │
│  this phone.              │
│                           │
│  ┌──────────────────────┐ │
│  │   Upload a file      │ │
│  └──────────────────────┘ │
│                  ➕       │
├───────────────────────────┤
│ ●Files  Uploads  Shares  ▸│
└───────────────────────────┘
```

## Interactions

- **Breadcrumbs**: tap any segment to jump up; long content collapses to `Root › … › Parent › Current` with the elided segment opening a picker.
- **FAB** → `source-picker` sheet (Camera / Files). Permission denied (camera/photos) → inline explanation with "Open Settings" — not a dead end.
- **Long-press a row** → selection mode with a contextual action bar (Download / Share / Move / Delete), or directly the action sheet — both acceptable; the action bar must include the same set as the sheet so no action is reachable only one way.
- **Upload flow**: pick → returns to `uploads` with items queued → progress lives there; on this screen only the FAB's availability changes (and quota banners).
- **Quota**: ≥80% → warning banner under breadcrumbs; exceeded → error banner, FAB disabled with "Quota exceeded — reads still work", plus a link to `usage`.
- **Permission-aware**: a viewer taps a disabled action → the row shows a caption ("Requires `write` on this folder") instead of an error toast; disabled rows aren't hidden (PRD MV4).
- **Offline**: pull-to-refresh shows "Offline — showing last loaded" with a banner; queued uploads persist in the `uploads` queue and sync on reconnect.
- **Loading**: skeleton rows; never a spinner over the whole list. **Error**: inline banner + Retry, list area blank (no stale rows).
- **Empty vs filtered-empty**: "This folder is empty · Drop files in from the picker" vs "No matches in this folder" — different copy, same layout (`design-collection/mobile-screen/patterns.md` Empty state).
- **Pull-to-refresh** on the list; the tab bar is always live (Files/Uploads/Shares/More) so navigation never dead-ends.
