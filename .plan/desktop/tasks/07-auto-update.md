# Task group — 07 auto-update

**Depends on:** `04-packaging-macos.md` (signing), `06-packaging-linux.md` (the update command). **Blocks:** `08-release-ci`.

Per `desktop/README.md` §Auto-update: **Wails ships no updater, so this is a decision — manual check, not silent auto-update.** It spans the shell (check/prompt), packaging (signed bundles), and CI (publishing the manifest), which is why it's its own file.

- [ ] **Manual check** — Help → Check for Updates… and a startup check (once per 24h) query the version manifest at `https://bloberry.example.com/desktop/latest.json`.
- [ ] **Non-blocking banner** — a newer version shows a banner with release notes and a **Download** link that opens the release page. Never an interruptive dialog.
- [ ] **No in-place binary replacement** — same reasoning as the CLI: overwriting a package-managed binary breaks the package database, and on macOS it invalidates the signature. **On Linux the banner prints `apt install --only-upgrade bloberry-desktop` instead of a download link** (the install-method-aware rule, matching the CLI's `version --check`).
- [ ] **The manifest** — `latest.json` with version, release notes, artifact URLs per OS, checksums. Published alongside the release artifacts by `08-release-ci`.
- [ ] **The gap is a recorded decision** — `.deb`-only Linux releases can't self-update; that's why the banner prints the apt command. Silent auto-update is explicitly a post-v1 candidate needing a signed-manifest mechanism.

**tests:** a real previous version's banner appears when `latest.json` is bumped; the Download link opens the release page; Linux renders the apt command; a signed-macOS app's update doesn't break its signature.
