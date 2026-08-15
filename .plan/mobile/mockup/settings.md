# Screen — settings

## Purpose & context

- **User goal**: app-level settings — the PIN/biometric unlock toggle, notifications, locale, and the security posture toggles. Tenant/account config is NOT here (that's web-only for tenant; account password/sessions live in `account-settings` on web and aren't replicated on mobile beyond what's needed to secure the device).
- **Entry points**: Settings from `more`.
- **Exit points**: back → `more`; Log out → `welcome`.
- **Data needed**: `user_settings` — notifications_enabled, biometric_unlock_enabled, locale, default_tenant_id. Platform capability check (biometric available?).

## States

- [x] Loading (skeleton toggles)
- [x] Populated
- [x] Error (a toggle failed to save — inline)
- [x] Domain-specific — biometric unavailable (toggle disabled with reason)

## Style reference

- **Components used**: settings-list pattern (`design-collection/mobile-screen/patterns.md`) — grouped headers, right-aligned values/toggles, destructive actions bottom. The PIN/biometric toggle lives under a SECURITY group per the pattern (not floating alone).
- No token deltas.

## Wireframe — mobile (populated)

```
┌───────────────────────────┐
│  ← Settings               │
├───────────────────────────┤
│  SECURITY                 │
│  ┌──────────────────────┐ │
│  │ PIN / biometricunlock│ │
│  │  ● ─ ─ ─ ─ ○         │ │  ← toggle
│  └──────────────────────┘ │
│  Unlock with Face ID or   │
│  your PIN after 5 min.    │
│  ──────────────────────── │
│  NOTIFICATIONS            │
│  ┌──────────────────────┐ │
│  │ Upload failures      │ │
│  │  ● ─ ─ ─ ─ ○         │ │
│  └──────────────────────┘ │
│  ┌──────────────────────┐ │
│  │ Key revoked alerts   │ │
│  │  ● ─ ─ ─ ─ ○         │ │
│  └──────────────────────┘ │
│  ──────────────────────── │
│  GENERAL                  │
│  ┌──────────────────────┐ │
│  │ Language      English ▾│
│  └──────────────────────┘ │
│  ┌──────────────────────┐ │
│  │ Default tenant   Acme ▾│
│  └──────────────────────┘ │
│  ──────────────────────── │
│  ABOUT                    │
│  Version 0.1.0            │
│  Log out              ⏻   │
└───────────────────────────┘
```

## Interactions

- **PIN/biometric unlock**: toggle off by default (`templates/flutter-defaults.md`); on enable → the PIN-setup flow (`unlock` screen's PIN pad in "set PIN" mode) or the platform biometric prompt, depending on what the toggle's subtitle says. When the device has no biometric sensor, the toggle is disabled with the caption "Biometric unavailable on this device — PIN still works."
- **Toggles** save immediately with a subtle toast (not a Save button); a failed save shows an inline error and reverts the toggle.
- **Language / Default tenant** are selects (right-aligned current value, tap → bottom sheet).
- **Log out** at the bottom, visually separated; revokes the refresh token server-side and lands on `welcome`. (The separate "Log out" on `more` is the same action — this one is the settings-list convention.)
- **A11y**: toggles carry text labels (not color alone); the toggle's on/off state is announced on change.
