# Task group — 20 guards

**Depends on:** `01-setup` (targets exist), the domains so they can be linted. Cross-ref `infra/tasks/` for actually wiring these into `ci.yml` (per `detail-infra`, CI checks are manually-triggered workflow_dispatch).

- [ ] **`make lint` in CI** — `golangci-lint` with `.golangci.yml` pinning at minimum: `errcheck`, `gosec`, `gocritic`, `revive`, `govet`, `staticcheck`, `unused`, `ineffassign` (`domains.md` §10).
- [ ] **`make security` in CI** — `gosec` standalone as its own gate. Must stay green on the Bloberry-specific items (`domains.md` §10): pinned JWT algorithm, no hardcoded credentials, no unhandled `io.Copy` errors on the streaming paths, no path traversal in the disk driver or the extraction worker.
- [ ] **`make mocks` wired** — `mockgen` auto-discovering `handler.go`/`usecase.go`/`repository.go`; the uniform interface naming (§2) is what makes one invocation cover every domain. Don't rename an interface for clarity.
- [ ] **`make generate` in CI** — `oapi-codegen` from `api/openapi.yaml`; **CI fails if the working tree changes** (spec-first enforcement, `19-openapi`).
- [ ] **Authz coverage gate** — `internal/authz/` requires **100% branch coverage** (PRD G5, `16-authz`), enforced in CI. Everything else has no hard gate; one function does, because one function is the security boundary.
- [ ] **Conformance suite CI** — the storage conformance suite runs against MinIO on every push; the **real-provider** suite (S3/R2/OSS/GCS/Azure Blob) runs on a schedule (costs money + real credentials, `17-storage`).
- [ ] **Cross-ref `infra/tasks/`** — the `ci.yml` that runs these is written by `build-infra`; this file defines *what* runs, infra owns *where*.

**tests:** this file IS the CI test plan — each checkbox is a workflow step whose green is a gate.
