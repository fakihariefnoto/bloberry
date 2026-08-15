# Task group — 07 desktop CI (the OS-native exception)

**Depends on:** `01-vps-runner.md`, `desktop/tasks/08-release-ci.md` (the packaging files). **Blocks:** the first desktop release.

Per `infra/README.md` §Workflows → `release-desktop.yml` and `desktop/README.md` §CI. **No desktop framework cross-compiles** — three runners, two billable. Infra owns the runners + secrets; `desktop/tasks/04–06` own the artifacts.

- [ ] **`release-desktop.yml` generated** — `workflow_dispatch` (manual-trigger-only keeps the billable minutes bounded — a real benefit of the no-automatic-triggers convention), three jobs:
  - **macOS `.dmg`** — `macos-14` (GitHub-hosted, arm64), **10× minute multiplier**: Xcode CLT, `create-dmg`, the signing + notarization secrets.
  - **Windows NSIS** — `windows-latest` (GitHub-hosted), **2× multiplier**: NSIS, WebView2 SDK. No signing cert secret yet (deferred — the absence is recorded, not accidental).
  - **Linux `.deb`/`.rpm`** — **the self-hosted VPS runner** (free), targeted via its `bloberry` label: `libwebkit2gtk-4.1-dev`, `libgtk-3-dev`, `nfpm`.
- [ ] **macOS signing secrets set** in a GitHub Environment: `APPLE_DEVELOPER_ID`, `APPLE_ID`, `APPLE_APP_PASSWORD`, `APPLE_TEAM_ID` (referenced only by the macOS job).
- [ ] **Per-OS installer toolchains installed on each runner** — `create-dmg` on macOS, NSIS on Windows, `nfpm` on the VPS runner (NSIS missing silently produces only a bare `.exe` — the toolchain task is what prevents the failure that looks like success).
- [ ] **Billable-minutes note recorded** — this is the one place the project pays GitHub for CI; the README says so and the manual trigger bounds it.

**verification:** a manual dispatch produces all three OS artifacts with the right toolchains; the macOS job signs + notarizes; the Linux job runs on the self-hosted runner (label match); each artifact then passes its `desktop/tasks/04–06` clean-install check.
