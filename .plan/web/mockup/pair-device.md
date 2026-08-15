# Screen — pair-device

## Purpose & context

- **User goal**: set up another device — the mobile app (by scanning a QR) or a desktop machine (by importing an encrypted login file) — without typing a password there (PRD M22, M23).
- **Entry points**: user menu → "Pair a device" (any authenticated user). The QR path requires the *current* session: the token is minted from it, so the QR signs the phone in as this user.
- **Exit points**: back → previous page. Scan complete → the phone is signed in (this page just shows "paired" + the token dies). File download → the file goes to the desktop's import flow.
- **Data needed**: a one-time pairing token (Redis, ~2 min TTL, single-use — `backend/domains.md` §4.8); a server-signed config payload (import-window claim + refresh token — §4.9). No new collections (`ERD.md` audit note).

## States

- [x] Loading (token fetch)
- [x] Populated (QR active + config-file card)
- [x] QR expired (auto-refresh with notice)
- [x] QR paired (scanned once → "paired" + token dies)
- [x] Config export — passphrase entry + validation
- [x] Error (token issue failed, download failed)

## Style reference

- **Components used**: `PageHeader`, two `Card`s, QR render (token payload), `FormField` (passphrase), `CopyableCode`-style download, `Toast`. `elevation.none` cards, `radius.lg`.
- No token deltas. Both cards fit one page — no tab switching needed.

## Wireframe — desktop (populated, both cards)

```
┌──────────┬───────────────────────────────────────────────────────────┐
│ (sidebar)│  Acme Inc ▾         Jane Doe ▾                            │
│          │  Pair a device                                            │
│          │  Sign in on another device without typing a password.     │
│          │  ┌──────────────────────────┐  ┌────────────────────────┐ │
│          │  │  📱  MOBILE APP          │  │  🖥  DESKTOP APP        │ │
│          │  │                          │  │                        │ │
│          │  │  [   QR code here   ]    │  │  Download a login file │ │
│          │  │  Scan with the Bloberry  │  │  to import on a new    │ │
│          │  │  app: Login → Scan QR.   │  │  desktop.              │ │
│          │  │                          │  │                        │ │
│          │  │  ⚠ This code signs you   │  │  Passphrase (for the   │ │
│          │  │  in as Jane Doe. It      │  │  file's encryption):   │ │
│          │  │  expires in 1:42.        │  │  ┌──────────────────┐  │ │
│          │  │  [Refresh]               │  │  │ •••••••••••••    │  │ │
│          │  │  [Paired ✓] (after scan) │  │  └──────────────────┘  │ │
│          │  │                          │  │  Min 12 chars. The     │ │
│          │  └──────────────────────────┘  │  passphrase never      │ │
│          │                             │  │  leaves this browser.  │ │
│          │                             │  │  [Download .bloberry]  │ │
│          │                             │  │  Usable for 24h. The   │ │
│          │                             │  │  session stays         │ │
│          │                             │  │  revocable via logout. │ │
│          │                             │  └────────────────────────┘ │
└──────────┴───────────────────────────────────────────────────────────┘
```

## Wireframe — mobile (populated)

```
┌────────────────────────────┐
│ ☰ ⬡    Acme Inc ▾   👤     │
├────────────────────────────┤
│ ‹ Pair a device            │
│ ─────────────────────────  │
│ 📱 MOBILE APP              │
│ ┌────────────────────────┐ │
│ │      [QR code]         │ │
│ │ Scan in the app:       │ │
│ │ Login → Scan QR.       │ │
│ │ ⚠ Signs you in as Jane.│ │
│ │ expires in 1:42        │ │
│ │ [Refresh]  [Paired ✓]  │ │
│ └────────────────────────┘ │
│ ─────────────────────────  │
│ 🖥 DESKTOP APP              │
│ Download a login file to   │
│ import on a new desktop.   │
│ Passphrase (for encryption)│
│ ┌────────────────────────┐ │
│ │ •••••••••••••          │ │
│ └────────────────────────┘ │
│ Min 12 chars. Never leaves │
│ this browser.              │
│ [Download .bloberry]       │
│ Usable for 24h. Session    │
│ stays revocable.           │
└────────────────────────────┘
```

## Wireframe — QR expired (desktop card only)

```
┌──────────────────────────┐
│  📱  MOBILE APP          │
│  [  QR expired —        ]│
│  [  refreshing…         ]│
│  ✓ New code ready (1:59) │
│  ⚠ This code signs you   │
│  in as Jane Doe.         │
└──────────────────────────┘
```

## Interactions

- **QR card**: on load, `POST /v1/auth/pair/issue` → render the payload as a QR. A live countdown shows the remaining TTL (~2 min). On expiry, auto-refresh (the old token is dead even if unscanned — `pair_invalid`). **Manual Refresh** mints a fresh one. Once the phone exchanges it, the card flips to "Paired ✓" with a note "the code has been used — no longer valid" (a used code must not still render as scannable).
- **Capability warning is permanent copy**, not a dismissible toast: "⚠ This code signs you in as Jane Doe. It expires in X:XX." A screenshot of this QR is a login credential for its TTL — the UI says so.
- **Config-file card**: passphrase entry with **min 12 chars** + a strength hint (TRD R14). On download: server issues the signed payload, the browser derives the AES key client-side and encrypts — **the passphrase never transits the server** (ADR-14). The button label + footer state the 24h import window and the revocability.
- **Wrong passphrase / expired file** are desktop-import concerns; this page only *issues* — the failure copy lives in the desktop import UI (`desktop/` first-run).
- **A11y**: the QR is an `<img>`/canvas with the payload as `aria-label`; the countdown announces once at expiry, not on every tick; both cards are keyboard-reachable.

**tests:** QR token is single-use (scan → `pair_invalid` on a second try); expiry auto-refreshes; the paired state no longer renders a scannable code; config download requires the passphrase strength; the passphrase is never sent to the server (network-inspector check).
