# Shared components — Bloberry web

The reusable-component inventory, so 26 pages don't produce 26 variants of the same table. Visual specs live in [`../design/style-guide.md`](../design/style-guide.md); this file records **what exists, where it's used, and what it must handle** — the parts a style guide doesn't cover.

**Stack:** Vue 3 + Vite + TypeScript, Tailwind, **Reka UI** (unstyled headless primitives — behavior only, no styles — wrapped into our own Tailwind components in `src/components/ui/`), **TanStack Table** for anything tabular, Bun for install/build. Icons: **Lucide** only.

---

## Why this list matters more than usual here

Six of Bloberry's pages are fundamentally *a table of things with row actions and a filter bar* — files, shares, applications, members, audit, jobs, plus four admin pages. Built independently they will diverge in row height, empty state, loading behaviour, selection semantics and pagination, and the app will read as ten different products. `DataTable` below exists specifically to prevent that.

---

## Layout

| Component | Used by | Must handle |
|---|---|---|
| **`AppShell`** | every authenticated route | Sidebar (260px) ↔ nav rail (72px) ↔ mobile drawer. Role-driven section visibility. Persists collapse state per user. |
| **`TenantSwitcher`** | `AppShell` header | Current tenant always visible. Switching **always lands on `files` at the new root** (`navigation.md`). Single-tenant users see a static label, not a dropdown with one item. |
| **`UserMenu`** | `AppShell` footer | Avatar, name, role badge, profile / account / log out. |
| **`PageHeader`** | every interior page | Title (`text.heading-lg`), optional description, right-aligned primary action slot, optional tabs. One component so 20 pages don't each invent their own top spacing. |
| **`Breadcrumbs`** | `files`, `file-detail` | **The primary navigation in the file browser** now that there's no tree pane. Collapses to `Root / … / Parent / Current` beyond 4 levels, ellipsis opens a dropdown of elided segments. Each segment is a drop target for drag-to-move. |

---

## Data display

### `DataTable`

The single most reused component. One implementation over TanStack Table, with Reka UI's table primitives only where they earn their keep (virtualized rows, focus management); most of the grid itself is plain Tailwind markup, since Reka's table is unstyled anyway.

**Used by:** `files`, `shares`, `jobs`, `applications`, `application-detail` (keys), `members`, `audit`, `admin-tenants`, `admin-backends`, `admin-usage` — **ten pages.**

**Must handle:**
- **Virtualized rows.** A folder can hold 10,000 objects (PRD G3). Rendering them all is not an option; neither is a 20-per-page pager for a file browser.
- **Server-side pagination, sort and filter**, since the data lives in MongoDB and the counts are large. Cursor-based, not offset — offset pagination over a collection being written to skips and repeats rows.
- **Four states, never improvised per page:** skeleton rows at `size.row-height` on first load; a distinct empty state per table (see `EmptyState`); an inline error banner with retry, leaving the table area blank rather than showing stale rows; and normal.
- **Row selection** (checkbox column) with a **bulk-action bar** that replaces the filter bar when anything is selected. Needed by `files` (bulk download/delete/move) and `shares` (bulk revoke).
- **Row actions** visible on hover for pointer input, and always present in an overflow menu — never hover-only, which is unusable on touch.
- **Permission-aware rows.** An action the principal can't perform is rendered disabled with a reason tooltip, not hidden (PRD MV4).
- **Sticky header** and a horizontally scrollable body — the page itself must never scroll horizontally.

### `FileTypeIcon`

Extension → Lucide icon + tint. **Defined once here and mirrored in Flutter** (`design/style-guide.md` says the mapping is shared). Folders, images, video, audio, archives, documents, code, and a generic fallback. Getting this inconsistent between web and mobile is small and very visible.

### `StatusPill`

`radius.full`, tinted per the semantic mapping in the style guide. One component, five meanings: key state (active/expiring/revoked), backend health, job state, object visibility (**public is a warning, not a success**), and member role.

### `ByteSize` / `RelativeTime`

Trivial, and worth extracting precisely because they're trivial — ten pages formatting bytes three different ways (`1.5 GB` / `1.5 GiB` / `1536 MB`) is a real and common inconsistency. `ByteSize` uses binary units with a decimal-unit tooltip, since storage providers bill in decimal and users compare against the invoice.

### `CopyableCode`

`text.mono` on `color.surface` with a copy button and a confirmation. Used for `file_id`, storage keys, short URLs, signed links, SDK snippets, CLI commands. **Masks by default** when the value is a secret, showing only the last 4 characters.

### `EmptyState`

Icon, headline, body, optional action. **Every table passes its own copy** — Bloberry has at least eight distinct empties (empty folder, no search results, no shares yet, no applications, no keys on this application, no audit events in range, no jobs, no tenants) and they must say different things. A shared "No data" is a bug, not a default.

---

## Files & upload

### `Dropzone`

Wraps the whole `files` table area, not just a small box — dropping onto the file list is what people expect. Drag-over state per the style guide. Rejects with a reason that reverts after 3s.

### `UploadQueue`

Docked bottom-right, `elevation.md`, collapsible, **persists across navigation** (which is why it's not a route — see `navigation.md`). Per item: type icon, middle-truncated name preserving the extension, determinate progress, byte count and rate, per-file retry.

**Progress is real, never decorative.** A stalled upload must visibly stall (`design/style-guide.md`) — a fake indeterminate bar makes a dead upload indistinguishable from a slow one, which is precisely when a user needs to know.

Handles the three upload paths (`architecture.md` §3.1): presigned PUT (default), multipart for large files with per-part progress, and direct as a fallback. Fires `beforeunload` while anything is in flight.

### `FolderTree`

**Scoped to the `move-picker` dialog only.** With the breadcrumbs-only browser layout there is no persistent tree pane, so this is the one surface that still needs hierarchy. Lazy-loads children on expand; disables the moved node and all its descendants as drop targets, which is how PRD TA-E2's cycle is prevented in the UI rather than only at the API.

---

## Forms & dialogs

| Component | Used by | Notes |
|---|---|---|
| **`ConfirmDestructive`** | folder delete, key revoke, member removal, backend change, make-public | Two tiers: a plain confirm for reversible actions, and **typed-name confirmation** for irreversible ones. States real consequences — "10,342 objects" not "this folder" (PRD TA-E1). |
| **`PermissionPicker`** | `grant-dialog`, key creation | Folder-subtree selector + permission checkboxes + optional expiry. **Must explain allow-only semantics inline** (PRD D7) — anyone with IAM habits will look for a deny and needs to be told there isn't one. |
| **`SecretRevealModal`** | `key-created` | Shows a secret exactly once. Explicit "you will not see this again" in `color.warning`, copy button, and dismissal requiring acknowledgement — **not** a backdrop click, which is how people lose keys. |
| **`FormField`** | every form | Label / input / helper / error, wired to the style guide's focus/error/disabled states. |
| **`DateRangePicker`** | `audit`, `usage`, `admin-usage` | Presets (24h, 7d, 30d, custom). |

---

## Feedback & state

| Component | Notes |
|---|---|
| **`Toast`** | Snackbar per the style guide. **Undo where the operation is reversible** (move, soft delete) and **no Undo where it isn't** (key revoke) — the absence is deliberate signal. |
| **`PermissionDenied`** | Not an error page. Renders the surface with actions disabled and a `text.caption` line naming what's needed ("Requires `write` on this folder"). A viewer should see the boundary, not a wall. |
| **`ErrorBoundary`** | Catches render errors per route, offers reload, never shows a stack trace. |
| **`JobProgress`** | Inline progress for extraction / bundling / subtree delete, polling `GET /v1/jobs/:id`. Used in `jobs` and as a toast-attached mini-view. |

---

## Data layer

| Module | Responsibility |
|---|---|
| **`lib/api/`** | **Generated from `api/openapi.yaml`** (`architecture.md` ADR-11) — not hand-written. Regenerated by `make generate`. |
| **`lib/envelope.ts`** | Unwraps `{data?, messages?}`, maps `messages[].code` to i18n keys, and **never surfaces raw provider or exception text**. One place, so 26 pages can't each invent error handling. |
| **`stores/`** | Pinia. `auth` (session, principal, role), `tenant` (current tenant, switcher), `upload` (the queue — survives navigation, which is why it's a store and not component state), `ui` (sidebar collapse, table density). |
| **`lib/permissions.ts`** | Client-side mirror of the resolver's *shape* for enabling/disabling UI. **Advisory only** — the server is the authority (`backend/domains.md` §5). Worth stating in the file itself, because a client-side permission helper is exactly the thing someone later mistakes for enforcement. |

---

## Environments

Per `templates/web-defaults.md`, three separated environments. Since the frontend is embedded in the Go binary (`architecture.md` ADR-1), "environment" here means **build-time** configuration only — there is no runtime Node server to hold secrets, and **nothing secret may ever live in these files** since they compile into shipped assets.

| File | `VITE_API_BASE_URL` | Notes |
|---|---|---|
| `.env.development` | `http://localhost:8080/v1` | Vite dev server proxies to the local Go binary |
| `.env.staging` | `/v1` | Same-origin — the binary serves both |
| `.env.production` | `/v1` | Same-origin |
| `.env.example` | — | Committed. The others are not. |

Staging and production are same-origin because of the embedded topology, which removes CORS from the dashboard's concerns entirely. The **SDKs** are cross-origin and do need CORS on the backend — a distinction worth keeping straight, since "we don't need CORS" is true of the dashboard and false of the product.

The desktop app is the exception: it loads the same bundle but points at a **user-configured server URL**, so the API base must be overridable at runtime there rather than baked in. `lib/api/` reads `window.__BLOBERRY_API_BASE__` when present and falls back to the build-time value.
