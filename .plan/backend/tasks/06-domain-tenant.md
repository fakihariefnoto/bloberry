# Task group — 06 domain: tenant

**Depends on:** `02-migrations`, `03-platform`, `16-authz` (role checks in the usecase layer per `domains.md` §1.2). Depends on `user.Reader`, `usage.Reader` via narrow interfaces.

- [ ] **`repository.go` / `repository/repository.go`** — tenant CRUD, membership CRUD, invitation CRUD, quota counters, `default_backend_id` assignment.
- [ ] **`usecase.go` / `usecase/usecase.go`** — tenant create (root folder created with it), update, suspend/reactivate; member add/remove/role-change (owner rules — at least one owner, §ERD memberships); invite issue/resend (replacing the old token); **backend reassignment** (applies to *new* objects only, ADR-4); quota change (platform-admin path).
- [ ] **`handler.go` / `handler/handler.go`** — `POST /v1/tenants`, `GET/PATCH /v1/tenants/:id`, `POST /v1/tenants/:id/members`, `DELETE /v1/tenants/:id/members/:uid`, `PATCH /v1/tenants/:id/members/:uid/role`, invitations create/list/resend, `PATCH /v1/tenants/:id/backend`, quota get/set (platform admin).
- [ ] **Exactly-one-owner business rule** — enforced in the usecase layer (not a DB constraint), per `ERD.md` memberships note.
- [ ] **Suspension behavior** — reads keep working, writes blocked with a clear reason (PRD PA-E2); the suspended state is a tenant field, not a wall of 500s.
- [ ] **Interface naming + mocks.**

**tests:** tenant create makes the root folder; owner rules (can't demote/remove the last owner); invite resend invalidates the old token; backend reassignment leaves existing objects' `backend_id` untouched; suspend blocks writes but not reads.
