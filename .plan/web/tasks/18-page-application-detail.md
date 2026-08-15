# Task group — 18 page: application-detail

**Depends on:** `01-setup`, `02-design-tokens`, `03-routing`, `05-core-infra`, `06-shared-components` (PageHeader, DataTable, StatusPill, SecretRevealModal, PermissionPicker, ConfirmDestructive, CopyableCode). **Blocks:** `33-flows` (key-issue flow). **Mockup:** [`web/mockup/application-detail.md`](../mockup/application-detail.md).

- [ ] **Layout — desktop** per the mockup: back link, PageHeader (app name + description + "Create key"), Access keys DataTable (masked key, scoped, perms, last used, `⋮`), the "keys shown once" footer.
- [ ] **Layout — mobile**: stacked; the `key-created` modal shown per its wireframe.
- [ ] **Key list** — masked always (`blob_live_••••4f2a`), never a full secret in a table; expiring < 7 days → warning pill; revoked → muted row with a Revoked pill (history retained); active → success pill.
- [ ] **Create key** — scope form with `PermissionPicker` (folder-subtree selector + read/write/delete/share checkboxes + optional expiry) and the **allow-only explanation inline** (PRD D7: a scope narrows, never widens; empty scope = whole tenant).
- [ ] **`key-created` modal** — full secret shown **exactly once** (`text.mono`, copy button with confirmation, "You won't see this again" in `color.warning`); dismissal requires the acknowledgement button — **no backdrop click** (`web/components.md`); copy is primary.
- [ ] **Revoke** — `ConfirmDestructive` stating the key's `last_used_at`/`last_used_ip` — "This key was last used Mar 13, 09:12 from 203.0.113.8. Revoking takes effect on the next request and is irreversible." **No Undo toast** (deliberate).
- [ ] **Last-active-key consequence** — revoking the app's only active key adds "This is <app>'s only active key. Its pipeline will fail on the next call" (PRD TA-E3).
- [ ] **Row actions** — `⋮` → Revoke (active/expiring) or Copy ID; expired keys show no Revoke (already dead).
- [ ] **Empty state** — "No keys yet · Create one to let <app> authenticate".
- [ ] **Loading / error** — skeleton / inline banner.
- [ ] **Permission-aware** — admin-wide create/revoke; viewers never reach the route.

**tests:** masked keys in the list; once-only modal (no backdrop dismiss, ack required); revoke confirm shows last-used; last-active-key warning; expiring pill threshold.
