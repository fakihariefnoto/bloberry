# Build index — CLI

Build-order table, dependency edges, status, command count, and distribution channels.

## Build order

| File | Covers | Status |
|---|---|---|
| `01-project-setup.md` | Entrypoint in the existing module, cobra root, version stamping | ☐ |
| `02-core-infra.md` | Config precedence, output layer, error/exit-code mapping, shared API client | ☐ |
| `03-auth.md` | `auth login`/`logout`/`status` — device flow, keychain, server-side revoke | ☐ |
| `04-commands-file-verbs.md` | `cp`, `ls`, `rm`, `mv`, `cat`, `stat`, `sync` (7) | ☐ |
| `05-commands-folder-share.md` | `folder create/tree`, `share link/short/public/list/revoke` (7) | ☐ |
| `06-commands-key-app-grant.md` | `key create/list/revoke`, `app create/list/delete`, `grant create/list/revoke` (9) | ☐ |
| `07-commands-tenant-job-archive.md` | `tenant list/use/usage`, `job list/status/watch`, `archive extract/bundle` (8) | ☐ |
| `08-commands-admin.md` | `admin backend {add,list,test,rate-card}`, `admin tenant {create,list,quota}`, `admin usage` (8) | ☐ |
| `09-commands-config-misc.md` | `init`, `config get/set/path`, `completion`, `version` (6) | ☐ |
| `10-flows.md` | First-run, CI/non-interactive, and the publish/rotate journeys | ☐ |
| `11-completions-man.md` | bash/zsh/fish/powershell + install by packages; man page for .deb/.rpm | ☐ |
| `12-testing.md` | Golden-file tests per command (stdout/stderr/exit code) | ☐ |
| `13-release-pipeline.md` | GoReleaser config, cross-compile matrix, archives/checksums, snapshot dry run | ☐ |
| `14-distribution.md` | 5 channels, each ending at *verified installed on a clean machine* | ☐ |

## Command count

**48 commands**, each with a designed file in `cli/commands/` — the smallest unit that can be built and verified against one artifact. `init` is the first-run entry.

## Dependency edges

- **`01-project-setup` blocks everything** — the binary must exist before any command is implemented.
- **`02-core-infra` blocks every command file** — config precedence, the output layer and error→exit-code mapping are shared; retrofitting them after commands exist means editing all of them.
- **`03-auth` blocks the commands that need a session** (everything except `init`, `config`, `completion`, `version`) — no command works until login/BLOBERRY_TOKEN is real.
- **`10-flows` depends on the commands it chains.**
- **`13-release-pipeline` depends on the commands** (the build must compile them all).
- **`14-distribution` depends on `13`** (channels ship release artifacts).

## External edges (from `architecture.md` §7)

- **Companion CLI — no init.** **`01-project-setup` Depends on `backend/tasks/01-setup.md`**: the Go module already exists (`github.com/<org>/bloberry`). This file adds `cmd/bloberry/main.go` + `internal/cli/` to the existing module. Nothing here runs `go mod init`.
- **Shared client, not a second one** — the CLI consumes `sdk/go` (generated from `api/openapi.yaml` by `backend/tasks/19-openapi.md`), parsing the standard envelope. Hand-rolling a second HTTP client is the failure this check prevents.
- **Test fixtures shared** — `testdata/golden/` holds the golden outputs, sourced from `cli/commands/*.md` sample blocks (`architecture.md` §7).

## Distribution channels in scope

GitHub Releases · Homebrew tap · Scoop bucket · `.deb`/`.rpm` · `go install`. Each ends at *verified installed* per the skill — never just *artifact produced*.

## Gaps flagged

None. All 47 commands have designed files with real help text, per-state output, and exit codes.
