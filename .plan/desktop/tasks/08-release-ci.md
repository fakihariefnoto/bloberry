# Task group — 08 release CI

**Depends on:** `04`/`05`/`06` (the packaging files), `07-auto-update.md`. **Blocks:** nothing (final).

Per `desktop/README.md` §CI and `templates/release-tooling.md` — **no desktop framework cross-compiles**, so this is the exception to the self-hosted-VPS-runner default. Three runners, two of them billable.

## The build matrix

- [ ] **macOS `.dmg`** — `macos-14` (GitHub-hosted, arm64): Xcode CLT, `create-dmg`, signing cert + notarization credentials.
- [ ] **Windows NSIS** — `windows-latest` (GitHub-hosted): NSIS, WebView2 SDK.
- [ ] **Linux `.deb`/`.rpm`** — the **self-hosted VPS runner** (the only non-billable one): `libwebkit2gtk-4.1-dev`, `libgtk-3-dev`, `nfpm`.
- [ ] **Per-OS build jobs, one publish job** — the GoReleaser-with-CGO pattern (per `release-tooling.md`): each OS job builds + uploads its artifacts; one publish job assembles the release, checksums and manifests. GoReleaser's single-run model can't build CGO targets from one machine.

## The release tool

- [ ] **GoReleaser config** — `.goreleaser.yaml` covers the desktop artifacts as extra builds alongside the CLI's (same file, `cli/tasks/13`). The publish half: checksums, release, channel manifests.
- [ ] **`goreleaser --snapshot` dry run before the first tag** — finds config errors against a real release (which means deleting a public tag).

## Secrets (GitHub Environment)

- [ ] **macOS**: `APPLE_DEVELOPER_ID`, `APPLE_ID`, `APPLE_APP_PASSWORD`, `APPLE_TEAM_ID` in a GitHub Environment, referenced only by the macOS job.
- [ ] **Releases**: `GITHUB_TOKEN`.
- [ ] **Windows signing secrets absent** by the deferral — an explicit recorded absence, not an accident.

## Version + artifacts

- [ ] **Version stamped from the git tag at build time** (never a source constant).
- [ ] **Artifacts + checksums attached to a GitHub Release**, including the `latest.json` manifest for `07-auto-update`.
- [ ] **Triggered by tag or `workflow_dispatch`** — the manual-trigger convention (cross-ref `infra/tasks/` if it exists; if Infra is out of scope, these tasks live here — a desktop app that deploys nothing still needs CI to package).

**tests:** a snapshot dry run produces all three OS artifacts; the macOS job signs + notarizes; the publish job assembles a Release with checksums + `latest.json`; version on each artifact matches the tag.
