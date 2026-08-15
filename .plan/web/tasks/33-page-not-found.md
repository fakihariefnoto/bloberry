# Task group — 32 page: not-found

**Depends on:** `01-setup`, `02-design-tokens`, `03-routing`. **Blocks:** none (terminal). **Mockup:** [`web/mockup/not-found.md`](../mockup/not-found.md).

- [ ] **Catch-all `/*` route** — any unknown path renders this. Two auth variants: authenticated users get the `AppShell` with the 404 panel; unauthenticated get the bare public layout (logo mark) per the mockup.
- [ ] **Content** — decorative 404 numeral (`role="presentation"`), "That page doesn't exist.", the attempted path shown in `text.mono` muted (enough to see the typo, not loud enough to be an error wall), one action.
- [ ] **Action by auth state** — authenticated → "Back to Files" → `files`; unauthenticated → "Back to sign in" → `login` (a visitor who hit a dead `/s/xyz` gets login, not a dashboard they can't reach).
- [ ] **Style** — `color.text-muted` (a dead end, not a failure the app is signaling).
- [ ] **A11y** — the real heading is "That page doesn't exist." (the numeral is decorative); keyboard-reachable action.

**tests:** both auth variants render the right action; attempted path shown muted; decorative numeral not announced.
