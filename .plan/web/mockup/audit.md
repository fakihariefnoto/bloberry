# Screen — audit

## Purpose & context

- **User goal**: answer "where did this file go / who did what, when, from where" (PRD TA6). The tenant's append-only record of object and permission mutations.
- **Entry points**: sidebar Audit (`tenant_admin`+); "see raw events" from `usage`; a deep link `?target=<file_id>` from `file-detail`.
- **Exit points**: click a target → `file-detail` (when the target is an object); filter/sort in place. Nothing here navigates away by default.
- **Data needed**: `audit_events` — `action` (object.upload/delete/share, grant.create, key.revoke, member.join…), `principal_type`/`principal_id`, `target_type`/`target_id`, `metadata`, `ip`, `user_agent`, `created_at`. Filter by `?from,to,action,principal`.

## States

- [x] Loading (skeleton rows)
- [x] Empty (no events in range)
- [x] Populated
- [x] Error
- [x] Domain-specific — "no events match this filter" vs "no events at all"
- [x] Domain-specific — the redirect-download limitation (audit records link *issuance*, not each byte-read — PRD D1/ADR-3) — must be stated, not discovered

## Style reference

- **Components used**: `AppShell`, `DataTable`, `DateRangePicker` (24h/7d/30d/custom), `StatusPill`-adjacent action tags, `RelativeTime`, `CopyableCode` (`file_id`s, principal IDs), `EmptyState`.
- The date range is a real `?from&to` URL state. This is a read-heavy, filter-heavy screen — the DataTable's server-side cursor pagination applies.
- No token deltas.

## Wireframe — desktop (populated)

```
┌──────────┬───────────────────────────────────────────────────────────┐
│ (sidebar)│  Acme Inc ▾         Jane Doe ▾                            │
│          │  Audit log                                                │
│          │  [24h ▾]  [Action ▾]  [Principal ▾]  [Clear]  [Export]    │
│          │  ┌──────────────────────────────────────────────────────┐ │
│          │  │ Time        Action          Principal   Target     ⋮ │ │
│          │  │ 14:03:12    object.upload   jane        hero.png   ⋮ │ │
│          │  │ 14:02:55    object.visibility jane       hero.png   ⋮│ │
│          │  │ 13:58:41    share.create    jane        kW9        ⋮ │ │
│          │  │ 13:41:09    key.revoke      sam         blob_••c9e7 ⋮│ │
│          │  │ 13:20:33    grant.create    jane        f_8Kd2p    ⋮ │ │
│          │  │ 12:58:04    object.delete   acme-cms    index.html ⋮ │ │
│          │  └──────────────────────────────────────────────────────┘ │
│          │  Showing 6 of 1,284          ‹ 1 2 3 … 214 ›              │
│          │                                                           │
│          │  Downloads on the redirect path record link issuance,     │
│          │  not each byte read (default is 5-min signed URLs).       │
└──────────┴───────────────────────────────────────────────────────────┘
```

## Wireframe — mobile (populated)

```
┌────────────────────────────┐
│ ☰ ⬡    Acme Inc ▾   👤     │
├────────────────────────────┤
│ Audit log   [24h] [Filter] │
│ ─────────────────────────  │
│ ┌────────────────────────┐ │
│ │ 14:03:12  upload       │ │
│ │ jane → hero.png        │ │
│ │ ▸ details              │ │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │ 13:41:09  key.revoke   │ │
│ │ sam → blob_••c9e7      │ │
│ │ ▸ details              │ │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │ 12:58:04  delete       │ │
│ │ acme-cms → index.html  │ │
│ │ ▸ details              │ │
│ └────────────────────────┘ │
│ ─────────────────────────  │
│ Downloads on the redirect  │
│ path record link issuance. │
└────────────────────────────┘
```

## Interactions

- **Filtering**: date range, action type, principal. All filter controls produce URL state (`?from&to&action&principal`) so an investigation can be shared/bookmarked. Export produces CSV of the *current filter*, with a caption of the range.
- **Row expand** (mobile cards, desktop row) shows the full event: `metadata` keyed per action, IP, user-agent. An `object.upload` expands to size/hash/backend; a `share.create` expands to the link + TTL.
- **Target link**: object targets link to `file-detail`; a deleted target opens the "file no longer exists" state but keeps the event context in the breadcrumb.
- **The redirect-download limitation is stated in a footer caption on every width** (not buried in a tooltip) — this is the documented honesty of ADR-3/`ERD.md` share-links note, and the one thing about this screen that will surprise an operator relying on it during an incident.
- **Empty states**: no events at all → "Nothing has happened yet · Mutations of files, keys and grants appear here"; no events in range → "No events in this range · Widen the date window or clear a filter". Distinct, with a Clear-filter action on the second.
- **Retention**: events older than the configured window (default 365 days) are hard-deleted by the monthly retention job (`ERD.md` Q2 resolution) — a muted footnote "Events are retained for 365 days" so a long look-back isn't silently empty.
- **Permission-aware**: `member`/`viewer` never see the route; `tenant_admin` sees the full tenant log; platform admins see only their own tenant's here — install-wide audit is a `platform_admin` surface (`admin-usage`/admin tenant detail).
