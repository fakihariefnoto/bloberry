# Task group — 19 screen: profile

**Depends on:** `02-design-tokens`, `03-routing`, `04-core-infra`. **Blocks:** none. **Mockup:** [`mobile/mockup/profile.md`](../mockup/profile.md).

- [ ] **Layout per the mockup** — avatar with Change, display name field, email field (✓ verified), Save, MEMBER OF list (tenant + role rows), member-since/last-login lines.
- [ ] **Save** — dirty-field save; changing email triggers re-verification (the `✓` becomes "verification sent" in `color.warning`).
- [ ] **Unverified-email banner** — "Verify jane@acme.dev to receive reset and invite emails. [Resend]" on load when `email_verified` is false.
- [ ] **Member of** — tapping a row switches tenant (`files` at the new root, `navigation.md` rule).
- [ ] **Avatar change** — OS image picker; byte write via the same presigned-PUT path as any object (hidden system folder). Placeholder acceptable in v1 if the plumbing isn't there.
- [ ] **A11y** — labels above fields; avatar `role="presentation"`; member rows ≥48dp.

**tests:** email-change re-verification; unverified banner; member row switches tenant.
