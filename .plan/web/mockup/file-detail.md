# Screen — file-detail

## Purpose & context

- **User goal**: see everything about one object — its stable `file_id` (PRD M4/G4), how it's shared, its permissions, its size/type/uploader — and act on it (download, share, move, delete, change visibility).
- **Entry points**: click a file row in `files`; click a linked file from `shares`; open a share link's source object from `audit`. Also reachable from mobile but this is the web detail view.
- **Exit points**: back/breadcrumb → `files` at the file's folder; Share → `share-dialog`; Move → `move-picker`; Delete → `confirm-destructive`; make public → `confirm-destructive` variant; "see raw events" links to `audit?target=<file_id>`.
- **Data needed**: from `ERD.md` `objects` — `name`, `id` (the public `file_id`), `size_bytes`, `content_type`, `content_hash`, `visibility`, `backend_id`, `uploaded_by`, `created_at`, `updated_at`, `deleted_at`. Plus current grants, share links (`share_links`), and audit events for this object.

## States

- [x] Loading (skeleton detail)
- [x] Populated (happy path)
- [x] Error (file missing/deleted/not-found 404)
- [x] Domain-specific — viewer without permission (download works, mutations disabled with reason)
- [x] Domain-specific — public visibility (warning banner — public is a caution)
- [x] Domain-specific — soft-deleted (S5)

## Style reference

- **Components used**: `PageHeader` with back, `Breadcrumbs`, `CopyableCode` (`file_id`, signed link, short URL), `StatusPill` (visibility, backend health), `ByteSize`, `RelativeTime`, `ConfirmDestructive`, `Toast`. Detail layout: left info column + right action column on desktop, stacked on mobile.
- No token deltas.

## Wireframe — desktop (populated, public file)

```
┌──────────┬───────────────────────────────────────────────────────────┐
│ (sidebar)│  ← Back to files                                          │
│          │  Root › Projects › Web › 2026 › hero.png                  │
│          │                                                           │
│          │  🖼 hero.png                          [Download]  [Share]  │
│          │  2.4 MB · PNG image · Mar 12, 2026  [Move] [Delete]       │
│          │  ⚠ Public — anyone with the link can view.                │
│          │  ───────────────────────────────────────────────────────  │
│          │  ┌──────────────────────────┐  ┌────────────────────────┐ │
│          │  │  Details                 │  │  Sharing               │ │
│          │  │  file_id  f_8Kd2pQxL31A  │  │  Short URL             │ │
│          │  │  ┌────────────────────┐  │  │  https://blb.io/s/kW9  │ │
│          │  │  │  copy  mask  copy  │  │  │  [copy]                │ │
│          │  │  └────────────────────┘  │  │  12 hits · last 2h     │ │
│          │  │  Storage   s3-eu-prod    │  │  [Revoke]              │ │
│          │  │  Backend   healthy       │  │  Signed links          │ │
│          │  │  Hash      sha256:9f2a…  │  │  · exp in 4h — [revoke]│ │
│          │  │  Uploaded  Jane Doe      │  │  + New signed link     │ │
│          │  │  Modified  Mar 12, 14:03 │  │                        │ │
│          │  │  Created   Mar 12, 09:41 │  │ [Makepublic][Unpublish]│ │
│          │  └──────────────────────────┘  └────────────────────────┘ │
│          │  ───────────────────────────────────────────────────────  │
│          │  Permissions on this file                                 │
│          │  acme-cms (application) — read · write — grants on 2026/  │
│          │  [Manage permissions]                                     │
└──────────┴───────────────────────────────────────────────────────────┘
```

## Wireframe — mobile (populated)

```
┌────────────────────────────┐
│ ☰ ⬡    Acme Inc ▾   👤     │
├────────────────────────────┤
│ ‹ Root › 2026 › hero.png   │
│ ─────────────────────────  │
│ 🖼 hero.png                 │
│ 2.4 MB · PNG · Mar 12      │
│ ⚠ Public                   │
│                            │
│ [Download]  [Share]        │
│ [Move]      [Delete]       │
│ ─────────────────────────  │
│ file_id                    │
│ f_8Kd2pQxL31A   [copy]     │
│ Storage  s3-eu-prod        │
│ Backend  healthy           │
│ Uploaded  Jane Doe         │
│ ─────────────────────────  │
│ Sharing                    │
│ https://blb.io/s/kW9 [copy]│
│ 12 hits · last 2h [revoke] │
│ Signed: exp in 4h  [revoke]│
│ [New signed link]          │
│ [Make public]  [Unpublish] │
└────────────────────────────┘
```

## Interactions

- **Download**: if the file's backend is a signed-redirect backend and the caller has read → streams the provider URL (or proxies for disk, `architecture.md` §3.2). A 410 from an expired link surfaces the human-readable message, never raw provider XML.
- **Share**: opens `share-dialog` with this object pre-filled.
- **Move**: opens `move-picker`; the file keeps its `file_id` forever (PRD M4).
- **Delete**: `confirm-destructive`; soft-delete by default (S5), hard-delete if the retention window is disabled. If the file is public, the delete flow warns that the public URL dies immediately.
- **Make public**: `confirm-destructive` variant — public is effectively irreversible once the URL has been copied (`TRD.md` R11), so the confirm states that consequence explicitly. Unpublish is the reverse and does **not** need typed-name confirmation.
- **CopyableCode**: `file_id`, short URL, signed link all use `text.mono` + copy button with confirmation; short URL and signed link copy the full URL.
- **Permissions**: "Manage permissions" opens `grant-dialog`; grants apply to a folder subtree, so this panel explains that the file inherits grants from `2026/` (deepest ancestor wins, allow-only, `backend/domains.md` §5.2).
- **Audit link**: a "who touched this" affordance jumps to `audit?target=<file_id>`.
- **Permission-aware**: a viewer without write sees Move/Delete/Share/Unpublish **disabled with a tooltip reason**; Download stays enabled. A 404 on the object shows a dedicated "This file no longer exists" state, not the error banner.
- **Soft-deleted (S5)**: banner "This file is in the trash · Restore or permanently delete", both actions in place; restoring re-enables the normal actions.
