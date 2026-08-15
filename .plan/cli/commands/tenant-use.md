# Command — `bloberry tenant use <slug>`

## Purpose & context

- **User goal**: set the default tenant for subsequent commands — persisted to the config file (`cli/README.md` precedence).
- **When they reach for it**: switching between projects in a terminal session; the script equivalent is `--tenant <slug>` per command or a fully-qualified `bloberry://<tenant>/…` path (the README's "CI scripts should use it" rule).
- **Needs**: auth; membership in the target tenant.
- **Data**: `memberships` — the tenant must be one the user belongs to.

## Signature

```
bloberry tenant use <slug> [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `<slug>` | string | — | yes | Tenant slug. |

## Help text

```
Set the default tenant for subsequent commands.

Persists to the config file. For scripts, prefer --tenant or the
fully-qualified bloberry://<tenant>/path form — ambient state is
exactly what a CI script shouldn't rely on.

Examples:
  bloberry tenant use folio
  bloberry tenant use acme
```

## Output states

**Success**

```
✓ Default tenant is now folio (Folio Notes)
```

**Not a member (exit 4)**

```
Error: you're not a member of tenant "secret-tenant"

Run 'bloberry tenant list' to see which tenants you can use.
```

**`--json`**

```json
{"tenant":"folio","name":"Folio Notes","current":true}
```

## Exit codes

| Code | When |
|---|---|
| `0` | Switched. |
| `4` | Not a member of that tenant. |
| `5` | No tenant with that slug. |
| `2` | Missing slug. |

## Behavior notes

- **Writes the config file** (`~/.config/bloberry/config.yaml` → `tenant:`); the file is managed via `config set` discipline, and this command is the ergonomic wrapper for the `tenant` key specifically.
- **stdout**: the result. **stderr**: nothing on success.
- **Membership-checked**: refuses (exit 4) if the user isn't a member — a `tenant use` into a tenant you can't see would make every later command fail with a confusing 403; this fails early and clearly.
- **Idempotent**: switching to the current tenant is a success no-op (still reports the tenant).
