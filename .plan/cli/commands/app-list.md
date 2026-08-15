# Command — `bloberry app list`

## Purpose & context

- **User goal**: list the tenant's applications and their key counts.
- **When they reach for it**: interactive review; discovering an app's ID before a `key create`/`app delete`.
- **Needs**: auth, `tenant_admin`+ role.
- **Data**: `applications` — name, description, `created_at`; derived active/revoked key counts.

## Signature

```
bloberry app list [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| — | — | — | — | No args. |

## Help text

```
List the tenant's applications and how many keys each has.

Shows the app ID you'll need for 'bloberry key create --app'.

Examples:
  bloberry app list
  bloberry app list --json
```

## Output states

**Success**

```
ID          NAME        KEYS     LAST USED        CREATED
app_3d9f    acme-cms    3         Mar 13 09:12    Mar 01
app_7b2c    ci-deploy   1         Mar 13 08:44    Mar 01
app_5f1d    legacy-api  0         never            Feb 12
```

**No applications**

```
No applications yet. Create one:
  bloberry app create <name>
```

**`--json`**

```json
[{"id":"app_3d9f","name":"acme-cms","keys":3,"last_used_at":"2026-03-13T09:12:00Z","created_at":"2026-03-01T10:00:00Z"}]
```

## Exit codes

| Code | When |
|---|---|
| `0` | Listed (including empty). |
| `4` | Forbidden — not `tenant_admin`+. |

## Behavior notes

- **`KEYS` column is a count** (total across states), consistent with the web's `applications` list; the "0 keys" row is not flagged differently here — that nuance belongs to `key list`.
- **stdout**: the table / JSON. **stderr**: nothing on success.
- **Empty state** names the create command — same rule as every list.
