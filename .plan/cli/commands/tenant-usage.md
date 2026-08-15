# Command — `bloberry tenant usage`

## Purpose & context

- **User goal**: see the current tenant's storage footprint and estimated monthly cost (PRD G7, M18) without opening the dashboard — the CLI's `usage` screen.
- **When they reach for it**: a quick "how much am I using / what's this going to cost" check; CI reporting.
- **Needs**: auth, tenant context, `tenant_admin`+ (reads the tenant's `usage_snapshots` and rate card).
- **Data**: `usage_snapshots` — bytes, object count, egress, estimated cost; the tenant's quota and backend rate card.

## Signature

```
bloberry tenant usage [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `--period <p>` | string | `month` | no | `month`, `7d`, `24h`. |

## Help text

```
Show the current tenant's usage and estimated monthly cost.

Reads the hourly metering snapshots and the backend's rate card. With
no rate card, cost shows as 'unknown' — never $0.

Examples:
  bloberry tenant usage
  bloberry tenant usage --period 7d --json
```

## Output states

**Success**

```
Tenant:  acme (Acme Inc)
Backend: s3-eu-prod (healthy)

Stored:        312 GB        ████████████░░░░  62% of 500 GB
Objects:       48,912
Egress (est):  1.2 GB        ⓘ ±10% — downloads redirect
Est. cost:     $21.40 / month
  storage  $0.023/GB-mo · egress $0.09/GB · 1k reqs $0.01
```

**No rate card**

```
Est. cost:     unknown ⚠
  No rate card configured for s3-eu-prod. Ask the platform admin:
  bloberry admin backend rate-card sb_3d9f --set …
```

**`--json`**

```json
{"tenant":"acme","backend":"s3-eu-prod","stored_bytes":312000000,"objects":48912,"egress_bytes":1200000000,"est_cost":{"amount":"21.40","currency":"USD"},"quota_bytes":500000000,"rate_card_missing":false}
```

## Exit codes

| Code | When |
|---|---|
| `0` | Shown. |
| `4` | Forbidden — not `tenant_admin`+. |
| `3` | Not authenticated. |
| `2` | Bad `--period`. |

## Behavior notes

- **The "unknown, never $0" rule is shared with the web** (PRD M18): a missing rate card renders `unknown ⚠` with the exact command to fix it — the CLI turns a display gap into the remediation step.
- **Egress is labeled estimated** (±10%, `ERD.md` usage-snapshots note) — the CLI must not present an estimate as a bill.
- **stdout**: the report / JSON. **stderr**: nothing on success.
- **Quota bar** mirrors the web (`files` footer): warning ≥80%, error at quota — but always with the cost context that makes an overage legible.
