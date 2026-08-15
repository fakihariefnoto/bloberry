# Screen — application-detail

## Purpose & context

- **User goal**: see one application's keys, issue a new scoped key, revoke a leaked one, and understand the blast radius of each key (PRD TA5, TA-E3). This is where the "urgent, away from a laptop" trust decision happens on the web.
- **Entry points**: click an application in `applications`; after creating a new app (lands here directly).
- **Exit points**: back → `applications`; Create key → scope form → `key-created` modal (shown once, then gone forever); Revoke → `confirm-destructive` with last-used info; key actions in the list.
- **Data needed**: `applications` (name, description, created_at); `access_keys` for the app — `prefix` (blob_live_/blob_test_), `last_four`, `scope_folder_ids`, `permissions`, `expires_at`, `last_used_at`, `last_used_ip`, `revoked_at`, `created_at`.

## States

- [x] Loading (skeleton)
- [x] Empty (no keys yet)
- [x] Populated
- [x] Error
- [x] Domain-specific — key expiring soon (< 7 days) → warning pill
- [x] Domain-specific — key revoked → muted row, audit trail retained (`ERD.md` access-key lifecycle)
- [x] Domain-specific — `key-created` modal (the once-only secret reveal)

## Style reference

- **Components used**: `PageHeader`, `DataTable`, `StatusPill`, `SecretRevealModal` (`key-created`), `PermissionPicker` (in the create form), `ConfirmDestructive` (revoke), `CopyableCode`, `RelativeTime`.
- The `key-created` modal is `elevation.lg` and **cannot be dismissed by backdrop click** (`web/components.md`).
- No token deltas.

## Wireframe — desktop (populated)

```
┌──────────┬───────────────────────────────────────────────────────────┐
│ (sidebar)│  ← Back to applications                                   │
│          │  acme-cms                           [Create key]          │
│          │  Machine account for the website. Created Mar 01, 2026.   │
│          │  ───────────────────────────────────────────────────────  │
│          │  Access keys                                              │
│          │  ┌──────────────────────────────────────────────────────┐ │
│          │  │ Key               Scoped    Perms      Last used  ⋮│   │
│          │  │ blob_live_••••4f2a whole-tenant all     Mar 13   ⋮ │   │
│          │  │ blob_test_••••8b1d 2026/     read/write Mar 12   ⋮ │   │
│          │  │ ⚠ blob_live_••••c9e7 2026/    read      exp in 5d ⋮│   │
│          │  │ blob_live_••••3a9c scripts/   write     revoked  ⋮ │   │
│          │  └──────────────────────────────────────────────────────┘ │
│          │  Keys are shown once at creation and never again.         │
│          │  A revoked key's history is kept for the audit trail.     │
└──────────┴───────────────────────────────────────────────────────────┘
```

## Wireframe — mobile (`key-created` modal)

```
┌──────────────────────────────┐
│ ☰ ⬡    Acme Inc ▾   👤       │
├──────────────────────────────┤
│ ‹ Applications               │
│ acme-cms        [New key]    │
│ ───────────────────────────  │
│ ┌──────────────────────────┐ │
│   │ Key            ⋮     │   │
│   │ blob_live_••••4f2a   │   │
│   │ whole-tenant · all   │   │
│ └──────────────────────────┘ │
│ ───────────────────────────  │
│   ┌──────────────────────┐   │
│   │  Access key created  │   │
│   │                      │   │
│   │  blob_live_9fK2mQx...│   │
│   │  [Copy secret]       │   │
│   │                      │   │
│   │  ⚠ You won't see this│   │
│   │  again. Anyone with  │   │
│   │  this key can act as │   │
│   │  acme-cms until it   │   │
│   │  expires.            │   │
│   │                      │   │
│   │  [  I've saved it  ] │   │
│   └──────────────────────┘   │
└──────────────────────────────┘
```

## Interactions

- **Create key**: opens the scope form — `PermissionPicker` (folder-subtree selector + read/write/delete/share checkboxes + optional expiry) with the **allow-only explanation inline** (PRD D7: no deny rules; a scope narrows, it never widens). Empty scope = whole tenant.
- **On save** → `key-created` modal shows the **full secret exactly once** (`text.mono`, copy button with confirmation, "You won't see this again" in `color.warning`). Dismissal requires the acknowledgement button — **no backdrop click** (`web/components.md`). Copy is the primary action; closing without copying is explicit and stated.
- **Revoke**: `confirm-destructive` stating the key's `last_used_at`/`last_used_ip` — "This key was last used Mar 13, 09:12 from 203.0.113.8. Revoking takes effect on the next request and is irreversible." No Undo toast (deliberate, `web/components.md`).
- **Revoking the last active key of an application**: the confirm adds the consequence — "This is acme-cms's only active key. Its pipeline will fail on the next call" (PRD TA-E3).
- **Key list**: masked always (`blob_live_••••4f2a`), never a full secret in a table. Expiring < 7 days → warning pill; revoked → muted with a Revoked pill; active → success pill.
- **Row actions**: `⋮` → Revoke (active/expiring) or "Copy ID". Expired keys show no Revoke (already dead).
- **Empty state**: "No keys yet · Create one to let acme-cms authenticate" with the Create button.
- **Permission-aware**: a `tenant_admin` sees everything; creating/revoking is admin-wide, not owner-only. Viewers never reach this route.
