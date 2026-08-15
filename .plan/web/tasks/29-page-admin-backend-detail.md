# Task group — 28 page: admin-backend-detail

**Depends on:** `01-setup`, `02-design-tokens`, `03-routing`, `05-core-infra`, `06-shared-components` (PageHeader, FormField, StatusPill, ByteSize, ConfirmDestructive). **Blocks:** none. **Mockup:** [`web/mockup/admin-backend-detail.md`](../mockup/admin-backend-detail.md).

- [ ] **Layout — desktop** per the mockup: back link, PageHeader (backend name + driver/bucket/since + Test connection), the unreachable banner, CONFIG card (name, bucket, endpoint, **write-only credentials**), RATE CARD card, CAPABILITIES card (read-only), Used-by panel.
- [ ] **Layout — mobile**: stacked; the credentials field visible without scrolling in the unreachable state.
- [ ] **Test connection** — re-runs the health check immediately and re-renders the banner — the fix loop ("replace key → test → healthy") must not wait for the 5-minute ticker.
- [ ] **Credentials write-only** — the field shows "(current: set)" with a Replace-file action; saving without replacing leaves the existing ciphertext untouched (`ERD.md` storage-backends note, R7). Never a masked echo of the real value (the server can't return it).
- [ ] **Rate card edit** — the input side of `usage`'s "est. $21.40"; saving re-renders historical estimates from the then-current card (`ERD.md` usage-snapshots note — old months don't silently change).
- [ ] **Capabilities read-only** — they come from the driver + conformance suite (`backend/domains.md` §6.1), not this form. Displaying them answers "why doesn't R2 offer storage classes" (TRD R2).
- [ ] **Used by** — read-only, links to `admin-tenants`; reassignment happens on the tenant, so this panel just states the blast radius of deleting.
- [ ] **Delete** — only enabled for 0-tenant backends; typed-name confirm.
- [ ] **Unreachable state** — real error in the banner; the config card (where the fix lives) highlighted.
- [ ] **A11y** — error banner `role="alert"`; credentials replacement keyboard-reachable.

**tests:** test-connection re-renders health without the ticker; credentials unchanged-on-save; capabilities read-only; delete disabled with tenants.
