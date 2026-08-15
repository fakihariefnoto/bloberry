# TRD — Bloberry

## Architecture overview

Bloberry is **one Go binary** that serves the JSON API, the embedded Vue dashboard, and the public/short-URL surfaces, backed by **MongoDB** (metadata, tenancy, permissions, audit) and **Redis** (sessions, refresh tokens, short-URL cache, extraction job queue). Object bytes never live in MongoDB — they live in whichever **storage backend** a tenant is configured for, reached through a single `StorageDriver` interface with **six** implementations (S3, R2, OSS, GCS, Azure Blob, local disk).

Everything else is a client of that binary: the Vue dashboard (embedded via `go:embed`), the Wails desktop shell (same Go code, same Vue frontend, wrapped natively), the Flutter mobile app, the Go CLI, and the three published SDKs — all speaking the same OpenAPI contract.

The full system-level design — context diagram, container breakdown, cross-platform flow sequences, deployment topology, architecture decisions, and the implementation layout — lives in [architecture.md](architecture.md). This section is a summary; don't re-derive a second system diagram here.

## Stack decisions

| Area | Choice | Notes |
|------|--------|-------|
| Mobile | Flutter | Always Flutter. Capture/pick/upload, browse, share, offline upload queue. |
| Web framework | **Vue** | Tailwind CSS + **Reka UI** (unstyled headless primitives — Vue's Radix equivalent). Nuxt-vs-plain-Vue is a `detail-web` sub-choice, but "embedded" topology points at plain Vue + Vite (static build, no Node server at runtime). |
| Web runtime / package manager | **Bun** | Build-time only — the dashboard compiles to static assets consumed by `go:embed`, so nothing Bun-specific ships to production. Reka UI is plain Vue 3 components — no Vite-plugin compatibility risk. |
| Backend language | **Go** | Modular monorepo. Chosen for concurrent streaming of large multipart uploads, first-party S3/GCS/OSS SDKs, cheap CPU-bound compression, and a single static binary for self-hosting. Also the CLI's language and Wails' host language — one language across three of six surfaces. |
| Backend web framework | **chi** | stdlib `http.Handler` all the way down, which matters here: streaming multipart uploads and HTTP range-request downloads are simplest when nothing wraps the raw `ResponseWriter`. |
| Database | **MongoDB** | Specified in the source note. Good fit for the dominant shape: object metadata is document-like with per-driver variable fields, and the folder tree maps naturally onto a materialized path + ancestor array. See R4 for the tradeoff on relational-shaped RBAC data. |
| User roles | **Yes** — `platform_admin`, `tenant_owner`, `tenant_admin`, `member`, `viewer` | Plus a non-human **application** principal type holding access keys. Folder-level grants (principal × folder subtree × permission set × expiry) layer on top of the role — the role is the floor, the grant is the reach. |
| Web + backend topology | **Embedded (same deployable)** | Vue builds to static assets embedded in the Go binary via `go:embed` and served by the same chi router. One binary to self-host — the whole point for this product. The Wails desktop shell reuses the identical frontend build. |
| CLI language | **Go** | Shares the API client package, envelope types and OpenAPI-generated models with the backend in one monorepo. GoReleaser drives every install channel from one config. |
| CLI role | **Companion to Backend** | Talks to a Bloberry server over the same API the SDKs and dashboard use; shares its auth (device login for humans, access keys for CI) and its committed OpenAPI contract. |
| Desktop framework | **Wails** | Go host + the same Vue frontend. Native file dialogs, filesystem watching and the sync engine are Go code that can reuse the backend's own client package. |
| Desktop wraps/reuses | **Web** | Wraps `web/`'s build output. Chrome deltas only (menu bar, tray, drag-and-drop target, sync status) — no parallel screen set. |
| App identifier (mobile) | `com.bloberry.app` | Must be the exact same string as Android `applicationId` and iOS `PRODUCT_BUNDLE_IDENTIFIER`. |
| VPS process model | **systemd** | The Bloberry binary as a systemd unit; MongoDB and Redis installed natively alongside. Matches the "one binary, few moving parts" self-hosting story. A Docker Compose file is still published as a convenience for self-hosters — but the reference deployment is systemd. |
| Reverse proxy | **Caddy** | Automatic TLS. Terminates the API/dashboard domain and the short-URL domain (see Q6 in `PRD.md`). Must be configured to **not** buffer request bodies, or large uploads will be buffered to the proxy's disk before reaching Go. |

## Tech stack

**Backend (Go)**

| Concern | Choice | Why |
|---|---|---|
| HTTP router | `go-chi/chi/v5` | stdlib-compatible; no `ResponseWriter` wrapping to fight during streaming. |
| MongoDB driver | `go.mongodb.org/mongo-driver/v2` | Official; the v2 API. |
| Redis client | `redis/go-redis/v9` | Sessions, refresh tokens, short-URL cache, job queue. |
| S3 / R2 | `aws/aws-sdk-go-v2` | One driver serves both — R2 is S3-compatible via endpoint override, as are MinIO, B2, Spaces and Wasabi. R2 presigning needs its own account-scoped endpoint and does **not** support all S3 features (see R2 below). |
| Alibaba OSS | `aliyun/aliyun-oss-go-sdk` | Separate SDK, separate presign semantics. |
| Google Cloud Storage | `cloud.google.com/go/storage` | Presigning uses a service-account signer, not a static key pair — a different credential shape from the other three. |
| Azure Blob Storage | `github.com/Azure/azure-sdk-for-go/sdk/storage/azblob` | SharedKey/SAS/AAD auth, container-vs-bucket, block-blob staging for multipart — a separate SDK and primitive model (like OSS/GCS, not an S3 endpoint override). |
| Local disk | stdlib `os` + `http.ServeContent` | Presigning is Bloberry-issued HMAC tokens against its own download endpoint, since there is no external signer. |
| JWT | `golang-jwt/jwt/v5` | Access tokens; refresh tokens are opaque and Redis-backed. |
| Archive | stdlib `archive/zip`, `archive/tar`, `compress/gzip` | Streaming where possible; extraction is queued, not inline (PRD D3). |
| OpenAPI | `oapi-codegen` | Spec-first. The committed spec drives the server interfaces, the Go SDK, and — via `openapi-generator` — the TypeScript and Dart SDKs. |
| Config | `spf13/viper` or stdlib + `caarlos0/env` | Confirm in `detail-backend`. |
| Lint / security | `golangci-lint`, `gosec` | `make lint`, `make security`. |

**Web (Vue)** — Vue 3 + Vite, TypeScript, Tailwind CSS, **Reka UI** (unstyled headless primitives: dialog, dropdown, popover, tabs, tooltip, toast, table — styled with Tailwind into the app's own components, not a themed kit), Bun as package manager and build runner. Built to static assets; `go:embed` picks up `web/dist`. Nuxt is not expected to be used (it would want its own Node server, contradicting the embedded topology) — confirmed in `detail-web`.

**Mobile (Flutter)** — per `templates/flutter-defaults.md`. UI kit picked in `detail-mobile`. Snake_case JSON deserialization configured explicitly (the Go backend emits snake_case; Dart defaults to camelCase).

**CLI (Go)** — `spf13/cobra` + `spf13/viper`, secrets in the OS keychain via `zalando/go-keyring`, releases via GoReleaser.

**Desktop (Wails v3)** — Go host, Vue frontend build shared with `web/`. Cross-compilation is **not** available (every desktop framework needs a native webview per OS), so its CI needs macOS and Windows runners — an explicit exception to the self-hosted-VPS-runner default, recorded in `infra/README.md`.

**SDKs** — Go (module in the monorepo, tagged separately), TypeScript (npm, generated from OpenAPI + a hand-written upload helper), Dart (pub.dev — pending Q8; if it stays internal, it lives under `mobile/` instead).

## Data model

MongoDB. Baseline auth/user collections plus Bloberry's own: **tenant**, **membership** (user × tenant × role), **application** (non-human principal), **access_key** (bearer, hashed at rest), **folder** (materialized path + ancestor-ID array, per PRD D2), **object** (the `file_id` record — metadata, its *own* backend pointer, visibility, content hash), **grant** (principal × folder subtree × permissions × expiry, allow-only per PRD D7), **share_link** (signed/short URL with TTL and revocation), **storage_backend** (driver type + encrypted credentials + rate card; nullable `tenant_id` for install-level vs BYO, per PRD D4), **usage_snapshot** (metered bytes/objects/bandwidth per tenant, feeding cost estimation), **audit_event**, and **job** (queued extraction, bundle generation, large subtree delete).

Two shapes worth flagging up front: `object` carries **its own** storage-backend reference rather than inheriting the tenant's current one — that's what makes PRD NG7 (no bulk migration in v1) survivable, since old objects keep resolving after a tenant switches backend. And `storage_backend` carries the **rate card** that PRD M18's cost estimation reads.

The authoritative ER diagram, collection shapes, **index design** and per-entity notes live in **[ERD.md](ERD.md)**. Don't paste a second copy of the diagram here.

Worth knowing without opening it: three of the PRD's measurable goals (G3 listing latency, G4 stable `file_id`, G5 authz latency) are **index designs**, not code optimizations — every index is compound and prefixed with `tenant_id`, which is simultaneously the isolation boundary and the highest-selectivity key. The two shapes carrying the most design weight are `ancestors[]` (denormalized onto objects as well as folders, so authorization never walks a tree) and `objects.backend_id` (per-object, so a tenant backend switch strands nothing).

## APIs & integrations

**Internal API** — REST/JSON over the standard envelope `{data?, messages?: [{code, content}]}`, both `omitempty`, no `error` boolean; the HTTP status code is the success/failure signal and `messages` carries both errors and success confirmations. Committed OpenAPI spec at `api/openapi.yaml`, shared by mobile, web, CLI and all three SDKs.

Surface groups (fleshed out in `detail-backend`):

| Group | Covers |
|---|---|---|
| `auth` | signup, login, refresh, logout, forgot-password, OTP login, Google login, **QR pairing (M22 — issue one-time token, verify by scan)**, **config-file login (M23 — export a signed, passphrase-encrypted login file)**, **TOTP 2FA (M24 — provision secret, confirm enable, verify on every human login, backup codes)** |
| `users` | profile, settings |
| `tenants` | CRUD, quota, storage backend assignment, members, invitations |
| `folders` | create, rename, move, delete-subtree, list children, tree |
| `objects` | direct upload, presigned-PUT init, multipart init/part/complete, download, stat, move, delete, visibility |
| `shares` | signed link create/revoke, short URL create/resolve |
| `applications` | register, access-key issue/list/revoke, last-used |
| `grants` | folder-level grant create/list/revoke |
| `archives` | extract-on-upload, bundle-download |
| `jobs` | status polling by job ID — extraction, bundle generation, large subtree delete (PRD M21) |
| `audit` | query per tenant |
| `usage` | per-tenant and install-wide bytes/objects/bandwidth, plus **estimated monthly cost** from the backend's rate card (PRD M18) |
| `admin` | storage backend registration + rate cards, backend health, tenant management, install stats |

**Two auth schemes on the same API** — human sessions (JWT access + Redis-backed opaque refresh, platform-aware TTLs: long-lived for mobile, short for web) and application access keys (bearer, prefixed, hashed at rest — PRD D5; the key-scope narrowing lives in `backend/domains.md` §5). Both resolve to a principal that the same RBAC middleware evaluates.

**Third-party integrations** — AWS S3, Cloudflare R2, Alibaba OSS, Google Cloud Storage (as storage backends); Google OAuth (login); an SMTP provider (verification, password reset, OTP, invitations). No payment provider in v1.

## Key technical risks

| # | Risk | Why it's real here | Mitigation |
|---|---|---|---|
| **R1** | **The six drivers don't actually share one interface.** Presigning differs most: S3/R2 use static key-pair signing, GCS needs a service-account signer (and IAM `signBlob` if running without a key file), OSS has its own signature version, Azure Blob uses SharedKey/SAS, and local disk has no external signer at all — Bloberry must issue its own HMAC tokens against its own endpoint. Multipart part-size minimums and ETag semantics also differ. | The whole product is this abstraction. If it leaks, every caller ends up branching on backend type and the value proposition is gone. | Define the interface against the **hardest** driver (local disk for presigning, GCS for credentials), not the easiest. Write one conformance test suite every driver must pass. Run it against real S3/R2/OSS/GCS/Azure Blob in CI, plus MinIO locally. |
| **R2** | **R2 is not fully S3-compatible.** No storage classes, different multipart part-size behavior, no `GetObjectAttributes`, and presigned URLs require the account-scoped endpoint. | Treating "R2 = S3 with a different endpoint" will produce runtime failures only under load or only on multipart. | Keep R2 as a *configuration* of the S3 driver but with an explicit capability flag set; make the conformance suite assert capabilities rather than assume them. |
| **R3** | **Folder moves at scale.** Moving a subtree with a materialized path rewrites the path of every descendant. A 100k-object subtree is a very large update. | "Reorganize after the fact" is a headline goal (M3, and the `file_id` stability promise in M4). | Store `ancestors: [folder_id]` alongside `path`, so permission checks and listings key off indexed IDs, not string prefixes. Path becomes a display concern that can be rebuilt in the background. Bound the operation and make it a job if it exceeds a threshold. |
| **R4** | **RBAC in MongoDB is the relational-shaped part of a document database.** Resolving "can principal P do action A on folder F" spans membership, grant, access_key scope and folder ancestry — a join in a store without joins. | This check runs on **every** request. Getting it wrong is a security bug; getting it slow is a latency bug on the hot path. | Denormalize the ancestor array onto the object document so a check is one indexed query, not a tree walk. Cache resolved principal→permissions in Redis with explicit invalidation on grant/key change. Write the permission resolver as one pure function with an exhaustive table-driven test suite — it is the single highest-risk function in the codebase. |
| **R5** | **Streaming large uploads through Caddy and Go.** Default reverse-proxy buffering will spool multi-GB bodies to the proxy's disk; a naive handler will read the body into memory. | Uploads are the core operation. This fails only at the sizes that matter. | Explicitly disable request buffering in Caddy. Stream `io.Copy` to the driver, never `io.ReadAll`. Enforce a hard body-size ceiling. Prefer presigned-PUT (bytes bypass Bloberry entirely) as the default path for browser uploads. |
| **R6** | **Server-side archive extraction is a denial-of-service surface.** Zip bombs, path-traversal entries (`../../etc/passwd`), symlink entries, and quota exhaustion mid-extraction. | M11 is a listed must-have and it's the most attacker-friendly feature in the product. | Queue it, never do it inline. Enforce a decompressed-size ceiling and a ratio ceiling; reject absolute paths, `..` segments and symlinks; check quota per entry, not once up front; extract to a staging prefix and commit atomically. |
| **R7** | **Storage credentials at rest.** Bloberry holds the keys to every tenant's bucket — it is a high-value target in a way a normal app database isn't. | A dump of the `storage_backend` collection is a full compromise of every connected bucket. | Envelope-encrypt credentials with a key from the environment (never in MongoDB), decrypt only in memory at driver construction, never log or return them through the API, and never echo them back to the dashboard (write-only fields). |
| **R8** | **Orphaned objects and the metadata/bytes split.** Bloberry writes metadata to Mongo and bytes to a remote store; there is no transaction across the two. A crash between them leaves either an unreferenced blob (silent cost) or a metadata row pointing at nothing (broken download). | Guaranteed to happen, not hypothetical — it's a two-phase write with no coordinator. | Write metadata first in a `pending` state, promote to `active` only after the byte write confirms. A reconciliation sweep hard-deletes `pending` records past a TTL and reports blobs with no metadata. Accept eventual consistency explicitly rather than pretending it's atomic. |
| **R9** | **Six surfaces, one contract.** Backend, web, mobile, desktop, CLI and three SDKs all consume the same API. A breaking change ripples eight ways. | The plan's breadth is its main schedule risk. | Spec-first: the committed OpenAPI file is the gate. Generate the Go/TS/Dart SDKs from it rather than hand-writing three clients. Version the API path from day one (`/v1/`). Desktop and web share literally the same frontend build, which removes one whole surface's drift. |
| **R10** | **Wails cannot cross-compile.** macOS and Windows desktop builds each need a native runner — the self-hosted-VPS-runner default doesn't cover them. | Discovered late, this stalls the desktop release entirely. | Recorded now in `infra/README.md` as an explicit exception; `detail-infra`/`build-desktop` plan the OS-native CI matrix. |
| **R11** | **A signed URL outlives its revocation.** On the default redirect download path (`architecture.md` ADR-3), once a presigned URL is issued Bloberry is out of the loop — revoking the grant, the key or the share link does **not** invalidate an already-issued URL until its TTL expires. | This directly contradicts the intuition "revoke means revoked", which PRD G5 establishes for everything else in the system. Someone will rely on it during an incident. | Keep presign TTLs short (5 minutes for API-issued downloads). Say so plainly in the sharing UI rather than burying it. Where a tenant genuinely needs instant revocation, that grant forces the proxy path — which is why the branch exists rather than being a local-disk special case. |
| **R12** | **The job worker shares a process with the API** (`architecture.md` ADR-8). A large extraction is CPU- and memory-heavy and competes directly with request handling. | The one-binary story (PRD G8) is a headline goal, so this is a deliberate trade — but it means a single 2 GB zip can degrade dashboard latency for every tenant on the install. | Bound worker concurrency to well below `GOMAXPROCS`; stream rather than buffer during extraction; enforce decompressed-size ceilings (R6) which also bound memory. Because the boundary is already a Redis queue, splitting the worker into its own process later is a deployment change rather than a rewrite — make that the documented escape hatch. |
| **R13** | **A pairing token or config file is a login capability.** The mobile QR (M22) and the desktop config file (M23) both mint a session for the scanned/imported identity — whoever holds one can log in as that user. | These are the first credentials in the system that live *outside* the user's own control (a screenshot, a downloaded file) before they're consumed. | QR token: one-time, ~2 min TTL, rate-limited verify, and **never logs the user in more than once** (single-use DEL). Config file: passphrase-encrypted **client-side** so the passphrase never transits the server, a server-signed **import window** so an old file is dead, and the resulting session stays revocable via normal `auth logout` (DT-E2). State the capability nature in the web UI ("this code signs you in — expires in 2 minutes"). |
| **R14** | **The config file passphrase is the user's last line of defense.** A leaked file with a weak passphrase can be brute-forced offline (AES-GCM decrypt attempts are cheap once the file is in hand). | The passphrase is user-chosen, and the whole point of the feature is that the file can travel (USB, cloud, email to yourself). | Require a minimum passphrase strength at export time; derive the key with a slow KDF (argon2id) so each offline guess is expensive; and make the session revocable (R13) so a compromise is contained even if the file is cracked. |
| **R15** | **2FA introduces a lockout surface.** A lost authenticator app with no usable backup codes is a locked-out user — and because 2FA gates *every* human login, support has to have a recovery path that isn't itself the thing 2FA protects. | TOTP recovery is the classic "recovery that defeats the security" problem: whoever can reset 2FA effectively bypasses it. | Backup codes (single-use, hashed, regenerable — MV-E4) are the primary recovery. The last-resort reset path is **platform-admin only** (a documented `admin` action, audited `auth.totp_reset`, requiring the admin's own authentication) — never self-serve, never email-based, so a compromised email can't bypass 2FA. Secret shown exactly once + confirm-before-enable (MV-E5) prevents lockout-by-misconfiguration. |

## Links

- Architecture: [architecture.md](architecture.md)
- Data model: [ERD.md](ERD.md)
- Mobile plan: [mobile/README.md](mobile/README.md)
- Backend plan: [backend/README.md](backend/README.md)
- Web plan: [web/README.md](web/README.md)
- CLI plan: [cli/README.md](cli/README.md)
- Desktop plan: [desktop/README.md](desktop/README.md)
- Infra plan: [infra/README.md](infra/README.md)
