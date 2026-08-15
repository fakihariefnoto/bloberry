# Task group — 15 screen: search

**Depends on:** `02-design-tokens`, `03-routing`, `04-core-infra`. **Blocks:** none. **Mockup:** [`mobile/mockup/search.md`](../mockup/search.md).

- [ ] **Pre-query state** — recent searches + suggested files the user can read (never a blank waiting screen, per the mobile-search pattern).
- [ ] **Field** — autofocus with keyboard up, back in the leading position.
- [ ] **Debounce ~300ms** — results replace recent/suggested content entirely (no stale mix); loading is a subtle inline indicator, not a blocking spinner.
- [ ] **Result row** — type icon + name + **parent folder path** (a hit without context is "found it but where") + size/date.
- [ ] **Empty** — "No results for 'hero' · Try a different name, type or date" (real copy, never "No data").
- [ ] **Query too short** (<2 chars) — hint caption, no request fired.
- [ ] **Back** → `files`, query cleared; search state not persisted across tabs.
- [ ] **A11y** — results announce count ("2 results"); rows are ≥48dp targets.

**tests:** pre-query shows recents; debounce fires once; result shows parent path; empty copy; short-query hint.
