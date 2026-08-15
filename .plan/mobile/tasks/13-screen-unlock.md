# Task group — 12 screen: unlock

**Depends on:** `02-design-tokens`, `03-routing`, `04-core-infra` (local_auth + FlutterFragmentActivity). **Blocks:** `25-flows` (biometric-unlock chain). **Mockup:** [`mobile/mockup/unlock.md`](../mockup/unlock.md).

- [ ] **Lock screen per the mockup** — brand mark, lock glyph, biometric button, PIN fallback, "Use password". Replaces the shell entirely (even the app bar — a stray tap can't leak a file row behind the lock).
- [ ] **Biometric prompt** — fires `local_auth`; success returns **exactly where the user was** (not `files`); rejected → shake/error pulse on the glyph, retry, escalation to PIN/password (never an infinite biometric loop).
- [ ] **PIN entry** — 4-digit pad per the mockup; correct → unlock in place; wrong → digits clear, one "PIN incorrect", cap 5 attempts then force "Use password". PIN stored encrypted (never plaintext).
- [ ] **Use password** → `login` (full login clears the lock state; success returns to the pre-lock destination).
- [ ] **Biometric unavailable** — renders PIN-first with a caption explaining why the sensor isn't offered.
- [ ] **Android `FlutterFragmentActivity`** — MainActivity extends it (the crash-if-missed requirement), verified here not buried in infra.
- [ ] **A11y** — single `Semantics` group; biometric button labeled ("Unlock with Face ID"); PIN digits announce on press; no haptic-only feedback.

**tests:** biometric success returns in place; failure escalates to PIN then password; 5-attempt PIN cap; unlock screen covers the app bar.
