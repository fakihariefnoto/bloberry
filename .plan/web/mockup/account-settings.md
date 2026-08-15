# Screen — account-settings

## Purpose & context

- **User goal**: manage the *account* — password, sessions, security. (Identity display/email lives in `profile`; tenant config in `tenant-settings` — three settings surfaces, each with one job.)
- **Entry points**: user menu → Account (any authenticated user).
- **Exit points**: save → toast; change password → success toast + all other sessions invalidated; log out → `login`.
- **Data needed**: no tenant data — the user's password (change), active sessions (from Redis refresh-token records), email/preferences for notification toggles.

## States

- [x] Loading (skeleton)
- [x] Populated
- [x] Error (wrong current password on change)
- [x] Domain-specific — active sessions list with per-session revoke
- [x] Domain-specific — 2FA off (enable button)
- [x] Domain-specific — 2FA on (backup-code count, disable + regenerate)
- [x] Domain-specific — 2FA provisioning modal (secret shown once, confirm code required — MV-E5)

## Style reference

- **Components used**: `PageHeader`, `FormField`, sectioned card, `ConfirmDestructive`-lite (session revoke is a plain confirm — the user is revoking their *other* sessions).
- No token deltas.

## Wireframe — desktop (populated)

```
┌──────────┬───────────────────────────────────────────────────────────┐
│ (sidebar)│  Acme Inc ▾         Jane Doe ▾                            │
│          │  Account settings                                         │
│          │  ┌──────────────────────────────────────────────────────┐ │
│          │  │  PASSWORD                                            │ │
│          │  │  Current password                                    │ │
│          │  │  ┌─────────────────────────────────────────────────┐ │ │
│          │  │  │ ••••••••••                                👁   │ │   │
│          │  │  └─────────────────────────────────────────────────┘ │ │
│          │  │  New password                                        │ │
│          │  │  ┌─────────────────────────────────────────────────┐ │ │
│          │  │  │ •••••••••••••                            👁    │ │   │
│          │  │  └─────────────────────────────────────────────────┘ │ │
│          │  │  At least 8 characters.                    [Save]    │ │
│          │  └──────────────────────────────────────────────────────┘ │
│          │  ┌──────────────────────────────────────────────────────┐ │
│          │  │  TWO-FACTOR AUTHENTICATION                           │ │
│          │  │  Status: Off                              [Enable]   │ │
│          │  │  Requires a code from your authenticator app at      │ │
│          │  │  every login — web, mobile and desktop.              │ │
│          │  │  ── (when enabled) ──                                │ │
│          │  │  Status: On ✓ · backup codes: 10remaining[Regenerate]│ │
│          │  │  [Disable]  (requires a current code or backup  code)│ │
│          │  └──────────────────────────────────────────────────────┘ │
│          │  ┌──────────────────────────────────────────────────────┐ │
│          │  │  ACTIVE SESSIONS                                     │ │
│          │  │  This browser · Web · San Francisco · now      (this)│ │
│          │  │  Chrome on Windows · Web · Jakarta · 2h ago [Revoke] │ │
│          │  │  Bloberry mobile · Mobile · Jakarta · 1d ago [Revoke]│ │
│          │  │                                          [Revoke all]│ │
│          │  │  Changing your password ends every other session.    │ │
│          │  └──────────────────────────────────────────────────────┘ │
│          │  ┌──────────────────────────────────────────────────────┐ │
│          │  │  NOTIFICATIONS                                       │ │
│          │  │  Email me when…                                      │ │
│          │  │  [✓] An upload fails in a queue I started            │ │
│          │  │  [✓] A key I manage is revoked                       │ │
│          │  │  [ ] Weekly usage summary                            │ │
│          │  └──────────────────────────────────────────────────────┘ │
└──────────┴───────────────────────────────────────────────────────────┘
```

## Wireframe — mobile (populated)

```
┌────────────────────────────┐
│ ☰ ⬡    Acme Inc ▾   👤     │
├────────────────────────────┤
│ ‹ Account settings         │
│ ─────────────────────────  │
│ PASSWORD                   │
│ Current                    │
│ ┌────────────────────────┐ │
│ │ ••••••••••          👁 │  │
│ └────────────────────────┘ │
│ New                        │
│ ┌────────────────────────┐ │
│ │ •••••••••••••       👁 │  │
│ └────────────────────────┘ │
│ [Save password]            │
│ ─────────────────────────  │
│ ACTIVE SESSIONS            │
│ This browser · now   (this)│
│ Windows · 2h ago  [Revoke] │
│ Mobile · 1d ago    [Revoke]│
│ [Revoke all others]        │
│ ─────────────────────────  │
│ NOTIFICATIONS              │
│ [✓] Upload fails in queue  │
│ [✓] Key I manage revoked   │
│ [ ] Weekly usage summary   │
└────────────────────────────┘
```

## Interactions

- **Change password**: current + new + confirm. On save → loading → success toast; **every other session is revoked server-side** (`backend/domains.md` §4.4 behavior on reset — same semantics on a voluntary change), which the caption under Sessions states before the user commits. Wrong current password → field error, no hint about how wrong.
- **2FA enable**: opens the provisioning modal — the `otpauth_url` QR + the secret **shown exactly once** (MV-E5), a code field, and "enable" disabled until a confirming code is entered. On success: **10 backup codes shown once** with a "you won't see these again" warning. This is the `SecretRevealModal` discipline applied to TOTP.
- **2FA disable**: requires a current TOTP code or an unused backup code (R15 — a disable can't be a bare click, or 2FA protects nothing).
- **Regenerate backup codes**: invalidates the old set immediately (a lost code list is replaced, not appended).
- **Backup-code count** shows `X remaining` so a user near zero regenerates before they're locked out.
- **Active sessions**: listed from the refresh-token store (device/platform/approximate location/last-active). **Revoke** ends one session; **Revoke all** ends every other session, with a plain confirm (this is "kick everyone else out", not destructive of data). The current browser session is marked "(this)" and cannot be revoked.
- **Notification toggles**: save-on-toggle (checkbox changes persist immediately, subtle toast). These are preferences on `user_settings` (`ERD.md`) — no tenant-level implication.
- **Log out of this device** is in the user menu, not here — deliberate: the menu's Log out is the single obvious exit, and putting a second one here invites confusion about which one "really" logs out.
- **Permission-aware**: account settings are personal — no role gating. A locked-out user's 2FA reset is platform-admin only (`backend/domains.md` §4.10, R15) and does NOT live here.
