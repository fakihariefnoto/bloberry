# Command — `bloberry admin tenant list`

## Purpose & context

- **User goal**: the platform admin's tenant overview — every tenant, its usage, quota and cost (PRD PA2/PA3). The "who's causing this bill" view in a terminal.
- **When they reach for it**: ops review; before a quota change.
- **Needs**: auth as `platform_admin`.
- **Data**: `tenants` + derived `usage_snapshots` — status, used bytes/objects, quota, backend, est. cost.

## Signature

```
bloberry admin tenant list [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `--sort <col>` | string | `cost` | no | `cost`, `bytes`, `name`, `created`. |
| `--json` | — | — | — | Global flag. |

## Help text

```
List all tenants with usage, quota and estimated cost.

The default sort is by estimated cost descending — the first row is
the answer to "who's burning money". Cost is 'unknown' for tenants
whose backend has no rate card.

Examples:
  bloberry admin tenant list
  bloberry admin tenant list --sort bytes --json
```

## Output states

**Success**

```
SLUG     NAME        STATUS     USED             QUOTA    EST. COST
folio    Folio Notes active     1.1 TB · 284k    2 TB     $64.10
acme     Acme Inc    active     312 GB · 48.9k   500 GB   $21.40
masjid   Masjid App  active     812 MB · 3.2k    10 GB    $0.28
legacy   Legacy      suspended  2.4 GB · 12k     5 GB     $0.90
kercis   Kercis      over quota 4.9 GB · 22k     5 GB     $4.10
```

**`--json` (partial)**

```json
[{"slug":"folio","name":"Folio Notes","status":"active","used_bytes":1100000000000,"quota_bytes":2000000000000,"est_cost":"64.10"}]
```

## Exit codes

| Code | When |
|---|---|
| `0` | Listed. |
| `4` | Forbidden — not `platform_admin`. |

## Behavior notes

- **Cost-first default sort** — the platform admin's actual question; every other sort is a deliberate choice away from it.
- **`over quota` / `suspended` statuses are text, not color alone** (TTY may tint, labels carry meaning) — a script greps the status column and gets a plain word.
- **stdout**: the table / JSON. **stderr**: nothing on success.
- **Empty**: "No tenants yet — create one: `bloberry admin tenant create …`."
