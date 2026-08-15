# Screen — usage

## Purpose & context

- **User goal**: check the tenant's storage footprint and estimated monthly cost from a phone — the admin's "am I near quota / what's the bill" glance (PRD G7, M18).
- **Entry points**: Usage from `more` (`tenant_admin`+ gate).
- **Exit points**: back → `more`; tap a stat → its own detail view (per the mobile analytics pattern — a stat card jumps to its metric's history); "see raw events" isn't here (audit is web-only for mobile in v1 — this screen is metering, not investigation).
- **Data needed**: `usage_snapshots` + rate card — bytes stored, object count, egress, est. monthly cost, quota position, trend.

## States

- [x] Loading (skeleton stat cards + chart placeholder — never a "0" flash)
- [x] Populated
- [x] Error
- [x] Domain-specific — cost "unknown" when no rate card (never $0 — PRD M18)
- [x] Domain-specific — no data in range (fresh tenant)

## Style reference

- **Components used**: analytics-dashboard pattern (`design-collection/mobile-screen/patterns.md`) — swipeable stat cards (2-3 per screen, horizontal scroll if more), **one primary full-width chart**, compact breakdown below. `DateRangePicker` (segmented control, global).
- Mobile-specific rules: stat cards scroll horizontally (don't push the chart below the fold); one chart per width; tap a stat → its detail.
- No token deltas.

## Wireframe — mobile (populated)

```
┌────────────────────────────┐
│  ← Usage     [This mo ▾]   │
├────────────────────────────┤
│  ┌─────────────┐ ┌───────┐ │
│ │ 312 GB stored│ │$21.40 │ │
│ │  ▁▃▅▂▇  ↑6%  │ │ est.mo│ │
│  └─────────────┘ └───────┘ │
│  ───────────────────────── │
│  Stored bytes · this month │
│  ┌───────────────────────┐ │
│  │   ╭───╮        ╭─╮    │ │
│  │ ╭─╯  ╰─╮  ╭────╯  ╰╮  │ │
│  │ Mar 1          Mar  31│ │
│  └───────────────────────┘ │
│  Rate card: s3-eu-prod     │
│  ⓘ egress estimated  (±10%)│
│  ───────────────────────── │
│  Egress     1.2 GB ↑41%    │
│  Objects    48,912  ↑3%    │
│  ───────────────────────── │
│  Quota 312/500 GB  ████ 62%│
│  ───────────────────────── │
│  Est. cost  $21.40         │
│  storage $0.023/GB-mo      │
│  egress  $0.09/GB          │
│  (this month, from rate    │
│   card)                    │
└────────────────────────────┘
```

## Wireframe — mobile (no rate card)

```
┌────────────────────────────┐
│  ← Usage     [This mo ▾]   │
├────────────────────────────┤
│  Stored     312 GB  ↑6%    │
│  Objects    48,912  ↑3%    │
│  ───────────────────────── │
│  ⚠ Cost unknown            │
│  No rate card configured   │
│  for s3-eu-prod. Ask the   │
│  platform admin to add one.│
│  ───────────────────────── │
│  Quota 312/500 GB    62%   │
└────────────────────────────┘
```

## Interactions

- **Period control** (segmented: 24h / 7d / 30d / This month) is global — chart, cards, breakdown all respond together.
- **Stat cards** scroll horizontally; tapping a card opens that metric's full chart + history (a detail screen) — the analytics pattern's "tap a stat to jump to detail".
- **Rate card rule**: no rate card → "Cost unknown" in `color.warning`, never "$0.00" (`usage`/`admin-usage` share this rule; `admin-backend-detail` is where the platform admin fixes it).
- **Egress figures** carry the ⓘ estimated (±10%) note (`ERD.md` usage-snapshots).
- **Quota bar** mirrors `files`/`more` but with cost context — near-limit shows the projected overage cost.
- **Loading**: skeletons at the real dimensions; the chart is a flat placeholder, not stale data or a blank axis.
- **Empty**: "No usage data for this period" — a brand-new tenant's first hour, not an error.
- **A11y**: the chart has a data-point summary for screen readers; stat cards are announced with label + value + direction.
