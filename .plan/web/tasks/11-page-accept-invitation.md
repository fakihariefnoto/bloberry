# Task group — 11 page: accept-invitation

**Depends on:** `01-setup`, `02-design-tokens`, `03-routing`, `05-core-infra`, `06-shared-components`. **Blocks:** `33-flows` (auth chain). **Mockup:** [`web/mockup/accept-invitation.md`](../mockup/accept-invitation.md).

- [ ] **Invite validation on load** — skeleton while checking; wireframe's role badge ("Member") is **read-only** (the invite carries the role, never editable here).
- [ ] **New-user branch** — display name + password + confirm; "Join <tenant>" primary; "Already have an account? Sign in instead" → `login`.
- [ ] **Existing-email branch** — display-name hidden; heading "Invited as Member of <tenant> — set a password to finish joining"; accepting adds the membership and logs in (`backend/domains.md` §4.1 — a user legitimately belongs to several tenants; not a duplicate-account conflict).
- [ ] **Join** → loading → success → `files` at the tenant root.
- [ ] **`invite_invalid`** — the expired wireframe (⚠, "Invitations expire after 7 days and can be used once", "Ask the admin who invited you for a new one", Sign in action). No retry — the token is dead.
- [ ] **A11y** — tenant + role announced on load; the error state's "Sign in" is a link.

**tests:** new-user vs existing-email branch; role read-only; invalid invite shows the expired state with no retry.
