# Task group — 16 screen: uploads

**Depends on:** `02-design-tokens`, `03-routing`, `04-core-infra` (SQLite queue store). **Blocks:** `25-flows` (upload flow). **Mockup:** [`mobile/mockup/uploads.md`](../mockup/uploads.md).

- [ ] **Layout — populated (mixed states)** per the mockup: rows for uploading (type icon, middle-truncated name, **determinate** progress bar in `color.primary`, byte count + rate), paused (⏸), completed (✓), failed (⚠ + reason + Retry). Header "3 uploading · 2 done · 1 failed".
- [ ] **Layout — empty** — "Nothing uploading · Pick a file from the Files tab to start an upload." with a Go to Files action.
- [ ] **Determinate, honest progress** — a stalled upload visibly stalls (`color.warning` paused state); never a fake bar (`design/style-guide.md` Motion).
- [ ] **Per-file retry** — re-queues only the failed file (a 30-file batch isn't lost to one failure, PRD MV-E1); the reason is the real one (quota/size/network).
- [ ] **Name collision** — inline replace/keep-both/cancel per file, queue stays open (PRD MV-E2).
- [ ] **Quota failure** — `quota_exceeded` reason + a `usage` link; other items continue.
- [ ] **Pause/resume per item**; backgrounded uploads resume where the platform allows.
- [ ] **Resume reconciliation** — on resume, the queue reconciles against `GET /v1/objects/:id/multipart/status` and re-sends **only missing parts** (PRD MB2, `architecture.md` §3.7). Both the local record and the server's are used (README's queue note).
- [ ] **Offline** — items show "Waiting for connection"; the queue is **persisted in SQLite**, so killing the app doesn't lose them.
- [ ] **Tap completed** → `file-detail`.
- [ ] **Footer caption** — "Uploads survive app restarts and connection loss" (the persistence guarantee is product copy, not a footnote).
- [ ] **Background upload** — Android foreground service while in flight (API 34 `FOREGROUND_SERVICE_DATA_SYNC`); iOS background processing task, **stated as best-effort in the UI** (README).

**tests:** queue survives kill+relaunch (SQLite); progress is determinate; resume re-sends only missing parts (mock server); per-file retry; collision replace/keep-both; offline waiting state.
