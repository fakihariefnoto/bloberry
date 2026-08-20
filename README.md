<div align="center">

<img src="bloberry.png" alt="Bloberry" width="640" />

# bloberry

**One self-hosted object service. Six storage backends. Every device you own.**

<br/>

Bloberry gives your apps a single S3-style API, real folders, per-folder permissions, application access keys, and temporary signed links — while the actual bytes live on **S3, Cloudflare R2, Alibaba OSS, Google Cloud Storage, Azure Blob Storage, or a plain VPS disk**. Swap vendors without touching a line of app code.

<br/>

Go · Vue 3 · Wails · Flutter · OpenAPI · MongoDB · Redis

</div>

---

## Why bloberry?

Object storage is powerful but annoying. Every provider has a different console, different SDK, different CORS setup, and no notion of "folders" or "who can see what." Bloberry wraps them all behind one contract:

- **Storage-agnostic by design** — your app talks to bloberry, bloberry talks to the vendor. Register one bucket or six; presigned uploads/downloads keep big files off your server.
- **Real folders + per-folder grants** — object storage has no directories. Bloberry gives you a real folder tree with materialized paths and RBAC down to `principal × folder subtree × permission × expiry`.
- **Application access keys** — `blob_live_…` bearer keys, hashed at rest, shown exactly once, revocable in seconds. Built for machines, scoped to a tenant.
- **Temporary signed links** — share anything for 1 hour or 1 day, revocable before expiry, with a friendly expiry page instead of a provider error.
- **Self-hosted, yours forever** — no lock-in, no per-GB SaaS bill. Run it on a $5 VPS and point it at whatever storage you already pay for.

## Feature highlights

| Area | What you get |
| --- | --- |
| **Storage** | 6 drivers — S3, Cloudflare R2, Alibaba OSS, Google Cloud Storage, Azure Blob, local disk. Capability-aware (presign, multipart, ranges), health-checked. |
| **Uploads** | Direct, browser presigned-PUT, and multipart/resumable for large files. Explicit replace / keep-both / cancel on name conflicts. |
| **Tenancy** | Tenant-isolated at the repository layer — no query can cross a tenant boundary. Per-tenant byte + object quotas with graceful read-only when over. |
| **Permissions** | 5 human roles (`platform_admin` → `viewer`) plus folder-level grants. Allow-only, most-specific-wins. Permission-aware UI, not error walls. |
| **API keys** | Prefixed, hashed at rest, one-time show, optional folder scoping, expiry, last-used tracking, instant revocation. |
| **Sharing** | Temporary signed links with TTL + revocation; public visibility with stable URLs; short slugs. |
| **Auth** | Email/password, refresh tokens, forgot-password, email-OTP, Google OAuth. **TOTP 2FA**, QR device pairing for mobile, and passphrase-encrypted config login for desktop. |
| **Ops** | Usage metering + estimated cost per backend, audit log, storage health, Redis-backed job queue for subtree deletes. Credentials envelope-encrypted, never returned by the API. |
| **Surfaces** | Web dashboard (Vue 3), **Go CLI**, **Wails desktop app**, **Flutter mobile app**, and SDK clients generated from one OpenAPI spec. |

## Surfaces

```
┌─────────────┐   ┌─────────────┐   ┌──────────────┐
│  Web (Vue3) │   │ CLI (cobra) │   │ Desktop/Wails│
└──────┬──────┘   └──────┬──────┘   └──────┬───────┘
       │                 │                 │
       └───────────────┐ │ ┌───────────────┘
                       ▼ ▼ ▼
               ┌──────────────┐
               │ bloberry API │   OpenAPI · one contract
               │  (Go server) │
               └──────┬───────┘
      ┌───────────────┼──────────────────┐
      ▼       ▼       ▼       ▼       ▼   ▼
     S3     R2       OSS     GCS    AZBlob  Disk
      └──────┴───────┴───────┴───────┴──────┴─────┘
                 (or MinIO, any S3-compatible)
```

## Quickstart

### 1. Server

```bash
# Requirements: Go 1.25+, MongoDB, Redis
cp .env.example .env            # set MONGODB_URI, REDIS_URI, JWT_SECRET,
                                # CREDENTIAL_ENCRYPTION_KEY, DISK_SIGNING_SECRET
make build
./bin/bloberry-server           # listens on $SERVER_ADDR (default 127.0.0.1:8080)
```

Open `http://localhost:8080`, complete the one-time setup wizard (platform admin + first storage backend), and you're live. The web dashboard is embedded in the binary — no separate frontend to serve.

### 2. CLI

```bash
./bin/bloberry auth login         # device flow / token
./bin/bloberry cp  photo.jpg /albums/
./bin/bloberry ls  /albums/ --json
./bin/bloberry share link photo.jpg --ttl 3600
./bin/bloberry key create
```

### 3. Point an app at it

Use the HTTP API directly (85 endpoints) — or the generated SDK clients for Go, TypeScript, and Dart. Everything is auth'd with an application access key:

```bash
# Presign a browser upload — the browser PUTs straight to the storage vendor,
# big files never touch your server.
curl -X POST http://localhost:9095/v1/objects/presign-put \
  -H "Authorization: Bearer blob_live_..." \
  -H "Content-Type: application/json" \
  -d '{"folder_id":"root","name":"photo.jpg","content_type":"image/jpeg","size":123456}'

# → { "file_id": "...", "upload_url": "https://<vendor>...X-Amz-Signature=..." }
```

## Deploy

Everything you need to self-host is in `deploy/`:

- `bloberry.service` — systemd unit
- `Caddyfile` — automatic HTTPS, reverse proxy, WebSocket-ready
- `.env.example` — full production environment reference

GitHub Actions ship CI (spec-drift gate, lint, gosec, authz 100% branch coverage, MinIO conformance), CLI releases, and desktop builds.

## Development

```bash
make build       # openapi → web → go build (spec-first; generated code must be committed)
make generate    # regenerate API server + clients from api/openapi.yaml
make test        # unit + authz 100% branch gate
make lint        # golangci-lint
make security    # gosec
go test ./internal/storage/... -run Conformance -tags conformance  # drivers vs MinIO
```

> The OpenAPI spec is the source of truth — `make generate` fails CI if the tree drifts. Never hand-edit generated code.

## Repository layout

```
api/openapi.yaml        Spec-first contract (server + SDK clients generated from it)
cmd/                    bloberry-server · bloberry (CLI) · bloberry-desktop (Wails)
internal/               api, auth, authz, apikey, audit, folder, grant, job,
                        object, share, storage (6 drivers), tenant, usage, user…
web/                    Vue 3 dashboard (embedded into the server binary)
mobile/                 Flutter app (QR pairing login, offline upload queue)
deploy/                 systemd, Caddyfile, production env template
.github/workflows/      CI, release-cli, release-desktop, deploy
```

## License

[GPL-3.0](LICENSE)
