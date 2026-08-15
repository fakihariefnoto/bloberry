# Command — `bloberry admin usage`

## Purpose & context

- **User goal**: the install-wide usage and cost view (PRD G7/PA3) — the CLI equivalent of `admin-usage`, for scripts and ops terminals.
- **When they reach for it**: "what's the total bill / who's causing it", and CI reporting.
- **Needs**: auth as `platform_admin`.
- **Data**: `usage_snapshots` across tenants + each tenant's backend rate card — totals + per-tenant share.

## Signature

```
bloberry admin usage [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `--period <p>` | string | `month` | no | `month`, `7d`, `24h`. |

## Help text

```
Show install-wide usage and estimated cost (platform admin).

Totals plus a per-tenant breakdown sorted by cost. Same rules as the
dashboard: cost is 'unknown' where a backend has no rate card, and
egress is estimated (±10%).

Examples:
  bloberry admin usage
  bloberry admin usage --period 7d --json
```

## Output states

**Success**

```
Install usage (this month)
  Stored:   2.4 TB  (9 tenants)
  Egress:   5.9 GB  (est. ±10%)
  Est. cost:$87.64

By tenant (cost desc):
  folio    1.1 TB   $64.10
  acme     312 GB   $21.40
  kercis   4.9 GB   $4.10
  legacy   2.4 GB   unknown ⚠ (no rate card on sb_5f8a)
```

**`--json`**

```json
{"period":"month","stored_bytes":2400000000000,"egress_bytes":5900000000,"est_cost_total":"87.64","tenants":[{"slug":"folio","est_cost":"64.10"},{"slug":"legacy","est_cost":null}]}
```

## Exit codes

| Code | When |
|---|---|
| `0` | Shown. |
| `4` | Forbidden — not `platform_admin`. |
| `2` | Bad `--period`. |

## Behavior notes

- **Cost-desc per-tenant breakdown** answers the bill-attribution question in one pass; the `unknown` row names the backend so the fix (`admin backend rate-card sb_5f8a --set …`) is one step away.
- **Egress is estimated and labeled** (`±10%`, `ERD.md` usage-snapshots) — the CLI's totals never present an estimate as a bill.
- **stdout**: the report / JSON. **stderr**: nothing on success.
- **Same rules as the web** (`admin-usage`): one rule set, two renderings — the CLI and the dashboard must not disagree on what "unknown" means.
