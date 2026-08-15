# Task group — 09 page: forgot-password

**Depends on:** `01-setup`, `02-design-tokens`, `03-routing`, `05-core-infra`, `06-shared-components`. **Blocks:** `33-flows` (auth chain). **Mockup:** [`web/mockup/forgot-password.md`](../mockup/forgot-password.md).

- [ ] **Form** — email field + "Send reset link" primary + "Back to sign in".
- [ ] **Submit** → loading → **identical confirmation whether or not the email exists** (no account-enumeration oracle, `backend/domains.md` §4.4): "If that email exists, a reset link is on its way to <email>. It expires in 30 minutes."
- [ ] **Confirmation state** — success panel replaces the form (check icon, "Check your inbox", TTL stated), with "Return to sign in" → `login`.
- [ ] **Network failure** — banner with Retry (the request genuinely didn't go out).
- [ ] **A11y** — confirmation announced as "Sent"; focus moves to the confirmation heading.

**tests:** identical response for existing/nonexistent email; network-failure retry.
