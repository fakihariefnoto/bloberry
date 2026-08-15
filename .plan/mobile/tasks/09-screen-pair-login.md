# Task group — 09 screen: pair-login

**Depends on:** `02-design-tokens`, `03-routing`, `04-core-infra` (camera + envelope). **Blocks:** `26-flows` (QR pair chain). **Mockup:** [`mobile/mockup/pair-login.md`](../mockup/pair-login.md).

- [ ] **Camera overlay per the mockup** — full-bleed preview with a scan-frame overlay, back button, caption under the frame. No app bar; the camera is the canvas.
- [ ] **`mobile_scanner` wired** — camera live on entry (no tap-to-start); decodes `bloberry://pair/<token>` payloads.
- [ ] **Verify flow** — a decoded payload triggers `POST /v1/auth/pair/verify`; subtle spinner over the frame while in flight; the frame disables further scans until the result lands (`backend/domains.md` §4.8).
- [ ] **Success** → `files` at the tenant root (mobile-TTL session; the token is single-use and dead).
- [ ] **`pair_invalid`** (expired / used) — the expired state ("This code is no longer valid · Codes expire in 2 minutes and work once · Refresh it on the web dashboard") with a "Scan again" action returning to the live scanner.
- [ ] **Not-a-Bloberry code** — the hint state ("That doesn't look like a Bloberry login code"), scanner keeps running (no error toast, MB-E2).
- [ ] **Camera permission denied** — inline explanation + "Open Settings" (the `files` FAB pattern, reused); not a dead end.
- [ ] **Back** → `login`, cancelling any in-flight scan.
- [ ] **A11y** — the frame is `role="img"` with an accessible label; the expired/not-a-code states are announced; the caption states the alternative (password login is one tap back).

**tests:** scan → verify → `files`; a second scan of the same token fails (`pair_invalid`); a foreign QR shows the hint and stays live; permission-denied shows the Open Settings path.
