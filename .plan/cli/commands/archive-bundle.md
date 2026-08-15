# Command — `bloberry archive bundle <paths...>`

## Purpose & context

- **User goal**: request N objects as one generated archive (PRD M11/AP5) — "download all" as a single link rather than a client-side zip.
- **When they reach for it**: interactive ("give me everything under this folder as one zip"); CI generating a distributable bundle.
- **Needs**: auth, tenant context, `read` on all included objects (a bundle never includes what the caller can't read — the resolver filters).
- **Data**: `jobs` (kind `bundle`) + `objects`. Result is a URL to the generated archive; queued server-side (PRD D3).

## Signature

```
bloberry archive bundle <paths...> [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `<paths...>` | path (repeatable) | — | yes | Objects/folders to include. |
| `--format <f>` | string | `zip` | no | `zip` or `tar.gz`. |
| `--wait` | bool | false | no | Block until the bundle is ready. |

## Help text

```
Create an archive containing the given objects or folder subtrees.

The archive is built server-side and served from a link. Only objects
you can read are included. Pass --wait to block until it's ready
(CI).

Examples:
  bloberry archive bundle bloberry://assets/v2 bloberry://README.md
  bloberry archive bundle bloberry://release -r --wait
```

## Output states

**Success (queued)**

```
✓ Bundle job queued
  Job: job_3d1c  (watch: bloberry job watch job_3d1c)
```

**Success (`--wait`)**

```
✓ Bundle ready (2.1 GB, 41 files)
  URL: https://bloberry.example.com/bundles/job_3d1c
  (signed, expires in 24h)
```

**`--json` (`--wait`)**

```json
{"job_id":"job_3d1c","state":"succeeded","url":"https://bloberry.example.com/bundles/job_3d1c","files":41,"bytes":2100000000,"expires_at":"2026-03-13T14:00:00Z"}
```

## Exit codes

| Code | When |
|---|---|
| `0` | Queued, or (with `--wait`) ready. |
| `5` | One of the paths doesn't exist (nothing queued). |
| `4` | Forbidden on an included path (nothing queued). |
| `2` | Bad invocation (no paths, bad format). |

## Behavior notes

- **Empty input is refused up front** (exit 2) — no "bundle of nothing" job.
- **Permission filtering**: the bundle contains exactly what the caller can read; a path that exists but isn't readable fails the whole request with `forbidden` (exit 4) rather than silently producing a partial archive — silent omissions in a download are data loss.
- **The result URL is signed** and expires (default 24h) — same revocation/TTL story as `share link`; the output states the expiry.
- **stdout**: the result / JSON. **stderr**: progress.
- **Not idempotent**: each run queues a fresh bundle job; the result is a new signed URL.
