# Task group — 24 screen: tenant-suspended

**Depends on:** `02-design-tokens`, `03-routing`. **Blocks:** none (terminal). **Mockup:** [`mobile/mockup/tenant-suspended.md`](../mockup/tenant-suspended.md).

- [ ] **Single-state page per the mockup** — ⚠ "This tenant is suspended", "Acme Inc's storage is paused. Reads and writes are blocked until the platform admin reactivates the tenant.", "Your files are safe. Nothing has been deleted.", Contact-platform-admin action, "You can still switch to other tenants or log out." Tab bar stays visible.
- [ ] **Shell guard** — any authenticated request while the tenant is suspended redirects here before the screen renders its own failure.
- [ ] **No data-fetching or destructive actions** — Contact opens the support surface (mailto/URL), never an in-app form (a suspended tenant shouldn't write anything, including a ticket).
- [ ] **Tenant switcher stays live** (app bar + `more`) — a multi-tenant user leaves the suspended tenant, per the caption. Log out works.
- [ ] **Tabs render this panel** — the user sees the boundary, not a wall of failures.
- [ ] **Style** — `color.warning` accent (administrative, not a crash).
- [ ] **A11y** — the warning announced on arrival; the support button is a plain external link.

**tests:** suspended tenant lands here (not a failure wall); switcher/logout still work; support opens externally.
