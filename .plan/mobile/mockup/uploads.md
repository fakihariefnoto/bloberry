# Screen — uploads

## Purpose & context

- **User goal**: see and manage in-flight uploads — progress, retry failures, resume after interruption, tap a completed item. The **permanent home of the queue** (a long-running operation needs a tab, not a transient sheet — `mobile/navigation.md` Shell notes).
- **Entry points**: tab tap; automatically shown after a pick (capture/file-picker returns here with items queued). Uploads continue in the background; on resume, this screen reconciles against the server (`GET multipart/status`) and re-sends only missing parts (PRD MB2, `architecture.md` §3.7).
- **Exit points**: tap a completed item → `file-detail`; retry a failure → re-queues in place; quota/size errors → shown here with the real reason.
- **Data needed**: the local upload queue (persisted, survives app kill) — per item: target folder, filename, byte count, progress, state (queued/uploading/paused/completed/failed/waiting-for-connection), error reason.

## States

- [x] Empty (nothing queued)
- [x] Populated (mixed: uploading/completed/failed)
- [x] All-complete
- [x] Error — item failed with real reason (per-file, other files continue — PRD MV-E1)
- [x] Domain-specific — offline (items "Waiting for connection", persisted)
- [x] Domain-specific — name collision (replace/keep-both/cancel, PRD MV-E2)
- [x] Domain-specific — quota exceeded (failed with `quota_exceeded`, link to `usage`)

## Style reference

- **Components used**: list rows (type icon, middle-truncated name preserving extension, **determinate** progress bar in `color.primary`, byte count + rate), status per row, action per failed row. Progress is real, never decorative (`design/style-guide.md` Motion).
- Mobile is the primary queue surface (web has the docked `upload-queue` panel; this tab is its mobile counterpart).
- No token deltas.

## Wireframe — mobile (mixed states)

```
┌────────────────────────────┐
│  Acme Inc ▾        ⚙       │
├────────────────────────────┤
│  3 uploading ·2done·1failed│
│ ────────────────────────── │
│ ┌────────────────────────┐ │
│ │ 🎥 field_recording.mp4 │ │
│ │ ████████░░░░░░  57%    │ │
│ │ 1.2 GB · 4.8 MB/s      │ │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │ 🖼 IMG_0421.jpg    ⏸    │ │
│ │ ██████████░░░░  82%    │ │
│ │ 6.2 MB · paused        │ │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │ ✓ invoice-q3.pdf       │ │
│ │ completed · 2.1 MB     │ │
│ └────────────────────────┘ │
│ ┌────────────────────────┐ │
│ │ ⚠ big-catalog.zip      │ │
│ │ quota exceeded — 500 GB│ │
│ │ [Retry] [See usage]    │ │
│ └────────────────────────┘ │
│ ────────────────────────── │
│ Uploads survive app        │
│ restarts and connection    │
│ loss.                      │
├────────────────────────────┤
│  Files  ●Uploads  Shares  ▸│
└────────────────────────────┘
```

## Wireframe — mobile (empty)

```
┌───────────────────────────┐
│  Acme Inc ▾        ⚙      │
├───────────────────────────┤
│                           │
│       ⬆                   │
│                           │
│  Nothing uploading        │
│                           │
│  Pick a file from the     │
│  Files tab to start an    │
│  upload.                  │
│                           │
│  ┌──────────────────────┐ │
│  │   Go to Files        │ │
│  └──────────────────────┘ │
├───────────────────────────┤
│  Files ●Uploads  Shares  ▸│
└───────────────────────────┘
```

## Interactions

- **Progress is determinate and honest**: a stalled upload visibly stalls (`color.warning` paused state) — a fake bar makes a dead upload indistinguishable from a slow one (`design/style-guide.md` Motion).
- **Retry** re-queues only the failed file — a 30-file batch isn't lost to one failure (PRD MV-E1). The reason shown is the real one (quota/size/network).
- **Name collision** → inline replace/keep-both/cancel per file; the queue stays open.
- **Quota failure** shows the real code (`quota_exceeded`) + a `usage` link (PRD MV-E1); other items continue.
- **Pause/resume** per item (the `⏸` row); backgrounded uploads resume automatically where the platform allows; on resume the queue reconciles missing parts server-side (`architecture.md` §3.7) — resumption re-sends only missing parts, not the whole file.
- **Offline**: items show "Waiting for connection"; the queue is **persisted**, so killing the app doesn't lose them (PRD MB2).
- **Tap completed** → `file-detail`. A completed row's share action is reachable from detail.
- **The footer caption** ("Uploads survive app restarts and connection loss") states the persistence guarantee — the reassurance is part of the product, not a footnote.
- **Empty state** directs to the pick flow (the FAB lives on `files`), so the queue tab never dead-ends.
