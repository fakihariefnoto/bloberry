# Command — `bloberry folder create`

## Purpose & context

- **User goal**: create a folder (or nested path) in the tenant tree.
- **When they reach for it**: interactive scaffolding; occasionally scripted by setup scripts that want a known structure to exist before a sync/cp.
- **Needs**: auth, tenant context, `write` on the parent folder.
- **Data**: `folders` — name, parent, path. The `ancestors` array is built server-side.

## Signature

```
bloberry folder create <path> [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `<path>` | path | — | yes | `bloberry://` path; creates intermediate folders. |

## Help text

```
Create a folder, creating any missing parents along the way.

Creating bloberry://assets/v2/2026 creates assets/, assets/v2/ and
assets/v2/2026/ if they don't exist. Exits 8 if the final folder
already exists.

Examples:
  bloberry folder create bloberry://assets/v2
  bloberry folder create bloberry://projects/web/2026 --json
```

## Output states

**Success**

```
✓ Created folder assets/v2 (fd_4b8a2c1e)
```

**Created with intermediate folders**

```
✓ Created assets/v2/2026 (fd_9c3f1d2b)
  (also created 2 intermediate folders: assets/, assets/v2/)
```

**Already exists (exit 8)**

```
Error: folder "bloberry://assets/v2" already exists

No change made. (name_conflict)
```

**`--json` (success)**

```json
{"id":"fd_4b8a2c1e","path":"assets/v2","created_intermediates":2}
```

## Exit codes

| Code | When |
|---|---|
| `0` | Created (including intermediate parents). |
| `8` | Folder already exists — no change. |
| `5` | Parent path can't be created (e.g. a component is actually an object). |
| `4` | Forbidden — no `write` on the parent. |
| `2` | Bad invocation (empty path, missing path). |

## Behavior notes

- **Idempotent-ish**: `--json` callers can treat exit `8` as "already there" and proceed; the message says "no change made" so a human doesn't assume a fresh folder.
- **stdout**: the created-folder line. **stderr**: nothing on success.
- **Intermediate creation is a single API call** (the server handles the recursion), so a race between "check then create" isn't this CLI's problem — it creates and reconciles.
- **No confirmation** (creating a folder is reversible by `rm`); destructive operations own the prompts, creation doesn't.
