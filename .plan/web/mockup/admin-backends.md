# Screen — admin-backends

## Purpose & context

- **User goal**: the platform admin's storage infrastructure — register backends, enter credentials once, set rate cards, see health, see which tenants use each one (PRD PA1, PA4, PA-E1). The screen that makes "storage-agnostic" operational rather than theoretical.
- **Entry points**: sidebar Storage backends (`platform_admin` only).
- **Exit points**: click a backend → `admin-backend-detail`; Register → modal; delete → `confirm-destructive` (refused when tenants are assigned).
- **Data needed**: `storage_backends` — driver, name, bucket/prefix, `health_status`/`health_error`, `rate_card`, and derived: count of tenants assigned, last health check. Per `ERD.md` storage-backends note: **many backends per driver type** — this is a list of accounts, not one row per provider.

## States

- [x] Loading (skeleton)
- [x] Empty (no backends yet — first-run state)
- [x] Populated
- [x] Error
- [x] Domain-specific — unreachable backend (error row with the raw provider error, PA-E1)
- [x] Domain-specific — backend with tenants assigned (delete refused)

## Style reference

- **Components used**: `AppShell`, `DataTable`, `StatusPill` (health), register modal, `ConfirmDestructive`.
- Rows grouped by driver, not flattened — the ERD's "many backends per driver" consequence (`ERD.md` storage-backends note).
- No token deltas.

## Wireframe — desktop (populated)

```
┌──────────┬───────────────────────────────────────────────────────────┐
│ (sidebar)│  Platform ▾        [👤] Platform Admin ▾                  │
│          │  Storage backends                            [+ Register] │
│          │  ┌──────────────────────────────────────────────────────┐ │
│         │  │ Driver  Name         Bucket/prefix   Health  Tenants ⋮│ │
│         │  │ s3      s3-eu-prod    app-uploads/    healthy 3     ⋮ │ │
│         │  │ s3      s3-us-archive archive/       healthy 1     ⋮  │ │
│         │  │ r2      r2-main       uploads/        healthy 2     ⋮ │ │
│         │  │ gcs     gcs-foundry   foundry-bkt/    ⚠ unreachable  1│ │
│         │  │ azblob  az-eu-prod    bloberry-prod/ healthy 2     ⋮  │ │
│         │  │ disk    vps-volume    /data/blob/     healthy 0     ⋮ │ │
│          │  └──────────────────────────────────────────────────────┘ │
│          │  Showing 6 of 6                                           │
│          │  Health is checked every 5 minutes. The error text on     │
│          │  an unreachable backend is visible only here.             │
└──────────┴───────────────────────────────────────────────────────────┘
```

## Wireframe — mobile (populated)

```
┌────────────────────────────┐
│ ☰ ⬡    Platform ▾    👤    │
├────────────────────────────┤
│ Backends         [+ Add]   │
│ ─────────────────────────  │
│ S3                         │
│ ┌────────────────────────┐ │
│ │ s3-eu-prod     healthy │ │
│ │ app-uploads/ · 3tenants│ │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │ s3-us-archive  healthy │ │
│ │ archive/ · 1 tenant    │ │
│ └────────────────────────┘ │
│ R2                         │
│ ┌────────────────────────┐ │
│ │ r2-main        healthy │ │
│ │ uploads/ · 2 tenants   │ │
│ └────────────────────────┘ │
│ GCS                        │
│ ┌────────────────────────┐ │
│ │ ⚠ gcs-foundry  unreach.│ │
│ │ foundry-bkt/ · 1 tenant│ │
│ └────────────────────────┘ │
│ Health checked every 5  min│
└────────────────────────────┘
```

## Interactions

- **Register**: modal with the driver selector (s3/r2/oss/gcs/azblob/disk), name, config (endpoint/bucket/prefix per driver), **credentials fields** (Azure: SharedKey or SAS; write-only), and the rate card (storage $/GB-mo, egress $/GB, per-1k-requests). Credentials are write-only — never echoed back on edit (PRD M20, R7). The modal notes "credentials are envelope-encrypted at rest".
- **Driver grouping**: rows group by driver (S3 / R2 / OSS / GCS / Disk headers), reflecting that one install legitimately has several distinct S3 accounts (`ERD.md` note). Sorting within a group by name.
- **Health**: polled every 5 minutes (`infra` — the in-process health ticker). An unreachable backend expands to show the **raw provider error** (PA-E1), the one place it's legal.
- **Tenants column** answers the deletion question before the admin asks it: a backend with tenants assigned gets **delete refused** ("3 tenants use this — reassign them first"). Only a 0-tenant backend can be deleted, and that still confirms.
- **Rate card**: missing rate card on a backend shows "no rate card" (`color.warning`) — usage shows "unknown" for tenants on it, so the admin sees the gap in the place that motivates filling it.
- **Empty state**: "No storage backends · Register one to start assigning tenants" — this is the true first-run gate (nothing else works until a backend exists).
- **A11y**: the health error is revealed on expand and announced; the register modal traps focus and returns it to the Register button.
