# Screen — link-expired

## Purpose & context

- **User goal**: open a shared link and find out, plainly, that it's dead. This is the `/s/:slug` 410 render (`TRD.md` error code `link_expired`; `architecture.md` §3.3) — a **public, human-readable HTML page**, deliberately not a JSON envelope, because the consumer is a person in a chat window (PRD MV-E3).
- **Entry points**: `/s/<slug>` where the link is expired, revoked, or the slug is unknown/never existed.
- **Exit points**: none actionable (the share is gone). The page is terminal.
- **Data needed**: none — the reason is derived from the link state server-side.

## States

- [x] Expired (the TTL elapsed)
- [x] Revoked (someone revoked the link before expiry)
- [x] Unknown slug (410 either way — don't reveal whether a slug ever existed)

## Style reference

- **Components used**: none beyond a logo mark and body text. This is the one screen that deliberately breaks the product's interior chrome — no sidebar, no nav, no "sign in" (the viewer may have no account at all).
- `color.warning`/`color.text-muted` only — never `color.error`, which would imply the site is broken rather than the link being dead.

## Wireframe — desktop (expired)

```
┌──────────────────────────────────────────────────────────────────┐
│                                                                  │
│                             ⬡                                    │
│                                                                  │
│                     This link has expired                        │
│                                                                  │
│                 The link you opened was set to expire.           │
│                 Ask whoever shared it to send a new one.         │
│                                                                  │
│                          ─ ─ ─ ─ ─                               │
│                                                                  │
│                    Powered by  ⬡ Bloberry                        │
└──────────────────────────────────────────────────────────────────┘
```

## Wireframe — mobile (revoked)

```
┌────────────────────────────┐
│                            │
│            ⬡               │
│                            │
│      This link was revoked │
│                            │
│   The owner stopped sharing│
│   this file.               │
│                            │
│   If you still need it,    │
│   ask them to send a new   │
│   one.                     │
│                            │
│        ─ ─ ─ ─ ─           │
│                            │
│   Powered by ⬡ Bloberry    │
└────────────────────────────┘
```

## Interactions

- **No interactions** — this page has nothing to click. That is the point: a dead link should be a dead end, not a funnel into signup.
- Copy differs by state: expired vs revoked vs "this link isn't working" (unknown slug). The distinction is honest — expired and revoked are both worth stating, since they imply different asks ("send a new one" vs "the owner deliberately stopped this").
- Never reveal whether a slug existed (`backend/domains.md` error table) — an attacker probing `/s/` must not learn which slugs are real.
- **A11y**: single `main` landmark, the logo is `role="presentation"`, and the "Powered by" line links to the install's landing page for anyone curious — a small, legitimate exit.
