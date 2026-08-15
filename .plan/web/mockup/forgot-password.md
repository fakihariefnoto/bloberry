# Screen — forgot-password

## Purpose & context

- **User goal**: request a password-reset email.
- **Entry points**: "Forgot password?" from `login`.
- **Exit points**: after submitting → success confirmation in place ("If that email exists, a reset link is on its way" — identical response whether or not the account exists, per `backend/domains.md` §4.4); "Back to sign in" → `login`. The emailed link opens `reset-password`.
- **Data needed**: email only.

## States

- [x] Loading (submit in flight)
- [x] Populated (form)
- [x] Submitted — success confirmation (same message regardless of account existence)
- [x] Error — network/unreachable

## Style reference

- **Components used**: primary button, secondary button (back), text input. Same centered-card public layout.
- No deltas.

## Wireframe — desktop (form)

```
┌──────────────────────────────────────────────────────────────────┐
│                                                                  │
│                 ┌──────────────────────────┐                     │
│                 │  Reset your password     │                     │
│                 │                          │                     │
│                 │  Enter the email on your │                     │
│                 │  account and we'll send  │                     │
│                 │  you a reset link.       │                     │
│                 │                          │                     │
│                 │  Email                   │                     │
│                 │  ┌────────────────────┐  │                     │
│                 │  │ you@project.dev    │  │                     │
│                 │  └────────────────────┘  │                     │
│                 │                          │                     │
│                 │  ┌────────────────────┐  │                     │
│                 │  │  Send reset link   │  │                     │
│                 │  └────────────────────┘  │                     │
│                 │                          │                     │
│                 │  ← Back to sign in       │                     │
│                 └──────────────────────────┘                     │
└──────────────────────────────────────────────────────────────────┘
```

## Wireframe — mobile (submitted)

```
┌────────────────────────────┐
│                            │
│   ← Back to sign in        │
│                            │
│        ✓                   │
│                            │
│   Check your inbox         │
│                            │
│   If that email exists, a  │
│   reset link is on its way │
│   to you@project.dev. It   │
│   expires in 30 minutes.   │
│                            │
│   [  Return to sign in  ]  │
└────────────────────────────┘
```

## Interactions

- **Send reset link**: on submit → loading → **same confirmation shown whether or not the email exists** — this endpoint must not be an account-enumeration oracle (`backend/domains.md` §4.4). The confirmation gives the link's TTL (30 min).
- **Return to sign in** → `login`.
- **Network failure**: banner above the card with Retry, so the user knows the request genuinely didn't go out.
- **A11y**: the confirmation icon uses a check inside a `radius.full` circle, announced as "Sent"; focus moves to the confirmation heading after submit.
