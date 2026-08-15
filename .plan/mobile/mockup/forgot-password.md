# Screen — forgot-password

## Purpose & context

- **User goal**: request a password-reset email.
- **Entry points**: "Forgot password?" from `login`.
- **Exit points**: submitted → success confirmation in place (identical whether or not the account exists — `backend/domains.md` §4.4); back → `login`. The emailed link opens `reset-password`.
- **Data needed**: email only.

## States

- [x] Loading (submit in flight)
- [x] Populated (form)
- [x] Submitted (confirmation)
- [x] Error — network

## Style reference

- **Components used**: form screen pattern, single input, primary button, back. No tab bar.
- No token deltas.

## Wireframe — mobile (form)

```
┌───────────────────────────┐
│  ← Back                   │
│                           │
│  Reset your password      │
│  Enter the email on your  │
│  account and we'll send   │
│  you a reset link.        │
│                           │
│  Email                    │
│  ┌──────────────────────┐ │
│  │ you@project.dev      │ │
│  └──────────────────────┘ │
│                           │
│  ┌──────────────────────┐ │
│  │  Send reset link     │ │
│  └──────────────────────┘ │
└───────────────────────────┘
```

## Wireframe — mobile (submitted)

```
┌───────────────────────────┐
│  ← Back                   │
│                           │
│        ✓                  │
│                           │
│  Check your inbox         │
│                           │
│  If that email exists, a  │
│  reset link is on its way │
│  to you@project.dev. It   │
│  expires in 30 minutes.   │
│                           │
│  ┌──────────────────────┐ │
│  │  Return to sign in   │ │
│  └──────────────────────┘ │
└───────────────────────────┘
```

## Interactions

- **Send reset link**: submit → loading → **same confirmation whether or not the email exists** (no account-enumeration oracle, `backend/domains.md` §4.4). Confirmation states the 30-min TTL.
- **Return to sign in** → `login`.
- **Network failure**: banner with Retry — the user must know the request genuinely didn't go out.
- **A11y**: confirmation is a `radius.full` check announced as "Sent"; focus moves to the confirmation heading after submit.
