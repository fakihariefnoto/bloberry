# Task group — 28 testing

**Depends on:** the screens it tests (`05`–`25`), `26-flows.md`. Named high-risk flows, not blanket "add tests."

- [ ] **Auth widget tests** — login (invalid-credentials identical message, button disabled-while-processing), OTP (rate-limit, attempt cap), reset (expired-token state, no form), invite (new-user vs existing-email), **pair-login (single-use token, foreign-QR hint, permission-denied)**. Covers `07`–`12`.
- [ ] **Biometric unlock** — prompt success returns in place; failure escalates PIN → password; 5-attempt PIN cap; `MainActivity` is `FlutterFragmentActivity` (`13`, `04`).
- [ ] **Upload queue** — survives kill+relaunch (SQLite); resume re-sends only missing parts (mock server); per-file retry; collision replace/keep-both; offline waiting (`17`, `26.1`).
- [ ] **Revoke-a-key flow end-to-end** — confirm shows last-used; typed-name refuses a wrong value; last-active-key warning; no Undo toast (`24`, `26.6`).
- [ ] **QR pairing end-to-end** — scan → verify → `files`; a second scan of the same token fails; foreign QR stays live (`09`, `26.3`).
- [ ] **Permission-aware rendering** — a viewer's disabled actions show the caption, not an error wall (PRD MV4) on `files`/`file-detail` (`14`, `15`).
- [ ] **Envelope + snake_case** — model fields populate (not all-null, the silent-break test); `messages[].code` branching; never-raw errors (`04`).
- [ ] **Overflow safety** — a 200-char filename middle-truncates without overflow (guard rule 3) on `files`/`uploads` (`14`, `17`).
- [ ] **Rate-card rule** — unknown-not-$0 on `usage` (`22`).
- [ ] **Golden strings** — snackbar/empty-state/status-pill copy asserted exactly, so copy changes are reviewable.

**tests:** this group IS the test plan — each checkbox maps to a widget/integration test under `flutter test`.
