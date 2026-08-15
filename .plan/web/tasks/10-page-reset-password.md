# Task group — 10 page: reset-password

**Depends on:** `01-setup`, `02-design-tokens`, `03-routing`, `05-core-infra`, `06-shared-components`. **Blocks:** `33-flows` (auth chain). **Mockup:** [`web/mockup/reset-password.md`](../mockup/reset-password.md).

- [ ] **Token validation on load** — skeleton while validating; invalid/expired token → the expired-state wireframe (⚠ "This link has expired", "Reset links last 30 minutes and can be used once", "Request a new link" → `forgot-password`). No password fields in that state.
- [ ] **Form** — New password + Confirm, each with a visibility toggle, `autocomplete="new-password"`.
- [ ] **Validation** (centralized validators) — min 8 chars (helper text, not discovered on submit); confirm mismatch is a field error.
- [ ] **Submit** → loading → success ("Password set — sign in with your new password") → auto-redirect to `login`. Success caption notes **every other session is invalidated** (`backend/domains.md` §4.4) so the desktop app logging out on purpose isn't a surprise.
- [ ] **The token never renders** — route consumes it and replaces the path immediately (leak vector otherwise).

**tests:** expired-token state shows no form; mismatch field error; success invalidates-other-sessions caption present; token absent from the rendered URL.
