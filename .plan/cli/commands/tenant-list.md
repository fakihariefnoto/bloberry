# Command — `bloberry tenant list`

## Purpose & context

- **User goal**: list the tenants the current user belongs to — the multi-tenant orientation command.
- **When they reach for it**: interactive ("which tenants am I in"), and before `tenant use` in a script.
- **Needs**: auth. Lists only tenants the user is a member of (never cross-tenant leaks, PRD M2).
- **Data**: `memberships` + `tenants` — slug, name, role, usage summary.

## Signature

```
bloberry tenant list [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| — | — | — | — | No args. |

## Help text

```
List the tenants you belong to, with your role in each.

The 'current' marker shows which tenant commands act on by default.
Switch with 'bloberry tenant use <slug>'.

Examples:
  bloberry tenant list
  bloberry tenant list --json
```

## Output states

**Success**

```
SLUG     NAME        ROLE          CURRENT
acme     Acme Inc    tenant_owner  * (default)
folio    Folio Notes member
masjid   Masjid App  viewer
```

**`--json`**

```json
[{"slug":"acme","name":"Acme Inc","role":"tenant_owner","current":true}]
```

## Exit codes

| Code | When |
|---|---|
| `0` | Listed. |
| `3` | Not authenticated. |
| `4` | — (a user is always in at least the tenants listed). |

## Behavior notes

- **stdout**: the table / JSON. **stderr**: nothing on success.
- **`* (default)` marks the current tenant** — the same information the `--tenant` flag and config carry; one source of truth displayed in the natural place.
- **Empty**: cannot occur in practice (a user has at least one membership); if it somehow does, "You belong to no tenants yet. Ask a tenant admin for an invitation."
- **Sort**: alphabetical by slug; stable across runs.
