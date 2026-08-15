# Screen — tenant-settings

## Purpose & context

- **User goal**: tenant-level configuration owned by the `tenant_owner`: name/slug, quota, and the storage backend the tenant writes new objects to (PRD TA7). Where a tenant switches from one backend to another and old objects keep resolving (`ERD.md` `objects.backend_id`, ADR-4).
- **Entry points**: sidebar Settings (`tenant_owner` only).
- **Exit points**: change backend → `confirm-destructive` (it changes where future bytes land); save name/quota → toast; account-level settings are `account-settings` (separate route).
- **Data needed**: `tenants` — name, slug, `default_backend_id`, `quota_bytes`, `quota_objects`, `used_bytes`/`used_objects`, `status`. Available backends (install-level pool + this tenant's BYO). NOTE: quota editing — decide whether tenant owners set their own quota or platform admins do; PRD PA2 has platform admins set quota on create. This screen shows quota read-only by default.

## States

- [x] Loading (skeleton form)
- [x] Populated
- [x] Error
- [x] Domain-specific — switching backend (confirm states old objects keep resolving, new ones go to the new backend)
- [x] Domain-specific — tenant suspended (settings read-only, banner)

## Style reference

- **Components used**: `PageHeader`, `FormField`, sectioned settings card (per `design-collection/web-screen/patterns.md` Settings pattern — groups, not a flat list), `ConfirmDestructive`.
- No token deltas.

## Wireframe — desktop (populated)

```
┌──────────┬───────────────────────────────────────────────────────────┐
│ (sidebar)│  Acme Inc ▾         Jane Doe ▾                            │
│          │  Tenant settings                                          │
│          │  ┌──────────────────────────────────────────────────────┐ │
│          │  │  GENERAL                                             │ │
│          │  │  Tenant name    Acme Inc                      [Save] │ │
│          │  │  Slug           acme           (read-only, shown)    │ │
│          │  │                                                      │ │
│          │  │  STORAGE                                             │ │
│          │  │  Storage backend  [s3-eu-prod ▾]          [Change]   │ │
│          │  │  New objects go to the selected backend. Existing    │ │
│          │  │  objects stay where they are and keep resolving.     │ │
│          │  │                                                      │ │
│          │  │  QUOTA                                               │ │
│          │  │  Quota           500 GB          (set by platform    │ │
│          │  │  Objects         1,000,000        admin)             │ │
│          │  │  Used            312 GB · 48,912 objects             │ │
│          │  │  ──────────────────────────────────────────────      │ │
│          │  │  Storage  312 GB of 500 GB ████████████░░░░  62%     │ │
│          │  └──────────────────────────────────────────────────────┘ │
│          │  DANGER ZONE                                              │
│          │  ┌──────────────────────────────────────────────────────┐ │
│          │  │  Delete tenant   Deletes folders, objects, keys and  │ │
│          │  │  memberships. Bytes remain in the bucket (orphans).  │ │
│          │  │                                        [Delete…]     │ │
│          │  └──────────────────────────────────────────────────────┘ │
└──────────┴───────────────────────────────────────────────────────────┘
```

## Wireframe — mobile (populated)

```
┌────────────────────────────┐
│ ☰ ⬡    Acme Inc ▾   👤     │
├────────────────────────────┤
│ Tenant settings            │
│ ─────────────────────────  │
│ GENERAL                    │
│ Tenant name                │
│ ┌────────────────────────┐ │
│ │ Acme Inc               │ │
│ └────────────────────────┘ │
│ [Save]                     │
│ ─────────────────────────  │
│ STORAGE                    │
│ Storage backend            │
│ ┌────────────────────────┐ │
│ │ s3-eu-prod        ▾    │ │
│ └────────────────────────┘ │
│ [Change backend]           │
│ New objects only. Existing │
│ objects keep resolving.    │
│ ─────────────────────────  │
│ QUOTA                      │
│ 500 GB · 1,000,000 objects │
│ set by platform admin      │
│ 312 GB · 62% ████████░░░░  │
│ ─────────────────────────  │
│ DANGER ZONE                │
│ Delete tenant        [Del] │
└────────────────────────────┘
```

## Interactions

- **Save (name)**: local-save pattern — edits mark the field, Save persists and toasts; a dirty field shows a subtle unsaved indicator rather than auto-saving.
- **Change backend**: `confirm-destructive` — the message is the ADR-4 contract: "New uploads will go to r2-main. Existing objects stay on s3-eu-prod and keep resolving. This does not move data." Switching back later is the same operation, no migration needed (PRD NG7/ADR-4). The dropdown lists install-level backends + the tenant's own BYO backend, with the health pill on each.
- **Quota is read-only here**: platform admins set quotas (PRD PA2) — this screen shows the number and a "set by platform admin" caption. Tenant owners don't cap themselves.
- **Delete tenant**: typed-name confirmation (DANGER ZONE). The confirm states what it does and does not do: "folders, objects and keys are deleted; **bytes in the bucket are not** — they become orphans" (`architecture.md` reconciliation sweep finds them later).
- **Suspended**: full-width banner "This tenant is suspended — reads and writes are blocked. Contact the platform admin." Form fields disabled, no Save.
- **Permission-aware**: route is `tenant_owner` only. A `tenant_admin` sees it in the sidebar but is refused by the guard (hidden section + server-side 403 → `forbidden`). Admins who need backend visibility use `usage`.
