# Command — `bloberry cp`

## Purpose & context

- **User goal**: copy objects/folders in any direction — local → remote, remote → local, or remote → remote (server-side) — with structure preserved.
- **When they reach for it**: interactive file moves, and **CI publish** (`bloberry cp ./dist bloberry://assets/v2 -r`). The CI case drives the design: `--dry-run`, `--exclude`, per-file progress to stderr, and code 7 on partial failure (PRD CL1, CL-E1).
- **Needs**: auth (`auth login` or `BLOBERRY_TOKEN`), tenant context (`--tenant` or config).
- **Data**: `objects` (name, `size_bytes`), `folders` (path, `ancestors`), `multipart_uploads` for large files.

## Signature

```
bloberry cp <src> <dst> [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `<src>` | path | — | yes | Local path or `bloberry://` remote path. |
| `<dst>` | path | — | yes | Local path or `bloberry://` remote path. |
| `-r, --recursive` | bool | false | no | Copy directories recursively. |
| `--exclude <pattern>` | string (repeatable) | — | no | Glob of paths to skip. |
| `--dry-run` | bool | false | no | Show what would happen, change nothing. |
| `--no-clobber` | bool | false | no | Skip existing destinations instead of overwriting. |

Global flags apply — `../README.md`.

## Help text

```
Copy files or folders between local paths and Bloberry.

One of <src> or <dst> may be a bloberry:// remote path; when both are
remote the copy happens server-side without the bytes crossing this
machine. Directories require -r. Fails with exit 7 if some items
copy and some don't — CI should treat 7 as "retryable", not "failed".

Flags:
  -r, --recursive       Copy directories recursively
      --exclude glob    Skip matching paths (repeatable)
      --dry-run         Show what would happen, change nothing
      --no-clobber      Skip existing destinations
  -y, --yes             Skip confirmations

Examples:
  # Publish a build (up)
  bloberry cp ./dist bloberry://assets/v2 -r

  # Pull a folder (down)
  bloberry cp bloberry://assets/v2 ./dist -r

  # Server-side move between folders
  bloberry cp bloberry://a/x.zip bloberry://b/x.zip

  # What a CI script actually runs
  bloberry cp ./dist bloberry://assets/v2 -r --json
```

## Output states

**Success (interactive, single file up)**

```
✓ assets/hero.png → bloberry://assets/hero.png  2.4 MB in 0.8s
```

**Success (recursive up, piped — one line per file on stderr, summary on stdout)**

```
copied 41 files, 312 MB, 0 skipped, 0 failed  → bloberry://assets/v2
```

**`--dry-run`**

```
would copy 41 files, 312 MB → bloberry://assets/v2
would skip 3 (already up to date)
```

**Partial failure (exit 7)**

```
✗ 3 of 41 files failed → bloberry://assets/v2
  README.md:        413 payload_too_large
  node_modules/…:   quota_exceeded
  build.log:        backend_unreachable
Successful files are uploaded. Re-run to retry only what failed.
```

**`--json` (success)**

```json
{"src":"dist/","dst":"bloberry://assets/v2","copied":41,"bytes":312000000,"skipped":3,"failed":0}
```

## Exit codes

| Code | When |
|---|---|
| `0` | All items copied. |
| `7` | Some copied, some failed — summary printed; successful items stay uploaded. |
| `1` | Everything failed. |
| `2` | Bad invocation (missing src/dst, `-r` required but absent for a dir). |
| `3` | Not authenticated. |
| `6` | A write was rejected by quota. |
| `9` | The storage backend is unreachable. |

## Behavior notes

- **stdout**: the summary line (data). **stderr**: per-file progress, failures, warnings. `--json` prints only the payload on stdout.
- **Prompt/`--yes`**: none by default — `cp` overwrites silently (it's `cp`). Use `--no-clobber` to skip existing. Overwriting a name-collision destination follows the same replace semantics as the web (PRD MV-E2); there is no keep-both in a non-interactive copy.
- **Partial failure is the design center**: exit `7` (PRD CL-E1), a per-file summary on stderr, successful files left in place, and a rerun is idempotent — this is the documented behavior, not an accident of iteration order.
- **TTY**: progress bar only when stderr is a terminal; a piped run prints one line per file instead of redraws.
- **Server-side remote→remote** never streams through the client (`architecture.md` §3.2) — size and speed are reported from the server's job/result.
- **Dry-run** performs all reads (stat, list) but zero writes; safe to run anywhere.
