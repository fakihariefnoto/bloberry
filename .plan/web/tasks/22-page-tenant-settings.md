# Task group — 22 page: tenant-settings

**Depends on:** `01-setup`, `02-design-tokens`, `03-routing`, `05-core-infra`, `06-shared-components` (PageHeader, FormField, ConfirmDestructive). **Blocks:** `33-flows` (backend-change flow). **Mockup:** [`web/mockup/tenant-settings.md`](../mockup/tenant-settings.md).

- [ ] **Layout — desktop** per the mockup: PageHeader "Tenant settings", sectioned settings card (GENERAL / STORAGE / QUOTA per the settings pattern), DANGER ZONE panel.
- [ ] **Layout — mobile**: stacked sections.
- [ ] **General** — tenant name editable + Save (dirty-field local save, unsaved indicator); slug read-only with a "shown, not editable" caption.
- [ ] **Storage backend** — dropdown of install-level backends + the tenant's own BYO backend, each with a health pill; **Change** → `ConfirmDestructive` with the ADR-4 contract: "New uploads will go to <backend>. Existing objects stay on <old> and keep resolving. This does not move data." Switching back is the same operation, no migration (PRD NG7).
- [ ] **Quota is read-only here** — the number + "set by platform admin" caption (PRD PA2; owners don't cap themselves). No edit affordance.
- [ ] **Used/limits** — the quota progress bar (312 GB of 500 GB · 62%).
- [ ] **DANGER ZONE — Delete tenant** — typed-name confirmation stating what it does and does not do: "folders, objects and keys are deleted; **bytes in the bucket are not** — they become orphans" (the reconciliation sweep finds them later).
- [ ] **Suspended state** — full-width banner "This tenant is suspended — reads and writes are blocked. Contact the platform admin."; form disabled, no Save.
- [ ] **Permission-aware** — route is `tenant_owner` only; a `tenant_admin` is refused by the guard (hidden section + server 403 → `forbidden`).

**tests:** backend-change confirm carries the ADR-4 contract; quota read-only; delete-tenant typed-name + orphan-bytes consequence; suspended disables the form.
