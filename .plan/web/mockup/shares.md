# Screen — shares

## Purpose & context

- **User goal**: see every active/expired/revoked share link in the tenant and revoke or re-share them. The "is this link still being used" answer lives here (`share_links.hit_count`, `last_accessed_at`).
- **Entry points**: sidebar Shares; "link created" from `share-dialog`.
- **Exit points**: click a linked file → `file-detail`; Revoke → `confirm-destructive` (or a non-typed confirm — revoking a share isn't irreversible-destructive in the same class, but it does kill the link, so a plain confirm). Share action → `share-dialog`.
- **Data needed**: `share_links` — target object (name + `file_id`), `kind` (signed/short), expiry, `hit_count`, `last_accessed_at`, `created_by`, `revoked_at`. Filter by `?status: active|expired|revoked`.

## States

- [x] Loading (skeleton rows)
- [x] Empty (no shares yet)
- [x] Populated
- [x] Filtered empty ("no expired links" ≠ "no links")
- [x] Error
- [x] Domain-specific — a share's target was deleted (row flagged, link effectively dead)

## Style reference

- **Components used**: `AppShell`, `DataTable` (with selection for bulk revoke), `StatusPill` (active/expired/revoked), `ByteSize`/`RelativeTime`, `CopyableCode` (short URL), `EmptyState`, `ConfirmDestructive`.
- Status filter tabs: All / Active / Expired / Revoked.
- No token deltas.

## Wireframe — desktop (populated)

```
┌──────────┬───────────────────────────────────────────────────────────┐
│ (sidebar)│  Acme Inc ▾         Jane Doe ▾                            │
│          │  Shares                                        [+ Share]  │
│          │  [All] [Active] [Expired] [Revoked]                       │
│          │  ┌──────────────────────────────────────────────────────┐ │
│         │  │ File           Kind      Created    Hits  Status  ⋮   │ │
│         │  │ 🖼 hero.png    short      Mar 12    12    active   ⋮   │ │
│         │  │ 📄 q3-report   signed 4h  Mar 13    3     active   ⋮  │ │
│         │  │ 📦 archive.zip short      Mar 08    45    expired  ⋮  │ │
│         │  │ 📄 roadmap.md  signed     Feb 27    0     revoked  ⋮  │ │
│          │  └──────────────────────────────────────────────────────┘ │
│          │  Showing 4 of 18          ‹ 1 2 ›                         │
│          │                                                           │
│          │  Active links expire or can be revoked any time.          │
│          │  Public objects are not listed here — visibility lives    │
│          │  on the file itself.                                      │
└──────────┴───────────────────────────────────────────────────────────┘
```

## Wireframe — mobile (populated, stacked cards)

```
┌────────────────────────────┐
│ ☰ ⬡    Acme Inc ▾   👤     │
├────────────────────────────┤
│ Shares          [+ Share]  │
│ [All] [Active] [Expired]   │
│ ─────────────────────────  │
│ ┌────────────────────────┐ │
│ │ 🖼 hero.png   active    │ │
│ │ short · Mar 12         │ │
│ │ 12 hits · last 2h   ⋮  │ │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │ 📄 q3-report   active  │ │
│ │ signed 4h · Mar 13     │ │
│ │ 3 hits · last 1h    ⋮  │ │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │ 📦 archive.zip  expired│ │
│ │ short · Mar 08         │ │
│ │ 45 hits · last 3d   ⋮  │ │
│ └────────────────────────┘ │
│ ─────────────────────────  │
│ Active links expire or     │
│ can be revoked any time.   │
└────────────────────────────┘
```

## Interactions

- **Row actions** (`⋮` / hover): Copy URL, Open (→ `file-detail` of the target), Revoke. Revoking an active link is a **plain confirm, not typed-name** — but the confirm says exactly what dies ("13 people have opened this in the last week — revoking kills the link now").
- **Selection + bulk revoke**: checkbox column for the active tab; bulk revoke confirms once with the count.
- **Hits column** answers the revoke question: a link with `45 hits · last 3d` is load-bearing; one with `0 hits · never` is safe to revoke. Sorting by Hits is default.
- **Status filter** maps to `?status=` and is a real URL state, so a filter isn't lost on refresh.
- **Target deleted**: if a share's object is soft-deleted, the row shows a muted "File in trash" tag and Copy URL becomes disabled — copying a dead URL is worse than not offering it.
- **Empty states differ**: no shares → "Nothing shared yet · Press Share on any file to create a link"; expired filter → "No expired links · Nothing here has outlived its TTL" — the second must not read like the first.
- **Public objects note**: the footer caption (in both wireframes) tells users public visibility lives on the file, not in this list — otherwise people search here for files they made public and can't find them.
