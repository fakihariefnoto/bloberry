# Screen — reset-password

## Purpose & context

- **User goal**: set a new password from the emailed reset link.
- **Entry points**: emailed deep link → `reset-password?token=…` (replace navigation — the token is one-shot).
- **Exit points**: password set → success → `login` (all other sessions invalidated server-side, `backend/domains.md` §4.4). Invalid/expired token → inline error with a "Request a new link" action → `forgot-password`.
- **Data needed**: new password, confirm password. Token from the query string — never shown, never stored in history.

## States

- [x] Loading (token validation in flight)
- [x] Populated (password form)
- [x] Success (password set)
- [x] Error — `reset_token_invalid` (missing/expired/used)
- [x] Error — password rules not met

## Style reference

- **Components used**: primary button, text inputs (password × 2). Same centered-card public layout.
- No deltas.

## Wireframe — desktop (form)

```
┌──────────────────────────────────────────────────────────────────┐
│                                                                  │
│                 ┌──────────────────────────┐                     │
│                 │  Set a new password      │                     │
│                 │                          │                     │
│                 │  New password            │                     │
│                 │  ┌────────────────────┐  │                     │
│                 │  │ ••••••••••    👁    │  │                     │
│                 │  └────────────────────┘  │                     │
│                 │  At least 8 characters.  │                     │
│                 │                          │                     │
│                 │  Confirm new password    │                     │
│                 │  ┌────────────────────┐  │                     │
│                 │  │ ••••••••••    👁    │  │                     │
│                 │  └────────────────────┘  │                     │
│                 │                          │                     │
│                 │  ┌────────────────────┐  │                     │
│                 │  │  Set password      │  │                     │
│                 │  └────────────────────┘  │                     │
│                 └──────────────────────────┘                     │
└──────────────────────────────────────────────────────────────────┘
```

## Wireframe — mobile (token expired)

```
┌────────────────────────────┐
│                            │
│   Reset your password      │
│                            │
│        ⚠                   │
│                            │
│   This link has expired    │
│                            │
│   Reset links last 30      │
│   minutes and can be used  │
│   once.                    │
│                            │
│   [ Request a new link ]   │
└────────────────────────────┘
```

## Interactions

- **Token validation**: on load, the token is validated; while checking, the form area shows a skeleton. If invalid → the expired-state wireframe, no password fields shown.
- **Set password**: on submit → loading → success swaps to "Password set — sign in with your new password" → auto-redirect to `login` after a beat. A successful reset **invalidates every existing session** for that user (`backend/domains.md` §4.4) — worth a caption on the success state so the user isn't surprised when the desktop app logs them out.
- **Validation** (centralized validators): min 8 chars; confirm must match — mismatch is a field error, not a submit-blocker toast. Password rules shown as helper text, not discovered on submit.
- **The token never renders** — a reset link that shows its own token in the address bar is a leak vector; the route consumes it and replaces it with `/reset-password` immediately.
- **A11y**: password fields use `autocomplete="new-password"`; both have independent visibility toggles; focus lands on the first field once the form renders.
