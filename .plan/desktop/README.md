# Desktop — Bloberry

**Wails v3.** Wraps [`../web/`](../web/README.md)'s build output — the same Vue bundle the browser gets — with a Go host process around it.

Follows [`templates/desktop-defaults.md`](../../../templates/desktop-defaults.md) and [`templates/desktop-packaging.md`](../../../templates/desktop-packaging.md). Release tooling per [`templates/release-tooling.md`](../../../templates/release-tooling.md).

---

## Not a separate project

`../architecture.md` §7 records this and it's the single most important thing on this page:

> **Do not run `wails3 init` in this repo.** The Go module already exists — `backend/tasks/01-setup.md` created it. Wails is added as a **third binary inside that module** (`cmd/bloberry-desktop/`), with the Wails dependency added to the existing `go.mod`.

The standard scaffold-owner rule in `templates/repo-layouts.md` says *Wails → `desktop/` owns the init*, because in the usual shape (`layout D`) the desktop app *is* the whole project. **Bloberry is layout I instead**: the server is a genuine separate deployable, and the desktop app is a network client of it (`architecture.md` ADR-9). Running `wails3 init` here would scaffold a competing project on top of an existing module.

If a Wails template is wanted as a reference, run `wails3 init` in a **scratch directory outside the repo** and port `main.go` and `build/` across by hand.

---

## What desktop adds

The entire justification for shipping a shell. If a capability isn't on this list, it belongs in `web/`, not here.

| Capability | Why the browser can't | Implementation |
|---|---|---|
| **Drag a folder in** | Browsers support directory upload poorly and inconsistently; structure is often lost | Native drag-drop event → walk the tree in Go → enqueue with paths preserved |
| **Native file dialogs** | The browser picker can't choose a save *destination* | Wails runtime dialogs |
| **One-way folder sync** | No filesystem watching in a browser | `fsnotify` in the Go host (`architecture.md` §3.6) |
| **Background upload queue** | Dies with the tab | Survives window close; continues from the tray |
| **Very large uploads** | Subject to tab lifecycle and memory limits | Streamed from disk in Go |
| **OS keychain credentials** | No browser equivalent | `zalando/go-keyring` — the same mechanism the CLI uses |

**Sync is one-way and additive.** A local delete never propagates (PRD NG4). The UI must say so plainly rather than implying mirror semantics — this is the difference between "convenience" and an implied backup guarantee the product doesn't make.

---

## Sub-choices

| Choice | Decision |
|---|---|
| **Frontend source** | `web/dist` — the *same* build artifact the server embeds. Not a copy, not a variant. Build order in `architecture.md` §7. |
| **API base URL** | **Runtime-configurable**, unlike the web build where it's same-origin. A first-run screen asks for the server URL; stored in the config file. `web/components.md` notes the `window.__BLOBERRY_API_BASE__` hook this depends on. |
| **IPC** | Wails bindings for the native-only surface — file dialogs, drag-drop, sync control, tray, keychain. Everything else goes over HTTPS to the API exactly as the browser does, so there is no second data path to keep in sync. |
| **Auth** | Two paths, same keychain destination: **browser device flow** (identical to the CLI, `../cli/README.md`) for interactive first-run, and **import an encrypted config file** (PRD M23) — a `.bloberry` file downloaded from the web dashboard, passphrase-decrypted locally, holding a refresh token. Both put the session in the OS keychain; both stay revocable via `auth logout`. |
| **Auto-updater** | See below — Wails ships none. |

---

## Window, menu, shortcuts

**Window:** 1280×800 default, 960×600 minimum. Position, size and maximized state persisted per user. Single window — no multi-window, since the app is one dashboard.

**Menu bar** (macOS gets the standard app menu; Windows/Linux get an in-window menu bar):

| Menu | Items |
|---|---|
| **Bloberry** *(macOS)* | About · Check for Updates… · Preferences ⌘, · Quit ⌘Q |
| **File** | Upload Files… ⌘O · Upload Folder… ⇧⌘O · New Folder ⌘N · — · Close Window ⌘W |
| **Edit** | Undo/Redo/Cut/Copy/Paste/Select All — standard, needed for form fields to behave natively |
| **View** | Back ⌘[ · Forward ⌘] · Reload ⌘R · Toggle Sidebar ⌘\ · Zoom In/Out/Reset |
| **Sync** | Add Sync Folder… · Pause All Syncing · Open Sync Folder · Sync Status |
| **Window** | Minimize ⌘M · Zoom · Bring All to Front |
| **Help** | Documentation · Report an Issue · View Logs |

**Global shortcuts:** ⌘/Ctrl+F focus search · ⌘/Ctrl+U toggle the upload queue panel · ⌘/Ctrl+, preferences · Delete on a selected row · Escape closes a dialog.

**System tray: yes** — and it earns its place rather than being reflexive. Uploads and folder sync continue after the window is closed, so the tray is the only way to see that they're running, show aggregate progress, pause syncing, and quit properly. An app with background work and no tray gives users no way to tell it's still doing something.

Tray menu: aggregate status line ("Uploading 3 of 12 · 45 MB/s") · Open Bloberry · Pause All Syncing · Quit.

**Closing the window does not quit** on macOS and Windows when uploads or syncs are active; it minimizes to tray with a first-time notification explaining that. Quitting with work in flight prompts for confirmation.

---

## Mockups — reuse, don't redraw

**No `mockup/` folder here.** Desktop screens *are* `../web/mockup/`'s desktop-width wireframes, loaded unchanged. Only the chrome differs, and those deltas are:

1. **Title bar / menu bar** above the web layout (native on macOS, in-window elsewhere).
2. **First-run server-URL screen** — the one surface with no web equivalent, since the browser knows its own origin. Wireframed in this file rather than as a mockup, below.
3. **Sync status** — a sidebar footer entry and a tray menu, neither of which exists on web.
4. **Native dialogs** replace the web file picker in the Upload flow.
5. **Tray**, with no web counterpart.

Everything else — the file browser, shares, applications, members, audit, usage, admin — is pixel-identical to web.

### First-run server URL

```
┌────────────────────────────────────────────────┐
│  ● ● ●                Bloberry                 │
├────────────────────────────────────────────────┤
│                                                │
│                  ◆ Bloberry                    │
│                                                │
│         Connect to your Bloberry server        │
│                                                │
│   Server URL                                   │
│   ┌──────────────────────────────────────┐     │
│   │ https://bloberry.example.com         │     │
│   └──────────────────────────────────────┘     │
│   The address of your organisation's install.  │
│                                                │
│            ┌──────────────────────┐            │
│            │       Continue       │            │
│            └──────────────────────┘            │
│                                                │
│   ⚠ Couldn't reach that server. Check the      │
│     address and that you're online.            │
│                                                │
└────────────────────────────────────────────────┘
```

States: empty (Continue disabled) · validating (spinner in the button, field disabled) · unreachable (the warning above) · reachable-but-not-Bloberry (distinct message — pointing at the wrong service is a real mistake) · success → device-flow login.

**Second option on first run: import a login file.** A secondary action under the URL field — "Or import a login file" → file picker → passphrase prompt → decrypt + validate signature + import window → store session in the keychain. Failure copy distinguishes the two cases (PRD DT-E1): **wrong passphrase** ("the passphrase didn't decrypt this file") vs **window passed** ("this file expired — download a fresh one from the web dashboard"). The imported session is an ordinary refresh token — revocable via `auth logout` on any device (DT-E2).

---

## Packaging

Per-OS, per `templates/desktop-packaging.md` → Wails. **Wails v3 is Taskfile-driven** — commands below are the v3 targets; pin the Wails version and take exact commands from that release's own Taskfiles, since v3's layout differs from v2's.

| OS | Artifact | Command | Output | Config / assets | Prerequisite |
|---|---|---|---|---|---|
| **macOS** | `.app` natively; **`.dmg`** via wrapper | `wails3 task darwin:package` then `create-dmg` | `bin/Bloberry.app` → `Bloberry-<ver>-universal.dmg` | `build/darwin/Info.plist`, `icon.icns` | `create-dmg` (brew). **Wails has no dmg target.** Build universal (`darwin/universal`) so Apple Silicon and Intel ship as one. |
| **Windows** | **NSIS setup `.exe`** + raw `.exe` | `wails3 task windows:package` | `bin/` | `build/windows/installer/project.nsi`, `info.json`, `icon.ico`, `wails.exe.manifest` | **NSIS must be installed on the build machine** — without it the packaging step silently produces only the bare `.exe`, which is a failure that looks like success |
| **Linux** | **`.deb` + `.rpm`**; AppImage opt-in | `wails3 task linux:build`, then `nfpm pkg --packager deb` / `--packager rpm` | `bin/` + nfpm output | `nfpm.yaml` listing the binary, `.desktop`, icon and MIME XML | `nfpm`. **Wails builds no Linux installer at all.** GTK/WebKit2GTK dev libs to build. |

### Linux desktop integration

Entirely hand-authored — Wails provides nothing here.

- **`bloberry.desktop`** → `/usr/share/applications/`: `Name=Bloberry`, `Exec=/usr/bin/bloberry-desktop %U`, `Icon=bloberry`, `Categories=Utility;Network;FileTransfer;`, `StartupWMClass=bloberry-desktop`, `MimeType=x-scheme-handler/bloberry;`
- **Icons** → `/usr/share/icons/hicolor/{16,32,48,128,256,512}x{...}/apps/bloberry.png`
- **URL scheme, not file associations.** Bloberry doesn't own a file type — it stores everyone else's. What it *does* want is the **`bloberry://` scheme** registered, so a share link can open the desktop app at that object. MIME XML declares `x-scheme-handler/bloberry`; macOS uses `CFBundleURLTypes` in `Info.plist`; Windows a registry key from the NSIS script.

### Signing

| OS | Plan |
|---|---|
| **macOS** | Required, not optional — unsigned apps are blocked by default on modern macOS. `codesign --deep --options runtime`, then `xcrun notarytool submit --wait`, then `xcrun stapler staple` on the `.dmg`. Needs a paid Apple Developer account ($99/yr) and an app-specific password in CI. |
| **Windows** | **Deferred for v1, explicitly.** An unsigned installer triggers a SmartScreen warning that users must click through — bad, but not blocking. A code-signing certificate is ~$200–400/yr and OV certs need business verification. Documented so users aren't surprised. When done: `signtool` on **both** the inner `.exe` and the NSIS installer — signing only the inner one still shows the warning on the installer. |
| **Linux** | No signing expectation. Optionally sign the `.deb` with `dpkg-sig`. |

### Auto-update

Wails ships none, so this is a decision rather than a default: **manual check, not silent auto-update.**

- Help → Check for Updates… and a startup check (once per 24h) query a version manifest at `https://bloberry.example.com/desktop/latest.json`.
- A newer version shows a non-blocking banner with release notes and a **Download** link that opens the release page.
- **No in-place binary replacement.** Same reasoning as the CLI (`../cli/README.md`): overwriting a package-managed binary breaks the package database, and on macOS it invalidates the signature. On Linux the banner prints `apt install --only-upgrade bloberry-desktop` instead of a download link.

Silent auto-update is a post-v1 candidate and would need a signed-manifest mechanism to be safe.

---

## CI — the exception to the self-hosted-runner default

**No desktop framework cross-compiles**, because each needs a native webview. This is the constraint that most often stalls a desktop release, so it's stated in three places (here, `../infra/README.md`, `../TRD.md` R10):

| Target | Runner | Needs installed |
|---|---|---|
| macOS `.dmg` | `macos-14` (GitHub-hosted, arm64) | Xcode CLT, `create-dmg`, signing cert + notarization credentials |
| Windows NSIS | `windows-latest` (GitHub-hosted) | NSIS, WebView2 SDK |
| Linux `.deb`/`.rpm` | **the self-hosted VPS runner** | `libwebkit2gtk-4.1-dev`, `libgtk-3-dev`, `nfpm` |

Only Linux uses the self-hosted runner. macOS and Windows use GitHub-hosted runners — which means **desktop CI minutes are billable** where the rest of this project's CI isn't. Worth knowing before it appears on an invoice.

**Release tool: GoReleaser**, the same one the CLI uses (`templates/release-tooling.md`), driven from one `.goreleaser.yaml` with the desktop artifacts as extra builds. Convenient, but note the asymmetry: the CLI's five platforms come off one Linux runner; the desktop's three need three.

Secrets: `APPLE_DEVELOPER_ID`, `APPLE_ID`, `APPLE_APP_PASSWORD`, `APPLE_TEAM_ID` (macOS); `GITHUB_TOKEN` (releases). Windows signing secrets are absent by the deferral above.

---

## Files

- `Makefile` — build/run/debug/test targets, from `templates/makefiles/desktop-wails.mk`
- `tasks/` — shell tasks plus **one packaging file per OS**, once `build-desktop` has run
