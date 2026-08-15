# Build index — Backend

Build-order table, dependency edges, and status for the Go backend task list.

## Build order

| File | Covers | Status |
|---|---|---|
| `01-setup.md` | Module init (the one `go mod init`), repo tree per §7, dev Mongo/Redis, migration tool, `.env.example`, Makefile fixes | ☐ |
| `02-migrations.md` | One migration per entity group from `ERD.md` (15 collections) in dependency order | ☐ |
| `03-platform.md` | `config`, `db`, `redis`, `httpx` (envelope/middleware/stream), `crypto`, `web/embed` | ☐ |
| `04-domain-auth.md` | auth domain (signup/login/refresh/forgot/OTP/Google/invitation) | ☐ |
| `05-domain-user.md` | user domain (profile, settings) | ☐ |
| `06-domain-tenant.md` | tenant domain (CRUD, quota, members, invitations, backend assignment) | ☐ |
| `07-domain-folder.md` | folder domain (tree CRUD, move, subtree delete, cycle prevention) | ☐ |
| `08-domain-object.md` | object domain (presign, multipart, direct, stat, move, visibility) | ☐ |
| `09-domain-share.md` | share domain (signed links, short URLs, resolution, revocation) | ☐ |
| `10-domain-apikey.md` | apikey domain (applications, access keys, principal resolution) | ☐ |
| `11-domain-grant.md` | grant domain (folder grants) | ☐ |
| `12-domain-job.md` | job domain (queue, worker, status — extract/bundle/subtree_delete) | ☐ |
| `13-domain-usage.md` | usage domain (hourly metering, cost estimation) | ☐ |
| `14-domain-audit.md` | audit domain (append-only event write/query) | ☐ |
| `15-domain-admin.md` | admin domain (backend registration, rate cards, health, install stats) | ☐ |
| `16-authz.md` | **The permission resolver** — not a domain, the highest-risk code (100% branch coverage) | ☐ |
| `17-storage.md` | **The driver abstraction** — not a domain; 5 drivers + conformance suite | ☐ |
| `18-flows.md` | Auth flows (7, from `domains.md` §4) + app-specific multi-step flows | ☐ |
| `19-openapi.md` | The committed `api/openapi.yaml` (spec-first) + `oapi-codegen` + SDK regeneration | ☐ |
| `20-guards.md` | lint/security/mocks/generate targets in CI, authz coverage gate | ☐ |

## Dependency edges

- **`01-setup` blocks everything.**
- **`02-migrations` blocks every domain file** (repositories read the collections).
- **`03-platform` blocks the domains that use it** — `httpx` (envelope/middleware) before any handler; `crypto` before `admin`/`storage`; `db`/`redis` before every repository.
- **`16-authz` blocks `folder`, `object`, `grant`, `apikey`, `share`** (their usecases call the resolver). Resolver-first: the pure function is testable before any domain depends on it.
- **`17-storage` blocks `object`, `share`, `job`, `admin`** (they call the driver/registry).
- **`19-openapi` is spec-first and gates everything**: server interfaces are generated from `api/openapi.yaml`. `make generate` must be green before domains compile — the spec is written up front, and CI fails if the working tree drifts.
- **`18-flows` depends on the domains it chains** (`auth` for the auth flows; `object`/`job`/`share` for the app-specific ones).

## External edges (from `architecture.md` §7)

- **Backend owns the module init** — `go mod init github.com/<org>/bloberry`. Every other platform depends on this file.
- **`01-setup` Blocks `web/tasks/01-setup.md`** (web scaffolds *inside* this repo), `cli/`, `desktop/`, `mobile/` (for the spec), `infra/`.
- **Build order `openapi → web build → go build`** — the frontend build must run **before** the backend build (the binary embeds `web/dist`). One `make build` target owns both in order; two independent CI jobs ships the previous frontend. Edge recorded in `00-index.md` and `web/tasks/00-index.md`.
- **`19-openapi` feeds everything** — `sdk/go`, `sdk/ts`, `mobile/lib/api/`, the CLI, and the server interfaces all derive from `api/openapi.yaml` (ADR-11).

## Gaps flagged

The backend Makefile's `swagger` target is the code-first template default (swaggo). `backend/domains.md` §1.1 **deviates deliberately**: OpenAPI is spec-first (`api/openapi.yaml` hand-authored, `oapi-codegen` generates server interfaces). `01-setup.md` carries a task to replace the `swagger` target with `generate` (oapi-codegen). Flagged here so it isn't assumed correct from the template.
