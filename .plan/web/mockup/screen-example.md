# Screen — <Screen Name>

## Purpose & context

- **User goal**: what the user is trying to accomplish here, one line.
- **Entry points**: which page/action leads here.
- **Exit points**: what happens after each primary action — which page next, or what state changes.
- **Data needed**: the actual fields/content this page shows or collects — pull from `../../PRD.md` / `../../ERD.md` (the data model). Not placeholders — if the real fields aren't decided yet, that's an open question to flag, not something to paper over with "Heading Text."

## States

Every state this page can actually be in, not just the happy path — check off what applies, delete what doesn't:

- [ ] Loading
- [ ] Empty (no data yet)
- [ ] Populated (happy path)
- [ ] Error
- [ ] Other domain-specific state (name it)

Draw a separate wireframe only for states that are *meaningfully different* in layout/content — a spinner over the same layout doesn't need its own box; an empty state with different copy and a different CTA does. Combine this with the mobile/desktop split below only where it actually matters (a loading state usually looks the same at both widths; don't draw four boxes when two would do).

## Style reference

See [`../../design/style-guide.md`](../../design/style-guide.md) for colors, typography, and spacing tokens (shared with mobile). Only list exceptions for this screen here — leave empty if none.

- **Components used**: nav bar, cards, buttons, etc.

## Wireframe — mobile (populated / happy path)

Use this page's *real* content — actual field labels, actual button copy, actual data shape from the Purpose & context section above.

```
┌───────────────────────────┐
│ ☰   Logo          🔍 👤   │
├───────────────────────────┤
│   Heading Text            │
│   Supporting copy line    │
│   goes here.              │
│                           │
│  ┌──────────────────────┐ │
│  │      Primary CTA    │  │
│  └──────────────────────┘ │
│                           │
│  [Card]  [Card]           │
│  [Card]  [Card]           │
└───────────────────────────┘
```

## Wireframe — desktop (populated / happy path)

```
┌─────────────────────────────────────────────────────────────┐
│  Logo     Nav Item   Nav Item   Nav Item        🔍  👤      │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   Heading Text                     [Card] [Card] [Card]     │
│   Supporting copy line             [Card] [Card] [Card]     │
│   goes here.                                                │
│   ┌────────────────┐                                        │
│   │  Primary CTA   │                                        │
│   └────────────────┘                                        │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

## Wireframe — [other state, e.g. empty / error]

```
[Only include mobile/desktop blocks for states that are genuinely different in layout — delete this section entirely if the happy path covers everything worth drawing.]
```

## Interactions

- What each clickable element does, and where it navigates or what it triggers.
- Form validation rules, if this page has a form (tie to the centralized validators from `templates/web-defaults.md`/backend defaults, don't restate the regex here).
- Loading/disabled behavior on the primary action — say what specifically happens on click, don't assume it's obvious.
- Breakpoint behavior beyond mobile/desktop, if there's a meaningful tablet/intermediate layout.

## Notes

Open questions, accessibility notes (focus order, contrast, keyboard navigation), anything not obvious from the above.
