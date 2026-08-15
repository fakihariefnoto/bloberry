# Task group — 07 screen: login

**Depends on:** `02-design-tokens`, `03-routing`, `04-core-infra` (envelope, auth). **Blocks:** `25-flows`. **Mockup:** [`mobile/mockup/login.md`](../mockup/login.md).

- [ ] **Form per the mockup** — email, password (visibility toggle), Log in primary, Google secondary, "Forgot password?", "Use a code instead". Scrollable form, `resizeToAvoidBottomInset` so the button stays reachable with the keyboard up.
- [ ] **Submit → loading** — button spinner replaces label, width held, disabled (guard rule 2).
- [ ] **Success** → `files` at the tenant root (or the preserved destination from a session-expiry return).
- [ ] **Invalid credentials** — inline field error under the password field, **identical for unknown email and wrong password** (`backend/domains.md` §4.2); password border `color.error`; focus moves to the errored field.
- [ ] **Network failure** — inline banner with Retry.
- [ ] **Google** — launches OAuth, disabled in flight; a valid identity with no account → `no_invitation` message (NG8).
- [ ] **"Use a code instead"** → `otp-login` (email pre-filled). **"Forgot password?"** → `forgot-password`.
- [ ] **Biometric-lock fallback** — "Use password" from `unlock` lands here; success returns to the destination the lock interrupted.
- [ ] **A11y** — field errors announced; proper labels, not placeholders-as-labels.

**tests:** invalid-credentials identical message; button disabled-while-processing; `?next`-style destination preserved; Google-failure landing.
