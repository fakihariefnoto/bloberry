# Screen — file-detail

## Purpose & context

- **User goal**: see one object's full metadata and act on it — share, download, move, change visibility, delete (PRD MV2, MV3).
- **Entry points**: tap a file row in `files`; tap a completed item in `uploads`; tap a share in `shares`; search result.
- **Exit points**: back → `files` at the file's folder; Share → `share-sheet`; Move → `folder-picker`; Delete → `confirm-destructive`; make public → `confirm-destructive`; created link → `shares` (per `navigation.md` flow).
- **Data needed**: `objects` — name, `file_id`, size, content type, visibility, uploader, created/modified, backend, hash; current grants (from parent folder), share links.

## States

- [x] Loading (skeleton detail)
- [x] Populated (happy path)
- [x] Error — file missing/deleted (404)
- [x] Domain-specific — viewer without write (actions disabled with reason)
- [x] Domain-specific — public (warning banner)
- [x] Domain-specific — soft-deleted (restore/delete banner, S5)

## Style reference

- **Components used**: detail screen pattern (hero content first, actions clearly separated), sticky bottom action bar, `share-sheet`/`folder-picker`/`confirm-destructive` sheets. `CopyableCode` for `file_id`.
- Mobile: the file-type icon is the hero (no preview in v1, NG3); metadata as a labeled list; actions in a sticky footer.
- No token deltas.

## Wireframe — mobile (populated, public)

```
┌───────────────────────────┐
│  ← Back            ⋮      │
├───────────────────────────┤
│      ┌───────────┐        │
│      │    🖼     │         │
│      └───────────┘        │
│  hero.png                 │
│  2.4 MB · PNG · Mar 12    │
│  ⚠ Public                 │
│ ───────────────────────── │
│  file_id                  │
│  f_8Kd2pQxL31A   [copy]   │
│  Storage   s3-eu-prod     │
│  Backend   healthy        │
│  Uploaded  Jane Doe       │
│  Created   Mar 12, 09:41  │
│  Modified  Mar 12, 14:03  │
│ ───────────────────────── │
│  Shared                   │
│ https://blb.io/s/kW9[copy]│
│  12 hits · last 2h[revoke]│
│  signed link ·expin4h[rev]│
├───────────────────────────┤
│  [Share]  [Download]  [⋯] │
└───────────────────────────┘
```

## Wireframe — mobile (file not found)

```
┌───────────────────────────┐
│  ← Back                   │
├───────────────────────────┤
│                           │
│      🗑                    │
│                           │
│  This file no longer      │
│  exists.                  │
│                           │
│  It may have been deleted │
│  by someone with access.  │
│                           │
│  ┌──────────────────────┐ │
│  │  Back to files       │ │
│  └──────────────────────┘ │
└───────────────────────────┘
```

## Interactions

- **Share** → `share-sheet`: pick signed link + TTL, or short URL, or make public. Created → OS share sheet opens with the URL pre-filled; the link lands in `shares`. Cancel → back to detail, nothing created. Make-public goes through `confirm-destructive` first (public is effectively irreversible once copied, `TRD.md` R11).
- **Download**: streams via the backend's signed path (redirect for cloud drivers, proxy for disk — `architecture.md` §3.2); progress + save handled by the platform downloader.
- **Move** → `folder-picker` (the one place the folder tree appears on mobile; moved node and descendants disabled as targets — PRD TA-E2 cycle prevention). The `file_id` survives the move (PRD M4).
- **Delete** → `confirm-destructive`; soft-delete by default (S5). A public file's delete warns the public URL dies immediately.
- **Copy `file_id`** via the `CopyableCode` affordance with a confirmation.
- **Sticky footer** holds Share / Download / ⋮ (Move, Rename, Delete, Make public/unpublish) — the primary action (Share) is always in thumb reach.
- **Permission-aware**: a viewer without write sees the actions disabled with a caption naming what's needed; Download stays enabled.
- **Soft-deleted**: banner "This file is in the trash · Restore or permanently delete", both actions in place.
- **A11y**: the file-type icon is `role="presentation"`; metadata rows are a labeled list, not div soup; actions are 48dp min.
