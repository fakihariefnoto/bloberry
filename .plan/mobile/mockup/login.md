# Screen — login

## Purpose & context

- **User goal**: authenticate to reach `files` at the tenant root.
- **Entry points**: "Log in" from `welcome`; the 401 session-expiry path (`Shell -.->|session expired| Login`); the biometric-lock fallback ("Use password"); after `reset-password` success.
- **Exit points**: authenticated → `files`; "Use a code instead" → `otp-login`; "Forgot password?" → `forgot-password`; "Continue with Google" → Google OAuth → `files`; **"Scan QR code" → `pair-login`** (sign in by scanning a QR from the web dashboard, M22). If the tenant is suspended → `tenant-suspended`.
- **Data needed**: email, password. `invalid_credentials` identical for unknown email and wrong password (`backend/domains.md` §4.2).

## States

- [x] Loading (submit in flight — button disabled, spinner)
- [x] Populated (form)
- [x] Error — invalid credentials
- [x] Error — network
- [x] Domain-specific — resuming with the biometric-unlock "Use password" path (same form, no layout change)
- [x] Domain-specific — **TOTP step** (M24): primary factor ok, `totp_required`, code field before any session is issued

## Style reference

- **Components used**: form screen pattern, primary button, secondary button (Google), text inputs, text links. No bottom tab bar (unauthenticated). Keyboard covers the form on small phones — the button must stay reachable (scrollable form, `Flutter: SingleChildScrollView` + `resizeToAvoidBottomInset`).
- No token deltas.

## Wireframe — mobile (happy path)

```
┌───────────────────────────┐
│  ←                        │
│                           │
│      ⬡  Bloberry          │
│                           │
│  Email                    │
│  ┌──────────────────────┐ │
│  │ you@project.dev      │ │
│  └──────────────────────┘ │
│                           │
│  Password                 │
│  ┌──────────────────────┐ │
│  │ ••••••••••     👁     │ │
│  └──────────────────────┘ │
│                           │
│  Forgot password?         │
│                           │
│  ┌──────────────────────┐ │
│  │      Log in          │ │
│  └──────────────────────┘ │
│                           │
│  ── or ──                 │
│ [⊙  Continue with Google] │
│                           │
│  Use a code instead       │
│  Scan QR code             │
└───────────────────────────┘
```

## Wireframe — mobile (invalid credentials)

```
┌───────────────────────────┐
│  Email                    │
│  ┌──────────────────────┐ │
│  │ you@project.dev      │ │
│  └──────────────────────┘ │
│                           │
│  Password                 │
│  ┌──────────────────────┐ │
│  │ ••••••••••     👁     │ │  ← border color.error
│  └──────────────────────┘ │
│  Wrong email or password. │
│                           │
│  ┌──────────────────────┐ │
│  │      Log in          │ │
│  └──────────────────────┘ │
└───────────────────────────┘
```

## Interactions

- **Log in**: on submit → button loading + disabled → success → `files`. On 401 → inline field error under password ("Wrong email or password." — identical for both causes). On network failure → inline banner with Retry.
- **Password visibility** toggle: standard, hit area ≥ 48dp.
- **Enter key** submits from either field.
- **Continue with Google**: launches OAuth; button disabled while in flight. A valid Google identity with no account → `no_invitation` message ("No account for this Google sign-in — ask your tenant admin to invite you", PRD NG8).
- **Use a code instead** → `otp-login` (email pre-filled when present).
- **Scan QR code** → `pair-login` (the camera overlay) — sign in by scanning a QR shown on the web dashboard; no password typed (M22). Requires the camera permission; see `pair-login.md` for the flow.
- **TOTP step (M24)**: when the primary factor succeeds but 2FA is enabled, the server returns `totp_required` with **no session issued** (`backend/domains.md` §4.10); the form swaps to a single 6-digit code field. Verify sends the code (or a backup code via a "Use a backup code" toggle); wrong/used → `totp_invalid`, attempt-capped. Success → `files`.
- **Biometric unlock fallback**: when the lock screen's "Use password" is tapped, this screen appears with the same form; success returns to the destination the lock interrupted.
- **A11y**: field errors are announced; focus moves to the first errored field; both inputs have proper labels (not placeholders-as-labels).
