# CLI — Bloberry

**Go.** Binary name `bloberry`. Role: **companion** to the Backend — a scriptable path to the same operations the API exposes, sharing its auth and its committed OpenAPI contract rather than hand-rolling a second client.

Follows [`templates/cli-defaults.md`](../../../templates/cli-defaults.md) and [`templates/cli-distribution.md`](../../../templates/cli-distribution.md). Lives in the same Go module as the server (`../architecture.md` §7) — `cmd/bloberry/` plus `internal/cli/`, consuming `sdk/go`. **No separate init.**

---

## Two users, one binary

The design has to serve both, and most of the conventions below exist because of the second:

- **A developer at a terminal** — moving files, issuing a key, checking a folder. Wants readable output, colour, progress, and a confirmation before anything destructive.
- **A CI job** — publishing build artifacts, syncing a directory, rotating a key. Wants `--json`, no colour, no prompts, an access key from an env var, and exit codes to branch on.

So: every prompt has a flag equivalent, every decoration is TTY-gated, and authentication has two shapes.

---

## Command tree

**Shape: short file verbs at the top level, `<noun> <verb>` for everything else.** This is a deliberate, documented exception to the one-shape rule in `cli-defaults.md`. File operations are the overwhelming majority of use and every comparable tool (`aws s3 cp`, `gsutil cp`, `rclone copy`) puts them within one word — `bloberry object copy ./dist …` in every CI script is friction users would alias away. Everything that isn't a file verb stays strictly noun-verb.

```
bloberry
├── init                       one-command first-run setup (server → tenant → auth)
├── cp <src> <dst>            copy up, down, or remote→remote  (-r, --exclude, --dry-run)
├── ls [path]                 list a folder                    (-l, --recursive)
├── rm <path>                 delete objects or a folder       (-r, --yes)
├── mv <src> <dst>            move / rename
├── cat <path>                stream an object to stdout
├── stat <path>               object or folder metadata
├── sync <local> <remote>     one-way mirror                   (--delete, --dry-run)
│
├── auth
│   ├── login                 browser device flow → OS keychain
│   ├── logout                revokes server-side, not just local
│   └── status                who am I, which tenant, token expiry
│
├── folder
│   ├── create <path>
│   └── tree [path]           the hierarchy as a tree
│
├── share
│   ├── link <path>           signed URL              (--ttl, --json)
│   ├── short <path>          short URL
│   ├── public <path>         make public             (--yes)
│   ├── list
│   └── revoke <id>
│
├── key
│   ├── create                scoped access key       (--folder, --permission, --expires)
│   ├── list
│   └── revoke <id>           (--yes)
│
├── app
│   ├── create <name>
│   ├── list
│   └── delete <id>           (--yes)
│
├── grant
│   ├── create                (--folder, --principal, --permission, --expires)
│   ├── list
│   └── revoke <id>
│
├── tenant
│   ├── list
│   ├── use <slug>            set the default tenant
│   └── usage                 bytes, objects, estimated cost
│
├── job
│   ├── list
│   ├── status <id>
│   └── watch <id>            block until terminal state — the CI-friendly one
│
├── archive
│   ├── extract <path>        server-side extraction into a folder
│   └── bundle <paths...>     N objects → one archive
│
├── admin                     platform_admin only
│   ├── backend {add,list,test,rate-card}
│   ├── tenant {create,list,quota}
│   └── usage
│
├── config {get,set,path}
├── completion <shell>
└── version
```

**48 commands**, each designed in its own file under [`commands/`](commands/) — purpose, flags, real `--help` text, sample output per state, exit codes. (The tree's `admin backend {add,list,test,rate-card}`, `admin tenant {create,list,quota}` and `config {get,set,path}` each expand to one file per verb.) `init` is the first-run entry — a fresh machine goes from install to a working `ls` in one command.

### Path syntax

Remote paths use a **`bloberry://` scheme**; local paths are bare. This is what lets `cp` copy in either direction from one argument order, the same way `aws s3 cp` does:

```bash
bloberry cp ./dist bloberry://assets/v2 -r        # up
bloberry cp bloberry://assets/v2 ./dist -r        # down
bloberry cp bloberry://a/x.zip bloberry://b/x.zip # remote → remote (server-side)
```

The tenant comes from context (`bloberry tenant use`, or the config file) and is overridable with `--tenant`. A fully-qualified form `bloberry://<tenant>/<path>` is accepted for scripts that shouldn't depend on ambient state — and CI scripts should use it.

---

## Global flags

On every command, per the defaults:

`--help/-h` · `--version` · `--config <path>` · `--json` · `--quiet/-q` · `--verbose/-v` · `--no-color` · `--yes/-y` · `--tenant <slug>`

`--tenant` is Bloberry's one addition and it earns its place: every operation is tenant-scoped, and a script that relies on whichever tenant happened to be current is a script that will one day write to the wrong one.

---

## Output conventions

- **stdout is data, stderr is everything else.** Progress bars, spinners, warnings, confirmations and errors all go to stderr, so `bloberry ls --json | jq` works. This is the convention most often broken by a stray status line.
- **`--json` emits only the payload** — no decoration, no progress. Every command that returns structured data supports it.
- **TTY gating**: colour, progress bars and prompts only when stdout is a terminal. `NO_COLOR` and `TERM=dumb` honoured. A piped `cp -r` prints one line per file, not 400 redraws of a progress bar.
- **Errors say what to do next.** `no folder 'assets/v3' (run 'bloberry ls bloberry://assets' to see what exists)` beats a 404. Panics and stack traces only under `--verbose`.
- **`--version` prints version, commit and build date**, injected via `-ldflags -X` at build time — never hardcoded, which silently lies after the next release.

---

## Exit codes

Scripts branch on these, so they are **API**. Stable once shipped.

| Code | Meaning |
|---|---|
| `0` | Success |
| `1` | The operation failed (generic) |
| `2` | Invocation error — bad flag, missing argument, malformed path |
| `3` | Not authenticated, or the token expired — run `bloberry auth login` |
| `4` | Authenticated but **forbidden** — the principal lacks the permission |
| `5` | Not found — object, folder, key or tenant doesn't exist |
| `6` | **Quota exceeded** — the write was rejected |
| `7` | **Partial failure** — some items succeeded, some didn't (`cp -r`, `sync`, `rm -r`) |
| `8` | Conflict — name collision, or a folder cycle |
| `9` | Storage backend unreachable |

**Code 7 is the one that matters most.** A recursive copy where 3 of 400 files fail must not exit `0` (CI would call it green) and must not exit `1` (which reads as "nothing happened"). It exits `7`, prints a per-file summary to stderr, and **leaves the successful files uploaded** so a rerun is idempotent rather than a full redo (PRD CL-E1).

Codes 3, 4, 5 and 6 are separated deliberately: `bloberry auth login` fixes 3, nothing the caller can do fixes 4, and 6 needs someone to raise a quota. A single code `1` for all four makes automated remediation impossible.

---

## Config, secrets, auth

**Precedence, highest first:** command flag → environment variable (`BLOBERRY_`) → config file → built-in default.

**Config file** at `~/.config/bloberry/config.yaml` (XDG) / `%APPDATA%\bloberry\config.yaml`, overridable with `--config`. Managed through `bloberry config get|set|path` rather than hand-editing.

```yaml
server: https://bloberry.example.com
tenant: acme
output: table          # table | json
color: auto
```

**Secrets never live in that file.** Tokens go in the **OS keychain** (macOS Keychain, libsecret, Windows Credential Manager) via `zalando/go-keyring`. In CI, where there is no keychain, `BLOBERRY_TOKEN` holds an access key instead — and the CLI tells you which one it used in `auth status`.

### Two auth shapes, both needed

| Shape | For | Mechanism |
|---|---|---|
| **Browser device flow** | A human at a terminal | `bloberry auth login` prints a code, opens the browser, polls until confirmed. Tokens land in the keychain. Works with Google login and any future SSO, because the browser does the auth. |
| **Access key via env** | CI | `BLOBERRY_TOKEN=blob_live_…`. No prompt, no keychain, no browser. The key is scoped to folders and permissions (PRD M10), so a leaked CI credential is bounded. |

Refresh follows the backend's platform-aware TTLs (`backend/domains.md` §7) with the CLI presenting as a non-web platform. **`auth logout` revokes server-side**, not just deleting the local token — a logout that leaves a live refresh token on the server isn't a logout.

---

## Shell completions and man page

`cobra completion bash|zsh|fish|powershell`, generated at build time into the release archives and **installed by the packages** — the Homebrew formula, `.deb` and `.rpm` all drop the completion file in the right place. A completion that ships in the tarball but isn't installed is one nobody discovers.

Completions are **dynamic where it matters**: `bloberry cp bloberry://<TAB>` completes remote folders by calling the API, and `bloberry key revoke <TAB>` completes key IDs. Static completion of subcommand names only would miss the part that's actually hard to type.

Man page via `cobra doc`, since `.deb`/`.rpm` are in scope and those users expect `man bloberry`.

---

## Testing

**Golden-file tests.** Each command's sample-output block in `commands/*.md` becomes a fixture in `testdata/golden/` — the design artifact and the test assertion describe the same thing, which is what keeps documented output honest.

Per command, at minimum: happy path, empty result, one input-error case, and `--json` where it exists. Plus, for the recursive commands, a **partial-failure** case asserting exit code 7 and the summary format.

---

## Distribution

One row per channel, with what a user actually types.

| Channel | Install command | Host | Automation | Credentials |
|---|---|---|---|---|
| **GitHub Releases** | download + `tar xzf` | this repo's Releases | GoReleaser on tag | `GITHUB_TOKEN` |
| **Homebrew tap** | `brew install <org>/tap/bloberry` | `<org>/homebrew-tap` repo | GoReleaser `brews:` | PAT with write access to the tap repo |
| **Scoop bucket** | `scoop bucket add <org> …` then `scoop install bloberry` | `<org>/scoop-bucket` repo | GoReleaser `scoops:` | same PAT |
| **`.deb` / `.rpm`** | `dpkg -i bloberry_*.deb` / `rpm -i` | attached to Releases | GoReleaser `nfpms:` | `GITHUB_TOKEN` |
| **`go install`** | `go install github.com/<org>/bloberry/cmd/bloberry@latest` | proxy.golang.org | none — free with a public module | none |

Archives per OS/arch: `darwin_arm64`, `darwin_amd64`, `linux_amd64`, `linux_arm64`, `windows_amd64`. Each contains the binary, LICENSE, README, completions and the man page. `checksums.txt` (SHA-256) alongside — every other channel's manifest embeds one of those hashes, so it's a hard dependency.

**Deliberately opt-in, not v1:** macOS `.pkg`, Snap, Flatpak, Winget, Chocolatey, a hosted apt/yum repo, and a `curl | sh` script. Each is maintenance forever; add one when someone asks.

**A pure-Go CLI cross-compiles every platform from one Linux runner** — so unlike the desktop app (`../desktop/README.md`), the CLI release needs no OS-native CI. Worth stating explicitly since the two ship from the same `.goreleaser.yaml`.

### Update story

`bloberry version` reports whether a newer release exists. There is **no `bloberry self-update`** that overwrites the binary in place, because most installs are package-managed and overwriting a Homebrew- or apt-managed binary breaks the package database. Instead, `version` detects the install method and prints the right command — `brew upgrade bloberry`, `scoop update bloberry`, `apt install --only-upgrade bloberry`, or a download link for a manual install. Getting this wrong is a common and genuinely destructive CLI bug.

---

## Files

- [`commands/`](commands/) — one `.md` per command: purpose, args/flags, real `--help` text, sample output per state, exit codes. Written by `generate-commands`. **The terminal output is this product's UI here**, so it gets designed rather than improvised.
- `Makefile` — build/run/test/lint/completions/release targets
- `tasks/` — the implementation task list, once `build-cli` has run
