# Screen — applications

## Purpose & context

- **User goal**: see the tenant's applications and their keys from a phone — primarily the "revoke a leaked key away from a laptop" case (PRD TA5, `navigation.md` flow chain) plus a read-only overview.
- **Entry points**: Applications from `more` (`tenant_admin`+ gate).
- **Exit points**: tap an application → `application-detail`; back → `more`. (Creating applications is web-only in v1 — key *management* is mobile-urgent, key *creation* is a desk task.)
- **Data needed**: `applications` — name, description, created; derived key counts + last-used.

## States

- [x] Loading (skeleton)
- [x] Empty (no applications — unusual; the tenant admin should create one on web)
- [x] Populated
- [x] Error
- [x] Domain-specific — an application with no active keys (flagged, likely broken pipeline)

## Style reference

- **Components used**: list rows (app icon/initial, name, keys + last-used line, warning pill for keyless apps), `⋮` overflow. No create affordance (web-only creation decision stated in the screen copy).
- No token deltas.

## Wireframe — mobile (populated)

```
┌───────────────────────────┐
│  ← Applications           │
├───────────────────────────┤
│ ┌───────────────────────┐ │
│ │ ⚙ acme-cms        ⋮   │ │
│ │ 3 keys · used Mar 13  │ │
│ └───────────────────────┘ │
│ ┌───────────────────────┐ │
│ │ ⚙ ci-deploy       ⋮   │ │
│ │ 1 key · used Mar 13   │ │
│ └───────────────────────┘ │
│ ┌───────────────────────┐ │
│ │ ⚠ legacy-api      ⋮   │ │
│ │ 0 keys · never used   │ │
│ │ all revoked           │ │
│ └───────────────────────┘ │
│ ───────────────────────── │
│ Creating applications is  │
│ a desk task — use the web │
│ dashboard.                │
└───────────────────────────┘
```

## Interactions

- **Tap a row** → `application-detail` (keys list with revoke).
- **`⋮`** → Revoke-last-key shortcut is NOT here — the revoke flow lives in detail where the last-used context is visible (PRD TA-E3 requires seeing what a key last accessed before revoking). The `⋮` menu only offers "Open detail".
- **Keyless app** flagged with the warning pill ("all revoked" / "never used") — a keyless app is usually a broken pipeline, not a cleaned-up one.
- **The footer caption** states the create-on-web decision, so a phone user hunting for a "New app" button understands why it's absent (deliberate, `mobile/navigation.md` platform-admin surfaces note).
- **Empty**: "No applications yet · Register one on the web dashboard."
- **A11y**: rows are 48dp targets; the warning pill carries a text label.
