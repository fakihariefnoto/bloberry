# Task group — 21 page: usage

**Depends on:** `01-setup`, `02-design-tokens`, `03-routing`, `05-core-infra`, `06-shared-components` (PageHeader, DateRangePicker, ByteSize, RelativeTime). **Blocks:** none. **Mockup:** [`web/mockup/usage.md`](../mockup/usage.md).

Per the analytics-dashboard pattern (`web-screen/patterns.md`): **one primary chart**, stat cards above with sparklines + direction, breakdown below.

- [ ] **Layout — desktop** per the mockup: PageHeader "Usage", global DateRangePicker (This month default), three stat cards (stored / objects / egress, each with sparkline + direction), one primary area chart (stored bytes over time) with the est. cost callout, the rate-card line + "estimated (±10%)" footnote, the weekly breakdown table, the quota bar.
- [ ] **Layout — mobile**: stat cards swipeable (never push the chart below the fold), chart full-width, breakdown stacked.
- [ ] **Period control is global** — chart, stat cards, breakdown, cost all respond together. Presets 24h/7d/30d/This month/custom.
- [ ] **Loading** — skeleton stat cards at real size + flat chart placeholder — never a "0" flash (`web-screen/patterns.md`; stale-data chart is banned).
- [ ] **The rate-card rule** — no rate card → "Cost unknown — no rate card configured for s3-eu-prod" in `color.warning`, **never $0.00** (PRD M18: a zero reads as free, which is a lie).
- [ ] **Rate-card provenance** — the cost line always names the backend whose card produced it.
- [ ] **Egress estimate** — every egress figure carries the ⓘ "estimated (±10%) — downloads redirect past Bloberry" (`ERD.md` usage-snapshots note, ADR-3).
- [ ] **Breakdown** — rows expand to per-day values; "see raw events" → `audit` (usage is metering, audit is who).
- [ ] **Quota bar** — mirrors `files` but with cost context; ≥95% shows the projected overage cost.
- [ ] **Empty state** — "No usage data for this period" (a brand-new tenant's first hour), not a blank chart.
- [ ] **Permission-aware** — `member`/`viewer` never see the route; `tenant_admin` sees usage + rate card but cannot edit the card (that's platform-admin).

**tests:** period control drives all widgets; unknown-cost vs $0; egress label; skeleton-not-zero on load; empty state for a fresh tenant.
