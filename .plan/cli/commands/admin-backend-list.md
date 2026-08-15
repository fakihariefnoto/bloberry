# Command — `bloberry admin backend list`

## Purpose & context

- **User goal**: see every registered storage backend, its health, and how many tenants use it (PRD M19/PA-E1) — the platform view, from a terminal.
- **When they reach for it**: ops checks; finding the backend ID a `rate-card` or `test` command needs.
- **Needs**: auth as `platform_admin`.
- **Data**: `storage_backends` — driver, name, health, assigned-tenant counts, `tenant_id` null = install-level.

## Signature

```
bloberry admin backend list [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `--driver <d>` | string | all | no | Filter by driver (`s3`, `r2`, …). |

## Help text

```
List registered storage backends with health and usage.

Health is checked periodically by the server; the raw provider error
on an unreachable backend is shown here (platform admins only).

Examples:
  bloberry admin backend list
  bloberry admin backend list --driver s3 --json
```

## Output states

**Success**

```
ID          DRIVER  NAME             BUCKET/PREFIX   HEALTH      TENANTS
sb_3d9f     s3      s3-eu-prod       app-uploads/    healthy     3
sb_7b2c     s3      s3-us-archive    archive/        healthy     1
sb_1d4e     r2      r2-main          uploads/        healthy     2
sb_5f8a     gcs     gcs-foundry      foundry-bkt/    ⚠ unreachable  1
```

**`--json` (partial)**

```json
[{"id":"sb_3d9f","driver":"s3","name":"s3-eu-prod","bucket":"app-uploads","health":"healthy","tenants":3}]
```

## Exit codes

| Code | When |
|---|---|
| `0` | Listed. |
| `4` | Forbidden — not `platform_admin`. |

## Behavior notes

- **Health carries the raw provider error** for unreachable backends — the one documented exception to never-passing-provider-errors (PRD PA-E1, `backend/domains.md` §8), visible only here and in `admin backend test`.
- **TENANTS column** shows how many tenants are assigned — the "what breaks if I remove this" answer before a delete is even attempted (`admin backends` web rule, `ERD.md`).
- **stdout**: the table / JSON. **stderr**: nothing on success.
- **Grouping**: the table lists by driver when `--driver` is absent is not required — sorting is by name; the CLI's flat table is fine where the web groups, because the web's grouping exists for visual scanning that a sortable table already provides.
