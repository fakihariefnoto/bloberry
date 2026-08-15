# Backend — Bloberry

**Go + chi + MongoDB + Redis.** Modular monorepo — one module to start, organized so a piece (the extraction worker, the download proxy) can split into its own service later without a rewrite.

This is the centre of Bloberry: it serves the JSON API, the embedded Vue dashboard, the public-object and short-URL surfaces, and it owns the storage-driver abstraction that the whole product exists to provide.

Follows [`templates/backend-go-defaults.md`](../../../templates/backend-go-defaults.md) — layered handler/usecase/repository with interface injection, the standard response envelope `{data?, messages?}`, OpenAPI-as-contract, JWT + Redis-backed refresh tokens with platform-aware expiry, the standard auth/user domains, mock-friendly interface naming, and `golangci-lint`/`gosec` guards. Use `detail-backend` to flesh out domains, the data model and the OpenAPI contract against that reference.

## Embedded web layer

The web layer is **embedded** (see [`../TRD.md`](../TRD.md) Stack decisions): `web/` builds to static assets with Bun/Vite, and this backend serves them via `go:embed` from the same chi router. One binary, one deployable. There is **no separate frontend server**, and no separate Go project for the desktop app either — Wails wraps this same Go code and the same Vue build. `../architecture.md` §7 records which platform's setup task owns the single `go mod init`; read it before scaffolding anything.

## What's specific to this backend

- **Storage drivers** — one `StorageDriver` interface, **six** implementations (S3, Cloudflare R2, Alibaba OSS, Google Cloud Storage, Azure Blob Storage, local VPS disk). Design the interface against the *hardest* driver, not the easiest: local disk has no external presigner, and GCS signs with a service account rather than a key pair. One conformance test suite that every driver must pass. See `TRD.md` risks R1/R2.
- **The permission resolver** — "can principal P do action A on folder F" spans membership, role, folder-grant, folder ancestry and access-key scope, runs on every request, and is the highest-risk function in the codebase. One pure function, table-driven tests, Redis-cached with explicit invalidation. See `TRD.md` risk R4.
- **Two auth schemes, one middleware** — human sessions (JWT + opaque Redis refresh token) and application access keys both resolve to a principal the same RBAC middleware evaluates.
- **Streaming, never buffering** — `io.Copy` to the driver, never `io.ReadAll`. Presigned-PUT is the preferred browser upload path so bytes bypass this process entirely.
- **Queued jobs** — archive extraction and bundle-download generation are Redis-queued, never inline. Extraction is also the product's main DoS surface (zip bombs, path traversal, symlinks) — see `TRD.md` risk R6.
- **Two-phase writes** — metadata to Mongo, bytes to a remote store, with no transaction across them. Write metadata `pending`, promote to `active` after the byte write confirms, reconcile the rest on a sweep. See `TRD.md` risk R8.

## Files

- `domains.md` — the domain structure, once `detail-backend` has run
- `../ERD.md` — the single authoritative data model (don't keep a second copy here)
- `Makefile` — build/run/debug/test/lint/security commands, from `templates/makefiles/go.mk`. Targets assume a `cmd/bloberry` entrypoint; adjust once the real layout exists per `../architecture.md` §7.
- `tasks/` — the implementation task list, once `build-backend` has run
