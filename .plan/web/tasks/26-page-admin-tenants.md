# Task group — 25 page: admin-tenants

**Depends on:** `01-setup`, `02-design-tokens`, `03-routing`, `05-core-infra`, `06-shared-components` (DataTable, StatusPill, ByteSize, EmptyState, FormField, ConfirmDestructive). **Blocks:** `26-page-admin-tenant-detail`. **Mockup:** [`web/mockup/admin-tenants.md`](../mockup/admin-tenants.md).

- [ ] **Layout — desktop** per the mockup: PageHeader "Tenants" with `+ New tenant`, DataTable (tenant, status, used, quota, est. cost, `⋮`), the rate-card footer caption.
- [ ] **Layout — mobile**: stacked cards.
- [ ] **New tenant modal** — name, slug (auto-suggested from name, editable), quota, default backend (from the install-level pool). On save → `admin-tenant-detail` so the admin can immediately invite the owner (a tenant with no owner is inert, `ERD.md`).
- [ ] **Default sort = Est. cost desc** — the platform admin's actual question is "what's burning money"; sortable server-side.
- [ ] **Statuses** — `suspended` → warning pill; **over quota** → used/quota pair in `color.error` ("4.9 GB / 5 GB") — reads+write-blocked is a state to notice before the tenant's users do (PRD PA-E2).
- [ ] **Row click** → `admin-tenant-detail`; `⋮` → Suspend/Reactivate, Delete.
- [ ] **Delete tenant** — typed-name confirm with the orphan-bytes consequence (bytes stay in the bucket; reconciliation finds them).
- [ ] **Footer caption** — distinguishes this per-tenant screen from `admin-usage` (install-wide) so the two cost views don't read as duplicates.
- [ ] **Empty state** — "No tenants yet · Create the first tenant to provision storage".
- [ ] **Permission-aware** — platform_admin only; PLATFORM hidden for everyone else, routes server-guarded.

**tests:** cost-desc default sort; over-quota color pair; new-tenant lands on detail; delete orphan-bytes wording.
