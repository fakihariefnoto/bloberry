# Screen — pair-login

## Purpose & context

- **User goal**: sign in on this phone by scanning a QR shown on the web dashboard — no password typing (PRD M22, MB3).
- **Entry points**: "Scan QR code" from `login` (push, camera overlay).
- **Exit points**: scan + verify succeed → `files` (the phone's own session, mobile TTLs); back → `login`; permission denied → explanation + "Open Settings"; not-a-Bloberry code → ignored with a hint (MB-E2).
- **Data needed**: none typed — the camera decodes a `bloberry://pair/<token>` payload, exchanged via `POST /v1/auth/pair/verify` (`backend/domains.md` §4.8). One-time token, ~2 min TTL, rate-limited verify.

## States

- [x] Camera preview (scanning)
- [x] Permission denied (camera) — explanation + Open Settings
- [x] Not-a-Bloberry code — ignored with a hint, scanner stays live
- [x] Token expired / used (`pair_invalid`) — "This code is no longer valid — refresh it on the web"
- [x] Verify in flight (loading)
- [x] Success → `files`
- [x] Error (network)

## Style reference

- **Components used**: full-bleed camera preview with a scan-frame overlay (the frame is the only "chrome"), back button, a caption under the frame. `local_auth`-style minimal screen — no app bar, the camera is the canvas (map/location pattern's floating-chrome discipline applied to a scanner).
- No token deltas.

## Wireframe — mobile (scanning)

```
┌───────────────────────────┐
│ ◀  Back                   │
│  ┌─────────────────────┐  │
│  │                     │  │
│  │   ┌───────────────┐ │  │
│  │   │  ╔════════════╗ │  │
│  │   │  ║  [camera] ║ │   │
│  │   │  ║  scan here║ │   │
│  │   │  ╚════════════╝ │  │
│  │   └───────────────┘ │  │
│  │                     │  │
│  └─────────────────────┘  │
│  Align the QR from the    │
│  web dashboard inside     │
│  the frame.               │
│                           │
│  Sign in by scanning —    │
│  no password needed.      │
└───────────────────────────┘
```

## Wireframe — mobile (not a Bloberry code)

```
┌───────────────────────────┐
│ ◀  Back                   │
│  ┌─────────────────────┐  │
│  │                     │  │
│  │   ┌───────────────┐ │  │
│  │   │  ╔════════════╗ │  │
│  │   │  ║  [camera] ║ │   │
│  │   │  ║  scan here║ │   │
│  │   │  ╚════════════╝ │  │
│  │   └───────────────┘ │  │
│  └─────────────────────┘  │
│  That doesn't look like   │
│  a Bloberry login code.   │
│  Keep scanning — codes    │
│  come from the web        │
│  dashboard.               │
└───────────────────────────┘
```

## Wireframe — mobile (expired / used)

```
┌───────────────────────────┐
│ ◀  Back                   │
│      ⚠                    │
│  This code is no longer   │
│  valid.                   │
│                           │
│  Codes expire in 2        │
│  minutes and work once.   │
│  Refresh it on the web    │
│  dashboard and scan the   │
│  new one.                 │
│                           │
│  [  Scan again  ]         │
└───────────────────────────┘
```

## Interactions

- **Scan**: the camera is live on entry (no tap-to-start). A decoded `bloberry://pair/<token>` triggers `POST /v1/auth/pair/verify`. While verifying, a subtle spinner over the frame; the frame disables further scans until the result lands.
- **Success** → `files` at the tenant root (the phone now has its own mobile-TTL session; the pairing token is single-use and dead).
- **`pair_invalid`** (expired or already used) → the expired state, with "Scan again" returning to the live scanner. The copy names the fix (refresh on web) — the token can't be re-minted from the phone.
- **Not-a-Bloberry code** → the hint state; the scanner keeps running (no error toast, no modal — a wrong QR shouldn't stop the flow, MB-E2).
- **Camera permission denied** → inline explanation + "Open Settings" (the `files` FAB flow's pattern, reused here); not a dead end.
- **Back** → `login`, cancelling any in-flight scan.
- **A11y**: the frame is `role="img"` with an accessible label; the expired/not-a-code states are announced; the camera is never the only input (the caption states the alternative — password login is one tap back).

**tests:** scan → verify → `files`; a second scan of the same token fails (`pair_invalid`); a foreign QR shows the hint and stays live; permission-denied shows the Open Settings path; the expired state offers Scan again.
