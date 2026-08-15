# Task group — 05 packaging: Windows

**Depends on:** `02-window-shell.md`, `web/tasks/01-setup.md`. **Blocks:** `08-release-ci`.

Per `desktop/README.md` §Packaging — Wails v3 Taskfile-driven. **Windows signing is deferred for v1, explicitly** (SmartScreen warning is bad but not blocking; a cert is $200–400/yr + business verification). Recorded so users aren't surprised.

## Config & assets

- [ ] **`build/windows/installer/project.nsi`** — the NSIS script: app name, install dir, Start-menu shortcut, **the `bloberry://` registry key** (URL scheme, not file associations — Bloberry owns no file type, README §Linux note applies here too), version from the tag.
- [ ] **`info.json`** (Wails metadata), **`icon.ico`**, **`wails.exe.manifest`** present and filled.

## The artifact

- [ ] **Raw `.exe`** — `wails3 task windows:package` produces the bare executable.
- [ ] **NSIS installer `.exe`** — the package step also produces the setup `Bloberry-<ver>.exe`.

## Prereq tooling — the silent failure

- [ ] **NSIS installed on the build machine** — **without it, the packaging step silently produces only the bare `.exe`** (a failure that looks like success, README's explicit warning). Verified by the installer artifact actually existing.
- [ ] **WebView2 SDK** on the build machine (per `desktop/README.md` §CI).

## Signing — deferred, but made visible

- [ ] **Deferred item recorded** — an explicit task noting the installer is unsigned for v1, the SmartScreen consequence ("users must click through a warning"), and that when signing is added, `signtool` must sign **both** the inner `.exe` and the NSIS installer (signing only the inner one still shows the warning on the installer). Not dropped — visible at release time.

## Verification — installed, not just built

- [ ] **Clean-machine install** — a Windows VM (or clean CI runner) runs the NSIS setup; the app installs to Program Files, appears in the Start menu with its icon.
- [ ] **The app launches** and the WebView2 runtime renders the web UI.
- [ ] **`bloberry://` scheme** opens the app at the linked object (the registry key worked).
- [ ] **Uninstall** removes the app cleanly.

**tests:** the installer artifact EXISTS (NSIS was present); fresh-install launches; Start-menu icon present; URL scheme opens the app; uninstall is clean.
