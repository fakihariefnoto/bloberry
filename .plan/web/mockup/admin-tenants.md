# Screen — admin-tenants

## Purpose & context

- **User goal**: the platform admin's tenant overview — create tenants, see each one's storage/quota/cost footprint, and drill in (PRD PA2, PA3). The "who's causing this bill" view.
- **Entry points**: sidebar Tenants (`platform_admin` only; visually separated PLATFORM group).
- **Exit points**: click a tenant → `admin-tenant-detail`; New tenant → modal (name, slug, quota, default backend); from `admin-usage` a tenant row → this list's detail.
- **Data needed**: `tenants` + derived `usage_snapshots` — name, slug, status, `used_bytes`, `used_objects`, quota, default backend, est. monthly cost (rate card).

## States

- [x] Loading (skeleton rows)
- [x] Empty (first install — no tenants yet)
- [x] Populated
- [x] Error
- [x] Domain-specific — suspended tenant (flagged row)
- [x] Domain-specific — tenant over quota (error-tinted used/quota)

## Style reference

- **Components used**: `AppShell`, `DataTable` (with the platform's most important sortable columns — used bytes, est. cost), `StatusPill`, `ByteSize`, `EmptyState`, create-tenant modal.
- This is a data-dense table page: `size.row-height`, `space.sm` row padding, sticky header.
- No token deltas.

## Wireframe — desktop (populated)

```
┌──────────┬───────────────────────────────────────────────────────────┐
│ (sidebar)│  Platform ▾        [👤] Platform Admin ▾                  │
│          │  Tenants                                   [+ New tenant] │
│          │  ┌──────────────────────────────────────────────────────┐ │
│          │  │ Tenant     Status   Used          Quota    Est. cost ⋮ │
│          │  │ Acme Inc   active   312 GB · 48.9k  500 GB  $21.40  ⋮  │
│          │  │ Folio Notes active   1.1 TB · 284k   2 TB   $64.10  ⋮  │
│          │  │ Masjid App  active   812 MB · 3.2k   10 GB  $0.28   ⋮  │
│          │  │ ⚠ Legacy    suspended 2.4 GB · 12k   5 GB   $0.90  ⋮   │
│          │  │ ⚠ Kercis    over quota 4.9 GB · 22k  5 GB   $4.10  ⋮   │
│          │  └──────────────────────────────────────────────────────┘ │
│          │  Showing 5 of 9            ‹ 1 2 ›                        │
│          │  Est. cost is computed from each tenant's backend rate    │
│          │  card and metering — see admin-usage for the install.     │
└──────────┴───────────────────────────────────────────────────────────┘
```

## Wireframe — mobile (populated)

```
┌────────────────────────────┐
│ ☰ ⬡    Platform ▾    👤    │
├────────────────────────────┤
│ Tenants          [+ New]   │
│ ─────────────────────────  │
│ ┌────────────────────────┐ │
│ │ Acme Inc    active     │ │
│ │ 312 GB / 500 GB  · 62% │ │
│ │ est. $21.40         ⋮  │ │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │ Folio Notes  active    │ │
│ │ 1.1 TB / 2 TB    · 55% │ │
│ │ est. $64.10        ⋮   │ │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │ ⚠ Legacy    suspended  │ │
│ │ 2.4 GB / 5 GB   · 48%  │ │
│ │ est. $0.90         ⋮   │ │
│ └────────────────────────┘ │
│ Est. cost from rate  cards.│
└────────────────────────────┘
```

## Interactions

- **New tenant**: modal — name, slug (auto-suggested from name, editable), quota, default backend (from install-level pool). On save → `admin-tenant-detail` so the admin can immediately invite the owner (`ERD.md`: a tenant with no owner is inert).
- **Sorting**: default by Est. cost desc — the platform admin's actual question is "what's burning money". Used bytes sortable; both server-side.
- **Statuses**: `suspended` flagged with the warning pill; **over quota** shows the used/quota pair in `color.error` ("4.9 GB / 5 GB") — reads + write-blocked is a state the admin should notice before the tenant's users do (PRD PA-E2).
- **Row click** → `admin-tenant-detail`. Row `⋮` → Suspend/Reactivate, Delete.
- **Delete tenant**: typed-name confirm with the orphan-bytes consequence (same as `tenant-settings` DANGER ZONE — bytes stay in the bucket, reconciliation finds them).
- **Empty state**: "No tenants yet · Create the first tenant to provision storage" with the New button.
- **The footer caption** distinguishes this screen (per-tenant) from `admin-usage` (install-wide) — the two cost views must not read as duplicates.
- **Permission-aware**: `platform_admin` only; the PLATFORM section is hidden for everyone else and the routes are server-guarded.
