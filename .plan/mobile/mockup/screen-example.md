# Screen — <Screen Name>

## Purpose & context

- **User goal**: what the user is trying to accomplish here, one line.
- **Entry points**: which screen/action leads here.
- **Exit points**: what happens after each primary action — which screen next, or what state changes.
- **Data needed**: the actual fields/content this screen shows or collects — pull from `../../PRD.md` / `../../ERD.md` (the data model). Not placeholders — if the real fields aren't decided yet, that's an open question to flag, not something to paper over with "Heading Text."

## States

Every state this screen can actually be in, not just the happy path — check off what applies, delete what doesn't:

- [ ] Loading
- [ ] Empty (no data yet)
- [ ] Populated (happy path)
- [ ] Error
- [ ] Other domain-specific state (name it — e.g. "pending approval," "expired")

Draw a separate wireframe only for states that are *meaningfully different* in layout/content — a spinner over the same layout doesn't need its own box; an empty state with different copy and a different CTA does.

## Style reference

See [`../../design/style-guide.md`](../../design/style-guide.md) for colors, typography, and spacing tokens (shared with web). Only list exceptions for this screen here — leave empty if none.

- **Components used**: buttons, cards, nav bar, etc.

## Wireframe — populated (happy path)

Use this screen's *real* content — actual field labels, actual button copy, actual data shape from the Purpose & context section above. A wireframe full of generic placeholder text hasn't actually been thought through yet.

```
┌───────────────────────────┐
│  ← Back        Title      │
├───────────────────────────┤
│                           │
│   [ Image / Icon ]        │
│                           │
│   Heading Text            │
│   Supporting copy line    │
│   goes here.              │
│                           │
│  ┌──────────────────────┐ │
│  │      Primary CTA    │  │
│  └──────────────────────┘ │
│                           │
│  Secondary link           │
│                           │
├───────────────────────────┤
│  [Home] [Search] [You]    │
└───────────────────────────┘
```

## Wireframe — [other state, e.g. empty / error]

```
[Only include this block for states that are genuinely different in layout — delete this section entirely if the happy path covers everything worth drawing.]
```

## Interactions

- What each tappable element does, and where it navigates or what it triggers.
- Form validation rules, if this screen has a form (tie to the centralized validators from `templates/flutter-defaults.md`, don't restate the regex here).
- Loading/disabled behavior on the primary action, per the UX guard rules in `templates/flutter-defaults.md` (loading indicator + disabled-while-processing) — say what specifically happens on tap, don't assume it's obvious.

## Notes

Open questions, accessibility notes (tap target sizes, contrast), anything not obvious from the above.
