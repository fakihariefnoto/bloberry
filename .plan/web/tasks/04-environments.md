# Task group — 04 environments

**Depends on:** `01-setup.md`. **Blocks:** `05-core-infra.md` (reads `VITE_API_BASE_URL`) and every data-dependent page.

**Source:** `web/components.md` §Environments and `templates/web-defaults.md`. Embedded topology means these are **build-time** config only — nothing secret ever lives in them (they compile into shipped assets).

- [ ] **`.env.development`** — `VITE_API_BASE_URL=http://localhost:8080/v1` (Vite dev server proxies to the local Go binary).
- [ ] **`.env.staging`** — `VITE_API_BASE_URL=/v1` (same-origin — the binary serves both).
- [ ] **`.env.production`** — `VITE_API_BASE_URL=/v1` (same-origin).
- [ ] **`.env.example`** committed with the keys and placeholder values; the three real env files are **not** committed (gitignore).
- [ ] **Build/run scripts wired per target** — `bun run dev` uses dev, a staging build uses staging, prod build uses prod. No manual env-file swapping; the scripts pick the file per target (Vite's `--mode`).
- [ ] **Cross-check against the backend** (backend is in scope): the staging/prod API base URLs agree with `backend/README.md`/`infra` staging/prod config — same meaning of "staging" and "prod" on both sides.
- [ ] **No secrets rule stated in the file** — a comment noting these compile into `web/dist` and are served to every visitor, so credentials/keys must never appear (matches `web/components.md` §Environments).

**tests:** a build per mode (`bun run build --mode staging`) producing `web/dist` with the right base baked in; a grep that no real secret pattern appears in `.env.*` (committed ones).
