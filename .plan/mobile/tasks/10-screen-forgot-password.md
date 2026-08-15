# Task group — 09 screen: forgot-password

**Depends on:** `02-design-tokens`, `03-routing`, `04-core-infra`. **Blocks:** `25-flows`. **Mockup:** [`mobile/mockup/forgot-password.md`](../mockup/forgot-password.md).

- [ ] **Form per the mockup** — email field, "Send reset link" primary, back.
- [ ] **Submit** → loading → **identical confirmation whether or not the email exists** (no account-enumeration, `backend/domains.md` §4.4): "If that email exists, a reset link is on its way to <email>. It expires in 30 minutes."
- [ ] **Confirmation state** — success panel replaces the form (check icon, TTL stated), "Return to sign in" → `login`.
- [ ] **Network failure** — banner with Retry.
- [ ] **A11y** — confirmation announced as "Sent"; focus to the confirmation heading.

**tests:** identical response for existing/nonexistent email; retry on network failure.
