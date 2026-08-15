# Screen — otp-login

## Purpose & context

- **User goal**: sign in with a 6-digit email code instead of a password — the "forgot my password but want in now" path, also handy on a device without the password manager.
- **Entry points**: "Use a code instead" from `login`.
- **Exit points**: code verified → `files` (or `?next` destination); "back to sign in" → `login`. Rate-limited or exhausted attempts block with the real reason.
- **Data needed**: email, 6-digit code. Backend: `POST /v1/auth/otp/request`, `POST /v1/auth/otp/verify` (`backend/domains.md` §4.5).

## States

- [x] Step 1 — email entry
- [x] Step 2 — code entry (after request succeeds)
- [x] Loading (request in flight; verify in flight)
- [x] Error — `otp_rate_limited`
- [x] Error — `otp_invalid` (wrong/expired/exhausted attempts)
- [x] Domain-specific — resend cooldown

## Style reference

- **Components used**: primary button, secondary button (back), text input, single-code input. Same centered-card public-route layout as `login`.
- No deltas from the style guide. The code input is a single field, not N separate boxes — N boxes make paste and autofill miserable, and this is a laptop-first surface.

## Wireframe — desktop (step 2, code entry)

```
┌──────────────────────────────────────────────────────────────────┐
│                                                                  │
│                 ┌──────────────────────────┐                     │
│                 │  Sign in with a code     │                     │
│                 │                          │                     │
│                 │  We emailed a 6-digitcode│                     │
│                 │  to you@project.dev.     │                     │
│                 │                          │                     │
│                 │  Code                    │                     │
│                 │  ┌────────────────────┐  │                     │
│                 │  │  4 2 8 1 9 6       │  │                     │
│                 │  └────────────────────┘  │                     │
│                 │                          │                     │
│                 │  Resend code (00:42)     │                     │
│                 │                          │                     │
│                 │  ┌────────────────────┐  │                     │
│                 │  │     Verify         │  │                     │
│                 │  └────────────────────┘  │                     │
│                 │                          │                     │
│                 │  ← Back to sign in       │                     │
│                 └──────────────────────────┘                     │
└──────────────────────────────────────────────────────────────────┘
```

## Wireframe — mobile (step 1, email entry)

```
┌────────────────────────────┐
│                            │
│   ← Back to sign in        │
│                            │
│   Sign in with a code      │
│   No password needed —     │
│   we'll email you a code.  │
│                            │
│   Email                    │
│   ┌─────────────────────┐  │
│   │ you@project.dev     │  │
│   └─────────────────────┘  │
│                            │
│   ┌─────────────────────┐  │
│   │     Send code       │  │
│   └─────────────────────┘  │
└────────────────────────────┘
```

## Wireframe — error (invalid code, desktop)

```
┌──────────────────────────────────────────────────────────────────┐
│                 ┌──────────────────────────┐                     │
│                 │  Sign in with a code     │                     │
│                 │                          │                     │
│                 │  Code                    │                     │
│                 │  ┌────────────────────┐  │                     │
│                 │  │  4 2 8 1 9 6       │                        │  ← border color.error
│                 │  └────────────────────┘  │                     │
│                 │  That code is wrong or   │                     │
│                 │  expired. Try again.     │                     │
│                 │                          │                     │
│                 │  Resend code (01:23)     │                     │
│                 │  ┌────────────────────┐  │                     │
│                 │  │     Verify         │  │                     │
│                 │  └────────────────────┘  │                     │
│                 └──────────────────────────┘                     │
└──────────────────────────────────────────────────────────────────┘
```

## Interactions

- **Step 1 → 2**: submit email → loading → success swaps to step 2 in place (no navigation). `otp_sent` confirmation. Failure (`otp_rate_limited`) shows the real message: "Too many requests for this email. Try again in about an hour." and does not advance.
- **Resend code**: disabled during the cooldown (shown as `(00:42)` counting down, `motion` none — a live countdown, not a fake interval). After cooldown it's an enabled link; tapping it re-runs the request and restarts the timer. Rate-limit failures here show the same message as step 1.
- **Verify**: on submit → loading → success navigates to `files` (or `?next`). `otp_invalid` shows the inline error; attempts are capped at 5 server-side, and the message on the 5th failure reads "Too many failed attempts — request a new code." 
- **Code entry**: accepts paste and autofill in one field; digits only, 6 characters, uppercase/lowercase n/a. Enter submits.
- **Back to sign in** → `login`, discarding this session's state.
- **A11y**: the code input has `autocomplete="one-time-code"`; the error is announced on arrival.
