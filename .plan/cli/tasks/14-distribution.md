# Task group — 14 distribution

**Depends on:** `13-release-pipeline` (artifacts), `11-completions-man`. Each channel **ends at *verified installed on a clean machine*** — never *artifact produced*.

## 14.1 GitHub Releases

- [ ] Artifacts + `checksums.txt` attached to the Release (`GITHUB_TOKEN`).
- [ ] **Verified:** download + `tar xzf` on a clean machine, binary on `$PATH`, `--version` reports the real version, completions work in a fresh shell.

## 14.2 Homebrew tap

- [ ] **`homebrew-tap` repo created and public** — `<org>/homebrew-tap` with a `bloberry.rb` formula. The formula's `test do` block passes (`version` or `--help`).
- [ ] **Formula bumped by automation** — GoReleaser `brews:` config generates the formula on release (a SHA updated by hand is the step that gets skipped; a stale tap installs an old binary while looking healthy).
- [ ] **A token with write access to the tap repo exists** — a repo-scoped `GITHUB_TOKEN` can't push there; this is the most common first-release failure (`release-tooling.md`).
- [ ] **Completions installed by the formula** — the formula drops bash/zsh/fish completion files in place (`11`).
- [ ] **Verified:** `brew install <org>/tap/bloberry` on a clean machine → `bloberry version` → completions work in a fresh zsh.

## 14.3 Scoop bucket

- [ ] **`scoop-bucket` repo created and public** — `<org>/scoop-bucket` with a `bloberry.json` manifest (URL, hash from `checksums.txt`, persist config path).
- [ ] **Manifest bumped by automation** — GoReleaser `scoops:` config; same token as the tap.
- [ ] **Verified:** `scoop bucket add <org> …` then `scoop install bloberry` on a clean Windows machine → `bloberry version`.

## 14.4 `.deb` / `.rpm`

- [ ] **Built by GoReleaser `nfpms:`** — `.deb`/`.rpm` attached to Releases, with `completions/` and `man/` at the right paths and a post-install that refreshes the shell completion cache.
- [ ] **Verified:** `dpkg -i bloberry_*.deb` / `rpm -i` on clean Linux → binary on `$PATH`, `man bloberry` renders, completions load.

## 14.5 `go install`

- [ ] The module is public; `go install github.com/<org>/bloberry/cmd/bloberry@latest` resolves (no automation — free with a public module).
- [ ] **Verified:** the install command succeeds on a machine without the binary and `bloberry --version` runs.

## 14.6 Install docs

- [ ] **README install section** — one block per channel with the exact command (`cli/README.md` §Distribution is the plan; the README's install section is how users find out the tap exists). Written and linked from the README.

**tests:** each channel's verified-installed step is the test — run the channel's install command on a clean machine/container, confirm version + completions. A channel that was published but never installed from is the standard way this stage looks done and isn't.
