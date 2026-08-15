# Command — `bloberry grant revoke`

## Purpose & context

- **User goal**: remove a grant. Takes effect on the next request (explicit cache invalidation, `backend/domains.md` §5.3).
- **When they reach for it**: access review ("bot@acme.dev shouldn't write to assets/ anymore"), deprovisioning.
- **Needs**: auth, `tenant_admin`+ role.
- **Data**: `grants` — set `revoked_at` (kept for audit, `ERD.md` grants note); the principal loses that folder's access immediately.

## Signature

```
bloberry grant revoke <id> [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `<id>` | string | — | yes | Grant ID (`gr_…`). |
| `--yes, -y` | bool | false | no | Skip the confirmation. |

## Help text

```
Revoke a folder grant. Takes effect on the next request.

The revoked grant is kept for the audit trail but stops granting
immediately. This never removes a role — a tenant_admin keeps their
role's access; it only removes this grant's added permissions.

Examples:
  bloberry grant revoke gr_6a1e
  bloberry grant revoke gr_6a1e --yes
```

## Output states

**Success**

```
✓ Revoked grant gr_6a1e (user:bot@acme.dev · write on assets/v2/)
```

**Already revoked**

```
Grant gr_6a1e is already revoked. (exit 0)
```

**`--json`**

```json
{"id":"gr_6a1e","revoked":true,"revoked_at":"2026-03-12T14:10:00Z"}
```

## Exit codes

| Code | When |
|---|---|
| `0` | Revoked (or already revoked — idempotent). |
| `5` | No grant with that ID. |
| `4` | Forbidden. |
| `2` | Missing ID. |

## Behavior notes

- **Immediate** (explicit invalidation, not TTL) — the help states "next request".
- **The role floor is untouched** — the help explicitly notes revoking a grant doesn't reduce a role, preventing the "I revoked the grant, why does the admin still have access" confusion.
- **stdout**: the result. **stderr**: the confirmation prompt.
- **Idempotent**: revoking an already-revoked grant is exit 0.
- **No Undo**: there is no un-revoke; recreating is the path (`grant create` with the same args) — the message doesn't pretend otherwise.
