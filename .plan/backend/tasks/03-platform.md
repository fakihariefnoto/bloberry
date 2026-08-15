# Task group — 03 platform packages

**Depends on:** `01-setup.md`. **Blocks:** every domain file (they all consume these). The non-domain infrastructure in `internal/platform/` and `internal/authz/`/`internal/storage/` scaffolding (the resolver and drivers get their own files — `16`, `17`).

- [ ] **`platform/config`** — env-driven config struct (`caarlos0/env` or stdlib, per TRD); typed `JWT_*` TTLs, `CREDENTIAL_ENCRYPTION_KEY`, rate-limit defaults, `SERVER_ADDR`, SMTP/Google/OAuth. Decided in `detail-backend`; wire what it named.
- [ ] **`platform/db`** — Mongo client constructor, `migrate.go` (runs `02` migrations), `indexes.go` (idempotent index ensure).
- [ ] **`platform/redis`** — client constructor + the key helpers the rest of the plan assumes (`refresh:`, `otp:`, `reset:`, `principal:`, `apikey:`, `slug:`, `job:queue`).
- [ ] **`platform/httpx/envelope.go`** — the standard response envelope `{data?, messages?: [{code, content}]}`, both omitempty, **no `error` boolean**; HTTP status is the signal.
- [ ] **`platform/httpx/errors.go`** — the 17 error codes from `backend/domains.md` §8 (API, never invented per handler); `forbidden` carries a `content` naming what was needed.
- [ ] **`platform/httpx/middleware.go`** — the auth middleware: discriminates JWT (human) vs access-key (application) by prefix, resolves to the unified `Principal` (`domains.md` §4.7), nothing downstream branches on scheme. Includes the **per-access-key rate limiter** (Redis token bucket, 429 + `Retry-After`, PRD-Q5).
- [ ] **`platform/httpx/stream.go`** — the streaming helpers: `io.Copy` to drivers (never `io.ReadAll`), HTTP range-request support, hard body-size ceiling.
- [ ] **`platform/crypto`** — envelope encryption for `storage_backends.credentials_encrypted` (TRD R7), key from env, decrypt-in-memory only. The rotation seam (`--rotate-credential-key`) is `infra`/CLI; this is the primitive it rotates.
- [ ] **`platform/web/embed.go`** — `//go:embed all:../../../web/dist` + the SPA catch-all (deep links like `/files/abc123` don't 404) + serving `web/dist` assets with correct cache headers.

**tests:** envelope shape (data only / messages only / both / empty 204); error-code mapping; auth-middleware scheme discrimination (JWT vs key → same Principal); rate-limiter 429 + `Retry-After`; stream ceiling rejects oversized bodies.
