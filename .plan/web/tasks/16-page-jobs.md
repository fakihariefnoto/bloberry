# Task group — 16 page: jobs

**Depends on:** `01-setup`, `02-design-tokens`, `03-routing`, `05-core-infra`, `06-shared-components` (DataTable, StatusPill, JobProgress, RelativeTime, Toast). **Blocks:** `33-flows` (extract, delete flows). **Mockup:** [`web/mockup/jobs.md`](../mockup/jobs.md).

- [ ] **Layout — desktop** per the mockup: PageHeader "Jobs", filter tabs (All / Running / Failed / Done), DataTable (job, kind, started, state, `⋮`), expandable running/failed rows with the progress bar and failure reason.
- [ ] **Layout — mobile**: stacked cards.
- [ ] **`?state=` filter** — Running is the default view when arriving from a toast (lands the user on what they were told about).
- [ ] **Live polling** — running rows poll `GET /v1/jobs/:id`; a state flip (running→done/failed) animates in place with the StatusPill change, no reload. **Polling pauses when the tab is hidden** (respects the user and the server).
- [ ] **Determinate progress only** — `progress_done`/`progress_total` drive the bar; never a fake indeterminate crawl (`design/style-guide.md` Motion; PRD M21).
- [ ] **Failure detail** — `failure_message` human-readable, `failure_code` copyable (`text.mono`); the extraction failure line **always** states "Target folder unchanged" (PRD AP-E2).
- [ ] **Retry** — only on failed jobs; re-enqueues with a fresh `attempts` counter. **Terminal error codes offer no Retry** (`archive_rejected` — a zip bomb isn't going to un-bomb; retrying on loop is a self-inflicted DoS, TRD R6) — only "View folder".
- [ ] **Bundle success** — result shows the generated archive size + Copy URL.
- [ ] **subtree_delete** — progress counts objects ("X of Y objects") so the scale is legible after the fact (PRD TA-E1).
- [ ] **View folder** → `files` at the job's target (the graph edge added in closure).
- [ ] **Empty states** — no jobs → "Nothing queued · Extraction and large deletes will appear here"; failed filter → "No failures · Clean run". Distinct.
- [ ] **A11y** — `role="progressbar"` with `aria-valuenow`; state changes announced via a polite live region, not an alert.

**tests:** polling flips state without reload; hidden-tab polling stops; determinate progress values; terminal-code rows have no Retry; target-folder-unchanged line on extraction failure.
