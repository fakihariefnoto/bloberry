# Task group — 19 openapi

**Depends on:** `01-setup` (make `generate` target). **Blocks:** everything that consumes the contract — `sdk/go`, `sdk/ts`, `mobile/lib/api/`, the CLI, and the web's `lib/api/`. **Spec-first** (the deliberate deviation in `domains.md` §1.1, ADR-11): `api/openapi.yaml` is hand-authored and is the source; `oapi-codegen` generates server interfaces + the Go SDK; `openapi-generator` generates TS + Dart clients.

- [ ] **`api/openapi.yaml` authored** — the full contract, per-surface groups from `TRD.md` §APIs: `auth`, `users`, `tenants`, `folders`, `objects`, `shares`, `applications`, `grants`, `archives`, `jobs`, `audit`, `usage`, `admin`. Versioned from day one (`/v1/`).
- [ ] **The response envelope is the contract** — every response shape wraps `{data?, messages?: [{code, content}]}` with the right `omitempty` semantics; error responses carry the documented codes from `domains.md` §8, not ad-hoc ones.
- [ ] **Server interfaces generated** — `oapi-codegen` emits the handler interfaces the `*domain/handler` packages implement; `make generate` regenerates, and **CI fails if the working tree changes** (the check that keeps spec-first honest).
- [ ] **`sdk/go` generated** — the Go SDK, tagged separately from the server (same module, its own version tag per TRD).
- [ ] **`sdk/ts` generated** — the TypeScript SDK for npm, plus the hand-written upload helper (multipart, per `TRD.md` §Tech stack).
- [ ] **`mobile/lib/api/` generated** — the Dart client (internal to mobile in v1, PRD D8), snake_case on the wire.
- [ ] **Breaking-change gate** — CI fails on a breaking spec diff before any client sees it (the reason spec-first was chosen over code-first, `domains.md` §1.1).
- [ ] **Contract consumed by the web client** — web's `lib/api/` (`web/tasks/05-core-infra.md`) regenerates against this same file; the two sides can't drift.

**tests:** `make generate` is idempotent (tree unchanged after a second run); a spec change that breaks a client fails CI; the envelope shape holds across every documented response.
