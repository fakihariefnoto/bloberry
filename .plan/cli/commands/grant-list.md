# Command — `bloberry grant list`

## Purpose & context

- **User goal**: see who has what access to which folder subtree — the "who can see this" answer (the `grants` index `{tenant_id, principal_type, principal_id}` and `{tenant_id, folder_id}`, `ERD.md`).
- **When they reach for it**: interactive audit before changing access; understanding why a principal can reach a folder.
- **Needs**: auth, `tenant_admin`+ role.
- **Data**: `grants` — principal, folder, permissions, expiry, granted_by, revoked state.

## Signature

```
bloberry grant list [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `--folder <path>` | path | all | no | Grants on this folder subtree. |
| `--principal <ref>` | string | all | no | `user:<email>` or `app:<id>`. |

## Help text

```
List folder grants.

Optionally filter by folder or principal. Remember grants only add —
a principal's effective access is role floor + most-specific grant
(see docs, D7).

Examples:
  bloberry grant list
  bloberry grant list --folder bloberry://assets
  bloberry grant list --principal user:bot@acme.dev --json
```

## Output states

**Success**

```
GRANT ID   PRINCIPAL             FOLDER            PERMS         EXPIRES
gr_6a1e    user:bot@acme.dev     assets/v2/        write         never
gr_2c8f    app:app_3d9f          shared/           read          in 30d
gr_9d4b    user:jane@acme.dev    assets/           read,write    revoked
```

**No grants**

```
No grants on this folder. Access here is role-based only.
```

**`--json`**

```json
[{"id":"gr_6a1e","principal_type":"user","principal":"bot@acme.dev","folder":"assets/v2","permissions":["write"],"expires_at":null,"revoked":false}]
```

## Exit codes

| Code | When |
|---|---|
| `0` | Listed (including empty). |
| `4` | Forbidden — not `tenant_admin`+. |
| `5` | The `--folder` path doesn't exist. |

## Behavior notes

- **Revoked grants are shown** (muted, `revoked` in the state) — the audit trail survives (`ERD.md` grants note), and the list is honest that a grant existed rather than erasing it.
- **stdout**: the table / JSON. **stderr**: nothing on success.
- **Empty state explains the model** — "Access here is role-based only" tells the caller why a grant isn't there (it's not that listing failed; there are none).
- **Effective-access is not computed here** — this lists raw grants; the resolver's allow-only combination is the backend's job (`backend/domains.md` §5), and the help points at that model rather than pretending to re-derive it.
