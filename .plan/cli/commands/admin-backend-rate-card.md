# Command — `bloberry admin backend rate-card`

## Purpose & context

- **User goal**: view or set a backend's rate card — the input side of PRD M18/G7's estimated cost (storage $/GB-mo, egress $/GB, requests per 1k).
- **When they reach for it**: configuring cost visibility; auditing what a tenant's estimate is based on.
- **Needs**: auth as `platform_admin`.
- **Data**: `storage_backends.rate_card`. A missing rate card makes tenant cost "unknown" — this is the fix command the `usage` screens point at.

## Signature

```
bloberry admin backend rate-card <id> [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `<id>` | string | — | yes | Backend ID (`sb_…`). |
| `--set <s>` | string | — | no | Set in `storage,egress,requests` order: `0.023,0.09,0.01`. |

## Help text

```
Show or set a backend's rate card.

The rate card drives estimated monthly cost (usage commands). The
format is storage$/GB-mo, egress$/GB, requests$per-1k:
  --set 0.023,0.09,0.01

Examples:
  bloberry admin backend rate-card sb_3d9f
  bloberry admin backend rate-card sb_3d9f --set 0.023,0.09,0.01
```

## Output states

**Show**

```
s3-eu-prod (sb_3d9f)
  Storage:  $0.023 / GB-month
  Egress:   $0.09   / GB
  Requests: $0.01   / 1,000
```

**Set**

```
✓ Rate card updated for s3-eu-prod
  Storage $0.023/GB-mo · Egress $0.09/GB · Requests $0.01/1k
```

**No rate card**

```
s3-eu-prod (sb_3d9f) has no rate card.

Tenant cost shows as "unknown" until one is set:
  bloberry admin backend rate-card sb_3d9f --set 0.023,0.09,0.01
```

**`--json` (show)**

```json
{"id":"sb_3d9f","name":"s3-eu-prod","rate_card":{"storage_per_gb_month":"0.023","egress_per_gb":"0.09","per_1k_requests":"0.01"}}
```

## Exit codes

| Code | When |
|---|---|
| `0` | Shown or set. |
| `2` | Bad `--set` format (must be three decimal numbers). |
| `5` | No backend with that ID. |
| `4` | Forbidden — not `platform_admin`. |

## Behavior notes

- **No rate card → the output gives the fix command** — this is the command `tenant usage`'s "ask the platform admin" points at; the CLI completes the loop itself.
- **Historical honesty**: estimates are computed at snapshot time from the then-current card (`ERD.md` usage-snapshots note) — changing the card does not rewrite past months, worth stating because an admin will notice old numbers not moving.
- **stdout**: the card / result / JSON. **stderr**: nothing on success.
- **Setting is immediate and non-destructive** (no confirmation — a rate card is a number, easily re-set).
