# Task group — 05 screen: onboarding

**Depends on:** `02-design-tokens`, `03-routing`. **Blocks:** `25-flows` (first-launch chain). **Mockup:** [`mobile/mockup/onboarding.md`](../mockup/onboarding.md).

- [ ] **3-pane carousel** per the mockup — folders / sharing / access keys, one idea per slide, `flutter_spinkit`-free (static content).
- [ ] **Skip top-right** — always available, jumps to `welcome`, persists the seen-flag.
- [ ] **Progress dots** — current filled; dots indicate position, not tappable.
- [ ] **Next → Get started** — label swaps on the final slide; Get started → `welcome`, flag persisted.
- [ ] **Swipe** advances/retreats panes.
- [ ] **A11y** — each slide's headline is a heading; dots announce "slide 2 of 3"; Skip/Get started keyboard-focusable.

**tests:** flag persisted after skip/get-started (never shows again); 3 slides with the right copy.
