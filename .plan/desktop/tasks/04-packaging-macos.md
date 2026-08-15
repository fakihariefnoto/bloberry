# Task group — 04 packaging: macOS

**Depends on:** `02-window-shell.md`, `web/tasks/01-setup.md` (the dist). **Blocks:** `07-auto-update`, `08-release-ci`.

Per `desktop/README.md` §Packaging — Wails v3 Taskfile-driven; pin the version and take exact commands from its Taskfiles. **macOS signing is required** (unsigned apps are blocked by default).

## Config & assets

- [ ] **`build/darwin/Info.plist`** — app metadata, bundle id, `CFBundleURLTypes` for `bloberry://` (so a share link opens the app at that object), version from the git tag. **`CFBundleIconFile` → `icon.icns`.**
- [ ] **`icon.icns`** generated at every required size (16/32/64/128/256/512 + @2x).

## The artifact

- [ ] **Universal `.app`** — `wails3 task darwin:package` building `darwin/universal` (Apple Silicon + Intel ship as one, per the README). Output `bin/Bloberry.app`.
- [ ] **`.dmg`** via `create-dmg` (Wails has no dmg target — the README is explicit). Output `Bloberry-<ver>-universal.dmg`.

## Prereq tooling

- [ ] **`create-dmg` installed** on the macOS build machine (`brew install create-dmg`).

## Signing + notarization (required)

- [ ] **`codesign --deep --options runtime`** the `.app` — Developer ID cert (`APPLE_DEVELOPER_ID` secret).
- [ ] **Notarize** — `xcrun notarytool submit --wait` with `APPLE_ID`/`APPLE_APP_PASSWORD`/`APPLE_TEAM_ID`.
- [ ] **Staple** — `xcrun stapler staple` on the `.dmg`.

## Verification — installed, not just built

- [ ] **Clean-machine install** — a macOS VM (or clean CI runner) mounts the `.dmg`, drags to Applications, launches; the app appears in Launchpad with its icon.
- [ ] **`spctl -a -vv` passes** on the installed app (Gatekeeper accepts it unsigned-flagged).
- [ ] **`bloberry://` scheme opens the app** at the linked object.
- [ ] **First launch** shows the first-run server-URL screen (not a Gatekeeper "damaged" error).

**tests:** `spctl -a -vv` clean; a fresh install launches past Gatekeeper; the dmg is double-clickable to a working app.
