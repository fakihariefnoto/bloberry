# Screen — admin-tenant-detail

## Purpose & context

- **User goal**: platform admin's full view of one tenant — its config (quota, backend), its members/keys at a glance, its usage trend, and the suspend/reactivate control (PRD PA2, PA-E2).
- **Entry points**: click a tenant in `admin-tenants`; click a tenant row in `admin-usage`.
- **Exit points**: back → `admin-tenants`; member/key rows → nothing (this screen summarizes; full management lives in the tenant's own admin surfaces); Suspend → confirm; usage range → in place.
- **Data needed**: `tenants`, `memberships` (count, owner), `access_keys` (count by state), `usage_snapshots`, `storage_backends` (the tenant's assigned backend + health).

## States

- [x] Loading (skeleton)
- [x] Populated
- [x] Error (tenant deleted between list and detail)
- [x] Domain-specific — suspended (banner, controls disabled)
- [x] Domain-specific — backend unreachable (warning panel with the real provider error — PRD PA-E1; the one documented exception to never passing provider errors through)

## Style reference

- **Components used**: `PageHeader` with back, stat cards (usage), small breakdown table, `StatusPill`, `ByteSize`, `ConfirmDestructive` (suspend/reactivate), sectioned panels.
- Read-mostly summary — the "manage" actions that live here are suspend/reactivate and backend reassignment; everything else links out to the tenant's own admin surfaces or the platform-level screens.
- No token deltas.

## Wireframe — desktop (populated)

```
┌──────────┬───────────────────────────────────────────────────────────┐
│ (sidebar)│  ‹ Back to tenants                                        │
│          │  Acme Inc                [Suspend]   [Reassign backend]   │
│          │  acme · active · since Mar 01, 2026                       │
│          │  ⚠ s3-eu-prod unreachable — SigV4 auth failed, credentials│
│          │    may be revoked.                                        │
│          │  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐       │
│         │  │ 312 GB stored│ │ 48,912 objects│ │ $21.40 est.  │       │
│         │  │  ▁▃▅▂▇  ↑6%  │ │  ▂▄▆▅▇   ↑3%  │ │ this month   │       │
│          │  └──────────────┘ └──────────────┘ └──────────────┘       │
│          │  ┌──────────────────────┐ ┌─────────────────────────────┐ │
│         │  │  CONFIG                │ │  PEOPLE & ACCESS           │ │
│         │  │  Backend    s3-eu-prod │ │  Members         5 (1owner)│ │
│         │  │  Health     unreachable│ │  Pending invites 1         │ │
│         │  │  Quota      500 GB     │ │  Applications    3         │ │
│         │  │  Objects    1,000,000  │ │  Keys     3 active / 1 rev │ │
│         │  │  Used       312 GB     │ │                            │ │
│          │  └──────────────────────┘ │  Full management lives in   │ │
│          │                           │  the tenant's own surfaces. │ │
│          │  ┌──────────────────────┐ └─────────────────────────────┘ │
│          │  │  USAGE (this month)    │                             │ │
│          │  │   Mar 01  284 GB $17.80│                             │ │
│          │  │   Mar 08  301 GB $19.10│                             │ │
│          │  │   Mar 15  312 GB $21.40│                             │ │
│          │  └─────────────────────┘                                │ │
└──────────┴───────────────────────────────────────────────────────────┘
```

## Wireframe — mobile (populated)

```
┌────────────────────────────┐
│ ☰ ⬡    Platform ▾    👤    │
├────────────────────────────┤
│ ‹ Tenants                  │
│ Acme Inc      [Suspend]    │
│ acme · active              │
│ ⚠ backend unreachable      │
│ ─────────────────────────  │
│ ┌────────────────────────┐ │
│ │  312 GB stored   ↑6%   │ │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │  $21.40 est this month │ │
│ └────────────────────────┘ │
│ ─────────────────────────  │
│ CONFIG                     │
│ Backend   s3-eu-prod  ⚠    │
│ Quota     500 GB           │
│ Used      312 GB (62%)     │
│ [Reassign backend]         │
│ ─────────────────────────  │
│ PEOPLE & ACCESS            │
│ Members 5 · 1 pending      │
│ Apps 3 · keys 3/1 active   │
│ ─────────────────────────  │
│ USAGE (month)              │
│ Mar 15  312 GB  $21.40     │
│ Mar 08  301 GB  $19.10     │
│ Mar 01  284 GB  $17.80     │
└────────────────────────────┘
```

## Interactions

- **Backend unreachable banner**: shows the **real provider error** (PRD PA-E1 — the single documented exception to the never-raw-provider-errors rule, `backend/domains.md` §8). Read-only surface for a `platform_admin`; the fix (re-enter credentials) lives in `admin-backend-detail`.
- **Suspend / Reactivate**: toggle with `ConfirmDestructive`-lite; the confirm states the tenant's users will see writes blocked (reads still work, PRD PA-E2). Suspended state renders the banner + disabled controls.
- **Reassign backend**: opens a backend picker; the confirm carries the ADR-4 contract ("new objects go to the new backend; existing objects keep resolving") — this is the platform-admin equivalent of `tenant-settings`'s backend change.
- **Stat cards**: skeleton on load; sparklines + direction per the analytics pattern.
- **People & Access panel** is summary-only with counts; "Full management lives in the tenant's own surfaces" is stated (a platform admin who tries to edit a member's role from here must be told where to go).
- **Usage breakdown**: same rate-card rule as `usage` — missing rate card renders "unknown", never $0.
- **Delete tenant** lives in `admin-tenants` (row action) and `tenant-settings` (owner's DANGER ZONE), not here — this screen is read/summary + suspend; deleting from a summary page invites fat-finger catastrophe.
- **A11y**: the unreachable-banner error is announced on arrival; the Suspend control returns focus to the banner when confirmed.
