# Screen — accept-invitation

## Purpose & context

- **User goal**: join a tenant from an emailed invitation link. The only way *into* the system — there is no self-serve signup (PRD NG8).
- **Entry points**: emailed link → `/invite/:token` (replace navigation). Also reachable from `welcome`/`login` when the visitor says "I have an invitation" (mobile-only concern; on web the link is the entry).
- **Exit points**: accepted → `files` at the new tenant's root. Already-have-an-account → `login`. Invalid/expired invite → error state with a message to contact the inviting admin. Invite for an existing email adds a membership rather than erroring (`backend/domains.md` §4.1).
- **Data needed**: invitation token (from path), display name, password + confirm (new users), or just a password if the email already has an account. Role and tenant come from the invite, never chosen here.

## States

- [x] Loading (invite validation in flight)
- [x] New user — display name + set password
- [x] Existing user — set password only (account exists, membership being added)
- [x] Success
- [x] Error — `invite_invalid` (expired / used / bad token)

## Style reference

- **Components used**: primary button, text inputs. Same centered-card public layout.
- No deltas.

## Wireframe — desktop (new user)

```
┌──────────────────────────────────────────────────────────────────┐
│                                                                  │
│                 ┌──────────────────────────┐                     │
│                 │  You're invited          │                     │
│                 │                          │                     │
│                 │  [👥]  Acme Inc. invited │                     │
│                 │  you to join as  Member  │                     │
│                 │                          │                     │
│                 │  Display name            │                     │
│                 │  ┌────────────────────┐  │                     │
│                 │  │ Jane Doe           │  │                     │
│                 │  └────────────────────┘  │                     │
│                 │                          │                     │
│                 │  Password                │                     │
│                 │  ┌────────────────────┐  │                     │
│                 │  │ ••••••••••    👁    │  │                     │
│                 │  └────────────────────┘  │                     │
│                 │  At least 8 characters.  │                     │
│                 │                          │                     │
│                 │  ┌────────────────────┐  │                     │
│                 │  │  Join Acme Inc.    │  │                     │
│                 │  └────────────────────┘  │                     │
│                 │                          │                     │
│                 │  Already have an account?│                     │
│                 │  Sign in instead         │                     │
│                 └──────────────────────────┘                     │
└──────────────────────────────────────────────────────────────────┘
```

## Wireframe — mobile (invite expired)

```
┌────────────────────────────┐
│                            │
│   Join your tenant         │
│                            │
│        ⚠                   │
│                            │
│   This invitation is no    │
│   longer valid.            │
│                            │
│   Invitations expire after │
│   7 days and can be used   │
│   once.                    │
│                            │
│   Ask the admin who        │
│   invited you for a new    │
│   one.                     │
│                            │
│   [  Sign in  ]            │
└────────────────────────────┘
```

## Interactions

- **Token validation**: on load, validate the invite; skeleton while checking. The wireframe's role badge ("Member") is read-only — the invite carries the role and it must not be editable here.
- **Existing-email branch**: when the email already has an account, the display-name field is hidden and the heading becomes "Invited as Member of Acme Inc. — set a password to finish joining." Accepting adds the membership and logs the user in.
- **Join**: on submit → loading → success → `files` at the tenant root. `invite_invalid` → error wireframe; no retry button (the token is dead), just the guidance.
- **Sign in instead** → `login` (for the visitor who already has an account but is mid-invite; the invite survives in the browser history if they need it).
- **A11y**: tenant name and role are announced on load; the error state links the words "Sign in" rather than bare text.
