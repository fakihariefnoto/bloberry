# Command — `bloberry sync`

## Purpose & context

- **User goal**: one-way mirror a local folder into a Bloberry folder (the desktop app's sync engine, scripted — PRD S1/DT2). Uploads newer/changed files, and with `--delete`, removes remote objects no longer present locally.
- **When they reach for it**: `rsync`-style publishing of a build or working directory; the CI/desktop-provisioning story. **One-way only** — this is explicitly not a backup/restore sync (PRD NG4; the direction is local → remote).
- **Needs**: auth, tenant context. `--delete` needs delete permission on the target.
- **Data**: `objects` (name, `size_bytes`, `content_hash`, `modified_at`) — sync keys off hash + size to avoid re-uploading unchanged files (PRD S3 gives us the hash).

## Signature

```
bloberry sync <local> <remote> [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `<local>` | local path | — | yes | Directory to mirror from. |
| `<remote>` | `bloberry://` path | — | yes | Directory to mirror into. |
| `--delete` | bool | false | no | Delete remote objects that no longer exist locally. |
| `--dry-run` | bool | false | no | Show what would change, change nothing. |
| `--exclude <pattern>` | string (repeatable) | — | no | Glob to skip (e.g. `.git`, `node_modules`). |

## Help text

```
Mirror a local directory into a Bloberry folder (one-way).

Uploads files that are new or changed (compared by size + content
hash), and with --delete removes remote objects that no longer exist
locally. Never deletes local files — this is publish, not two-way
sync. Dry-run is free and recommended before the first real run.

Flags:
      --delete        Remove remote objects missing locally
      --dry-run       Show what would change, change nothing
      --exclude glob  Skip matching paths (repeatable)

Examples:
  bloberry sync ./dist bloberry://assets/v2 --dry-run
  bloberry sync ./dist bloberry://assets/v2 --delete
  bloberry sync ./site bloberry://site -r --json
```

## Output states

**`--dry-run`**

```
would upload 12 (new), update 3 (changed), delete 2 (with --delete), unchanged 41
```

**Success**

```
✓ 15 uploaded · 41 unchanged · 2 deleted · 0 failed → bloberry://assets/v2
```

**Partial failure (exit 7)**

```
✗ 2 of 15 failed → bloberry://assets/v2
  dist/icon-2x.png: quota_exceeded
  dist/vendor.js:   payload_too_large
Re-run to upload only what failed (already-synced files are skipped).
```

**`--json` (success)**

```json
{"uploaded":15,"updated":3,"unchanged":41,"deleted":2,"failed":0}
```

## Exit codes

| Code | When |
|---|---|
| `0` | All changes applied. |
| `7` | Some changed, some failed — summary on stderr; successful uploads persist. |
| `1` | Nothing succeeded (fatal error). |
| `2` | Bad invocation (missing dirs, `--delete` on the tenant root refused). |
| `6` | Quota blocked the upload. |

## Behavior notes

- **Comparison is by size + `content_hash`** — a file with identical bytes is skipped even if mtime differs. This makes reruns cheap and idempotent.
- **`--delete` is the dangerous half**: refuses when `<remote>` is the tenant root or a top-level folder without explicit confirmation; in a non-TTY without `--yes`, it fails rather than guessing.
- **stdout**: the summary. **stderr**: per-file progress, delete warnings, errors.
- **Direction is one-way local → remote, stated in the help** — nothing ever deletes a local file (PRD NG4: "must be described as such in the UI").
- **TTY**: a live progress bar when stderr is a terminal; one line per file when piped.
- **Name collisions** follow replace semantics (a changed local file replaces the remote); there's no keep-both in a mirror.
