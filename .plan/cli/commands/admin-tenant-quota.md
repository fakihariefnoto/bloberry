# Command — `bloberry admin tenant quota`

## Purpose & context

- **User goal**: view or change a tenant's byte/object quota (PRD M17/PA2) — the remediation when a tenant hits `quota_exceeded` (exit 6).
- **When they reach for it**: ops ("kercis is over quota — raise it or tell them"), review.
- **Needs**: auth as `platform_admin`.
- **Data**: `tenants.quota_bytes`/`quota_objects`, `used_bytes`/`used_objects` (denormalized counters, `ERD.md`). Raising a quota is instant; the metering job reconciles.

## Signature

```
bloberry admin tenant quota <slug> [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `<slug>` | string | — | yes | Tenant slug. |
| `--set-bytes <n>` | string | — | no | New byte quota (`500GB`, `2TB`, `0` = unlimited). |
| `--set-objects <n>` | int | — | no | New object quota. |

## Help text

```
Show or change a tenant's quota (platform admin).

With no --set flags, prints current usage against quota. Setting 0
means unlimited. A tenant over quota keeps reads working; only writes
are rejected (see docs, PA-E2).

Examples:
  bloberry admin tenant quota kercis
  bloberry admin tenant quota kercis --set-bytes 10GB
```

## Output states

**Show**

```
kercis (Kercis)
  Used:  4.9 GB · 22,114 objects
  Quota: 5 GB   · 1,000,000 objects
  ⚠ over quota — writes rejected until raised or usage drops
```

**Set**

```
✓ Quota updated for kercis
  Bytes: 5 GB → 10 GB
  Objects: 1,000,000 (unchanged)
  Reads still work; writes unblock immediately.
```

**`--json` (set)**

```json
{"slug":"kercis","quota_bytes":10000000000,"quota_objects":1000000,"used_bytes":4900000000,"over_quota":false}
```

## Exit codes

| Code | When |
|---|---|
| `0` | Shown or set. |
| `5` | No tenant with that slug. |
| `4` | Forbidden — not `platform_admin`. |
| `2` | Bad quota size. |

## Behavior notes

- **Writes unblock immediately on raise** — the message says so (reads never stop; the fix is instant, which is why this is the remediation command for exit-6 callers).
- **Over-quota state is called out** in the show view, with the reads-still-work caveat (PRD PA-E2).
- **stdout**: the report / result / JSON. **stderr**: nothing on success.
- **No confirmation** — a quota is a number and re-settable; the change is visible in the output.
