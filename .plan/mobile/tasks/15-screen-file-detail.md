# Task group — 14 screen: file-detail

**Depends on:** `02-design-tokens`, `03-routing`, `04-core-infra`. **Blocks:** `25-flows` (share flow). **Mockup:** [`mobile/mockup/file-detail.md`](../mockup/file-detail.md).

- [ ] **Layout — populated** per the mockup: back + `⋮`, file-type icon hero (no preview in v1, NG3), name/size/type/date, the public warning banner, metadata list (file_id with copy, storage + backend health, uploader, timestamps), Shared section (short URL + hits + revoke, signed link + TTL + revoke), sticky footer actions (Share / Download / ⋮).
- [ ] **Layout — file not found** — 🗑 "This file no longer exists" + Back to files.
- [ ] **Soft-deleted** — "This file is in the trash · Restore or permanently delete" (S5).
- [ ] **Loading** — skeleton detail, not a spinner.
- [ ] **Share** → `share-sheet`: signed link + TTL / short URL / make public; created → OS share sheet with the URL pre-filled; link lands in `shares`. Cancel → back, nothing created. Make-public goes through `confirm-destructive` first (effectively irreversible once copied, TRD R11).
- [ ] **Download** — streams via the backend's signed path (redirect for cloud, proxy for disk — `architecture.md` §3.2); a dead link shows the human-readable message, never raw provider error.
- [ ] **Move** → `folder-picker` (the one folder-tree surface on mobile; moved node + descendants disabled as targets, PRD TA-E2). The `file_id` survives (PRD M4).
- [ ] **Delete** → `confirm-destructive`; soft-delete default; a public file's delete warns the public URL dies.
- [ ] **Copy file_id** via the CopyableCode affordance.
- [ ] **Sticky footer** — Share (primary) always in thumb reach; ⋮ holds Move, Rename, Delete, Make public/unpublish.
- [ ] **Permission-aware** — a viewer without write sees the actions disabled with a caption; Download stays enabled.
- [ ] **A11y** — type icon `role="presentation"`; metadata a labeled list; actions ≥48dp.

**tests:** 404 vs soft-deleted states; share-sheet create/cancel; make-public confirm; permission-disabled actions; sticky footer present.
