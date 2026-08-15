# Task group — 03 chrome mockups (reuse, don't redraw)

**Depends on:** `02-window-shell.md`. **Blocks:** the packaging files.

Per `desktop/README.md` §Mockups: desktop screens ARE `web/mockup/`'s desktop-width wireframes, loaded unchanged. This file covers only the chrome deltas — the five items the README lists.

- [ ] **Title bar / menu bar above the web layout** — native on macOS, in-window on Windows/Linux. The web bundle's own top area is not doubled; the native chrome replaces/augments where the README specifies.
- [ ] **First-run server-URL screen** — the one surface with no web equivalent (the browser knows its own origin); already built as part of `02`, this task confirms it renders correctly above the shell chrome per the wireframe.
- [ ] **Sync status in the sidebar footer** — a web `web/mockup/` page has no such element; the Go host injects it (a sidebar footer entry showing active sync folders + one-way state). Per `desktop/README.md` — sync is **one-way and additive**, the UI states it plainly (PRD NG4).
- [ ] **Native dialogs replace the web file picker** in the Upload flow — confirmed against `web/mockup/files.md`'s upload interaction, now routed through Wails dialogs.
- [ ] **Tray** (no web counterpart) — the tray surface from `02`, with its aggregate-progress status line.

**tests:** every screen still renders pixel-identical to its web desktop-width mockup except the listed chrome deltas; the sync footer and tray show the same live state (one-way, no local-delete propagation).
