# Task group — 01 project setup

**Depends on:** `backend/tasks/01-setup.md` (the module exists — **no init**). **Blocks:** everything in this folder.

- [ ] **Entrypoint added, not created** — `cmd/bloberry/main.go` in the existing Go module (a second entrypoint alongside `cmd/bloberry-server/`). No `go mod init`.
- [ ] **`internal/cli/` tree** per `architecture.md` §7 — `commands/` (one package per command group), `output/` (the shared renderer from `02`), plus the root command wiring.
- [ ] **cobra root command** — the `bloberry` root with the global flags from `cli/README.md` §Global flags: `--help/-h`, `--version`, `--config <path>`, `--json`, `--quiet/-q`, `--verbose/-v`, `--no-color`, `--yes/-y`, `--tenant <slug>`.
- [ ] **Version stamped at build time** — `-ldflags "-X main.version=… -X main.commit=… -X main.date=…"` per the CLI Makefile (already in the template); the `version` command reads them. **Never hardcoded** — a hardcoded version lies after the next release, exactly when someone files a bug against it.
- [ ] **`make build`/`run`/`test` working** — the CLI Makefile's targets produce `bin/bloberry` and run its tests from `cmd/bloberry`.
- [ ] **`--help` renders** the root help with the command tree; `bloberry` with no args shows help + a usage hint, exits 2.
- [ ] **Consumes `sdk/go`** — the CLI imports the generated Go SDK (from `backend/tasks/19-openapi.md`), not a hand-rolled client.

**tests:** `go build ./cmd/bloberry` clean; `bin/bloberry --version` reports the injected version; `bin/bloberry` with no args exits 2 with help.
