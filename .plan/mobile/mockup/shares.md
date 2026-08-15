# Screen — shares

## Purpose & context

- **User goal**: see the tenant's share links from the phone — active ones, their hit counts, and revoke an urgent one (the "shared something and need to un-share it now" case). Sharing from a phone is a primary use case (PRD MV2/MV3).
- **Entry points**: tab tap; a link created in `share-sheet` lands the user here.
- **Exit points**: tap a linked file → `file-detail`; revoke → `confirm-destructive`-lite; + Share → opens a file picker → `share-sheet` (share a *new* file).
- **Data needed**: `share_links` — target file (name + `file_id`), kind (signed/short), expiry, `hit_count`, `last_accessed_at`, `created_by`, revoked state.

## States

- [x] Loading (skeleton rows)
- [x] Empty (nothing shared yet)
- [x] Populated
- [x] Error
- [x] Domain-specific — expired/revoked links shown muted in a separate section (not mixed with active)

## Style reference

- **Components used**: list rows (type icon, name, kind + expiry line, hit count, status pill), action sheet on `⋮` (Copy URL / Open / Revoke), `+` share affordance. Active links first, expired/revoked grouped below under a muted header.
- No token deltas.

## Wireframe — mobile (populated)

```
┌────────────────────────────┐
│  Acme Inc ▾        [+]     │
├────────────────────────────┤
│  ACTIVE (2)                │
│ ┌────────────────────────┐ │
│ │ 🖼 hero.png   active    │ │
│ │ short · 12 hits · 2h   │ │
│ │ https://blb.io/s/kW9 ⋮ │ │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │ 📄 q3-report   active  │ │
│ │ signed · exp in 4h     │ │
│ │ 3 hits · last 1h     ⋮ │ │
│ └────────────────────────┘ │
│  EXPIRED / REVOKED (2)     │
│ ┌────────────────────────┐ │
│ │ 📦 archive.zip  expired│ │
│ │ short · 45 hits · 3d   │ │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │ 📄 roadmap.md  revoked │ │
│ │ signed · Mar 02        │ │
│ └────────────────────────┘ │
│ ────────────────────────── │
│ Public objects live on     │
│ the file, not here.        │
├────────────────────────────┤
│  Files  Uploads  ●Shares  ▸│
└────────────────────────────┘
```

## Wireframe — mobile (empty)

```
┌───────────────────────────┐
│  Acme Inc ▾        [+]    │
├───────────────────────────┤
│                           │
│        🔗                 │
│                           │
│  Nothing shared yet       │
│                           │
│  Share any file to create │
│  a link — signed links    │
│  expire, short URLs stick │
│  until revoked.           │
│                           │
│  ┌──────────────────────┐ │
│  │  Share a file        │ │
│  └──────────────────────┘ │
├───────────────────────────┤
│  Files Uploads  ●Shares  ▸│
└───────────────────────────┘
```

## Interactions

- **Revoke**: `⋮` → Revoke with a plain confirm stating the consequence ("12 people have opened this this week — revoking kills the link now"). No Undo toast (revocation is deliberate, `design/style-guide.md`).
- **Hits column** is the revoke decision: `45 hits · last 3d` is load-bearing; `0 hits · never` is safe to revoke. Sort active links by hits desc by default.
- **Copy URL** copies the full link with confirmation.
- **Tap the linked file** → `file-detail` (which shows the same link's management).
- **`+` Share** → file picker → `share-sheet` — sharing a *new* file from this tab.
- **Expired/revoked section**: muted, no actions (Copy is disabled for a dead link). Deliberately grouped so the active list is scannable.
- **Footer caption**: "Public objects live on the file, not here" — prevents the "I made it public, why isn't it in Shares" confusion (`web/mockup/shares.md` has the same rule).
- **A11y**: status pills carry a text label (not color alone); rows are 48dp targets.
