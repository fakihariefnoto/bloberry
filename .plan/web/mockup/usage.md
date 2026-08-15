# Screen — usage

## Purpose & context

- **User goal**: see the tenant's storage footprint and estimated monthly cost (PRD G7, M18) — the "am I about to overrun the quota / get a big bill" answer.
- **Entry points**: sidebar Usage (`tenant_admin`+); quota-exceeded banner's "see usage" link from `files`; upload-queue failure's usage link.
- **Exit points**: "see raw events" → `audit`; change period → in place.
- **Data needed**: `usage_snapshots` (hourly buckets) + rate card from the tenant's `storage_backends` — bytes stored, object count, egress, request count, `estimated_cost`. Trend for the selected period.

## States

- [x] Loading (skeleton stat cards + chart placeholder — never "0" flash)
- [x] Populated
- [x] Error
- [x] Domain-specific — "cost unknown" when no rate card configured (must render "unknown", never $0 — a zero reads as free, PRD G7/M18)
- [x] Domain-specific — egress is estimated (±10%, PRD G7) — labeled as such

## Style reference

- **Components used**: `AppShell`, `PageHeader`, `DateRangePicker` (global to page, per `design-collection/web-screen/patterns.md` analytics rule), stat cards with sparklines + direction, one primary area chart, breakdown table, `ByteSize`, `RelativeTime`.
- Per the analytics-dashboard pattern: **one primary chart** (storage over time), stat cards above, breakdown below. Charts get skeleton placeholders on load, not stale data.
- No token deltas.

## Wireframe — desktop (populated)

```
┌──────────┬───────────────────────────────────────────────────────────┐
│ (sidebar)│  Acme Inc ▾        Jane Doe ▾      [This month ▾]         │
│          │  Usage                                                    │
│          │  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐       │
│         │  │ 312 GB stored│ │ 48,912 objects│ │ 1.2 GB egress│       │
│         │  │  ▁▃▅▂▇   ↑6% │ │   ▂▄▆▅▇   ↑3% │ │  ▁▂▄▆▇  ↑41% │       │
│          │  └──────────────┘ └──────────────┘ └──────────────┘       │
│          │  ┌──────────────────────────────────────────────────────┐ │
│          │  │  Stored bytes · this month                est. $21.40│ │
│          │  │    ╭──╮            ╭─╮                               │ │
│          │  │  ╭─╯  ╰─╮    ╭─────╯  ╰╮                             │ │
│          │  │ ╭╯      ╰────╯        ╰──╮                           │ │
│          │  │ Mar 1          Mar 15          Mar 31                │ │
│          │  └──────────────────────────────────────────────────────┘ │
│          │  Estimated cost   $21.40   (rate card: s3-eu-prod)        │
│          │  ── storage $0.023/GB-mo · egress $0.09/GB · 1k reqs $0.01│
│          │  ⓘ Egress is estimated (±10%) — downloads redirect.       │
│          │  ┌──────────────────────────────────────────────────────┐ │
│          │  │ Period      Stored    Objects   Egress    Est. cost  │ │
│          │  │ Mar 01-07   284 GB    46,102    640 MB     $17.80    │ │
│          │  │ Mar 08-14   301 GB    47,590    1.02 GB    $19.10    │ │
│          │  │ Mar 15-21   312 GB    48,912    1.21 GB    $21.40    │ │
│          │  └──────────────────────────────────────────────────────┘ │
│          │  Quota  312 GB of 500 GB ████████████░░░░  62%            │
└──────────┴───────────────────────────────────────────────────────────┘
```

## Wireframe — mobile (populated)

```
┌───────────────────────────────┐
│ ☰ ⬡    Acme Inc ▾   👤        │
├───────────────────────────────┤
│ Usage          [This mo ▾]    │
│ ┌──────────────┐ ┌──────────┐ │
│ │312 GB stored│ │48,912 objs│ │
│ │ ↑6%         │ │ ↑3%       │ │
│ └──────────────┘ └──────────┘ │
│ ┌───────────────────────────┐ │
│ │ Stored · this month    │    │
│ │   ╭─────╮        ╭─╮   │    │
│ │ ╭─╯  ╰─╮  ╭──────╯  ╰╮ │    │
│ │ Mar 1      Mar 31      │    │
│ │ est. $21.40            │    │
│ └───────────────────────────┘ │
│ Egress this month  1.2 GB     │
│ ⓘ estimated (±10%)            │
│ ────────────────────────────  │
│ Quota 312/500 GB   62%        │
│ ────────────────────────────  │
│ Breakdown (weekly)            │
│ Mar 01-07  $17.80    ▸        │
│ Mar 08-14  $19.10    ▸        │
│ Mar 15-21  $21.40    ▸        │
│ [see raw events →]            │
└───────────────────────────────┘
```

## Interactions

- **Period control** is global: changes the chart, stat cards, breakdown, and cost figure together. Presets 24h/7d/30d/This month/custom.
- **Rate card**: the "est. $21.40" line always shows which backend's rate card produced it. **No rate card** → "Cost unknown — no rate card configured for s3-eu-prod" in `color.warning`, never "$0.00" (M18: a zero reads as free, which is a lie).
- **Egress estimate**: every egress figure carries the ⓘ "estimated (±10%) — downloads redirect past Bloberry" (`ERD.md` usage-snapshots note, ADR-3).
- **Stat cards**: sparkline + direction indicator; loading skeletons at the card's real size — never a "0" flash (`design-collection/web-screen/patterns.md`).
- **Breakdown** rows expand to per-day values; the "see raw events" link goes to `audit` (usage is metering, audit is who).
- **Quota bar** mirrors the `files` footer but with the cost context — a tenant at 95% sees the projected cost of the overage.
- **Empty state**: no snapshots in range → "No usage data for this period" (a brand-new tenant's first hour), not a blank chart.
- **Permission-aware**: `member`/`viewer` never see this route; a `tenant_admin` sees the tenant's usage and rate card but cannot edit the rate card (that's `admin-backend-detail`, `platform_admin` only).
