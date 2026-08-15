# Task group — 25 page: pair-device

**Depends on:** `01-setup`, `02-design-tokens`, `03-routing`, `05-core-infra` (auth session), `06-shared-components` (FormField, Toast). **Blocks:** none. **Mockup:** [`web/mockup/pair-device.md`](../mockup/pair-device.md).

The "set up another device" screen — two cards on one page: mobile QR pairing (PRD M22) and desktop encrypted config (PRD M23).

- [ ] **Layout — desktop + mobile** per the mockup: PageHeader "Pair a device", subtitle, two `Card`s side by side (stacked on mobile): 📱 MOBILE APP and 🖥 DESKTOP APP.
- [ ] **QR card — issue on load** — `POST /v1/auth/pair/issue` from the authenticated session; render the `bloberry://pair/<token>` payload as a QR.
- [ ] **Live countdown** — remaining TTL (~2 min); on expiry **auto-refresh** with a notice (the old token is dead even unscanned — `pair_invalid`).
- [ ] **Manual Refresh** — mints a fresh token.
- [ ] **Paired state** — once the phone exchanges the token, the card flips to "Paired ✓" with "the code has been used — no longer valid"; a used code must not still render scannable.
- [ ] **Capability warning is permanent copy** — "⚠ This code signs you in as Jane Doe. It expires in X:XX." (a screenshot of the QR is a login credential for its TTL — ADR-13). Never a dismissible toast.
- [ ] **Config-file card — passphrase entry** — min 12 chars + a strength hint (TRD R14).
- [ ] **Config download** — server issues the signed payload (`POST /v1/auth/config/issue`), the **browser encrypts client-side** with a passphrase-derived key (argon2id → AES-GCM); **the passphrase never transits the server** (ADR-14). Download a `.bloberry` file.
- [ ] **Config-file copy states the facts** — 24h import window, session stays revocable via logout (DT-E2), passphrase never leaves the browser.
- [ ] **Error states** — token issue failed, download failed → inline + retry.
- [ ] **A11y** — the QR has an `aria-label` with its payload; the countdown announces once at expiry, not per tick; both cards keyboard-reachable.

**tests:** QR token is single-use (scan → `pair_invalid` on a second try); expiry auto-refreshes; the paired state no longer renders a scannable code; config download enforces the passphrase strength; a network-inspector check confirms the passphrase is never sent to the server.
