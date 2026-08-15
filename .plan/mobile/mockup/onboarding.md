# Screen — onboarding

## Purpose & context

- **User goal**: a one-time, 3-pane intro to what Bloberry does — folders, sharing, access keys — shown only on genuine first launch (`navigation.md`: a persisted flag, shown once ever).
- **Entry points**: cold first launch. Never again after the flag is set.
- **Exit points**: Skip/Get started → `welcome`; swiping through all 3 panes reaches Get started on the last.
- **Data needed**: none — static marketing content. The tenant name isn't known yet (the user hasn't authenticated).

## States

- [x] Slide 1 — folders (the reason the app exists)
- [x] Slide 2 — sharing
- [x] Slide 3 — access keys / control
- [x] Loading — not applicable; static content, no network. (No skeleton needed.)

## Style reference

- **Components used**: onboarding carousel per `design-collection/mobile-screen/patterns.md` — progress dots, Skip top-right, Next/Get started bottom, `color.primary` primary button.
- Mobile width only (`size.appbar-height` free, no app bar on this screen — the illustration is the hero).
- No token deltas.

## Wireframe — slide 1 (folders)

```
┌───────────────────────────┐
│                   Skip →  │
│                           │
│      ┌─────────────┐      │
│      │   📁 📁     │      │
│      │  📁  ┌──┐   │      │
│     │      │📁 │   │      │
│      │      └──┘   │      │
│      └─────────────┘      │
│                           │
│  Real folders, not keys   │
│  Your files live in a     │
│  tree — create, rename,   │
│  move. It behaves like    │
│  a filesystem.            │
│                           │
│        ●  ○  ○            │
│                           │
│  ┌──────────────────────┐ │
│      │        Next │      │
│  └──────────────────────┘ │
└───────────────────────────┘
```

## Wireframe — slide 3 (access keys)

```
┌───────────────────────────┐
│                   Skip →  │
│                           │
│      ┌─────────────┐      │
│      │  🔑 key     │      │
│      │  blob_live_ │      │
│      │  ••••4f2a   │      │
│      └─────────────┘      │
│                           │
│  One key per application  │
│  Apps authenticate with   │
│  scoped keys — grant one  │
│  folder, never the tenant │
│  unless you mean to.      │
│                           │
│        ○  ○  ●            │
│                           │
│  ┌──────────────────────┐ │
│      │   Getstarted│      │
│  └──────────────────────┘ │
└───────────────────────────┘
```

## Interactions

- **Skip** (top-right, always available) jumps straight to `welcome` and persists the flag — onboarding is never a wall.
- **Next** advances the pane; the label becomes **Get started** on the final slide. Progress dots show position (current filled).
- **Get started** → `welcome`, flag persisted.
- **Swipe** advances/retreats panes; dots are not tappable (they indicate position, they don't navigate).
- **A11y**: each slide is a full-width element with its headline as a heading; dots announce "slide 2 of 3"; Skip and Get started are reachable by keyboard focus order.
