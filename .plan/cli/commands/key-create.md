# Command — `bloberry key create`

## Purpose & context

- **User goal**: issue a scoped access key for an application (PRD M10) — the CLI's way to provision CI credentials. This is the command that produces `BLOBERRY_TOKEN` values.
- **When they reach for it**: setting up CI, or when a leaked key was revoked and a replacement is needed.
- **Needs**: auth, `tenant_admin`+ role, `admin` permission in scope. The key's secret is shown **exactly once** (PRD D5/M10) — the CLI prints it once and it is unrecoverable.
- **Data**: `access_keys` — prefix (`blob_live_`/`blob_test_`), scope_folder_ids, permissions, expires_at. The `secret_hash` is argon2id server-side; only the one-time secret round-trips.

## Signature

```
bloberry key create [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `--app <id>` | string | — | yes | The application to issue the key for (`app_…`). |
| `--name <label>` | string | — | no | Human label shown in key lists. |
| `--folder <path>` | path (repeatable) | whole tenant | no | Restrict the key to these subtrees. |
| `--permission <p>` | string (repeatable) | all | no | `read`, `write`, `delete`, `share`, `admin`. |
| `--expires <when>` | string | never | no | `30d`, `2026-12-31`, or ISO timestamp. |
| `--prefix <p>` | string | `blob_live_` | no | `blob_live_` or `blob_test_`. |

## Help text

```
Create a scoped access key for an application.

The secret is printed once and never again — capture it immediately.
Scope narrows: with --folder the key only reaches those subtrees, and
--permission only those actions. An unscoped key is whole-tenant,
which you usually don't want for CI.

Examples:
  bloberry key create --app app_3d9f --folder bloberry://assets/v2 --permission read,write
  bloberry key create --app app_3d9f --expires 30d --json
```

## Output states

**Success**

```
✓ Access key created for acme-cms (app_3d9f)

  Key:    blob_live_9fK2mQxL31A8pQ4rT7vB2cD5eF8hJ1k
  Scope:  assets/v2/ (read, write)
  Expires: 2026-04-11 (in 30 days)

  ⚠ This secret is shown once and cannot be recovered. Store it now
  in your secrets manager or BLOBERRY_TOKEN in CI.
```

**`--json`** — note the secret is in this payload, so scripts must capture stdout before it's gone:

```json
{"id":"ak_5f2b","app_id":"app_3d9f","prefix":"blob_live_","secret":"blob_live_9fK2mQxL31A8pQ4rT7vB2cD5eF8hJ1k","last_four":"J1k","scope":["assets/v2"],"permissions":["read","write"],"expires_at":"2026-04-11T09:00:00Z"}
```

## Exit codes

| Code | When |
|---|---|
| `0` | Created; the secret was printed (or emitted) once. |
| `4` | Forbidden — not `tenant_admin`+ in the scope. |
| `2` | Bad invocation (missing `--app`, invalid `--expires`, unknown permission). |
| `5` | The app doesn't exist. |

## Behavior notes

- **The one-time secret is the entire point**: the human output labels it clearly and warns it's unrecoverable; `--json` output carries `secret` so a CI script can pipe it into a secrets manager in the same pipeline step — but once gone, it's gone (PRD D5).
- **Scoping is the default posture**: the help nudges `--folder`/`--permission`; an unscoped whole-tenant key prints a warning line "unscoped key — reaches the whole tenant".
- **`--expires` formats**: `30d`/`90d` durations or `YYYY-MM-DD`; the output states the expiry both ways.
- **stdout**: the result. **stderr**: nothing on success (the warning is part of the human result, still on stdout because a key's disclosure warning belongs with the key; `--quiet` drops it).
- **Test prefix** (`blob_test_`) exists so dev/CI staging can't touch production buckets.
