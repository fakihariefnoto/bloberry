# Screen — tenant-suspended

## Purpose & context

- **User goal**: understand that the tenant is suspended and what that means. A terminal, single-state page (per `mobile/navigation.md`, this route is shared conceptually with mobile).
- **Entry points**: any authenticated request while the user's tenant has `status: suspended` — a shell-level guard redirects here rather than letting every screen 500.
- **Exit points**: none actionable. "Contact the platform admin" is the only path. Log out still works (user menu stays available).
- **Data needed**: tenant name + status; nothing else — deliberately no data that a suspended tenant's reads would expose (reads keep working per PRD PA-E2, but this page doesn't fetch them).

## States

- [x] Single state — suspended. No loading/empty/error variants worth drawing (the page is the state).

## Style reference

- **Components used**: `AppShell` (sidebar visible but sections inert), a full-content status panel. `color.warning` accent, not `color.error` — suspension is administrative, not a crash.
- No token deltas.

## Wireframe — desktop

```
┌──────────┬──────────────────────────────────────────────────────────┐
│ (sidebar)│  Acme Inc ▾         Jane Doe ▾                           │
│          │                                                          │
│          │            ⚠  This tenant is suspended                   │
│          │                                                          │
│          │    Acme Inc's storage is paused. Reads and writes        │
│          │    are blocked until the platform admin reactivates      │
│          │    the tenant.                                           │
│          │                                                          │
│          │    Your files are safe. Nothing has been deleted.        │
│          │                                                          │
│          │    [  Contact platform admin  ]                          │
│          │                                                          │
│          │    You can still log out and switch to other tenants.    │
└──────────┴──────────────────────────────────────────────────────────┘
```

## Wireframe — mobile

```
┌────────────────────────────┐
│ ☰ ⬡    Acme Inc ▾   👤     │
├────────────────────────────┤
│      ⚠  Tenant suspended   │
│                            │
│  Acme Inc's storage is     │
│  paused. Reads and writes  │
│  are blocked until the     │
│  platform admin reactivates│
│  the tenant.               │
│                            │
│  Your files are safe.      │
│  Nothing has been deleted. │
│                            │
│ [  Contact platform admin ]│
│                            │
│  You can still log out and │
│  switch to other tenants.  │
└────────────────────────────┘
```

## Interactions

- **No destructive or data-fetching actions** — the page reassures and directs. "Contact platform admin" opens the configured support surface (mailto or a support URL from install config), never an in-app form (a suspended tenant shouldn't be able to write anything, including a support ticket).
- **Tenant switcher stays live**: a multi-tenant user can leave the suspended tenant — the caption says so explicitly (otherwise it reads as a full account lockout).
- **Log out** works from the user menu.
- The sidebar renders but every section link is disabled — the user is told the boundary, not shown a broken wall of 500s.
