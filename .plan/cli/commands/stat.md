# Command — `bloberry stat`

## Purpose & context

- **User goal**: see one object's or folder's full metadata — size, type, visibility, `file_id`, backend, uploader, timestamps.
- **When they reach for it**: interactively ("what is this file exactly"), and in scripts to fetch a `file_id` to store in another system.
- **Needs**: auth, tenant context, read permission.
- **Data**: `objects` (`id`, `name`, `size_bytes`, `content_type`, `visibility`, `content_hash`, `backend_id`, `uploaded_by`, timestamps), `folders` (`id`, `path`, `ancestors`, descendant counts).

## Signature

```
bloberry stat <path> [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `<path>` | path | — | yes | `bloberry://` object or folder path. |

## Help text

```
Show metadata for one object or folder.

Prints the stable file_id, size, type, visibility, storage backend,
hash and timestamps. Every field here is exactly what the web
file-detail page shows.

Examples:
  bloberry stat bloberry://assets/hero.png
  bloberry stat bloberry://assets --json
```

## Output states

**Success (object)**

```
Name:       hero.png
file_id:    f_8Kd2pQxL31A
Size:       2.4 MB (2,400,000 bytes)
Type:       image/png
Visibility: public
Backend:    s3-eu-prod (healthy)
Hash:       sha256:9f2a4c8e…
Uploaded:   Jane Doe (user_8f2a1c)
Modified:   2026-03-12 14:03:08 UTC
Created:    2026-03-12 09:41:00 UTC
```

**`--json` (object)**

```json
{"id":"f_8Kd2pQxL31A","name":"hero.png","size":2400000,"content_type":"image/png","visibility":"public","backend_id":"sb_3d9f","content_hash":"9f2a4c8e…","uploaded_by":"user_8f2a1c","created_at":"2026-03-12T09:41:00Z","modified_at":"2026-03-12T14:03:08Z"}
```

## Exit codes

| Code | When |
|---|---|
| `0` | Stat succeeded. |
| `5` | Path doesn't exist. |
| `4` | Forbidden to read. |
| `3` | Not authenticated. |

## Behavior notes

- **stdout**: the key/value lines (human) or the single JSON object (`--json`). Nothing else.
- **`file_id` is the point** — it's the stable identity (PRD M4) another system should store; the human output puts it second right after the name, and scripts use `--json` + `.id`.
- **Folder stats** report the descendant counts (objects + subfolders + bytes) so a script can decide whether to recurse.
- **Alias**: `stat` on a folder path with a trailing `/` and on a bare name behave identically (path resolution is the same).
