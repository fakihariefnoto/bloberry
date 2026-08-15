# Task group — 06 commands: key + app + grant

**Depends on:** `02-core-infra`, `03-auth`. One checkbox per command; acceptance = "output matches the sample in `cli/commands/<file>.md`".

## `key create` (`cli/commands/key-create.md`)

- [ ] Issue a scoped key: `--app`, `--folder` (repeatable), `--permission` (read/write/delete/share/admin), `--expires` (`30d`/`YYYY-MM-DD`), `--prefix` (`blob_live_`/`blob_test_`).
- [ ] **The secret prints exactly once** and is unrecoverable (PRD D5) — the human output labels it clearly with the warning; `--json` carries `secret` for a one-step pipe into a secrets manager.
- [ ] Scoping is the default posture — an unscoped whole-tenant key prints a warning line; the help nudges `--folder`/`--permission`.
- [ ] Argon2id-hash server-side; only `secret_hash` + `last_four` stored.

## `key list` (`cli/commands/key-list.md`)

- [ ] Masked keys only (`blob_live_••••4f2a` — `last_four` in `--json`, never a full secret); scope/perms/last-used/state columns; `expiring` = < 7 days; active sort by last-used desc; filter `--app`, `--status`.

## `key revoke` (`cli/commands/key-revoke.md`)

- [ ] Revoke by id; confirmation shows **last-used + IP** (PRD TA-E3); immediate (explicit cache invalidation, next request); **irreversible, no Undo**; last-active-key detection adds "This was <app>'s only active key" + the replacement command; idempotent (already-revoked = exit 0).

## `app create` (`cli/commands/app-create.md`)

- [ ] Register an application; name-conflict → exit 8 with the existing app's ID (scriptable continue); output chains to `key create --app <id>` (provisioning is two steps, made discoverable).

## `app list` (`cli/commands/app-list.md`)

- [ ] List applications with key counts + last-used; empty state names `app create`.

## `app delete` (`cli/commands/app-delete.md`)

- [ ] Delete an application; **refused with active keys** (exit 8, the exact revoke commands listed — deleting an app with live keys silently breaks CI, `ERD.md` lifecycle); confirmation unless `--yes`; already-deleted → exit 5 not silent success.

## `grant create` (`cli/commands/grant-create.md`)

- [ ] Grant principal × folder × permissions × expiry; `--principal user:<email>|app:<id>`; **allow-only semantics stated in the help** (PRD D7 — no deny flag exists by design); `admin` not a grant-level permission.

## `grant list` (`cli/commands/grant-list.md`)

- [ ] List grants, filter by `--folder`/`--principal`; **revoked grants shown muted** (audit trail survives); empty state explains "access here is role-based only".

## `grant revoke` (`cli/commands/grant-revoke.md`)

- [ ] Revoke by id; immediate; **role floor untouched** (stated — revoking a grant never reduces a role); idempotent; no Undo.

**tests:** per command — happy path, error, `--json`; `key create` secret-shown-once; `app delete` active-key refusal; `grant create` allow-only wording.
