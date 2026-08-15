# Build index — Web

Build order, dependency edges, and status for the web dashboard task list. Source of truth for what's done.

## Build order

| File | Covers | Status |
|---|---|---|
| `01-setup.md` | Project init inside the existing repo, Bun, Vite, Vue 3 + TS, Tailwind, Reka UI, deps, Makefile | ☐ |
| `02-design-tokens.md` | Tailwind theme + CSS custom properties from `design/style-guide.md`, every scale + component states | ☐ |
| `03-routing.md` | Router, 27 routes with typed params, AppShell (sidebar/rail/drawer), auth gates + redirects | ☐ |
| `04-environments.md` | `.env.{development,staging,production}`, `.env.example`, build scripts per target | ☐ |
| `05-core-infra.md` | HTTP client + envelope, snake_case mapping, auth token handling, runtime API-base hook | ☐ |
| `06-shared-components.md` | 23 shared components from `web/components.md`, each built once | ☐ |
| `07-page-login.md` … `33-page-not-found.md` | One file per page (27), grounded in that page's mockup states/interactions | ☐ |
| `34-flows.md` | Cross-page flow chains from `web/navigation.md` §Flow chains (9 flows) | ☐ |
| `35-testing.md` | Named high-risk-flow tests + golden-output discipline | ☐ |

## Dependency edges

- **`01-setup` blocks everything** — no file before it compiles.
- **`02-design-tokens` blocks every page** (and `06-shared-components`) — pages built before tokens exist accumulate one-off values (`p-[13px]`, hex literals). This ordering is deliberate (`build-web` step 5).
- **`03-routing` blocks every page** — pages need the router and shell to render in context.
- **`05-core-infra` blocks the data-dependent pages** (`files`, `file-detail`, `shares`, `jobs`, `applications`, `application-detail`, `members`, `audit`, `usage`, `tenant-settings`, all admin pages) and the auth pages' submit paths.
- **`06-shared-components` blocks the pages that use them** — `DataTable` before the ten table pages, `ConfirmDestructive` before `files`/`application-detail`/`members`/`tenant-settings`, etc.
- **`34-flows` depends on the pages it chains** — write it after 07–33 so its references resolve.
- **`35-testing` depends on the pages it tests.**

## External edges (from `architecture.md` §7)

- **`01-setup` Depends on `backend/tasks/01-setup.md`** — the Go module + repo already exist; this file creates a *package* at `web/`, not a repo. No `git init`.
- **Build order is `openapi codegen → web build → go build`** (`architecture.md` §7, edge 2). The root `Makefile`'s `build` target must run the web build **before** the backend build — a backend build that runs first embeds the *previous* frontend. **This is an edge, not a nicety**: two independent CI jobs shipping the backend first is exactly how a release ships stale UI.
- **`web/dist` is consumed twice** — by the server (`internal/platform/web/embed.go`) and by the desktop shell. Both depend on this build. Same edge as above.
- **Desktop depends on `web/tasks/01-setup.md` + this folder's output** — `desktop/tasks/` must not run `wails3 init`; it consumes `web/dist`.

## Gaps flagged

None. All 26 pages have real mockups with Interactions sections; the style guide is complete (every scale + component states); the navigation map exists with route table and flow chains.
