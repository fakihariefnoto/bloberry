# Task group — 12 page: link-expired

**Depends on:** `01-setup`, `02-design-tokens`, `03-routing`. **Blocks:** none (terminal page). **Mockup:** [`web/mockup/link-expired.md`](../mockup/link-expired.md).

- [ ] **The `/s/:slug` 410 render** — public HTML page, deliberately **not** a JSON envelope (the consumer is a person in a chat window, PRD MV-E3). Served by the Go router for expired/revoked/unknown slugs.
- [ ] **Three copy variants** by link state: expired ("The link you opened was set to expire"), revoked ("The owner stopped sharing this file"), unknown ("this link isn't working"). Honest distinction — expired and revoked imply different asks.
- [ ] **Never reveals whether a slug existed** — unknown slug renders identically-shaped copy (`backend/domains.md` error table; an attacker probing `/s/` must not learn which slugs are real).
- [ ] **Style** — no sidebar, no nav, no "sign in" (the viewer may have no account). `color.warning`/`text-muted` only — never `color.error` (a dead link isn't a broken site). "Powered by ⬡ Bloberry" footer links to the install landing page.
- [ ] **A11y** — single `main` landmark, logo `role="presentation"`.

**tests:** each state renders its copy; unknown-slug copy does not differ observably from expired/revoked in a way that leaks existence; HTML content-type.
