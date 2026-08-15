# Task group — 08 screen: otp-login

**Depends on:** `02-design-tokens`, `03-routing`, `04-core-infra`. **Blocks:** `25-flows`. **Mockup:** [`mobile/mockup/otp-login.md`](../mockup/otp-login.md).

- [ ] **Two-step flow** — step 1 email → step 2 code, advancing in place; email pre-filled from the route arg when present.
- [ ] **`otp_rate_limited`** — "Too many requests for this email. Try again in about an hour.", does not advance, re-request disabled during cooldown.
- [ ] **Single code field** — not N boxes; digits only, length 6, `autocomplete="one-time-code"`, paste + autofill work, Enter submits, numeric keyboard.
- [ ] **Resend cooldown** — live countdown (`00:42`), disabled until expiry, re-request restarts the timer.
- [ ] **Verify → success** → `files` (destination preserved).
- [ ] **`otp_invalid`** — inline error; 5th failure reads "Too many failed attempts — request a new code" (attempt cap, `backend/domains.md` §4.5).
- [ ] **Back to sign in** → `login`, state discarded.
- [ ] **A11y** — error announced; countdown announced once at expiry.

**tests:** rate-limited doesn't advance; wrong-code error; resend cooldown; one-field paste.
