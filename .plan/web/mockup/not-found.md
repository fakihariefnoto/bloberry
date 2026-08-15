# Screen — not-found

## Purpose & context

- **User goal**: hit an unknown path and get back to the app with one action. The catch-all `/*` route (`web/navigation.md`). Serves any auth state (public users get it too — the Go router's SPA catch-all renders the same shell and this route).
- **Entry points**: any path that matches no route — typos, stale links, API paths typed into the browser.
- **Exit points**: "Back to Files" → `files` (authenticated) or "Back to sign in" → `login` (unauthenticated). Both replace the dead route.
- **Data needed**: none. The attempted path is shown muted (helpful for typo diagnosis) but never echoed as a loud string.

## States

- [x] Single state, two auth variants (the action button differs; the layout is identical).

## Style reference

- **Components used**: `AppShell` for authenticated users (sidebar visible — a logged-in user who hits a bad path still has their navigation), a bare centered panel for unauthenticated. `color.text-muted` — a 404 is a dead end, not a failure the app is signaling.
- No token deltas.

## Wireframe — desktop (authenticated)

```
┌──────────┬──────────────────────────────────────────────────────────┐
│ (sidebar)│  Acme Inc ▾         Jane Doe ▾                           │
│          │                                                          │
│          │            404                                           │
│          │                                                          │
│          │    That page doesn't exist.                              │
│          │    /filesx/abc123 isn't a route.                         │
│          │                                                          │
│          │    [  Back to Files  ]                                   │
└──────────┴──────────────────────────────────────────────────────────┘
```

## Wireframe — mobile (unauthenticated)

```
┌────────────────────────────┐
│                            │
│          ⬡  Bloberry       │
│                            │
│          404               │
│                            │
│      That page doesn't     │
│      exist.                │
│      /s/xyz isn't a link.  │
│                            │
│      [  Back to sign in  ] │
└────────────────────────────┘
```

## Interactions

- **Back to Files / Back to sign in**: the only action, varying by auth state (an unauthenticated user who hits a dead `/s/xyz` gets `login`, not a dashboard they can't reach).
- The attempted path is shown in `text.mono`, muted — enough to see the typo, not so loud it reads as an error wall.
- **A11y**: `role="presentation"` on the big 404 numeral (it's decorative — the heading is the real "That page doesn't exist" text); keyboard reachable action.
- Unauthenticated variant reuses the public-layout chrome (logo mark), not the sidebar shell — matching what `login`/`link-expired` look like pre-auth.
