# Task group — 08 commands: admin

**Depends on:** `02-core-infra`, `03-auth` (platform_admin role). One checkbox per command; acceptance = "output matches the sample in `cli/commands/<file>.md`". These are the install-provisioning surface.

## `admin backend add` (`cli/commands/admin-backend-add.md`)

- [ ] Register a backend: `--config key=value` (bucket/endpoint/region), **`--credential-file` for secrets** (the process-list rule — `--config key=secret` risks `/proc` exposure), `--rate-card` optional. Credentials envelope-encrypted, **never echoed** (`credentials_set: true` only). Name unique per install → exit 8; health starts `unchecked`.

## `admin backend list` (`cli/commands/admin-backend-list.md`)

- [ ] List backends with health + assigned-tenant counts; **the raw provider error on unreachable backends** (the one place it's legal, PRD PA-E1).

## `admin backend test` (`cli/commands/admin-backend-test.md`)

- [ ] Immediate reachability probe; **exit 0 healthy / 9 unreachable** (scripts page on 9); the failure output shows the real error + the fix hint; not a full conformance run (fast enough for a fix loop).

## `admin backend rate-card` (`cli/commands/admin-backend-rate-card.md`)

- [ ] Show or set (`--set storage,egress,requests`); no rate card → "unknown" + the exact `--set` fix command; setting is immediate + non-destructive (no confirm).

## `admin tenant create` (`cli/commands/admin-tenant-create.md`)

- [ ] Create a tenant with `--backend` (required), `--quota-bytes` (`500GB`/`2TB`/`0` = unlimited), `--quota-objects`, `--slug`; root folder created with it; output chains to the owner-invite next step (a tenant with no owner is inert).

## `admin tenant list` (`cli/commands/admin-tenant-list.md`)

- [ ] All tenants with usage/quota/cost; **default sort = est. cost desc** (the "who's burning money" answer); `over quota`/`suspended` as plain text words (script-greppable).

## `admin tenant quota` (`cli/commands/admin-tenant-quota.md`)

- [ ] Show or set quota (`--set-bytes`, `--set-objects`); over-quota state called out with the reads-still-work caveat (PRD PA-E2); writes unblock immediately on raise (the exit-6 remediation).

## `admin usage` (`cli/commands/admin-usage.md`)

- [ ] Install-wide usage + cost with a cost-desc per-tenant breakdown; the `unknown` row names the backend + fix command; egress labeled estimated.

**tests:** per command — happy path, error, `--json`; `admin backend test` exit 9; `admin tenant create` size parsing; quota unblock.
