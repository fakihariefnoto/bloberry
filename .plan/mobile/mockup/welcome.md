# Screen — welcome

## Purpose & context

- **User goal**: the landing choice point for an unauthenticated user — log in, or follow an invitation.
- **Entry points**: after `onboarding` (first launch) or as the root public route on subsequent launches; after logout.
- **Exit points**: Log in → `login`; "I have an invitation" → `accept-invitation` (deep link flow; on mobile the invite email's link goes straight there). From `reset-password` success, back to `login` — not here.
- **Data needed**: none — the app brand + two actions. No network.

## States

- [x] Single state (static; no loading/error — nothing loads). First-launch vs returning looks identical; the difference is that returning users land here directly without `onboarding`.

## Style reference

- **Components used**: brand mark, `text.display` tagline, primary button (Log in), secondary button (I have an invitation). No app bar — full-bleed brand screen.
- No token deltas.

## Wireframe — mobile

```
┌───────────────────────────┐
│                           │
│                           │
│                           │
│         ⬡                 │
│                           │
│       Bloberry            │
│  Storage for your         │
│  projects, one API over   │
│  any bucket.              │
│                           │
│                           │
│                           │
│  ┌──────────────────────┐ │
│  │      Log in          │ │
│  └──────────────────────┘ │
│                           │
│  ┌──────────────────────┐ │
│  │  I have an invitation│ │
│  └──────────────────────┘ │
│                           │
│  No public signup —       │
│  your tenant admin sends  │
│  the invitation.          │
└───────────────────────────┘
```

## Interactions

- **Log in** → `login` (primary, full width).
- **I have an invitation** → `accept-invitation` — but on mobile the practical path is the emailed invite link (deep link), so this button covers the "I have a token I want to enter manually" case.
- The footer caption states NG8 plainly — a visitor hunting for a Sign up button is told where it went instead of hitting a dead end.
- **A11y**: the tagline is a heading; both buttons are keyboard-focusable; the footer is `text.caption` `color.text-muted`, never visually louder than the actions.
