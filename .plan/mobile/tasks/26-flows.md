# Task group — 26 cross-screen flows

**Depends on:** the screens it chains (`05`–`25`). **Blocks:** nothing (integration file).

The vertical-slice view a pure per-screen split misses. One checkbox group per flow from `mobile/navigation.md` §Flow chains, each listing the screens it touches and its cancel/failure landing.

## 26.1 Capture and upload

- [ ] `files` → FAB → `source-picker` → camera or file picker (OS) → returns to `uploads` with items queued → per-file progress → complete. Chain: `14` → `source-picker` sheet → `17`.
- [ ] **Permission denied** (camera/photos) → inline explanation with an "Open Settings" action, not a dead end.
- [ ] **Backgrounded mid-upload** → continues where the platform allows; on resume `uploads` reconciles missing parts (`17`).
- [ ] **Offline** → items stay queued with "Waiting for connection"; the SQLite queue survives a kill.
- [ ] **Quota exceeded** → that item fails with the real reason; the rest continue (PRD MV-E1).

## 26.2 Share from a phone

- [ ] `files` → tap file → `file-detail` → Share → `share-sheet` → signed link + TTL / short URL / make public → created → OS share sheet opens with the URL pre-filled. Chain: `14` → `15` → `share-sheet` → `18`.
- [ ] **Cancel** → back to `file-detail`, nothing created.
- [ ] **Make public** → `confirm-destructive` first (public is effectively irreversible once the URL has been copied — TRD R11).

## 26.3 QR pairing login (M22)

- [ ] `login` → "Scan QR code" → `pair-login` (camera) → scan a `bloberry://pair/<token>` payload → verify → `files`. Chain: `07` → `09` + the `04` camera/envelope plumbing.
- [ ] **Expired / used token** → `pair_invalid` state with "Scan again" (the token can't be re-minted from the phone — refresh on web, `backend/domains.md` §4.8).
- [ ] **Not-a-Bloberry code** → hint state, scanner stays live (MB-E2).
- [ ] **Permission denied** → Open Settings path.

## 26.4 Biometric unlock

- [ ] App resumed after the configured timeout → `unlock` over the shell → biometric prompt. Chain: `13` + the `04` biometric plumbing + `21` (toggle).
- [ ] **Success** → returns to exactly where the user was, not to `files`.
- [ ] **Failure / unavailable** → "Use password" → `login`.
- [ ] Toggle lives in `settings`, off by default.

## 26.5 Session expiry

- [ ] Any route → 401 → `login`, preserving the intended destination and returning to it after login. Chain: `03-routing` redirect + `07`.
- [ ] **Queued uploads survive** — they're presigned and go straight to the provider; only their `complete` call waits for re-login (`17`).

## 26.6 Revoke a key from a phone (the urgent case)

- [ ] `more` → `applications` → `application-detail` → key row → Revoke → `confirm-destructive` showing **when it was last used** (PRD TA-E3) → revoked. Chain: `19` → `23` → `24`.
- [ ] **No Undo** — key revocation is irreversible; the toast deliberately offers none.

**tests:** one end-to-end test per flow covering the happy path and its cancel/failure landing — 6 flows total.
