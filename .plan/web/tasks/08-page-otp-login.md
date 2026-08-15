# Task group — 08 page: otp-login

**Depends on:** `01-setup`, `02-design-tokens`, `03-routing`, `05-core-infra`, `06-shared-components`. **Blocks:** `33-flows` (auth chain). **Mockup:** [`web/mockup/otp-login.md`](../mockup/otp-login.md).

- [ ] **Two-step flow in place** — step 1 email entry → step 2 code entry, advancing in place (no navigation between steps). Email pre-filled from the `?email`/route arg when present.
- [ ] **Step 1 submit** → loading → on `otp_sent` advance to step 2.
- [ ] **`otp_rate_limited`** — "Too many requests for this email. Try again in about an hour." — does **not** advance.
- [ ] **Code field** — single field (not N boxes), digits only, 6 chars, `autocomplete="one-time-code"`, accepts paste, Enter submits.
- [ ] **Resend cooldown** — disabled while counting down (live `00:42`), re-enabled after; re-request restarts the timer.
- [ ] **Verify** → loading → success → `files` (or `?next`).
- [ ] **`otp_invalid`** — inline error; 5th failure reads "Too many failed attempts — request a new code" (attempt cap, `backend/domains.md` §4.5).
- [ ] **Back to sign in** → `login`, discarding state.
- [ ] **A11y** — error announced; countdown announced once at expiry, not on every tick.

**tests:** rate-limited doesn't advance; wrong-code error; resend cooldown enabled/disabled; the one-field paste path.
