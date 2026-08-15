# Task group — 06 packaging: Linux

**Depends on:** `02-window-shell.md`, `web/tasks/01-setup.md`. **Blocks:** `07-auto-update` (the Linux update command), `08-release-ci`.

Per `desktop/README.md` §Packaging — **Wails builds no Linux installer at all**: `wails3 task linux:build` produces the binary, then `nfpm` assembles `.deb`/`.rpm`. All `.desktop`/icon/MIME integration is hand-authored.

## Config & assets

- [ ] **`nfpm.yaml`** — describes the binary, dependencies, and the `.desktop`/icon/MIME files to install.
- [ ] **`bloberry.desktop`** → `/usr/share/applications/` — `Name=Bloberry`, `Exec=/usr/bin/bloberry-desktop %U`, `Icon=bloberry`, `Categories=Utility;Network;FileTransfer;`, `StartupWMClass=bloberry-desktop`, `MimeType=x-scheme-handler/bloberry;`.
- [ ] **Icons** → `/usr/share/icons/hicolor/{16,32,48,128,256,512}x{...}/apps/bloberry.png` (every size).
- [ ] **MIME XML for the URL scheme** — `x-scheme-handler/bloberry` declared (Bloberry registers its `bloberry://` scheme, not file associations — it owns no file type).

## The artifact

- [ ] **`wails3 task linux:build`** — the Linux binary (from the self-hosted VPS runner, `libwebkit2gtk-4.1-dev` + `libgtk-3-dev` installed).
- [ ] **`.deb`** — `nfpm pkg --packager deb`.
- [ ] **`.rpm`** — `nfpm pkg --packager rpm`.

## Prereq tooling

- [ ] **`nfpm` installed** on the Linux build machine (the self-hosted VPS runner).
- [ ] **GTK/WebKit2GTK dev libs** present to build.

## Signing

- [ ] **None expected** — per the README; optional `dpkg-sig` on the `.deb` noted but not required in v1.

## Verification — installed, not just built (the launcher check)

- [ ] **Clean-machine install** — `dpkg -i` (and `rpm -i` on a Fedora-family box); the app installs and **appears in the app launcher with its icon** (the `.desktop` entry worked — the check that catches a package that installs successfully and then shows up nowhere).
- [ ] **`bloberry://` scheme** opens the app (the MIME XML registered `x-scheme-handler/bloberry`).
- [ ] **`man bloberry-desktop`** renders if a man page is shipped.
- [ ] **Update path** — `apt install --only-upgrade bloberry-desktop` works (the command the auto-update banner prints on Linux).

**tests:** `dpkg -i` on a clean box → app in the launcher with the icon; URL scheme opens it; `apt install --only-upgrade` upgrades.
