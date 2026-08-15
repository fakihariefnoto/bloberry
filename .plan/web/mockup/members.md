# Screen — members

## Purpose & context

- **User goal**: manage who's in the tenant — invite people, change roles, remove members (PRD TA1). Role is the *floor*; folder grants layer on top (PRD D7).
- **Entry points**: sidebar Members (`tenant_admin`+).
- **Exit points**: Invite → `invite-member` modal; Remove member → `confirm-destructive`; role change → inline select (confirm for demoting an owner).
- **Data needed**: `memberships` + `users` (display_name, email), `invitations` (pending), roles (`tenant_owner`/`tenant_admin`/`member`/`viewer`). Derived: how many folders a member has been granted.

## States

- [x] Loading (skeleton)
- [x] Empty (no members — shouldn't happen, owner exists; treat as bug-state "invite people")
- [x] Populated
- [x] Error
- [x] Domain-specific — pending invitations (awaiting acceptance, with expiry)
- [x] Domain-specific — the owner row (cannot be demoted or removed by anyone)

## Style reference

- **Components used**: `AppShell`, `DataTable`, `StatusPill` (role, pending), `invite-member` modal, `ConfirmDestructive`, `RelativeTime`.
- Role changes use an inline select, not a second dialog — demote-owner is the exception that confirms.
- No token deltas.

## Wireframe — desktop (populated)

```
┌──────────┬───────────────────────────────────────────────────────────┐
│ (sidebar)│  Acme Inc ▾         Jane Doe ▾            [Invite]        │
│          │  Members                                                  │
│          │  ┌──────────────────────────────────────────────────────┐ │
│          │  │ Member         Email             Role     Joined   ⋮ │ │
│          │  │ 👤 Jane Doe     jane@acme.dev     owner    Mar 01  ⋮ │ │
│          │  │ 👤 Sam Khan     sam@acme.dev      admin    Mar 02  ⋮ │ │
│          │  │ 👤 Ana Torres   ana@acme.dev      member   Mar 05  ⋮ │ │
│          │  │ 👤 Lee Chen     lee@acme.dev      viewer   Mar 09  ⋮ │ │
│          │  └──────────────────────────────────────────────────────┘ │
│          │  ───────────────────────────────────────────────────────  │
│          │  Pending invitations (1)   — auto-expire after 7 days     │
│          │  ┌──────────────────────────────────────────────────────┐ │
│          │  │ ✉  priya@acme.dev   member  invited Mar 10  [Resend] │ │
│          │  └──────────────────────────────────────────────────────┘ │
│          │                                                           │
│          │  Roles set a floor. Folder grants add access on top —     │
│          │  a viewer can still be granted write on one folder.       │
└──────────┴───────────────────────────────────────────────────────────┘
```

## Wireframe — mobile (populated)

```
┌────────────────────────────┐
│ ☰ ⬡    Acme Inc ▾   👤     │
├────────────────────────────┤
│ Members          [Invite]  │
│ ─────────────────────────  │
│ ┌────────────────────────┐ │
│ │ 👤 Jane Doe   owner    │ │
│ │ jane@acme.dev          │ │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │ 👤 Sam Khan   admin ⋮  │ │
│ │ sam@acme.dev           │ │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │ 👤 Ana Torres member ⋮ │ │
│ │ ana@acme.dev           │ │
│ └────────────────────────┘ │
│ ─────────────────────────  │
│ Pending invitations (1)    │
│ ✉ priya@acme.dev  member   │
│   invited Mar 10 [Resend]  │
│ ─────────────────────────  │
│ Roles set a floor. Folder  │
│ grants add access on top.  │
└────────────────────────────┘
```

## Interactions

- **Invite**: `invite-member` modal — email + role select. On send, the invite appears in the pending list. Only admins/owners can invite; the owner can always be seen by all members.
- **Pending invites**: show the invited email, role, and expiry; **Resend** regenerates the invite (and replaces the old token — an old emailed link stops working). Expired invites drop off the list and self-delete server-side (`ERD.md` invitations TTL).
- **Remove member**: `confirm-destructive`. Removing the **last owner** is refused with "A tenant must have at least one owner." Removing a member who holds grants just removes the membership — grants for a user outside the tenant are inert but retained for the audit trail (`ERD.md` grants).
- **Role change**: inline select. **Demoting the owner** requires typed-name confirmation (the tenant's sole authority is changing hands); any other role change confirms once.
- **The owner row** shows no role select and no remove — this is UI lock, not hidden authority (PRD's role floor means the owner is always visible).
- **Footer caption** (both widths) states the allow-only model so nobody hunts for deny rules here (PRD D7).
- **Empty states**: no members → "Invite your team · People join through emailed invitations"; pending-empty → no separate panel at all.
- **A11y**: role select is keyboard-accessible with visible focus; the confirm dialogs trap focus and return it to the triggering row on dismiss.
