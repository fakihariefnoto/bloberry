# Task group — 20 screen: settings

**Depends on:** `02-design-tokens`, `03-routing`, `04-core-infra` (biometric plumbing). **Blocks:** `25-flows` (biometric chain). **Mockup:** [`mobile/mockup/settings.md`](../mockup/settings.md).

- [ ] **Layout per the mockup** — SECURITY (PIN/biometric unlock toggle + subtitle), NOTIFICATIONS (upload failures, key revoked alerts), GENERAL (Language select, Default tenant select), ABOUT (version, Log out).
- [ ] **PIN/biometric toggle** — off by default (`templates/flutter-defaults.md`); on enable → PIN-setup flow or the platform biometric prompt (per the subtitle); **disabled with "Biometric unavailable on this device — PIN still works"** when no sensor.
- [ ] **Toggles save immediately** — subtle toast, no Save button; a failed save reverts with an inline error.
- [ ] **Language / Default tenant** — selects with right-aligned current value, tap → bottom sheet.
- [ ] **Log out** — bottom, visually separated; revokes the refresh token server-side; lands on `welcome`.
- [ ] **A11y** — toggles carry text labels (not color alone); on/off announced on change.

**tests:** toggle persists; biometric-unavailable disables with the right caption; a failed save reverts; log out revokes + lands on welcome.
