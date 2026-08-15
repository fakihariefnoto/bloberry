# Task group — 22 screen: applications

**Depends on:** `02-design-tokens`, `03-routing`, `04-core-infra`. **Blocks:** `23-screen-application-detail`. **Mockup:** [`mobile/mockup/applications.md`](../mockup/applications.md).

- [ ] **Layout per the mockup** — list rows (⚙ initial, name, "N keys · used <time>" line, warning pill for keyless apps, `⋮`), the create-on-web footer.
- [ ] **Tap a row** → `application-detail`.
- [ ] **`⋮` only offers "Open detail"** — the revoke flow lives in detail where the last-used context is visible (PRD TA-E3); no revoke shortcut here.
- [ ] **Keyless app flagged** — warning pill ("all revoked" / "never used") — a keyless app is usually a broken pipeline, not a cleaned-up one.
- [ ] **Footer caption** — "Creating applications is a desk task — use the web dashboard." (explains the absent New-app button as a decision, `navigation.md` platform-admin note).
- [ ] **Empty** — "No applications yet · Register one on the web dashboard."
- [ ] **A11y** — rows ≥48dp; the warning pill carries a text label.

**tests:** row → detail; `⋮` has no revoke shortcut; keyless warning pill; footer caption.
