# Command — `bloberry app delete`

## Purpose & context

- **User goal**: delete an application. Guarded: **an app with active keys cannot be deleted** (deleting it would silently orphan the keys that authorize production CI — `ERD.md` access-key lifecycle).
- **When they reach for it**: decommissioning a service or cleaning up after a migration.
- **Needs**: auth, `tenant_admin`+ role.
- **Data**: `applications` — deletion is refused while active keys exist; keys must be revoked first.

## Signature

```
bloberry app delete <id> [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `<id>` | string | — | yes | Application ID (`app_…`). |
| `--yes, -y` | bool | false | no | Skip the confirmation. |

## Help text

```
Delete an application.

Refused while the app has active keys — revoke them first, because
deleting an app with live keys silently breaks whatever uses them.
After deletion the app's keys (revoked or not) are gone.

Examples:
  bloberry app delete app_5f1d
  bloberry app delete app_5f1d --yes
```

## Output states

**Success**

```
✓ Deleted application legacy-api (app_5f1d)
```

**Refused — active keys (exit 8)**

```
Error: cannot delete acme-cms — it has 2 active keys.

Revoke them first, then delete:
  bloberry key revoke ak_5f2b
  bloberry key revoke ak_3d1c
  bloberry app delete app_3d9f
```

**`--json`**

```json
{"id":"app_5f1d","deleted":true}
```

## Exit codes

| Code | When |
|---|---|
| `0` | Deleted. |
| `8` | Refused — active keys exist (a conflict, not a failure to find). |
| `5` | No app with that ID. |
| `4` | Forbidden. |
| `2` | Missing ID. |

## Behavior notes

- **The active-key guard is the design center** — the refused message lists the exact revoke commands, turning a refusal into the next step.
- **Confirmation** unless `--yes`: "Delete acme-cms? This removes the app and all its (revoked) keys." Non-TTY without `--yes` fails with "pass `--yes`".
- **stdout**: the result. **stderr**: the confirmation, the refusal message.
- **Idempotent**: deleting an already-deleted app exits `5` (not found) — no silent success, no fake failure; the message points at `app list`.
