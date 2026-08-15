# Screen — accept-invitation

## Purpose & context

- **User goal**: join a tenant from an emailed invitation — the only way *into* the system (PRD NG8, no self-serve signup). Fills the gap the standard Signup screen would occupy (`mobile/navigation.md` Required screens table).
- **Entry points**: emailed deep link → `accept-invitation` (replace navigation). No manual-entry path on mobile beyond pasting the invite URL into the OS (the app registers the deep link).
- **Exit points**: accepted → `files` at the tenant root; "Sign in instead" → `login` (existing-account case); invalid/expired → error state with guidance to contact the admin.
- **Data needed**: invite `token`; display name + password (new user) or password only (existing email — membership added rather than a duplicate account, `backend/domains.md` §4.1). Role and tenant come from the invite, never chosen here.

## States

- [x] Loading (invite validation)
- [x] New user — display name + set password
- [x] Existing user — set password only
- [x] Success
- [x] Error — `invite_invalid`

## Style reference

- **Components used**: form screen pattern, primary button, tenant badge (read-only role pill). No tab bar.
- No token deltas.

## Wireframe — mobile (new user)

```
┌───────────────────────────┐
│                           │
│  You're invited           │
│                           │
│  [👥]  Acme Inc invited   │
│  you to join as Member    │
│                           │
│  Display name             │
│  ┌──────────────────────┐ │
│  │ Jane Doe             │ │
│  └──────────────────────┘ │
│                           │
│  Password                 │
│  ┌──────────────────────┐ │
│  │ ••••••••••     👁     │ │
│  └──────────────────────┘ │
│  At least 8 characters.   │
│                           │
│  ┌──────────────────────┐ │
│  │   Join Acme Inc      │ │
│  └──────────────────────┘ │
│                           │
│  Already have an account? │
│  Sign in instead          │
└───────────────────────────┘
```

## Wireframe — mobile (invite expired)

```
┌───────────────────────────┐
│                           │
│      ⚠                    │
│                           │
│  This invitation is no    │
│  longer valid.            │
│                           │
│  Invitations expire after │
│  7 days and can be used   │
│  once.                    │
│                           │
│  Ask the admin who        │
│  invited you for a new    │
│  one.                     │
│                           │
│  ┌──────────────────────┐ │
│  │      Sign in         │ │
│  └──────────────────────┘ │
└───────────────────────────┘
```

## Interactions

- **Token validation** on load; skeleton while checking. The role badge is read-only — the invite carries the role, never editable here.
- **Existing-email branch**: display-name field hidden; heading becomes "Invited as Member of Acme Inc — set a password to finish joining." Accepting adds the membership and logs in.
- **Join**: submit → loading → `files` at the tenant root. `invite_invalid` → error wireframe (no retry — the token is dead), guidance to contact the admin.
- **Sign in instead** → `login` (the invite stays in the OS's open URL history if needed).
- **A11y**: tenant + role announced on load; the error state's "Sign in" is a link, not bare text.
