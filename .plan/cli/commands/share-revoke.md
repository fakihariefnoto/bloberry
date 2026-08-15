# Command — `bloberry share revoke`

## Purpose & context

- **User goal**: kill a share link before its expiry — the containment action when a link leaked or was shared too far.
- **When they reach for it**: interactive cleanup and incident response; scripts rotating an exposed link.
- **Needs**: auth, tenant context, `share` permission.
- **Data**: `share_links` — the record is set `revoked_at` rather than deleted (audit trail survives, `ERD.md` share-links note).

## Signature

```
bloberry share revoke <id> [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `<id>` | string | — | yes | The link ID (`sl_…`, from `share list`). |
| `--yes, -y` | bool | false | no | Skip the confirmation. |

## Help text

```
Revoke a share link. It stops working immediately.

The link is marked revoked (not deleted) so the audit trail and hit
counts survive. A revoked link renders the "this link has expired"
page to whoever opens it.

Examples:
  bloberry share revoke sl_7f2a91
  bloberry share revoke sl_7f2a91 --yes
```

## Output states

**Success (interactive, after confirm)**

```
✓ Revoked link sl_7f2a91 (https://…/s/v2-1-0)
  12 hits · last used 2h ago · now dead
```

**Already revoked**

```
Link sl_7f2a91 is already revoked. (exit 0)
```

**`--json`**

```json
{"id":"sl_7f2a91","revoked":true,"revoked_at":"2026-03-12T14:10:00Z"}
```

## Exit codes

| Code | When |
|---|---|
| `0` | Revoked (or already revoked — idempotent). |
| `5` | No link with that ID. |
| `4` | Forbidden. |
| `2` | Missing ID. |

## Behavior notes

- **Immediate**: revocation takes effect on the next request to the link (`backend/domains.md` §7 — the cache entry is invalidated explicitly, not TTL-out).
- **The confirmation states the hit count** ("12 hits · last used 2h ago · now dead") — the same "understand what you're breaking" context as the web's revoke flows (PRD TA-E3 discipline applied to links).
- **Idempotent**: revoking an already-revoked link is a success (exit 0) — scripts can revoke defensively without guarding state.
- **No Undo**: revoking is irreversible (there is no un-revoke); the confirmation is the only gate.
- **stdout**: the result. **stderr**: the confirmation prompt.
