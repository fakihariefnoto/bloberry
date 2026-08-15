# Screen — admin-usage

## Purpose & context

- **User goal**: the platform admin's install-wide cost view — total stored bytes, total estimated cost, and the per-tenant breakdown that attributes the bill (PRD G7, PA3). The screen the admin opens when the storage invoice is higher than expected.
- **Entry points**: sidebar Install usage (`platform_admin` only); a cost-attribution link from `admin-tenants`.
- **Exit points**: click a tenant row → `admin-tenant-detail`; "see raw events" → per-tenant audit (via detail); period → in place.
- **Data needed**: `usage_snapshots` across tenants + each tenant's backend rate card — total bytes/objects/egress, est. total cost, per-tenant share.

## States

- [x] Loading (skeleton stat cards + chart)
- [x] Populated
- [x] Error
- [x] Domain-specific — a tenant with no rate card ("unknown" — never $0, PRD M18)
- [x] Domain-specific — no data in range (fresh install)

## Style reference

- **Components used**: `AppShell`, `DateRangePicker` (global), stat cards with sparklines, one primary area chart (total stored), per-tenant breakdown `DataTable`, `ByteSize`.
- Per the analytics-dashboard pattern: primary chart largest, stat cards with direction indicators, breakdown table below. Server-side pagination on the tenant table.
- No token deltas.

## Wireframe — desktop (populated)

```
┌──────────┬───────────────────────────────────────────────────────────┐
│ (sidebar)│  Platform ▾       [👤] Platform Admin ▾      [This mo ▾]  │
│          │  Install usage                                            │
│          │  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐       │
│         │  │ 2.4 TB stored│ │ 9 tenants    │ │ $87.64 est.   │       │
│         │  │  ▁▃▅▂▇   ↑8% │ │  active 9/9  │ │  ▂▄▆▅▇   ↑12% │       │
│          │  └──────────────┘ └──────────────┘ └──────────────┘       │
│          │  ┌──────────────────────────────────────────────────────┐ │
│          │  │  Total stored bytes · all tenants                   │  │
│          │  │    ╭─╮        ╭──╮                                  │  │
│          │  │  ╭╯  ╰─╮  ╭────╯    ╰──╮                            │  │
│          │  │ ╭╯      ╰───╯         ╰──╮                          │  │
│          │  │ Mar 1          Mar 15          Mar 31               │  │
│          │  └──────────────────────────────────────────────────────┘ │
│          │  Cost by tenant (this month)      [Export CSV]            │
│          │  ┌──────────────────────────────────────────────────────┐ │
│          │  │ Tenant     Stored     Objects   Egress    Est. cost │  │
│          │  │ Folio Notes 1.1 TB   284k      3.4 GB    $64.10     │  │
│          │  │ Acme Inc   312 GB    48.9k     1.2 GB    $21.40     │  │
│          │  │ Kercis     4.9 GB    22k       90 MB     $4.10      │  │
│          │  │ ⚠ Legacy   2.4 GB    12k       12 MB     unknown    │  │
│          │  └──────────────────────────────────────────────────────┘ │
│          │  Total est. $87.64 · egress figures estimated (±10%)      │
│          │  A tenant with no rate card reports "unknown".            │
└──────────┴───────────────────────────────────────────────────────────┘
```

## Wireframe — mobile (populated)

```
┌────────────────────────────┐
│ ☰ ⬡    Platform ▾    👤    │
├────────────────────────────┤
│ Install usage  [This mo ▾] │
│ ┌────────────────────────┐ │
│ │  2.4 TB stored   ↑8%   │ │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │  $87.64 est     ↑12%   │ │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │ Total stored           │ │
│ │   ╭─╮        ╭──╮      │ │
│ │ Mar 1      Mar 31      │ │
│ └────────────────────────┘ │
│ Cost by tenant             │
│ Folio Notes  $64.10    ▸   │
│  1.1 TB · 3.4 GB egress    │
│ Acme Inc    $21.40    ▸    │
│  312 GB · 1.2 GB egress    │
│ ⚠ Legacy     unknown  ▸    │
│  2.4 GB · no rate card     │
│ ─────────────────────────  │
│ egress estimated (±10%)    │
│ Total est. $87.64          │
└────────────────────────────┘
```

## Interactions

- **Period control** is global (chart + stats + table together). Presets 24h/7d/30d/This month/custom.
- **Cost by tenant** is the load-bearing table — default sort by Est. cost desc, so the top row is the answer to "who's burning money".
- **The `unknown` row** (no rate card on the tenant's backend) renders "unknown" in `color.warning`, with a "no rate card" hint — and the admin-usage footer states the rule once so the whole table's semantics are legible (`usage` has the same rule; `admin-backend-detail` is where the gap is fixed).
- **Row click** → `admin-tenant-detail` (from there, the tenant's audit/usage deepens). "See raw events" for a tenant also routes through the detail — this page is attribution, not investigation.
- **Export CSV**: exports the current filter/period; caption carries the range + "estimates" note.
- **Egress** everywhere is footnoted estimated (±10%, `ERD.md` usage-snapshots note).
- **Empty state**: "No usage data yet · Snapshots are taken hourly" — a fresh install, not an error.
- **A11y**: chart has a data-table fallback for screen readers; stat cards are `dl` structure, not bare divs.
