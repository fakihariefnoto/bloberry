# Task group — 14 page: file-detail

**Depends on:** `01-setup`, `02-design-tokens`, `03-routing`, `05-core-infra`, `06-shared-components` (CopyableCode, StatusPill, Breadcrumbs, ConfirmDestructive, Toast). **Blocks:** `33-flows` (share, move, delete flows). **Mockup:** [`web/mockup/file-detail.md`](../mockup/file-detail.md).

- [ ] **Layout — desktop** per the mockup: back + breadcrumbs, `PageHeader` with the file name + size/type/modified line, the action row (Download / Share / Move / Delete), the public-object warning banner, Details card (file_id with copy, storage backend + health, hash, uploader, timestamps), Sharing card (short URL + hits + revoke, signed links + revoke, "New signed link", Make public/Unpublish), Permissions panel.
- [ ] **Layout — mobile**: stacked single column, actions first (Download/Share/Move/Delete), details below.
- [ ] **Loading** — skeleton detail (never a spinner).
- [ ] **404 / soft-deleted** — "This file no longer exists" state with Back to files; soft-deleted shows the trash banner with Restore / Permanently delete (S5).
- [ ] **Download** — streams via the backend's signed path (redirect for cloud drivers, proxy for disk — `architecture.md` §3.2); a dead link surfaces the human-readable message, never raw provider XML.
- [ ] **Share** → opens `share-dialog` pre-filled (see `33-flows`).
- [ ] **Move** → `move-picker`; the `file_id` survives the move (PRD M4). 
- [ ] **Delete** → `ConfirmDestructive`; soft-delete default; a public file's delete warns the public URL dies immediately.
- [ ] **Make public** → `ConfirmDestructive` variant — public is effectively irreversible once the URL is copied (TRD R11); the confirm states that consequence. Unpublish is the reverse and needs no typed-name.
- [ ] **CopyableCode** — `file_id`, short URL, signed link all copy with confirmation; the short/signed copy the full URL.
- [ ] **Permissions panel** — "Manage permissions" → `grant-dialog`; explains the file inherits grants from its deepest ancestor folder (allow-only, most-specific-wins, `backend/domains.md` §5.2).
- [ ] **Audit link** — "who touched this" → `audit?target=<file_id>` (the graph edge added in mockup closure).
- [ ] **Permission-aware** — a viewer without write sees Move/Delete/Share/Unpublish **disabled with a tooltip reason**; Download stays enabled.

**tests:** 404 vs soft-deleted states; public banner on a public file; make-public confirm consequence; permission-disabled actions; audit link carries `?target=`.
