# Command — `bloberry grant create`

## Purpose & context

- **User goal**: grant a principal access to a folder subtree (PRD M9) — the folder-level RBAC layer. Allow-only, most-specific-wins, **no deny rules** (PRD D7).
- **When they reach for it**: interactive ("let the build agent write to assets/"), and scripted onboarding.
- **Needs**: auth, `tenant_admin`+ role; the principal must be in the tenant (user or application).
- **Data**: `grants` — principal_type/id, folder_id, permissions, expires_at, granted_by. `ERD.md` grants note: `revoked_at` set rather than deleted.

## Signature

```
bloberry grant create [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `--folder <path>` | path | — | yes | Folder subtree the grant applies to. |
| `--principal <ref>` | string | — | yes | `user:<email>` or `app:<id>`. |
| `--permission <p>` | string (repeatable) | — | yes | `read`, `write`, `delete`, `share`. |
| `--expires <when>` | string | never | no | `30d` or `YYYY-MM-DD`. |

## Help text

```
Grant a principal access to a folder and everything under it.

Grants are allow-only: they add to the principal's role, never
subtract, and there is no deny rule (see docs, D7). A grant to
bloberry://assets covers assets/ and all descendants. Use the most
specific folder that gives what's needed.

Examples:
  bloberry grant create --folder bloberry://assets/v2 --principal user:bot@acme.dev --permission write
  bloberry grant create --folder bloberry://shared --principal app:app_3d9f --permission read --expires 30d
```

## Output states

**Success**

```
✓ Granted user:bot@acme.dev write on assets/v2 (and descendants)
  Grant ID: gr_6a1e · expires: never
```

**`--json`**

```json
{"id":"gr_6a1e","folder":"assets/v2","principal_type":"user","principal":"bot@acme.dev","permissions":["write"],"expires_at":null,"granted_by":"user_8f2a1c"}
```

## Exit codes

| Code | When |
|---|---|
| `0` | Created. |
| `5` | Folder or principal doesn't exist. |
| `4` | Forbidden — not `tenant_admin`+, or the principal isn't in this tenant. |
| `2` | Bad invocation (missing `--folder`/`--principal`, empty `--permission`). |

## Behavior notes

- **The help states the allow-only model** (PRD D7) — a CLI user with IAM habits will look for a deny flag; the absence is a decision, and the help says so.
- **Permission validation**: at least one of read/write/delete/share required; `admin` isn't a grant-level permission (grants are folder-scoped, `ERD.md`).
- **stdout**: the result / JSON. **stderr**: nothing on success.
- **No confirmation** — granting is additive and revocable (`grant revoke`); destructive commands own the prompts.
- **Expiry**: `--expires` absent = permanent until revoked; the output states it.
