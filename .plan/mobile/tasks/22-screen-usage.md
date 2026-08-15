# Task group — 21 screen: usage

**Depends on:** `02-design-tokens`, `03-routing`, `04-core-infra`, `06-shared-chart` (fl_chart). **Blocks:** none. **Mockup:** [`mobile/mockup/usage.md`](../mockup/usage.md).

Per the mobile analytics pattern: swipeable stat cards, one full-width primary chart, breakdown below.

- [ ] **Layout — populated** per the mockup: PageHeader "Usage" + period segmented control, swipeable stat cards (stored / est. cost, with sparkline + direction), the primary chart (stored bytes over time), egress + objects lines, quota bar, est.-cost card with the rate-card breakdown.
- [ ] **Layout — no rate card** — "⚠ Cost unknown · No rate card configured for s3-eu-prod. Ask the platform admin to add one." in `color.warning`, never $0 (PRD M18).
- [ ] **Period control** (segmented 24h / 7d / 30d / This month) — global: cards, chart, breakdown all respond together.
- [ ] **Stat cards swipeable** — never push the chart below the fold; tapping a card opens that metric's full chart + history (per the mobile analytics pattern).
- [ ] **Egress estimated** — every egress figure carries the ⓘ "(±10%)" note (`ERD.md` usage-snapshots).
- [ ] **Rate-card provenance** — the cost line names the backend whose card produced it.
- [ ] **Loading** — skeleton stat cards + flat chart placeholder (never a "0" flash or stale data).
- [ ] **Empty** — "No usage data for this period" (a fresh tenant's first hour).
- [ ] **A11y** — chart has a data-point summary; stat cards announce label + value + direction.

**tests:** period control drives everything; unknown-not-$0; egress label; skeleton-not-zero; empty state.
