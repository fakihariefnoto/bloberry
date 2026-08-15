# Screen — login

## Purpose & context

- **User goal**: sign in to the dashboard to reach `files` at the tenant root.
- **Entry points**: the app's only public route. Any authenticated route hit without a session redirects here with `?next=<path>`. Also the landing point after session expiry (401), after logout, and after a successful password reset.
- **Exit points**: authenticated → `files` (or `?next` destination); "Use a code instead" → `otp-login`; "Forgot password?" → `forgot-password`; "Continue with Google" → Google OAuth → `files`. A brand-new visitor with an invitation has no self-serve signup (PRD NG8) — they use the emailed `/invite/:token` link, not this screen.
- **Data needed**: email, password. Google OAuth button. Error cases: `invalid_credentials`, `key_revoked` isn't relevant here. If the tenant is suspended, the redirect after login lands on `tenant-suspended`.

## States

- [x] Loading (submit in flight — button disabled, spinner)
- [ ] Empty (fields blank, this is the happy path)
- [x] Populated
- [x] Error — invalid credentials (identical message for unknown email and wrong password)
- [x] Error — network/backend unreachable
- [x] Domain-specific — returning with `?next=` (subtle "Continue to where you were" — no separate layout, just a hint line)
- [x] Domain-specific — **TOTP step** (M24): primary factor succeeded, `totp_required` returned, a code field appears before any session is issued

## Style reference

- **Components used**: primary button, secondary button (Google), text inputs, link-style text buttons. Centered card on `color.background`. No sidebar — this is a public route.
- No deltas from the style guide.

## Wireframe — desktop (happy path)

```
┌──────────────────────────────────────────────────────────────────┐
│                                                                  │
│                    ⬡  Bloberry                                   │
│                    Storage for your projects                     │
│                                                                  │
│                 ┌──────────────────────────┐                     │
│                 │  Sign in                 │                     │
│                 │                          │                     │
│                 │  Email                   │                     │
│                 │  ┌────────────────────┐  │                     │
│                 │  │ you@project.dev    │  │                     │
│                 │  └────────────────────┘  │                     │
│                 │                          │                     │
│                 │  Password                │                     │
│                 │  ┌────────────────────┐  │                     │
│                 │  │ ••••••••••    👁    │  │                     │
│                 │  └────────────────────┘  │                     │
│                 │                          │                     │
│                 │  Forgot password?        │                     │
│                 │                          │                     │
│                 │  ┌────────────────────┐  │                     │
│                 │  │     Sign in        │  │                     │
│                 │  └────────────────────┘  │                     │
│                 │                          │                     │
│                 │  ── or continue with ──  │                     │
│                 │  [  ⊙  Continue with Google  ]                 │
│                 │                          │                     │
│                 │  Use a code instead      │                     │
│                 └──────────────────────────┘                     │
│                                                                  │
│                 No account? Ask your tenant admin                │
│                 for an invitation.                               │
└──────────────────────────────────────────────────────────────────┘
```

## Wireframe — mobile (happy path)

```
┌────────────────────────────┐
│                            │
│        ⬡  Bloberry         │
│   Storage for your projects│
│                            │
│        Sign in             │
│                            │
│  Email                     │
│  ┌──────────────────────┐  │
│  │ you@project.dev      │  │
│  └──────────────────────┘  │
│                            │
│  Password                  │
│  ┌──────────────────────┐  │
│  │ ••••••••••     👁     │  │
│  └──────────────────────┘  │
│                            │
│  Forgot password?          │
│                            │
│  ┌──────────────────────┐  │
│  │      Sign in         │  │
│  └──────────────────────┘  │
│                            │
│  ── or continue with ──    │
│  [  ⊙  Continue withGoogle]│
│                            │
│  Use a code instead        │
│                            │
│  No account? Ask your      │
│  tenant admin for an       │
│  invitation.               │
└────────────────────────────┘
```

## Wireframe — error (invalid credentials, desktop)

```
┌──────────────────────────────────────────────────────────────────┐
│                                                                  │
│                 ┌──────────────────────────┐                     │
│                 │  Sign in                 │                     │
│                 │                          │                     │
│                 │  Email                   │                     │
│                 │  ┌────────────────────┐  │                     │
│                 │  │ you@project.dev    │  │                     │
│                 │  └────────────────────┘  │                     │
│                 │                          │                     │
│                 │  Password                │                     │
│                 │  ┌────────────────────┐  │                     │
│                 │  │ ••••••••••    👁    │                        │   ← border color.error
│                 │  └────────────────────┘  │                     │
│                 │  Wrong email or password.                      │
│                 │  Check both and try again.                     │
│                 │                          │                     │
│                 │  ┌────────────────────┐  │                     │
│                 │  │     Sign in        │  │                     │
│                 │  └────────────────────┘  │                     │
│                 │                          │                     │
│                 │  ── or continue with ──  │                     │
│                 │  [  ⊙  Continue with Google  ]                 │
│                 └──────────────────────────┘                     │
└──────────────────────────────────────────────────────────────────┘
```

## Wireframe — TOTP step (M24, desktop)

The primary factor succeeded but `totp_required` was returned — the session is **not** issued yet. The card collapses to a single code field.

```
┌──────────────────────────────────────────────────────────────────┐
│                                                                  │
│                 ┌──────────────────────────┐                     │
│                 │  Two-factor check        │                     │
│                 │                          │                     │
│                 │  Enter the 6-digit code  │                     │
│                 │  from your authenticator │                     │
│                 │  app.                    │                     │
│                 │                          │                     │
│                 │  ┌────────────────────┐  │                     │
│                 │  │  4 2 8 1 9 6       │  │                     │
│                 │  └────────────────────┘  │                     │
│                 │                          │                     │
│                 │  Use a backup code       │                     │
│                 │                          │                     │
│                 │  ┌────────────────────┐  │                     │
│                 │  │     Verify         │  │                     │
│                 │  └────────────────────┘  │                     │
│                 └──────────────────────────┘                     │
└──────────────────────────────────────────────────────────────────┘
```

## Interactions

- **Sign in**: on submit, button goes loading (spinner replaces label, width held, non-interactive). On success → `files` or the `?next` destination. On 401 → inline error under the password field: "Wrong email or password." — **identical for unknown email and wrong password** (`backend/domains.md` §4.2, no account-enumeration). On network failure → banner above the card with a Retry action.
- **Password visibility toggle** (`👁`): shows/hides, standard icon-button pattern, hit area ≥ 40×40px.
- **Enter key** submits the form from either field.
- **Continue with Google**: launches OAuth, button disabled while in flight. Failure (token rejected, no account because NG8) lands back here with the relevant message.
- **Use a code instead** → `otp-login`, pre-filling the email when present.
- **TOTP step (M24)**: when the primary factor succeeds but the user has 2FA enabled, the server returns `totp_required` with **no session issued** (`backend/domains.md` §4.10). The card swaps to the code field. Verify sends the TOTP code (or a single-use backup code via the "Use a backup code" toggle); wrong/used → `totp_invalid` inline, attempt-capped. Success → session issued → `files` / `?next`. The `pending` token makes this one flow across all three login paths (password, email-OTP, Google).
- **Forgot password?** → `forgot-password`.
- When arriving with `?next=`, a muted caption line under the title says "You'll be returned to where you left off after signing in" — so the redirect isn't a surprise.
- **Session expiry return**: the 401 path redirects to `login?next=<path>`; after a successful login the user lands exactly where they were.
- **A11y**: the error message is linked to the password field via `aria-describedby` and announced; focus moves to the first errored field on submit failure.
