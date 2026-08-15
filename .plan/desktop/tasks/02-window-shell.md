# Task group — 02 window shell

**Depends on:** `01-framework-setup.md`. **Blocks:** `03-chrome-mockups`, the packaging files. Window/menu/shortcuts/tray per `desktop/README.md` §Window, menu, shortcuts — the exact recorded list, not "standard menu."

- [ ] **Window** — 1280×800 default, 960×600 minimum; position/size/maximized state persisted per user (Go-side config, the Wails way per `desktop-defaults.md`). Single window.
- [ ] **Menu bar** — macOS gets the standard app menu; Windows/Linux an in-window menu bar. Items exactly per the README table:
  - **Bloberry** (macOS): About · Check for Updates… · Preferences ⌘, · Quit ⌘Q
  - **File**: Upload Files… ⌘O · Upload Folder… ⇧⌘O · New Folder ⌘N · — · Close Window ⌘W
  - **Edit**: Undo/Redo/Cut/Copy/Paste/Select All (needed for form fields to behave natively)
  - **View**: Back ⌘[ · Forward ⌘] · Reload ⌘R · Toggle Sidebar ⌘\ · Zoom In/Out/Reset
  - **Sync**: Add Sync Folder… · Pause All Syncing · Open Sync Folder · Sync Status
  - **Window**: Minimize ⌘M · Zoom · Bring All to Front
  - **Help**: Documentation · Report an Issue · View Logs
- [ ] **Global shortcuts** — ⌘/Ctrl+F focus search · ⌘/Ctrl+U toggle the upload queue panel · ⌘/Ctrl+, preferences · Delete on a selected row · Escape closes a dialog. Wired to the menu actions above.
- [ ] **System tray: yes** (the recorded decision — uploads/sync continue after the window closes, so the tray is how users tell it's still working). Tray menu: aggregate status line ("Uploading 3 of 12 · 45 MB/s") · Open Bloberry · Pause All Syncing · Quit.
- [ ] **Close ≠ quit** — closing the window minimizes to tray (not quits) on macOS/Windows while uploads or syncs are active, with a **first-time notification** explaining that; quitting with work in flight prompts for confirmation.
- [ ] **Native notifications** — use the OS notification system (Wails `Notification` / a Go package) for anything that should reach the user while the window isn't focused (upload-complete, sync-paused, tray-level events), not in-app toasts alone.
- [ ] **The first-run server-URL screen** — the wireframe in `desktop/README.md`: empty (Continue disabled) · validating (spinner in the button, field disabled) · unreachable (the warning) · reachable-but-not-Bloberry (distinct message) · success → device-flow login. Stored to the config file, feeding the `__BLOBERRY_API_BASE__` hook.
- [ ] **First-run "import a login file" option** (PRD M23, `desktop/README.md` §Sub-choices) — a secondary action under the URL field: file picker → passphrase prompt → decrypt client-side (argon2id → AES-GCM) → validate the server signature + 24h import window → store the session in the OS keychain. Failure copy distinguishes **wrong passphrase** from **window passed** (DT-E1). The imported session is an ordinary revocable refresh token (DT-E2).
- [ ] **Native dialogs replace the web file picker** in the Upload flow (the browser picker can't choose a save *destination*; Wails runtime dialogs).

**tests:** window opens at the persisted size/position; every menu item triggers its action; global shortcuts fire; tray shows live aggregate progress; close-minimizes-to-tray with the first-time notification; first-run screen handles all 5 states; the sync caption states "one-way and additive — a local delete never propagates" (PRD NG4).
