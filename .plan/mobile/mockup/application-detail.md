# Screen — application-detail

## Purpose & context

- **User goal**: see an application's keys and **revoke a leaked one from a phone** — the urgent, away-from-a-laptop case the mobile app exists to serve (PRD TA5, TA-E3, `navigation.md` flow chain "Revoke a key from a phone").
- **Entry points**: tap an application in `applications`.
- **Exit points**: back → `applications`; Revoke → `confirm-destructive` stating last-used; revoked → toast (no Undo), row mutes. Key *creation* is not offered (web-only in v1 — a desk task).
- **Data needed**: `access_keys` for the app — `prefix`, `last_four`, scope, permissions, expiry, `last_used_at`, `last_used_ip`, `revoked_at`.

## States

- [x] Loading (skeleton)
- [x] Empty (no keys)
- [x] Populated
- [x] Error
- [x] Domain-specific — key expiring soon (warning pill)
- [x] Domain-specific — revoked (muted row, history retained)
- [x] Domain-specific — last active key (revoke confirm warns the app's pipeline dies)

## Style reference

- **Components used**: list rows (masked key `blob_live_••••4f2a` in `text.mono`, scope + permissions line, last-used line, status pill), Revoke via `⋮` → `confirm-destructive` sheet (typed-name confirmation — key revocation is irreversible, `design/style-guide.md`).
- No token deltas.

## Wireframe — mobile (populated)

```
┌────────────────────────────┐
│  ← Applications            │
├────────────────────────────┤
│  ⚙ acme-cms                │
│  Machine account for the   │
│  website. Created Mar 01.  │
│ ────────────────────────── │
│  KEYS (3)                  │
│ ┌────────────────────────┐ │
│ │ blob_live_••••4f2a  ⋮  │ │
│ │ whole-tenant · all     │ │
│ │ used Mar 13, 09:12     │ │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │ ⚠ blob_live_••••c9e7 ⋮ │ │
│ │ 2026/ · read           │ │
│ │ exp in 5 days          │ │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │ blob_live_••••3a9c  ⋮  │ │
│ │ scripts/ · write       │ │
│ │ revoked · Mar 02       │ │
│ └────────────────────────┘ │
│ ────────────────────────── │
│ Keys are shown once at     │
│ creation and never again.  │
└────────────────────────────┘
```

## Wireframe — mobile (revoke confirm)

```
┌────────────────────────────┐
│   [ app-detail behind,     │
│     dimmed ]               │
│ ┌─────────────────────────┐│
│ │        ───────          ││
│ │  Revoke this key?       ││
│ │                         ││
│ │  blob_live_••••4f2a     ││
│ │                         ││
│ │  Last used Mar 13,      ││
│ │  09:12 from 203.0.113.8.││
│ │                         ││
│ │  Type the key's last 4  ││
│ │  chars to confirm:      ││
│ │  ┌──────────────────┐   ││
│ │  │ 4f2a             │   ││
│ │  └──────────────────┘   ││
│ │                         ││
│ │  [ Cancel ]  [ Revoke ] ││
│ └─────────────────────────┘│
└────────────────────────────┘
```

## Interactions

- **Revoke** (`⋮` → Revoke) opens `confirm-destructive` showing the key's **last-used time + IP** — the TA-E3 context ("understand what you're about to break"). Typed-name confirmation (last-4-chars) for the irreversible act.
- **Revoking the last active key** adds the consequence line: "This is acme-cms's only active key. Its pipeline will fail on the next call."
- **Revoked** → toast (deliberately **no Undo** — key revocation is irreversible, `web/components.md` rule shared with mobile), row mutes to `color.text-muted` with a Revoked pill; the history (last_used_at) survives for the audit trail (`ERD.md` access-key lifecycle).
- **Expiring < 7 days** → warning pill; expired keys show no Revoke (already dead).
- **Create key is absent** (web-only, stated in the footer — the phone is for containment, not provisioning).
- **A11y**: the typed-name field confirms exact-char matching with case-insensitivity where sensible; the confirm sheet traps focus and returns it to the triggering row.
