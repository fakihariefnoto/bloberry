# Screen — otp-login

## Purpose & context

- **User goal**: sign in with a 6-digit emailed code instead of a password.
- **Entry points**: "Use a code instead" from `login` (`otp-login` route, `email: String?` arg pre-filled when available).
- **Exit points**: code verified → `files`; back → `login`. Rate-limited / attempts exhausted → real reason, no dead end.
- **Data needed**: email, 6-digit code. `otp_rate_limited` and `otp_invalid` error codes (`backend/domains.md` §4.5: hashed in Redis, 5-attempt cap, per-email rate limit).

## States

- [x] Step 1 — email entry
- [x] Step 2 — code entry (after request succeeds)
- [x] Loading (request / verify in flight)
- [x] Error — `otp_rate_limited`
- [x] Error — `otp_invalid`
- [x] Domain-specific — resend cooldown countdown

## Style reference

- **Components used**: form screen pattern, single code field (not N boxes — paste/autofill on a phone matters more than the aesthetic), primary button, back. `autocomplete="one-time-code"` on the field.
- No token deltas.

## Wireframe — mobile (step 2)

```
┌───────────────────────────┐
│  ← Back                   │
│                           │
│  Sign in with a code      │
│  We emailed a 6-digit     │
│  code to you@project.dev. │
│                           │
│  Code                     │
│  ┌──────────────────────┐ │
│  │  4 2 8 1 9 6         │ │
│  └──────────────────────┘ │
│                           │
│  Resend code (00:42)      │
│                           │
│  ┌──────────────────────┐ │
│  │      Verify          │ │
│  └──────────────────────┘ │
└───────────────────────────┘
```

## Wireframe — mobile (rate-limited)

```
┌───────────────────────────┐
│  Sign in with a code      │
│                           │
│  Code                     │
│  ┌──────────────────────┐ │
│  │                      │ │
│  └──────────────────────┘ │
│                           │
│  ⚠ Too many requests for  │
│  this email. Try again in │
│  about an hour.           │
│                           │
│  ┌──────────────────────┐ │
│  │  Request a new code  │ │  ← disabled (cooldown)
│  └──────────────────────┘ │
│                           │
│  ← Back to sign in        │
└───────────────────────────┘
```

## Interactions

- **Step 1 → 2**: submit email → loading → advance to code entry in place (no navigation). `otp_rate_limited` shows the cooldown message and disables re-request.
- **Resend code**: disabled during countdown (shows `00:42` live), re-enabled after; re-request restarts the timer.
- **Verify**: on submit → loading → success → `files` (preserving `?next`-style destination). `otp_invalid` → inline error; the 5th failure reads "Too many failed attempts — request a new code" (attempt cap, `backend/domains.md` §4.5).
- **Code entry**: paste + autofill in one field; digits only, length 6; Enter submits. Keyboard is numeric.
- **A11y**: the code field's error is announced; the countdown is announced once at expiry (a live-updating label is noise on a screen reader).
