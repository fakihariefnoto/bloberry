# Screen — admin-backend-detail

## Purpose & context

- **User goal**: one backend's full config — credentials (write-only), rate card, capabilities, health, and which tenants use it (PRD PA1, PA4, PA-E1). Where the platform admin fixes an unreachable backend by re-entering credentials.
- **Entry points**: click a backend in `admin-backends`.
- **Exit points**: back → `admin-backends`; save config → toast; test connection → health re-check; delete (0-tenants only) → `confirm-destructive`.
- **Data needed**: `storage_backends` — driver, name, config, `credentials_encrypted` (write-only: on edit, fields show "unchanged" unless replaced), `capabilities`, `rate_card`, `health_status`/`health_error`/`health_checked_at`, derived tenant assignments.

## States

- [x] Loading (skeleton)
- [x] Populated
- [x] Error (backend deleted between list and detail)
- [x] Domain-specific — unreachable (raw provider error shown, re-test action)
- [x] Domain-specific — credentials-unchanged-on-edit (write-only fields)

## Style reference

- **Components used**: `PageHeader`, `FormField`, sectioned card, `StatusPill`, `ByteSize`, test-connection action, `ConfirmDestructive`.
- Credentials render per the `Code / secret display` spec — `text.mono`, masked, never a full value (`web/components.md` write-only fields).
- No token deltas.

## Wireframe — desktop (populated, unreachable)

```
┌──────────┬───────────────────────────────────────────────────────────┐
│ (sidebar)│  ‹ Back to storage backends                               │
│          │  gcs-foundry                    [Test connection]         │
│          │  Google Cloud Storage · foundry-bkt/ · registered Mar 04  │
│          │  ⚠ Unreachable · checked 12:41 · SigV4/oidc auth failed:  │
│          │    service account key expired 2026-07-31.                │
│          │  ┌──────────────────────────────────────────────────────┐ │
│          │  │  CONFIG                                             │  │
│          │  │  Name         [gcs-foundry                     ]    │  │
│          │  │  Bucket       [foundry-bkt/                     ]   │  │
│          │  │  Endpoint     (gcs uses the default endpoint)       │  │
│          │  │  Credentials  service-account JSON — [Replace file] │  │
│          │  │               (current: set · unchanged on save)    │  │
│          │  └──────────────────────────────────────────────────────┘ │
│          │  ┌──────────────────────────────────────────────────────┐ │
│          │  │  RATE CARD                                          │  │
│          │  │  Storage  [0.023]  $/GB/month                       │  │
│          │  │  Egress   [0.09]   $/GB                             │  │
│          │  │  Requests [0.01]   per 1k                           │  │
│          │  └──────────────────────────────────────────────────────┘ │
│          │  ┌──────────────────────────────────────────────────────┐ │
│          │  │  CAPABILITIES (from the conformance suite)          │  │
│          │  │  Presign ✓  Multipart ✓  Storage classes ✗          │  │
│          │  │  Min part 256 KiB  Max parts 10,000                 │  │
│          │  └──────────────────────────────────────────────────────┘ │
│          │  Used by: Folio Notes (default), Masjid App (default)     │
│          │  (read-only — reassignment happens on the tenant)         │
└──────────┴───────────────────────────────────────────────────────────┘
```

## Wireframe — mobile (unreachable)

```
┌────────────────────────────┐
│ ☰ ⬡    Platform ▾    👤    │
├────────────────────────────┤
│ ‹ Backends                 │
│ gcs-foundry  [Test conn]   │
│ ⚠ Unreachable              │
│   service account key      │
│   expired 2026-07-31.      │
│ ─────────────────────────  │
│ CONFIG                     │
│ Name    [gcs-foundry ]     │
│ Bucket  [foundry-bkt/]     │
│ Credentials                │
│ [Replace file] (set)       │
│ ─────────────────────────  │
│ RATE CARD                  │
│ Storage  [0.023] $/GB-mo   │
│ Egress   [0.09]  $/GB      │
│ Reqs     [0.01]  per 1k    │
│ ─────────────────────────  │
│ Used by  Folio, Masjid     │
│ [Save changes]             │
└────────────────────────────┘
```

## Interactions

- **Test connection**: re-runs the health check immediately and re-renders the banner — the fix loop ("replace key → test → healthy") must not wait for the 5-minute ticker.
- **Credentials are write-only**: the field shows "(current: set)" with a Replace-file action; saving without replacing leaves the existing ciphertext untouched (`ERD.md` storage-backends note, R7). Never a masked echo of the real value — the server can't return it and this UI never pretends to.
- **Rate card** edit is here (platform-admin owns it) — this is the input side of `usage`'s "est. $21.40" figures; saving re-renders historical estimates from the then-current card (`ERD.md` usage-snapshots note).
- **Capabilities** are **read-only** — they come from the driver + conformance suite (`backend/domains.md` §6.1), not from this form. Displaying them here is for the admin who wonders why R2 doesn't offer storage classes (`TRD.md` R2).
- **Used by** is read-only and links to `admin-tenants` — reassigning a tenant's backend happens on the tenant, not here, so this panel just states the blast radius of deleting.
- **Delete**: only enabled for 0-tenant backends (see `admin-backends`); typed-name confirm.
- **Unreachable state** renders the banner with the real error and highlights the config card — the field most likely to fix it (credentials) is visible without scrolling.
- **A11y**: the error banner has `role="alert"`; credentials replacement is keyboard-reachable.
