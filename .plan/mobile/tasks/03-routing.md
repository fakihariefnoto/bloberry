# Task group — 03 routing

**Depends on:** `01-setup.md` (go_router), `02-design-tokens.md`. **Blocks:** every screen file.

**Source:** `mobile/navigation.md` — the authoritative route graph. Routes, args, presentation, back destinations and auth gates come from its route table verbatim.

- [ ] **All 21 routes registered in go_router**, each with its exact path and typed args from the route table:
  `onboarding` `/onboarding` · `welcome` `/welcome` · `login` `/login` · `otp-login` `/login/otp` (`email: String?`) · **`pair-login` `/login/pair` (push, camera overlay)** · `forgot-password` `/forgot-password` · `reset-password` `/reset-password` (`token: String`, replace via deep link) · `accept-invitation` `/invite/:token` (`token: String`, replace via deep link) · `unlock` `/unlock` (replace over shell) · `files` `/files/:folderId?` (tab, default; `folderId` null = root) · `file-detail` `/files/detail/:fileId` (push) · `search` `/search` (`q: String?`, push) · `uploads` `/uploads` (tab) · `shares` `/shares` (tab) · `more` `/more` (tab) · `profile` `/profile` (push) · `settings` `/settings` (push) · `usage` `/usage` (push, `tenant_admin`+) · `applications` `/applications` (push, `tenant_admin`+) · `application-detail` `/applications/:appId` (push, `tenant_admin`+) · `tenant-suspended` `/suspended` (replace).
- [ ] **The 4-tab shell** (`size.navbar-height` 64px): Files (default), Uploads, Shares, More — per the Shell section. No 5th tab; administrative surfaces stay under More.
- [ ] **Presentation per route** — tab / push / replace / modal, exactly per the route table's Presentation column. Back behavior matches the table's "Back goes to" (`file-detail` → `files` at the file's folder, etc.).
- [ ] **Auth gating per route** — public (onboarding, welcome, login, otp-login, forgot-password, reset-password, accept-invitation); session-held-locked (unlock); any-authenticated (files, file-detail, search, uploads, shares, more, profile, settings, tenant-suspended); `tenant_admin`+ (usage, applications, application-detail).
- [ ] **Default after login** → `files` at the tenant root; **unauthenticated** → `welcome`, or `onboarding` on genuine first launch (persisted flag, shown once ever).
- [ ] **Deep links** — `reset-password` and `accept-invitation` handle their tokens and consume them (replace the route so the token never stays in the URL).
- [ ] **Session-expiry redirect** — any 401 → `login`, preserving the intended destination and returning to it after login.
- [ ] **Tenant-suspended guard** — a suspended tenant's requests → `tenant-suspended` before the screen renders its own failure.
- [ ] **The 5 sheets** as route-less modals, per the "Deliberately sheets, not screens" table: `source-picker`, `share-sheet`, `folder-picker`, `confirm-destructive`, `tenant-switcher`. Listed so they're not built as routes.

**tests:** route-table conformance (every route has the right args + presentation + gate); deep-link token consumption; session-expiry destination preservation; tab shell active-state.
