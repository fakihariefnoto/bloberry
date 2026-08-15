# Task group — 06 shared components

**Depends on:** `02-design-tokens.md`, `05-core-infra.md` (DataTable needs the envelope). **Blocks:** the pages that use them — `DataTable` before the ten table pages, `ConfirmDestructive` before the destructive-action pages, etc.

**Source:** `web/components.md` (the full inventory) + `design/style-guide.md` (visual specs). Each component is built **once** and referenced by page tasks — never copy-pasted per page. A checkbox per component with its must-handle list from `components.md`.

## Layout

- [ ] **`AppShell`** — sidebar 260px ↔ nav rail 72px ↔ mobile drawer; role-driven section visibility; persists collapse state (from `03-routing.md`; shell skeleton and routing live there, this is the component's behavior).
- [ ] **`TenantSwitcher`** — current tenant always visible; switching always lands on `files` at the new root; single-tenant users see a static label, not a one-item dropdown.
- [ ] **`UserMenu`** — avatar, name, role badge, profile / account / log out.
- [ ] **`PageHeader`** — title (`text.heading-lg`), optional description, right-aligned primary action slot, optional tabs. One component so 20 pages don't each invent top spacing.
- [ ] **`Breadcrumbs`** — the file browser's primary navigation. Collapses to `Root / … / Parent / Current` beyond 4 levels, ellipsis opens a dropdown of elided segments. Every segment is a drop target for drag-to-move.

## Data display

- [ ] **`DataTable`** — the single most reused component, one implementation over TanStack Table (Reka's table primitives only where they earn their keep). Serves **ten pages**: `files`, `shares`, `jobs`, `applications`, `application-detail`, `members`, `audit`, `admin-tenants`, `admin-backends`, `admin-usage`. Must handle: virtualized rows (10,000-object folders); server-side cursor pagination/sort/filter (never offset); four states (skeleton rows at `size.row-height` / distinct empty per table / inline error banner leaving the table blank / normal); row selection + bulk-action bar; row actions hover-on-pointer + always in an overflow menu (never hover-only); permission-aware disabled rows with reason tooltip; sticky header, no page-level horizontal scroll.
- [ ] **`FileTypeIcon`** — extension → Lucide icon + tint, matching the shared web/mobile mapping (`design/style-guide.md`). Folders, images, video, audio, archives, documents, code, generic fallback.
- [ ] **`StatusPill`** — `radius.full`, one component, five meanings: key state, backend health, job state, object visibility (**public is a warning, not a success**), member role.
- [ ] **`ByteSize` / `RelativeTime`** — binary units with decimal-unit tooltip (providers bill decimal); one `ByteSize` so ten pages don't format three different ways. `RelativeTime` ("2h ago", "Mar 12").
- [ ] **`CopyableCode`** — `text.mono` on `color.surface`, copy button + confirmation. Masks by default when the value is a secret (last-4 only).
- [ ] **`EmptyState`** — icon + headline + body + optional action. Every table passes its own copy — Bloberry has at least eight distinct empties and they must say different things. A shared "No data" is a bug.

## Files & upload

- [ ] **`Dropzone`** — wraps the whole `files` table area; drag-over state per style guide; rejected (size/type) with a reason reverting after 3s.
- [ ] **`UploadQueue`** — docked bottom-right, `elevation.md`, collapsible, **persists across navigation** (why it's a store, not a route). Per item: type icon, middle-truncated name preserving extension, **determinate** progress, byte count + rate, per-file retry. Fires `beforeunload` while anything is in flight. Handles the three upload paths (presigned PUT default, multipart for large files with per-part progress, direct fallback) per `architecture.md` §3.1.
- [ ] **`FolderTree`** — **move-picker dialog only** (no persistent tree pane). Lazy-loads children; disables the moved node + all descendants as drop targets (PRD TA-E2 cycle prevention). Width `size.tree-width`.

## Forms & dialogs

- [ ] **`ConfirmDestructive`** — plain confirm for reversible actions; **typed-name confirmation** for irreversible ones. States real consequences ("10,342 objects" not "this folder", PRD TA-E1). `elevation.lg`.
- [ ] **`PermissionPicker`** — folder-subtree selector + permission checkboxes + optional expiry. **Explains allow-only semantics inline** (PRD D7 — anyone with IAM habits looks for a deny).
- [ ] **`SecretRevealModal`** (`key-created`) — shows a secret exactly once; "you won't see this again" in `color.warning`; copy button; dismissal requires acknowledgement, **no backdrop click**.
- [ ] **`FormField`** — label / input / helper / error, wired to the style guide's focus/error/disabled states; uses the centralized validators.
- [ ] **`DateRangePicker`** — presets (24h, 7d, 30d, custom); used by `audit`, `usage`, `admin-usage`.

## Feedback & state

- [ ] **`Toast`** — snackbar per style guide. **Undo where reversible** (move, soft delete), **none where not** (key revoke) — the absence is deliberate signal. Errors = `color.error` icon/accent stripe, not full background.
- [ ] **`PermissionDenied`** — not an error page; renders the surface with actions disabled and a `text.caption` naming what's needed ("Requires `write` on this folder").
- [ ] **`ErrorBoundary`** — catches render errors per route, offers reload, never a stack trace.
- [ ] **`JobProgress`** — inline determinate progress, polls `GET /v1/jobs/:id`; used in `jobs` and as a toast-attached mini-view.

**tests:** a `DataTable` test per its four states + selection/bulk-bar + permission-disabled rows; a `ConfirmDestructive` typed-name test (refuses wrong name, confirms right); `UploadQueue` persistence-across-navigation test; `SecretRevealModal` no-backdrop-dismiss test.
