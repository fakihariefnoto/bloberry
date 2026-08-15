# Screen — profile

## Purpose & context

- **User goal**: view/edit the user's own display name, email, avatar; see their memberships.
- **Entry points**: Profile from `more` (or user-header tap).
- **Exit points**: back → `more`; save → toast; tapping a different tenant's membership row → tenant switch → `files` at that tenant's root.
- **Data needed**: `users` — display_name, email, email_verified, last_login_at; memberships across tenants (role per tenant).

## States

- [x] Loading (skeleton)
- [x] Populated
- [x] Error (save failed — inline)
- [x] Domain-specific — unverified email banner

## Style reference

- **Components used**: form/list hybrid — edit fields inline, memberships as a list below. Avatar with change affordance.
- No token deltas.

## Wireframe — mobile (populated)

```
┌───────────────────────────┐
│  ← Profile                │
├───────────────────────────┤
│       👤  [Change]        │
│                           │
│  Display name             │
│  ┌──────────────────────┐ │
│  │ Jane Doe             │ │
│  └──────────────────────┘ │
│                           │
│  Email                    │
│  ┌──────────────────────┐ │
│  │ jane@acme.dev   ✓    │ │
│  └──────────────────────┘ │
│                           │
│  ┌──────────────────────┐ │
│  │        Save          │ │
│  └──────────────────────┘ │
│  ──────────────────────── │
│  MEMBER OF                │
│  ┌──────────────────────┐ │
│  │ Acme Inc   owner   › │ │
│  └──────────────────────┘ │
│  ┌──────────────────────┐ │
│  │ Folio Notes  member ›│ │
│  └──────────────────────┘ │
│  ──────────────────────── │
│  Member since Mar 01, 2026│
│  Last login today 09:12   │
└───────────────────────────┘
```

## Interactions

- **Save**: dirty-field save; changing email triggers re-verification (the `✓` becomes "verification sent" in `color.warning`).
- **Unverified email banner**: "Verify jane@acme.dev to receive reset and invite emails. [Resend]" on load when `email_verified` is false.
- **Member of**: each row shows tenant + role; tapping switches tenant (`files` at the new root, `navigation.md` rule).
- **Avatar change**: OS image picker; the byte write uses the same presigned-PUT path as any object (hidden system folder). Placeholder acceptable in v1 if plumbing isn't there.
- **A11y**: labels above fields; the avatar is `role="presentation"`; member rows are 48dp targets.
