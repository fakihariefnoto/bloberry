# Command — `bloberry admin tenant create`

## Purpose & context

- **User goal**: create a tenant with its quota and default backend (PRD M2/M17/PA2) — install provisioning from a terminal or a scripted bootstrap.
- **When they reach for it**: setting up a new project on the install.
- **Needs**: auth as `platform_admin`.
- **Data**: `tenants` — name, slug, `quota_bytes`/`quota_objects`, `default_backend_id`. The tenant's root folder is created with it.

## Signature

```
bloberry admin tenant create <name> [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `<name>` | string | — | yes | Tenant display name. |
| `--slug <slug>` | string | derived from name | no | URL slug; unique per install. |
| `--quota-bytes <n>` | string | `0` (unlimited) | no | Byte quota; `0` = unlimited (`ERD.md` tenants note). |
| `--quota-objects <n>` | int | `0` (unlimited) | no | Object quota. |
| `--backend <id>` | string | — | yes | Default backend (`sb_…`) for new objects. |

## Help text

```
Create a tenant (platform admin).

Every tenant gets a root folder; assign its default backend and quota
here. quota-bytes 0 means unlimited (not zero storage). The owner is
added later by invitation.

Examples:
  bloberry admin tenant create "Acme Inc" --backend sb_3d9f --quota-bytes 500GB
  bloberry admin tenant create Folio --backend sb_1d4e --quota-bytes 2TB --json
```

## Output states

**Success**

```
✓ Created tenant acme (Acme Inc)
  Slug:    acme
  Backend: s3-eu-prod (default)
  Quota:   500 GB · 1,000,000 objects
  Next: invite the owner
    bloberry admin tenant invite acme owner@acme.dev  (or the dashboard)
```

**`--json`**

```json
{"id":"tnt_8f2a","slug":"acme","name":"Acme Inc","backend_id":"sb_3d9f","quota_bytes":500000000000,"quota_objects":1000000}
```

## Exit codes

| Code | When |
|---|---|
| `0` | Created. |
| `8` | Slug taken. |
| `4` | Forbidden — not `platform_admin`. |
| `2` | Bad invocation (missing `--backend`, bad quota size). |

## Behavior notes

- **Sizes are human-readable**: `500GB`, `2TB` accepted (and plain bytes); the output and `--json` both state the canonical byte count.
- **The output chains to the next step** — a tenant with no owner is inert (`ERD.md`); the success message names the owner-invite next step rather than declaring victory at a half-provisioned tenant.
- **stdout**: the result / JSON. **stderr**: nothing on success.
- **No confirmation** (creation is reversible via `admin`/dashboard delete).
