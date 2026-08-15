# Task group — 18 screen: more

**Depends on:** `02-design-tokens`, `03-routing`, `04-core-infra`. **Blocks:** none (menu hub). **Mockup:** [`mobile/mockup/more.md`](../mockup/more.md).

- [ ] **Layout per the mockup** — user header (avatar, name, email, role badge), TENANT group (Acme Inc ▾ + inline storage bar 312/500 GB 62%), ACCOUNT group (Profile ›, Settings ›), ADMIN group (Usage ›, Applications ›), ABOUT group (version, Log out).
- [ ] **User header** — tap → `profile`.
- [ ] **Tenant row** → `tenant-switcher` sheet; switching always lands on `files` at the new tenant's root (`navigation.md` rule).
- [ ] **Inline usage summary** — the quota bar renders live so a glancing admin doesn't even need to open `usage`.
- [ ] **Admin group gated** — Usage/Applications shown only to `tenant_admin`+; a member/viewer sees the section header hidden entirely (not an empty group). Role badge is the same gate.
- [ ] **Log out** — bottom, visually separated (settings-list pattern), a plain action (reversible by logging back in); revokes the refresh token server-side.
- [ ] **Version row** — static `text.caption`; tappable only in debug builds (shows build commit).
- [ ] **A11y** — avatar `role="presentation"`; role badge has a text label; rows ≥48dp.

**tests:** admin group hidden for member/viewer; tenant switch lands on files; inline quota renders; log out revokes server-side.
