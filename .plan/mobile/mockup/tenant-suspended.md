# Screen — tenant-suspended

## Purpose & context

- **User goal**: understand that the tenant is suspended and what that means. A terminal, single-state page (shared conceptually with the web `tenant-suspended`).
- **Entry points**: any authenticated request while the tenant has `status: suspended` — the shell guard redirects here instead of letting every screen fail.
- **Exit points**: none actionable beyond Log out and (for multi-tenant users) switching tenants. "Contact platform admin" opens the install's support surface.
- **Data needed**: tenant name + status only — deliberately nothing that requires a read (reads work per PRD PA-E2, but this page fetches nothing to stay trivially reliable).

## States

- [x] Single state.

## Style reference

- **Components used**: minimal status panel over the tab shell (tabs remain visible so the user can still reach `more` → switch tenant / log out), `color.warning` accent (administrative, not a crash). No tab content renders.
- No token deltas.

## Wireframe — mobile

```
┌───────────────────────────┐
│  Acme Inc ▾               │
├───────────────────────────┤
│      ⚠  Tenant suspended  │
│                           │
│  Acme Inc's storage is    │
│  paused. Reads and writes │
│  are blocked until the    │
│ platform admin reactivates│
│  the tenant.              │
│                           │
│  Your files are safe.     │
│  Nothing has been deleted.│
│                           │
│  ┌──────────────────────┐ │
│ │ Contact platform admin│ │
│  └──────────────────────┘ │
│                           │
│  You can still switch to  │
│  other tenants or log out.│
├───────────────────────────┤
│  Files  Uploads  Shares  ▸│
└───────────────────────────┘
```

## Interactions

- **No data-fetching or destructive actions** — the page reassures and directs. "Contact platform admin" opens the configured support surface (mailto/URL), never an in-app form (a suspended tenant shouldn't write anything, including a ticket).
- **Tenant switcher stays live** (app bar + `more`) — a multi-tenant user leaves the suspended tenant, per the caption. **Log out** works.
- The four tabs render but their content areas show this panel — the user sees the boundary, not a broken wall of failures.
- **A11y**: the warning is announced on arrival; the support button is a plain external link.
