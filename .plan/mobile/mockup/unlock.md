# Screen — unlock

## Purpose & context

- **User goal**: quickly re-authenticate to a locked app after the configured idle timeout, via PIN or biometric — without a full login (PRD security posture; `templates/flutter-defaults.md` PIN/biometric unlock rule, off by default in `settings`).
- **Entry points**: app resumed after the timeout with a live session → `unlock` replaces the shell (session is held but locked). Also after cold start if the lock is enabled.
- **Exit points**: authenticated → returns **exactly where the user was** (not `files` — `navigation.md` flow chain); "Use password" → `login`; failure → stays on `unlock` with a retry + "Use password".
- **Data needed**: none from the network — the biometric/PIN check is local (`local_auth` + a local PIN stored in the keychain/encrypted storage). No token refresh on this screen.

## States

- [x] Biometric prompt (Face ID / fingerprint)
- [x] PIN entry fallback
- [x] Success (unlocks, returns in place)
- [x] Failure — biometric rejected / PIN wrong (retry, then "Use password")
- [x] Domain-specific — biometric unavailable on this device (PIN only, with the reason)

## Style reference

- **Components used**: minimal lock screen — brand mark, lock glyph, biometric button, PIN pad, "Use password" link. No tab bar, no app bar; the shell is fully replaced (even the app bar is covered so a stray tap can't leak a file row behind the lock).
- `color.primary` for the biometric button; `color.error` only for a rejected attempt.
- No token deltas.

## Wireframe — mobile (biometric prompt)

```
┌───────────────────────────┐
│                           │
│         ⬡                 │
│                           │
│      🔒  Locked           │
│                           │
│  Bloberry is locked.      │
│  Unlock to continue.      │
│                           │
│     ┌─────────────┐       │
│     │    👁  ╲_╱  │        │  ← Face ID / fingerprint
│     │   Unlock   │        │
│     └─────────────┘       │
│                           │
│  ─────────────────────    │
│  Use PIN                  │
│                           │
│  Use password             │
└───────────────────────────┘
```

## Wireframe — mobile (PIN entry)

```
┌───────────────────────────┐
│                           │
│         ⬡                 │
│                           │
│  Enter your PIN           │
│                           │
│       ●  ●  ○  ○          │
│                           │
│      1   2   3            │
│      4   5   6            │
│      7   8   9            │
│     ⌫   0   ✓             │
│                           │
│  Use password             │
└───────────────────────────┘
```

## Interactions

- **Unlock (biometric)**: fires the platform prompt; success returns to the exact pre-lock location. Rejected → shake/error pulse on the glyph, retry, and a "Use PIN"/"Use password" escalation path (never an infinite biometric loop).
- **PIN**: 4-digit pad; correct → unlock in place; wrong → digits clear, one "PIN incorrect" error, cap at 5 attempts before forcing "Use password" (`local_auth` conventions; the PIN itself is stored encrypted, never plaintext).
- **Use password** → `login` (the full login invalidates the lock state; success returns to the pre-lock destination).
- **Biometric unavailable** (no sensor / permission denied): the screen renders PIN-first with a caption explaining why the sensor isn't offered.
- **A11y**: the lock screen is a single `Semantics` group; the biometric button has a labeled name ("Unlock with Face ID"); PIN digits announce on press. No haptic-only feedback.
