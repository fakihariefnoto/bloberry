# Task group — 23 page: profile

**Depends on:** `01-setup`, `02-design-tokens`, `03-routing`, `05-core-infra`, `06-shared-components` (PageHeader, FormField). **Blocks:** none. **Mockup:** [`web/mockup/profile.md`](../mockup/profile.md).

- [ ] **Layout — desktop** per the mockup: PageHeader "Profile", identity card (avatar + Change photo, display name field, email field with verified badge), Save, "Member of" list, member-since/last-login line.
- [ ] **Layout — mobile**: single column.
- [ ] **Save** — dirty-field save; changing email triggers re-verification (the `✓ verified` becomes "verification sent" in `color.warning`).
- [ ] **Unverified-email banner** — "Verify jane@acme.dev to receive reset and invite emails. [Resend]" on load when `email_verified` is false.
- [ ] **Member of** — rows show tenant + role; clicking switches tenant (`files` at the new root — the tenant-switch rule from `navigation.md`).
- [ ] **Avatar** — OS picker; the byte write uses the same presigned-PUT path as any object (hidden system folder). Placeholder acceptable in v1 if the storage plumbing isn't there.
- [ ] **Permission-aware** — always the current user's own profile; no role gating; no "view another member's profile" route on web.

**tests:** email-change re-verification; unverified banner; member-of row switches tenant to `files` at the new root.
