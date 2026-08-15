# Task group — 10 screen: reset-password

**Depends on:** `02-design-tokens`, `03-routing` (deep link), `04-core-infra`. **Blocks:** `25-flows`. **Mockup:** [`mobile/mockup/reset-password.md`](../mockup/reset-password.md).

- [ ] **Deep link + token consumption** — the emailed link opens `reset-password` (replace navigation); the token is validated on load (skeleton while checking) and **never rendered** (route replaced immediately).
- [ ] **Form per the mockup** — New + Confirm password fields, each with a visibility toggle, `autocomplete="new-password"`.
- [ ] **Validation** (centralized validators) — min 8 chars (helper text); confirm mismatch is a field error.
- [ ] **Submit** → loading → success ("Password set — sign in with your new password") → auto-redirect to `login`. Success caption notes **every other session is invalidated** (`backend/domains.md` §4.4 — the mobile app's other session dying on purpose is expected).
- [ ] **Invalid/expired token** → the expired wireframe (⚠ "This link has expired", "Reset links last 30 minutes and can be used once", "Request a new link" → `forgot-password`); no password fields in that state.

**tests:** expired-token state shows no form; mismatch field error; success invalidates-other-sessions caption; token absent from the route after consumption.
