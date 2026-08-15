# Task group — 07 page: login

**Depends on:** `01-setup`, `02-design-tokens`, `03-routing`, `05-core-infra`, `06-shared-components` (FormField, Toast). **Blocks:** `33-flows` (auth chain). **Mockup:** [`web/mockup/login.md`](../mockup/login.md).

- [ ] **Layout — desktop + mobile** from the mockup's happy-path wireframes: centered card on `color.background`, brand mark, Email + Password fields, Sign in primary button, Google secondary button, "Forgot password?" and "Use a code instead" links, the "No account? Ask your tenant admin" caption (PRD NG8 — no signup link).
- [ ] **`?next=` hint** — when arriving with a destination, a muted caption says the user will return there after login (not a separate layout).
- [ ] **Submit → loading** — button spinner replaces label, width held, control non-interactive (`02` button loading state).
- [ ] **Success** → navigate to `?next` destination or `files`.
- [ ] **Invalid credentials** — inline error under the password field ("Wrong email or password."), **identical for unknown email and wrong password** (no account-enumeration, `backend/domains.md` §4.2); password field border `color.error`; focus moves to the errored field.
- [ ] **Network failure** — banner above the card with Retry.
- [ ] **Password visibility toggle** — `👁`, hit area ≥ 40×40px.
- [ ] **Enter submits** from either field.
- [ ] **Google** — launches OAuth, disabled while in flight; `no_invitation` → message about needing an invitation (NG8).
- [ ] **"Use a code instead"** → `otp-login` (email pre-filled). **"Forgot password?"** → `forgot-password`.
- [ ] **TOTP step (M24)** — when the primary factor succeeds but 2FA is enabled, the server returns `totp_required` with no session issued; the card swaps to the single code field per the mockup wireframe. Verify sends the TOTP code or a single-use backup code ("Use a backup code" toggle); `totp_invalid` inline + attempt-capped; success → `files`/`?next`. One flow across password/OTP/Google via the pending token (`backend/domains.md` §4.10).
- [ ] **A11y** — error linked to the password field via `aria-describedby` and announced; labels, not placeholders-as-labels.

**tests:** invalid-credentials identical-message test; `?next` round-trip; button-disabled-while-processing; Google-failure landing.
