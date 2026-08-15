# Screen — jobs

## Purpose & context

- **User goal**: watch queued/running work — archive extraction, bundle generation, large subtree deletes — and see why a job failed. The home of PRD AP-E2's "the folder is unchanged" promise: a failure shows the reason here, never a half-populated target folder.
- **Entry points**: sidebar Jobs; "Extraction queued"/"Delete started" toasts link here; `?state=` filter for CI/ops folks.
- **Exit points**: click a job → expands inline detail (or `job-detail` if it ever earns a route — v1 keeps it inline); Retry on a failed job → re-enqueues; the linked object/folder → `files` at that folder.
- **Data needed**: `jobs` — `kind` (extract/bundle/subtree_delete), `state` (queued/running/succeeded/failed), `progress_done`/`progress_total`, `failure_code`, `failure_message`, `attempts`, `created_at`/`started_at`/`finished_at`, payload/result. Poll `GET /v1/jobs/:id`.

## States

- [x] Loading (skeleton rows)
- [x] Empty (no jobs yet)
- [x] Populated (happy path, mixed states)
- [x] Failed (with real reason + retry)
- [x] Error (poll connection dropped)
- [x] Domain-specific — job finished while the page was open (state flips live, no refresh needed)

## Style reference

- **Components used**: `AppShell`, `DataTable`, `StatusPill` (queued=pending/warning, running=primary, succeeded=success, failed=error), `JobProgress` (determinate bar at the real rate — never fake indeterminate, `design/style-guide.md`), `RelativeTime`, `Toast`.
- Job rows are `size.row-height`; the progress bar lives inside the row's second line when expanded. Filter tabs: All / Running / Failed / Done.
- No token deltas.

## Wireframe — desktop (populated)

```
┌──────────┬───────────────────────────────────────────────────────────┐
│ (sidebar)│  Acme Inc ▾         Jane Doe ▾                            │
│          │  Jobs                                                     │
│          │  [All] [Running] [Failed] [Done]                          │
│          │  ┌──────────────────────────────────────────────────────┐ │
│          │  │ Job            Kind         Started      State    ⋮  │ │
│          │  │ ▸ Extract 2026  extract      now          25%     ⋮  │ │
│          │  │    ████████░░░░░░░░░ 47 / 184 files · 2 min left     │ │
│          │  │   Bundling release  bundle       03:12        done  ⋮│ │
│          │  │   Delete _old       subtree_del  03:05        done  ⋮│ │
│          │  │ ⚠ Extract marketing extract       02:58       failed⋮│ │
│          │  │    zip bomb: decompressed size exceeded 8 GB ceiling │ │
│          │  │    Target folder unchanged. [Retry] [View folder]    │ │
│          │  └──────────────────────────────────────────────────────┘ │
│          │  Showing 4 of 9                                           │
│          │                                                           │
│          │  Long-running jobs survive reloads and are safe to        │
│          │  close — status resumes from the server.                  │
└──────────┴───────────────────────────────────────────────────────────┘
```

## Wireframe — mobile (running + failed)

```
┌────────────────────────────┐
│ ☰ ⬡    Acme Inc ▾   👤     │
├────────────────────────────┤
│ Jobs   [All][Run][Fail]    │
│ ─────────────────────────  │
│ ▸ Extract 2026     25%     │
│   ████████░░░░░░░░         │
│   47/184 files · 2 min left│
│ ─────────────────────────  │
│ Bundling release    done   │
│   03:12 · 2.1 GB bundle    │
│ ─────────────────────────  │
│ ⚠ Extract marketing failed │
│   decompressed size        │
│   exceeded 8 GB ceiling    │
│   Target folder unchanged. │
│   [Retry] [View folder]    │
│ ─────────────────────────  │
│ Long-running jobs survive  │
│ reloads.                   │
└────────────────────────────┘
```

## Interactions

- **Live polling**: running rows poll `GET /v1/jobs/:id`; a state flip (running→done/failed) animates in place with the `StatusPill` change — no page reload. Polling pauses when the tab is hidden (respects the user, and saves the server).
- **Retry**: only on `failed` jobs; re-enqueues with a fresh `attempts` counter. Terminal error codes (`archive_rejected` — a zip bomb isn't going to un-bomb) offer **no Retry**, only "View folder" — retrying a bomb on loop is a self-inflicted DoS (TRD R6).
- **Extraction detail**: shows the payload summary ("Extract 2026/ into Projects/2026/") on expand; a successful extract links to the target folder.
- **Bundle**: result shows the generated archive's size and a Copy URL.
- **subtree_delete**: progress counts objects, and the confirm that started it is long gone — the row itself carries "X of Y objects" so the scale is legible after the fact (PRD TA-E1).
- **Errors**: `failure_message` is human-readable; `failure_code` is copyable (`text.mono`) for filing an issue. The extraction failure line **always** states "Target folder unchanged" (PRD AP-E2).
- **Filter tabs** map to `?state=`; Running is the default view when arriving from a toast so the user lands on the thing they were told about.
- **Empty states**: no jobs → "Nothing queued · Extraction and large deletes will appear here"; failed filter → "No failures · Clean run". Distinct.
- **A11y**: progress bars expose `role="progressbar"` with `aria-valuenow`; state changes are announced via a polite live region, not a jarring alert.
