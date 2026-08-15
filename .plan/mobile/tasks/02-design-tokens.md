# Task group — 02 design tokens

**Depends on:** `01-setup.md`. **Blocks:** every screen file. **Ordering is the point** — screens built before this file exists accumulate hardcoded values.

**Source:** `design/style-guide.md` — the single shared spec. Everything here maps 1:1 to a token in that file. **Mobile uses the system font** (SF Pro/Roboto), not Inter — the one documented mobile/web delta.

- [ ] **Color — light and dark.** All 15 tokens (`background`, `surface`, `surface-raised`, `border`, `text`, `text-muted`, `primary`, `on-primary`, `primary-subtle`, `accent`, `error`, `success`, `warning`, `disabled`, `on-disabled`) as Flutter `Color` constants in a token file, light + dark variants. The two `on-primary` values differ per mode (light `#FFFFFF`, dark `#17161D`) — the contrast caveat is a constraint.
- [ ] **Semantic storage-state mapping** — public = warning, private = muted, active/healthy = success, expiring/pending = warning, revoked/failed/over-quota = error — one mapping, used by every screen (the style-guide §Colors table).
- [ ] **Type scale** — `text.display` 29/36 bold, `heading-lg` 24/32 bold, `heading-md` 20/28 semibold, `heading-sm` 17/24 semibold, `body` 14/20, `body-strong` 14/20 semibold, `button` 14/20 semibold, `label` 12/20 medium, `caption` 12/20, `mono` 13/20 — with line-heights on the 4px grid. System font family.
- [ ] **Spacing** — `space.xs` 4 · `sm` 8 · `md` 16 · `lg` 24 · `xl` 32 · `2xl` 48. Screen padding `space.md`; `space.sm` row padding for the file list.
- [ ] **Radius** — `sm` 8 · `md` 12 · `lg` 20 · `full` 9999 (Soft language).
- [ ] **Elevation** — `none` (border), `sm`, `md`, `lg` — Soft language, border-first in dark mode.
- [ ] **Sizing / touch targets** — `touch-min` 48×48dp mobile, `button-height` 48, `button-height-sm` 36 (pointer-first only), `input-height` 48, `icon` 24, `icon-sm` 20, `avatar` 40, `appbar-height` 56, `navbar-height` 64, `row-height` 44.
- [ ] **Motion** — `fast` 150ms ease-out · `base` 250ms ease-in-out · `slow` 400ms ease-in-out; respect the OS reduce-motion setting with an instant change.
- [ ] **Component specs** — button (primary/secondary/destructive/ghost) with pressed/disabled/loading states (loading = spinner replaces label, width held); input with focus/error/disabled; card; data-row; dropzone; upload-queue item; status pill; code/secret display (masked); nav (bottom tabs active state); snackbar; empty state; skeleton loader; `ConfirmDestructive` typed-name.
- [ ] **Screens consume tokens only** — follow-on rule enforced in review: no hardcoded colors/radii/sizes in screen files.

**tests:** a token-coverage lint that fails if a screen file references a raw hex/color literal outside the theme file; light/dark contrast pairs match the style-guide AA table.
