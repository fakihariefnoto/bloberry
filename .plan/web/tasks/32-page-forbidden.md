# Task group — 31 page: forbidden

**Depends on:** `01-setup`, `02-design-tokens`, `03-routing`. **Blocks:** none (terminal). **Mockup:** [`web/mockup/forbidden.md`](../mockup/forbidden.md).

- [ ] **Single-state page** — server-side guards redirect here on a 403 (a `member` typing `/members` directly, a non-admin hitting a platform route with a stale URL). The route replaces the failed one.
- [ ] **The message names the required role** from the route guard, not a generic "forbidden": "Members requires the `tenant_admin` or `tenant_owner` role." (`design/style-guide.md` → Permission-denied state: a caption naming what's needed).
- [ ] **"If you need it, ask a tenant admin to change your role or grant you folder-level access."**
- [ ] **Back to Files** → `files`; the sidebar stays fully functional (the user isn't locked out of what they *can* do).
- [ ] **Style** — `color.warning` icon, not `color.error` (a permission boundary, not a crash).
- [ ] **A11y** — heading + action in a single `main` landmark.

**tests:** role name in the message matches the route guard; back-to-files works; sidebar intact.
