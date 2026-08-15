# Screen — applications

## Purpose & context

- **User goal**: manage the tenant's non-human principals (PRD TA4) — applications that hold access keys. The list view: name, description, key count, last used.
- **Entry points**: sidebar Applications (`tenant_admin`+).
- **Exit points**: click an application → `application-detail`; New application → modal (short form: name, description). 
- **Data needed**: `applications` — `name`, `description`, `created_at`; derived: number of active/expired/revoked `access_keys`, last key `last_used_at`.

## States

- [x] Loading (skeleton rows)
- [x] Empty (no applications yet)
- [x] Populated
- [x] Error
- [x] Domain-specific — application has keys but none active (all revoked/expired) — flag, since a silent keyless app is how CI starts 401ing

## Style reference

- **Components used**: `AppShell`, `DataTable` (no selection needed — no bulk action here), `StatusPill`, `RelativeTime`, `EmptyState`, modal form.
- No token deltas.

## Wireframe — desktop (populated)

```
┌──────────┬───────────────────────────────────────────────────────────┐
│ (sidebar)│  Acme Inc ▾         Jane Doe ▾           [+ New app]      │
│          │  Applications                                             │
│          │  Machine accounts that hold access keys.                  │
│          │  ┌──────────────────────────────────────────────────────┐ │
│          │  │ Name        Keys  Last used       Created       ⋮    │ │
│          │  │ acme-cms    3     Mar 13, 09:12   Mar 01        ⋮    │ │
│          │  │ ci-deploy   1     Mar 13, 08:44   Mar 01        ⋮    │ │
│          │  │ ⚠ legacy-api 0    never           Feb 12        ⋮    │ │
│          │  └──────────────────────────────────────────────────────┘ │
│          │  Showing 3 of 3                                           │
└──────────┴───────────────────────────────────────────────────────────┘
```

## Wireframe — mobile (populated)

```
┌────────────────────────────┐
│ ☰ ⬡    Acme Inc ▾   👤     │
├────────────────────────────┤
│ Applications   [+ New app] │
│ Machine accounts holding   │
│ access keys.               │
│ ─────────────────────────  │
│ ┌────────────────────────┐ │
│ │ acme-cms         ⋮    │  │
│ │ 3 keys · Mar 13       │  │
│ │ created Mar 01        │  │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │ ci-deploy         ⋮   │  │
│ │ 1 key · Mar 13        │  │
│ │ created Mar 01        │  │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │ ⚠ legacy-api      ⋮   │  │
│ │ 0 keys · never used   │  │
│ │ created Feb 12        │  │
│ └────────────────────────┘ │
└────────────────────────────┘
```

## Interactions

- **New app**: modal with name + description; on save → immediately `application-detail` for that app (keys are issued there, and issuing a key is the app's real purpose — don't drop the user back on the list to hunt for it).
- **Row click** → `application-detail`. Row actions: `⋮` → Delete application (only allowed when the app has **no active keys** — deleting an app with live keys is refused with a message, because it silently orphans the keys that authorize production CI; `ERD.md` access-key lifecycle).
- **The ⚠ state**: an application with zero active keys gets the warning pill ("0 keys · never" or "all revoked") — a keyless app is usually a broken pipeline, not a cleaned-up one. This is a deliberate deviation from a neutral "0".
- **Keys column** links into the app's detail filtered to keys.
- **Empty state**: "No applications yet · Register an app to issue it scoped access keys" with the New button.
- **Permission-aware**: `member`/`viewer` never see this route (sidebar hides the section); a `tenant_admin` who can view but is not the owner sees everything here — key management is admin-wide, not owner-only.
