# Command — `bloberry key list`

## Purpose & context

- **User goal**: see an application's access keys — masked, with scope, expiry, last-used, state. The audit-before-revoke view.
- **When they reach for it**: interactive review ("what keys exist, when were they last used"), and CI preflight.
- **Needs**: auth, `tenant_admin`+ role, `admin` permission.
- **Data**: `access_keys` — prefix + `last_four` (masked), `scope_folder_ids`, `permissions`, `expires_at`, `last_used_at`, `revoked_at`.

## Signature

```
bloberry key list [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `--app <id>` | string | all apps | no | Restrict to one application. |
| `--status <s>` | string | all | no | Filter: `active`, `expiring`, `revoked`, `expired`. |

## Help text

```
List access keys, masked — only the last 4 characters of each secret
are ever shown.

Shows scope, expiry and last-used, which is the information you need
before revoking. Active keys sort by last-used descending.

Examples:
  bloberry key list
  bloberry key list --app app_3d9f
  bloberry key list --status active --json
```

## Output states

**Success**

```
APP        KEY                  SCOPE          PERMS      LAST USED          STATE
acme-cms   blob_live_••••4f2a   whole-tenant   all        Mar 13 09:12       active
acme-cms   blob_live_••••c9e7   2026/          read/write Mar 12 14:03       expiring (5d)
ci-deploy  blob_live_••••3a9c   scripts/       write      never              revoked
```

**No keys**

```
No access keys yet. Create one:
  bloberry key create --app app_3d9f
```

**`--json`**

```json
[{"id":"ak_5f2b","app_id":"app_3d9f","prefix":"blob_live_","last_four":"4f2a","scope":null,"permissions":["read","write","delete","share"],"expires_at":null,"last_used_at":"2026-03-13T09:12:00Z","revoked":false}]
```

## Exit codes

| Code | When |
|---|---|
| `0` | Listed (including empty). |
| `4` | Forbidden — not `tenant_admin`+. |
| `5` | `--app` doesn't exist. |

## Behavior notes

- **Masking is absolute**: only `prefix + •••• + last_four` in human output, and `--json` emits `last_four`, never a full secret (`ERD.md` access-keys note; `design/style-guide.md` Code/secret display).
- **`expiring` state** is `expires_at` within 7 days — the same warning the web UI shows.
- **The `state` column is text, not color alone** (TTY may tint it, but the label always carries the meaning).
- **stdout**: the table / JSON. **stderr**: nothing on success.
- **Column layout is real**: the table is a fixed-width grid (see `ls`/`share list` for the same discipline); `check-wireframes.mjs` verifies `cli/` box output too.
