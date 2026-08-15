# Task group — 23 screen: application-detail

**Depends on:** `02-design-tokens`, `03-routing`, `04-core-infra`. **Blocks:** `25-flows` (revoke-a-key flow). **Mockup:** [`mobile/mockup/application-detail.md`](../mockup/application-detail.md).

- [ ] **Layout — populated** per the mockup: back, app name + description + created, KEYS list (masked key `blob_live_••••4f2a` in mono, scope + perms line, last-used line, status pill, `⋮`), the "shown once" footer.
- [ ] **Revoke** (`⋮` → Revoke) — `confirm-destructive` sheet per the mockup wireframe: shows the key's **last-used time + IP** (the TA-E3 context), typed-name (last-4-chars) confirmation for the irreversible act.
- [ ] **Last-active-key consequence** — "This is acme-cms's only active key. Its pipeline will fail on the next call."
- [ ] **Revoked** — toast with **no Undo** (deliberate, `web/components.md` rule shared); row mutes with a Revoked pill; history (last_used_at) retained for the audit trail (`ERD.md`).
- [ ] **Expiring < 7 days** → warning pill; expired keys show no Revoke (already dead).
- [ ] **Create key absent** — web-only (stated in the footer; the phone is for containment, not provisioning).
- [ ] **A11y** — the typed-name field matches last-4 case-insensitively; the confirm sheet traps focus and returns it to the triggering row.

**tests:** revoke confirm shows last-used; typed-name refuses a wrong value; last-active-key warning; no Undo toast; masked keys only.
