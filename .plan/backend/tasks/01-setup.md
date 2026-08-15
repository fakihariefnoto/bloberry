# Task group — 01 setup

**Depends on:** nothing (this is the scaffold owner). **Blocks:** `web/tasks/01-setup.md`, every file in this folder, and the cli/desktop/mobile/infra setup files.

**Context from `architecture.md` §7:** layout B (embedded web) + F (CLI companion) + I (desktop network client) composed into **one repo, one Go module**. This file runs the one `go mod init`. Everything else is downstream.

- [ ] **`go mod init github.com/<org>/bloberry`** at the repo root — the single init. Nothing else in the plan runs an init.
- [ ] **Repo tree created** per §7: `cmd/bloberry-server/`, `internal/{auth,user,tenant,folder,object,share,apikey,grant,job,usage,audit,admin,authz,storage,platform}/`, `api/`, `migrations/`, `sdk/go/`, `web/` (created by web's setup), `.plan/` (this folder, committed).
- [ ] **`cmd/bloberry-server/main.go`** — wiring only: config → db → repos → usecases → handlers → router. No business logic in main.
- [ ] **Dev MongoDB provisioned** — local instance (docker-compose or native) per `deploy/`; the `run-all` Makefile target starts infra then the app.
- [ ] **Dev Redis provisioned** — same story; sessions, authz cache, job queue all need it.
- [ ] **Migration tool initialized** — `golang-migrate` (per the Go defaults) wired to Mongo. First migration file is the baseline in `02-migrations.md`.
- [ ] **`.env.example` committed** — `MONGODB_URI`, `REDIS_URI`, `JWT_*` TTLs (web 48h/144h, mobile 720h/2160h), `CREDENTIAL_ENCRYPTION_KEY`, `SMTP_*`, `GOOGLE_OAUTH_*`, `SERVER_ADDR`, rate-limit defaults. Real `.env` never committed.
- [ ] **Makefile corrected for spec-first** — **replace the template's `swagger` (swaggo) target** with `generate` (`oapi-codegen` from `api/openapi.yaml`), per `backend/domains.md` §1.1's deliberate deviation. Keep `lint`/`security`/`mocks`/`test-coverage` as-is.
- [ ] **`make build` runs `openapi → web build → go build` in order** (`architecture.md` §7 edge 2) — the binary embeds `web/dist`. If `web/tasks/01-setup.md` hasn't produced a dist yet, `make build` fails with a clear message, not a stale embed.
- [ ] **`make generate` green** — `oapi-codegen` produces server interfaces + `sdk/go` from the (initially minimal) `api/openapi.yaml`. CI fails if the working tree changes (keeps spec-first honest).
- [ ] **`make tools`** installs `dlv`, `golangci-lint`, `gosec`, `mockgen`, `oapi-codegen` (drop `swag` from the list).
- [ ] **Smoke run** — `go run ./cmd/bloberry-server` boots, connects to Mongo+Redis, serves a health endpoint.

**tests:** `go build ./...` clean; `make generate` leaves the tree unchanged; the server boots against fresh Mongo/Redis.
