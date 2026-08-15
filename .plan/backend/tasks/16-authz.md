# Task group — 16 authz (the permission resolver)

**Depends on:** `02-migrations`, `03-platform` (Redis). **Blocks:** `06/07/08/09/10/11` (every domain whose usecase calls it). **Not a domain** — no endpoints, no repository of its own (`domains.md` §2 note). The highest-risk function in the codebase (TRD R4, PRD G5, ADR-6).

- [ ] **`internal/authz/principal.go`** — the unified `Principal` type: `{Type user|application, ID, TenantID, Role, IsPlatformAdmin, ScopeFolderIDs, KeyPermissions, Grants}` (`domains.md` §5.1). Built by the auth middleware from either scheme; nothing downstream branches on scheme (§4.7).
- [ ] **`internal/authz/resolver.go`** — the pure function `Resolve(p Principal, action Permission, folderID string, ancestors []string) Decision`. **No I/O.** Evaluation order, each step's reason in `domains.md` §5.2:
  1. Platform admin → allow (runs the install, above tenancy).
  2. Tenant mismatch → deny (the isolation boundary, checked first).
  3. Key scope (applications) — non-empty `ScopeFolderIDs` and neither folderID nor an ancestor in it → deny. **A key only ever narrows.**
  4. Key permissions — non-empty and lacking `action` → deny. Same rule.
  5. Role floor — owner/admin all actions; member read+write+share; viewer read.
  6. Grants — the grant whose folder is folderID or the **deepest** ancestor entry; if it includes `action` → allow. **Grants only add.**
  7. Otherwise → deny.
  **There is no deny rule** (PRD D7). Absence is the default; nothing overrides an allow.
- [ ] **`internal/authz/resolver_test.go`** — table-driven, **100% branch coverage required** (PRD G5), enforced in CI. Must cover at minimum (`domains.md` §5.4): each role × each action; platform admin crossing tenants; a key scoped to a subtree accessing its own folder, a descendant, an ancestor, a sibling; key permissions narrower than its role; an expired grant; a revoked grant; nested grants at two depths where the deeper wins; a grant on a non-ancestor folder.
- [ ] **`internal/authz/cache.go`** — Redis cache keyed by principal ID. **Principal assembly is cached; `Resolve` itself is never cached** (pure + fast; caching decisions would multiply the invalidation surface).
- [ ] **Explicit invalidation, never TTL** — key revoke, grant create/revoke, membership change, role change all delete the principal's cache entry synchronously (`domains.md` §5.3). PRD G5's "revocation takes effect on the next request" is a hard requirement; the invalidation callers in `06/10/11` wire into this.

**tests:** the 100%-branch-coverage table; a separate suite asserting the evaluation *order* (platform-admin before tenant-mismatch before scope before grants — a reordering bug is a security bug even with full coverage); cache invalidation on each mutation path.
