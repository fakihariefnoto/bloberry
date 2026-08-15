# Task group — 01 setup

**Depends on:** `backend/tasks/01-setup.md` (the Go module and repo already exist). **Blocks:** everything in this folder.

**Context from `architecture.md` §7:** layout B (embedded web). `web/` is a **package inside the existing repo**, not a repo and not a deployable. `bun create vite web --template vue-ts` scaffolds it in place; nothing here runs `git init` or creates a second project.

- [ ] **Depends-on note verified.** `backend/tasks/01-setup.md` has run: repo root with `go.mod` exists, `web/` does not yet. If the backend setup hasn't run, stop here.
- [ ] **Scaffold the package.** `bun create vite web --template vue-ts` from the repo root. This creates `web/` with Vite + Vue 3 + TypeScript. It is a package inside the existing repo — no `.git`, no separate init.
- [ ] **Bun confirmed as the package manager and build runner** (`TRD.md` Stack decisions). `bun install` works; `bun run dev/build/preview` run the Vite scripts. No npm/pnpm/yarn anywhere in `web/` — the `Makefile` and scripts use `bun …`.
- [ ] **Router installed.** `bun add vue-router@4`. Used by `03-routing.md`; no default routes yet.
- [ ] **State installed.** `bun add pinia`. Used by `05-core-infra.md` (`stores/`).
- [ ] **UI kit installed.** `bun add reka-ui` — the unstyled headless primitives (dialog, dropdown, popover, tabs, tooltip, toast, table). No themed kit, no vendored component layer beyond our own `src/components/ui/`.
- [ ] **Styling installed.** `bun add -D tailwindcss @tailwindcss/vite` and wire the Vite plugin per Tailwind v4 conventions. `02-design-tokens.md` fills the theme in.
- [ ] **Icons installed.** `bun add lucide-vue-next` — the only icon set (`design/style-guide.md` §Sizing).
- [ ] **Table installed.** `bun add @tanstack/vue-table` (+ any peer deps it declares). Used by the `DataTable` component and the ten table pages.
- [ ] **API client installed.** `bun add openapi-fetch` (or the generated-client runtime the OpenAPI toolchain needs) — consumed by `lib/api/` in `05-core-infra.md`. The *client code itself* is generated, not this dependency list.
- [ ] **Vite config set.** `vite.config.ts` with `outDir: dist` (matches `internal/platform/web/embed.go`'s `//go:embed all:../../../web/dist`), and a dev-server proxy for `http://localhost:8080` so the dev environment works against the local Go binary (`04-environments.md`).
- [ ] **TypeScript strict.** `tsconfig.json` at strict, `verbatimModuleSyntax` on — a data-dense dashboard with 26 pages doesn't survive `any`-drift.
- [ ] **`Makefile` wired to Bun.** `web/Makefile` uses `bun run dev/build/test/lint` (already converted from the template — verify each target actually runs under Bun).
- [ ] **Root `make build` ordering verified** (`architecture.md` §7 edge 2): the root Makefile's `build` target runs `openapi codegen → web build → go build`, **in that order**. Confirm the web build is a dependency of the backend build, not a parallel job.
- [ ] **Basic smoke build.** `bun run build` produces `web/dist/` without errors, and the Go binary serves it (via `internal/platform/web/embed.go`). Empty shell page at this point is fine.

**tests:** the smoke build (`bun run build` → `web/dist` exists) and `go build ./cmd/bloberry-server` succeeding with the embedded dist.
