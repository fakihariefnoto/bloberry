# Task group — 30 page: tenant-suspended

**Depends on:** `01-setup`, `02-design-tokens`, `03-routing`. **Blocks:** none (terminal). **Mockup:** [`web/mockup/tenant-suspended.md`](../mockup/tenant-suspended.md).

- [ ] **Single-state page** — the shell guard redirects here on any authenticated request while the tenant is suspended, before the page renders its own failure. Rendered inside the `AppShell` (sidebar visible but sections inert) per the mockup.
- [ ] **Content** — ⚠ "This tenant is suspended", "Acme Inc's storage is paused. Reads and writes are blocked until the platform admin reactivates the tenant.", "Your files are safe. Nothing has been deleted.", Contact-platform-admin action, "You can still log out and switch to other tenants."
- [ ] **No data-fetching or destructive actions** — the page reassures and directs. Contact opens the configured support surface (mailto/URL), never an in-app form (a suspended tenant shouldn't write anything, including a ticket).
- [ ] **Tenant switcher stays live** — a multi-tenant user can leave the suspended tenant (the caption says so explicitly, so it doesn't read as a full lockout). Log out works from the user menu.
- [ ] **Style** — `color.warning` accent, not `color.error` (suspension is administrative, not a crash).
- [ ] **A11y** — the warning announced on arrival; the support button is a plain external link.

**tests:** suspended tenant lands here (not a 500 wall); switcher/logout still work; support opens externally.
