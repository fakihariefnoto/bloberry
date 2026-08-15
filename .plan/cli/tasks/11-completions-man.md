# Task group — 11 completions + man page

**Depends on:** `01-project-setup` (cobra), the command groups. **Blocks:** `14-distribution` (packages carry these).

- [ ] **Completion generation** — bash/zsh/fish/powershell via `cobra completion`, per `completion` command design (`09`).
- [ ] **Dynamic completions implemented** — remote paths complete folders via the API (`bloberry cp bloberry://<TAB>`), key/app/job IDs complete from the API (`key revoke <TAB>`). Static-only completion would miss the part that's actually hard to type (`cli/README.md` §Completions).
- [ ] **`make completions`** produces the files into `completions/` (the Makefile target exists) — bash → `bloberry.bash`, zsh → `_bloberry`, fish → `bloberry.fish`, powershell → `bloberry.ps1`.
- [ ] **Installed by the packages** — the Homebrew formula, `.deb` and `.rpm` drop the completion file in the right place (`cli/README.md` §Completions: a completion that ships in the tarball but isn't installed is one nobody discovers). Wired in `14-distribution`.
- [ ] **Man page** — `cobra doc` → `man/bloberry.1`, shipped in `.deb`/`.rpm` (those users expect `man bloberry`).

**tests:** each completion script sources without error in its shell; the zsh completion is on `$fpath` after package install; `man bloberry` renders.
