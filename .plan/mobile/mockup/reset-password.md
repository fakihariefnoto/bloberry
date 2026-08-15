# Screen — reset-password

## Purpose & context

- **User goal**: set a new password from the emailed reset link.
- **Entry points**: emailed deep link → `reset-password` (replace navigation; `token: String` arg). The token is consumed once and never shown.
- **Exit points**: password set → `login` (all other sessions invalidated, `backend/domains.md` §4.4). Invalid/expired token → error state with "Request a new link" → `forgot-password`.
- **Data needed**: new password + confirm. Token from the deep link.

## States

- [x] Loading (token validation)
- [x] Populated (password form)
- [x] Success
- [x] Error — `reset_token_invalid`
- [x] Error — password rules not met

## Style reference

- **Components used**: form screen pattern, two password fields (each with visibility toggle), primary button. `autocomplete="new-password"` on both.
- No token deltas.

## Wireframe — mobile (form)

```
┌───────────────────────────┐
│                           │
│  Set a new password       │
│                           │
│  New password             │
│  ┌──────────────────────┐ │
│  │ ••••••••••     👁     │ │
│  └──────────────────────┘ │
│  At least 8 characters.   │
│                           │
│  Confirm new password     │
│  ┌──────────────────────┐ │
│  │ ••••••••••     👁     │ │
│  └──────────────────────┘ │
│                           │
│  ┌──────────────────────┐ │
│  │    Set password      │ │
│  └──────────────────────┘ │
└───────────────────────────┘
```

## Wireframe — mobile (token expired)

```
┌───────────────────────────┐
│                           │
│      ⚠                    │
│                           │
│  This link has expired    │
│                           │
│  Reset links last 30      │
│  minutes and can be used  │
│  once.                    │
│                           │
│  ┌──────────────────────┐ │
│  │ Request a new link   │ │
│  └──────────────────────┘ │
└───────────────────────────┘
```

## Interactions

- **Token validation** on load (skeleton while checking); invalid → expired wireframe, no password fields.
- **Set password**: submit → loading → success ("Password set — sign in with your new password") → auto-redirect to `login`. A successful reset **invalidates every existing session** — the success caption says so (the mobile app's other session dying on purpose is expected).
- **Validation** (centralized `lib/core/validators.dart`): min 8 chars; confirm must match (field error on mismatch).
- **The token never renders** — the deep-link handler consumes it and replaces the route immediately.
- **A11y**: both fields have visibility toggles and proper labels; the expired state is a heading + action, not a raw error string.
