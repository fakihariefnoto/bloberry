# Task group — 01 framework setup

**Depends on:** `backend/tasks/01-setup.md` (the module), `web/tasks/01-setup.md` (the dist). **Blocks:** everything in this folder.

**Context from `architecture.md` §7 / `desktop/README.md`:** layout I — desktop is a *network client*, not the project. **Do NOT run `wails3 init` in this repo.** Wails is a third binary inside the existing Go module.

- [ ] **`wails3 init` run in a scratch directory OUTSIDE the repo** (only if a reference is wanted) — the reference `main.go` and `build/` are ported across by hand. Never in the repo.
- [ ] **`cmd/bloberry-desktop/main.go` added** to the existing Go module — a third entrypoint alongside `bloberry-server` and `bloberry`. Wails dependency added to the existing `go.mod`.
- [ ] **Frontend source wired to `web/dist`** — the same build artifact the server embeds (`desktop/README.md` §Sub-choices). Not a copy, not a variant. The window loads the built assets; `web/tasks/01-setup.md`'s build must run first.
- [ ] **IPC layer set up** — Wails bindings for the native-only surface only (file dialogs, drag-drop, sync control, tray, keychain). **Everything else goes over HTTPS to the API exactly as the browser does** — no second data path to keep in sync (§Sub-choices IPC).
- [ ] **API base URL runtime-configurable** — the web bundle's `window.__BLOBERRY_API_BASE__` hook (`web/components.md` §Environments) is set by the Go host from the config file, defaulting to the first-run screen's answer. This is the desktop-specific delta to the same-origin web build.
- [ ] **Auth = browser device flow** — identical to the CLI (`cli/README.md`), tokens in the OS keychain via `zalando/go-keyring`. No separate login form — one implementation, SSO for free.
- [ ] **`internal/desktop/` packages** — `menu`, `tray`, `sync`, `dialogs` (per §7 tree), the only packages importing Wails APIs.
- [ ] **`Makefile`** from `templates/makefiles/desktop-wails.mk`, adjusted for this layout (v3's Taskfile is wrapped by the Makefile).
- [ ] **Smoke run** — `go run ./cmd/bloberry-desktop` opens a window loading `web/dist` with the shell chrome, no API calls needed yet.

**tests:** the window opens loading the real `web/dist`; a first-run server URL can be entered and stored to the config file.
