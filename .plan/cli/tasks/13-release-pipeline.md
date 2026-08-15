# Task group — 13 release pipeline

**Depends on:** `12-testing` (green), `11-completions-man` (the archives carry them). **Blocks:** `14-distribution`.

Per `templates/release-tooling.md` — Go CLI → **GoReleaser**, pure Go so one Linux runner cross-compiles every platform (`architecture.md` §7 / `cli/README.md` §Distribution: unlike desktop, the CLI needs no OS-native CI).

- [ ] **`.goreleaser.yaml` authored** — the cross-compile matrix from `cli/README.md` §Distribution: `darwin_arm64`, `darwin_amd64`, `linux_amd64`, `linux_arm64`, `windows_amd64`.
- [ ] **Archives contain the right contents** — binary, LICENSE, README, `completions/`, `man/` (a GoReleaser archive: `files:` list).
- [ ] **`checksums.txt` generated** — SHA-256 alongside the artifacts (every other channel's manifest embeds one of these hashes — a hard dependency).
- [ ] **Version stamping verified** — run `bin/bloberry version` on a built artifact and confirm the injected version/commit/date match the tag (`-ldflags`, not a source constant).
- [ ] **`goreleaser check` green** — the config validates without building.
- [ ] **`goreleaser --snapshot` dry run** — produces all artifacts locally *before the first tag* (`make snapshot`). Finding a config error against a real release means deleting a public tag.
- [ ] **Trigger convention** — release triggered by a tag or `workflow_dispatch`, consistent with this repo's manual-trigger default (per `infra-defaults.md`; the workflow itself is `build-infra`'s job — this file defines what it runs).
- [ ] **Update story wired** — `version --check` knows the install method and prints the right update command (`09`), and there is deliberately **no `self-update`**.

**tests:** a snapshot artifact of each OS/arch builds; `--version` on each reports the injected values; `checksums.txt` hashes match the archives.
