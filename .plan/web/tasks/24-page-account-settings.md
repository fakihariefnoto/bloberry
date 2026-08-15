# Task group — 24 page: account-settings

**Depends on:** `01-setup`, `02-design-tokens`, `03-routing`, `05-core-infra`, `06-shared-components` (PageHeader, FormField, ConfirmDestructive-lite). **Blocks:** none. **Mockup:** [`web/mockup/account-settings.md`](../mockup/account-settings.md).

- [ ] **Layout — desktop** per the mockup: PageHeader "Account settings", four sections — PASSWORD, **TWO-FACTOR AUTHENTICATION**, ACTIVE SESSIONS, NOTIFICATIONS.
- [ ] **Layout — mobile**: stacked sections.
- [ ] **2FA enable (M24)** — opens the provisioning modal: the `otpauth_url` QR + secret shown **exactly once** (MV-E5), a confirming-code field, Enable disabled until the code validates; on success **10 backup codes shown once** with the "you won't see these again" warning (the `SecretRevealModal` discipline).
- [ ] **2FA disable** — requires a current TOTP code or an unused backup code (R15 — never a bare click).
- [ ] **Regenerate backup codes** — invalidates the old set immediately; the count shows `X remaining` so a user near zero regenerates first.
- [ ] **Change password** — current + new + confirm; wrong current → field error (no hint about how wrong); on success **every other session is revoked server-side** (`backend/domains.md` §4.4 semantics); the caption under Sessions states this *before* the user commits.
- [ ] **Active sessions** — listed from the refresh-token store (device/platform/location/last-active); **Revoke** ends one session; **Revoke all** ends every other with a plain confirm; the current browser is marked "(this)" and can't be revoked.
- [ ] **Notifications** — save-on-toggle (persist immediately, subtle toast); preferences on `user_settings` (`ERD.md`), no tenant implication.
- [ ] **No "log out of this device" here** — the user menu's Log out is the single obvious exit (a second one invites confusion).
- [ ] **Permission-aware** — personal, no role gating.

**tests:** password change invalidates other sessions (server-verified); wrong-current error; revoke-all confirms; "(this)" session not revocable; toggle persistence; **2FA secret shown once + confirm-before-enable; disable requires a code/backup; regenerate invalidates the old set; backup-code count updates**.
