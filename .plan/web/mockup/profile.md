# Screen — profile

## Purpose & context

- **User goal**: see and edit the user's own identity — display name, email, avatar. Reachable from the user menu; distinct from `account-settings` (which is security/account-level).
- **Entry points**: user menu → Profile (any authenticated user).
- **Exit points**: back → previous page; save → toast.
- **Data needed**: `users` — `display_name`, `email`, `email_verified`, `last_login_at`, plus memberships across tenants (a user may belong to several). `created_at`.

## States

- [x] Loading (skeleton)
- [x] Populated
- [x] Error
- [x] Domain-specific — unverified email (banner + resend verification)

## Style reference

- **Components used**: `PageHeader`, `FormField`, avatar, sectioned card. Simple, single-column content card — no charts, no table.
- No token deltas.

## Wireframe — desktop (populated)

```
┌──────────┬───────────────────────────────────────────────────────────┐
│ (sidebar)│  Acme Inc ▾         Jane Doe ▾                            │
│          │  Profile                                                  │
│          │  ┌──────────────────────────────────────────────────────┐ │
│          │  │  👤  [Change photo]                                  │ │
│          │  │  Display name                                        │ │
│          │  │  ┌────────────────────────────────────────────────┐  │ │
│          │  │  │ Jane Doe                                       │  │ │
│          │  │  └────────────────────────────────────────────────┘  │ │
│          │  │  Email                                               │ │
│          │  │  ┌────────────────────────────────────────────────┐  │ │
│          │  │  │ jane@acme.dev   ✓ verified                     │  │ │
│          │  │  └────────────────────────────────────────────────┘  │ │
│          │  │                                        [Save]        │ │
│          │  └──────────────────────────────────────────────────────┘ │
│          │  Member of                                                │
│          │  ┌──────────────────────────────────────────────────────┐ │
│          │  │ Acme Inc        owner      since Mar 01, 2026   ▸    │ │
│          │  │ Folio Notes     member     since Mar 05, 2026   ▸    │ │
│          │  └──────────────────────────────────────────────────────┘ │
│          │  Member since Mar 01, 2026 · Last login today at 09:12    │
└──────────┴───────────────────────────────────────────────────────────┘
```

## Wireframe — mobile (populated)

```
┌────────────────────────────┐
│ ☰ ⬡    Acme Inc ▾   👤     │
├────────────────────────────┤
│ ‹ Profile                  │
│ ─────────────────────────  │
│      👤  [Change photo]    │
│ Display name               │
│ ┌────────────────────────┐ │
│ │ Jane Doe               │ │
│ └────────────────────────┘ │
│ Email                      │
│ ┌────────────────────────┐ │
│ │ jane@acme.dev  ✓       │ │
│ └────────────────────────┘ │
│ [Save]                     │
│ ─────────────────────────  │
│ Member of                  │
│ Acme Inc  owner  ▸         │
│ Folio Notes  member ▸      │
│ ─────────────────────────  │
│ Member since Mar 01, 2026  │
│ Last login today 09:12     │
└────────────────────────────┘
```

## Interactions

- **Save**: dirty-field save; email changes trigger re-verification — the field goes from `✓ verified` to "verification email sent" (`color.warning`) until confirmed.
- **Unverified email banner**: "Verify jane@acme.dev to receive reset and invite emails. [Resend]" — appears on load when `email_verified` is false.
- **Member of**: rows link to a tenant switch (`TenantSwitcher` behavior — switching lands on that tenant's `files` root, `navigation.md` tenant-switch rule). Clicking a row switches tenant.
- **Photo**: standard avatar picker; the upload path is the same presigned-PUT flow as any object but writes to a hidden system folder. (Not an M-level feature — a placeholder avatar is acceptable in v1 if the storage plumbing isn't there yet.)
- **Permission-aware**: always the *current* user's own profile — no role gating on viewing; there is no "view another member's profile" route on web (that's `members`, admin-scoped).
