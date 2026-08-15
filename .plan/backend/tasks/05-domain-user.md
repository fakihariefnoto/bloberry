# Task group — 05 domain: user

**Depends on:** `02-migrations`, `03-platform`. The standard auth-adjacent domain: profile + settings. Smallest domain in the system.

- [ ] **`repository.go` / `repository/repository.go`** — get/update profile, get/update settings (the embedded `users.settings` doc), password change.
- [ ] **`usecase.go` / `usecase/usecase.go`** — update-profile (email change → re-verification), update-settings, change-password (current-password check; **every other session invalidated** on success, `domains.md` §4.4 semantics).
- [ ] **`handler.go` / `handler/handler.go`** — `GET/PATCH /v1/users/me`, `GET/PATCH /v1/users/me/settings`, `POST /v1/users/me/password`. All scoped to the authenticated principal, never an arbitrary user id.
- [ ] **Interface naming + mocks** — `user.Repository`/`user.Usecase`/`user.Handler` + mockgen coverage.

**tests:** profile update; email-change triggers verification state; settings update; password change invalidates other sessions; handlers reject a non-`me` id.
