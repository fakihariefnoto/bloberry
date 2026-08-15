# Web — Bloberry

**Vue 3 + Vite + TypeScript + Tailwind CSS + Reka UI**, package-managed and built with **Bun**.

Reka UI is the unstyled headless primitives layer (Vue's Radix equivalent); all styling is Tailwind-driven, wrapped into the app's own components in `src/components/ui/`. Picked over a themed kit because Bloberry's data-dense surfaces (tables, dialogs, dropzones) live and die on the style-guide's exact tokens — a themed kit's defaults would fight `design/style-guide.md` everywhere it matters.

This is Bloberry's administration surface: the file browser, folder tree, upload, sharing, application/access-key management, member management, usage and audit log. It is also — unchanged — the UI the Wails desktop app ships, so nothing here may assume a browser-only capability without a fallback.

Follows [`templates/web-defaults.md`](../../../templates/web-defaults.md). Use `detail-web` to settle the Nuxt-vs-plain-Vue question, the env separation, the navigation shell and the reusable-component review.

## Topology

**Embedded** (see [`../TRD.md`](../TRD.md) Stack decisions). This app builds to static assets that the Go backend embeds via `go:embed` and serves from the same binary — there is no Node server at runtime, and Bun is a build-time dependency only. Practical consequences:

- **Plain Vue + Vite is the expected pick, not Nuxt.** Nuxt wants its own Node server, which contradicts the single-binary story. Confirm in `detail-web`.
- Client-side routing needs a catch-all fallback in the Go router so a deep link like `/files/abc123` doesn't 404.
- The API is same-origin, so no CORS config and cookies work naturally — but the SDKs' consumers *are* cross-origin, so the backend still needs CORS for them.
- `../architecture.md` §7 records which platform's setup task owns the project `init`; this folder does not scaffold its own repo.

## What's specific to this dashboard

- **Data-dense by nature.** File tables, key lists and audit logs are the core surfaces — `space.sm` row padding and `size.row-height` 44px, not comfortable consumer spacing. See [`../design/style-guide.md`](../design/style-guide.md).
- **Responsive is required, not optional.** Tenant admins check usage and revoke keys from a phone. Every mockup here carries both a mobile-width and a desktop-width wireframe.
- **Upload is the hard component.** Drag-and-drop with folder structure preserved, presigned-PUT direct to storage, multipart for large files, a persistent queue panel that survives navigation, real progress (never a fake indeterminate crawl), and per-file retry.
- **Secrets appear exactly once.** A newly issued access key is shown in a modal with a copy button and never again; everywhere else it's masked to the last 4 characters.
- **Permission-aware rendering.** A viewer sees the same layout with blocked actions visibly disabled and a reason, not a wall or a surprise error toast.

Colour, type, shape and component specs are shared with mobile — see [`../design/style-guide.md`](../design/style-guide.md), don't redefine them here.

## Files

- `mockup/` — one `.md` per page, each with mobile-width and desktop-width wireframes
- `navigation.md` — the authoritative route graph and route table, once `detail-web` has run
- `components.md` — the shared-component inventory, once `detail-web` has run
- `Makefile` — install/dev/build/test/lint commands, converted to Bun from `templates/makefiles/web.mk`. Targets may need adjusting once the real project layout exists.
- `tasks/` — the implementation task list, once `build-web` has run
