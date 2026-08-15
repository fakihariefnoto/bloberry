# Changelog — Bloberry

Format based on [Keep a Changelog](https://keepachangelog.com/), versions follow [SemVer](https://semver.org/).

## [Unreleased]

### Added

- **CLI `init` command** (48th) — one-command first-run: server URL → reachability → tenant → auth; non-interactive via `--server --tenant --token` for CI; `--token` never written to the config file.
- **TOTP 2FA (M24)** — web account-settings provisioning (secret shown once, confirm-before-enable, 10 single-use hashed backup codes), gating every human login (password/OTP/Google) on all surfaces; access keys exempt; platform-admin-only recovery (R15).
- **Azure Blob Storage — 6th storage driver** — own `azblob` SDK (SharedKey/SAS, container primitive), conformance suite + weekly real-provider CI gain a sixth target.
- **Cloud-backend scale-out model** (`infra/README.md` §Scaling out, ADR-15) — stateless API nodes + Mongo replica set + Redis sentinel, job-worker split, scheduler leader lease; local-disk driver single-node and absent from scaled installs.
- **Build task lists generated for every platform** (`build-backend/web/mobile/cli/desktop/infra`) — 120 task files total: `backend/tasks/` (21), `web/tasks/` (36), `mobile/tasks/` (29), `cli/tasks/` (15), `desktop/tasks/` (9), `infra/tasks/` (10). Each numbered in dependency order with Depends-on/Blocks headers, per-endpoint/per-screen/per-command depth grounded in the designed artifacts, and ending at verified-installed/verified-working where shipping matters. No gaps flagged in any pass.
- **`web/tasks/` build list generated** (`build-web`, 36 files) — setup (scaffold as a package inside the existing repo, depends on the backend's module), design tokens before any page, routing with role gates + `?next`, env separation, core infra (envelope, generated API client, snake_case once), 23 shared components, one file per page (27) grounded in its mockup, 9 cross-page flows, and a named testing file. No gaps flagged.
- **All 48 CLI commands designed** (`cli/commands/`) — `init` plus the 7 short file verbs, auth, folder, share, key, app, grant, tenant, job, archive, the 8 admin verbs, config, completion and version. Real `--help` text, per-state sample output and exit codes per command, ready to become golden-file test fixtures.
- **All 27 web mockups** (`web/mockup/`, 61 wireframes at mobile + desktop widths) — auth (incl. the TOTP step and `pair-device`), the breadcrumb+table file browser, file-detail, shares, jobs, applications/keys, members, audit, usage + cost, tenant settings, the platform-admin set, and the terminal states. All wireframes pass `check-wireframes.mjs`; the web route graph is closed.
- **All 21 mobile mockups** (`mobile/mockup/`, 38 wireframes) — the 4-tab shell with the persistent uploads queue, unlock, browse/share, `pair-login` (QR scan), and the revoke-a-key flow. Mobile route graph closed.
- Initial idea and plan captured, from `notes/agnostic-storage-drive.md`.
- Stack settled: Go + chi backend on MongoDB and Redis; Vue 3 + Vite + Tailwind + Reka UI (Bun) built to static assets and embedded in the Go binary; Flutter mobile (`com.bloberry.app`); Go CLI as a companion to the backend; Wails desktop wrapping the Vue build; systemd + Caddy on a self-hosted VPS.
- Scope decisions: **six** storage drivers in v1 (S3, Cloudflare R2, Alibaba OSS, Google Cloud Storage, Azure Blob Storage, local VPS disk); first-party SDKs in Go, TypeScript and Dart; five human roles (`platform_admin`, `tenant_owner`, `tenant_admin`, `member`, `viewer`) plus application principals holding access keys, with folder-level grants layered on top of the role.
- `design/style-guide.md` complete — blueberry-indigo palette generated from `#8B7DEB` with every WCAG AA pair passing, Soft shape language for radius/elevation/motion, Lucide icon set, and component specs covering the file table, folder tree, dropzone, upload queue, status pills and secret display.

- `ERD.md` written: the authoritative MongoDB data model — 15 collections, denormalization decisions, a full compound-index table (every index prefixed with `tenant_id`), and per-entity lifecycle/nullability notes.
- `architecture.md` written: system context, container breakdown, 8 multi-container sequence flows, deployment topology with trust boundaries, cross-cutting concerns, **15 ADRs**, and the §7 implementation layout.
- PRD detailed: 10 measurable goals, 9 explicit non-goals, ~40 user stories with edge cases traced to **25 Must requirements** (M22 QR pairing, M23 config login, M24 TOTP 2FA added later), and a D1–D10 decisions table recording every resolved open question with its reasoning.

### Changed

- **Web UI kit switched from shadcn-vue to Reka UI.** The dashboard's surfaces are styled almost entirely by the style-guide's exact tokens, so the unstyled headless primitives let those tokens be the theme rather than fighting a themed kit. TanStack Table stays for tabular data. Mobile keeps shadcn/ui for Flutter, now justified by shared style-guide tokens rather than a shared kit.
- **All six surfaces now target v1** (backend, web, CLI, SDKs, mobile, desktop) rather than phasing mobile/desktop to v1.1. Accepted as the plan's largest schedule risk; mitigated by contract-first OpenAPI with generated clients, and by desktop reusing the web build verbatim.
- **Per-tenant cost estimation moved from non-goal to Must (M18).** Billing is built in v1 but switched off — metering, rate cards per storage backend, and estimated monthly cost per tenant ship; invoicing, payments and plans do not.
- **Scaffold ownership settled: `backend/` owns the single `go mod init`, and `desktop/` must not run `wails3 init`.** Bloberry composes three repo layouts at once (embedded web + companion CLI + desktop-as-network-client), so the standard "Wails owns the scaffold" rule doesn't apply — here the Go module already exists and Wails is a third binary inside it.
- `object` records now carry their own storage-backend pointer rather than inheriting the tenant's current one, so switching a tenant's backend doesn't strand existing files (supports NG7 — no bulk migration in v1).

### Fixed
