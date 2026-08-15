# Screen — forbidden

## Purpose & context

- **User goal**: land on a 403 and be told plainly what happened and what would be needed — not an opaque wall (PRD MV4: blocked actions are explained, never error-walls). Server-side guards redirect here (`web/navigation.md` graph: `Shell -.->|"403 on a guarded route"| Forbidden`).
- **Entry points**: a guarded route rejected by the RBAC middleware — e.g. a `member` typing `/members` directly, or a `platform_admin` route reached by a non-admin with a stale URL.
- **Exit points**: "Back to Files" → `files`. (The sidebar still shows the sections the user *can* see.)
- **Data needed**: none from the server — the page knows the attempted route's guard from the client router and can name what's needed ("Requires `tenant_admin` or `tenant_owner`").

## States

- [x] Single state. (The 403 redirect replaces the current route, so there's no loading/error variant — the page *is* the error state.)

## Style reference

- **Components used**: `AppShell`, a centered status panel. `color.warning` icon — a permission boundary, not a crash. No `color.error` (the app didn't fail; the user asked for something they can't have).
- No token deltas.

## Wireframe — desktop

```
┌──────────┬──────────────────────────────────────────────────────────┐
│ (sidebar)│  Acme Inc ▾         Jane Doe ▾                           │
│          │                                                          │
│          │            🔒  You don't have access here                │
│          │                                                          │
│          │    Members requires the tenant_admin or tenant_owner     │
│          │    role.                                                 │
│          │                                                          │
│          │    If you need it, ask a tenant admin to change your     │
│          │    role or grant you folder-level access.                │
│          │                                                          │
│          │    [  Back to Files  ]                                   │
└──────────┴──────────────────────────────────────────────────────────┘
```

## Wireframe — mobile

```
┌────────────────────────────┐
│ ☰ ⬡    Acme Inc ▾   👤     │
├────────────────────────────┤
│      🔒  No access here    │
│                            │
│  Members requires the      │
│  tenant_admin or tenant    │
│  owner role.               │
│                            │
│  If you need it, ask a     │
│  tenant admin to change    │
│  your role or grant you    │
│  folder-level access.      │
│                            │
│  [  Back to Files  ]       │
└────────────────────────────┘
```

## Interactions

- **Back to Files** → `files` (the shell default). The sidebar remains fully functional — the user isn't locked out of what they *can* do; they're shown the boundary of what they can't.
- **The message names the required role** from the route guard, not a generic "forbidden" — "Requires `tenant_admin` or `tenant_owner`" tells the user what to ask for (`design/style-guide.md` → Permission-denied state: a `text.caption` naming what's needed).
- Never an error toast + stale page: the 403 replaces the route so the user doesn't see half a page behind a toast.
- **A11y**: `role="alert"`-free (this is a stable destination, not a transient alert); heading + action in a single `main` landmark.
