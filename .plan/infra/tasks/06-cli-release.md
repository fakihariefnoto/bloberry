# Task group — 06 CLI release workflow

**Depends on:** `01-vps-runner.md`, `cli/tasks/13-release-pipeline.md` (the GoReleaser config). **Blocks:** the first CLI release.

Per `infra/README.md` §Workflows → `release-cli.yml`. **The CLI cross-compiles all five platforms from one Linux runner** — the OS-native desktop exception does NOT apply here (a pure-Go CLI needs no per-OS runners). Infra owns the runner and credentials; `cli/tasks/14-distribution.md` owns the channel publishing.

- [ ] **`release-cli.yml` generated** — `workflow_dispatch`, GitHub-hosted Linux, GoReleaser driven from `cli/tasks/13`'s `.goreleaser.yaml`.
- [ ] **The cross-platform PAT secret set** — a token with **write access to the `homebrew-tap` and `scoop-bucket` repos**, stored as a secret. A repo-scoped `GITHUB_TOKEN` can't push a formula to another repo — the most common first-release failure (`templates/release-tooling.md`).
- [ ] **`GITHUB_TOKEN`** for the GitHub Release itself (attached to this repo's Releases).
- [ ] **Channel-publishing wiring confirmed** — the workflow invokes GoReleaser's `brews:`/`scoops:`/`nfpms:` blocks that `cli/tasks/14` verifies end at *installed on a clean machine*; this file's job is the runner + token that make them run.
- [ ] **Dry run before the first real tag** — `goreleaser --snapshot` in CI (or `make snapshot`) validates the config against a fake release before the first tag.

**verification:** a manual dispatch produces all five platform artifacts + checksums; the tap/bucket manifests are pushed (proving the PAT works); a snapshot dry run passes before the first tag.
