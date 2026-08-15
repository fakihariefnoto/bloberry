# Command — `bloberry key revoke`

## Purpose & context

- **User goal**: kill a leaked or unwanted access key immediately. Takes effect on the next request (`backend/domains.md` §7 explicit invalidation, PRD G5) — the containment action.
- **When they reach for it**: incident response (a key in a leaked log), and rotating credentials.
- **Needs**: auth, `tenant_admin`+ role, `admin` permission.
- **Data**: `access_keys` — the record is set `revoked_at` (not deleted) so the audit trail and `last_used_at` survive (PRD TA-E3).

## Signature

```
bloberry key revoke <id> [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `<id>` | string | — | yes | Key ID (`ak_…`, from `key list`). |
| `--yes, -y` | bool | false | no | Skip the confirmation. |

## Help text

```
Revoke an access key. It stops working on the next request.

The key's history (last-used time, IP) is kept for the audit trail.
Revoking the last active key of an application makes its pipeline fail
until a new key is issued. This is irreversible — no undo.

Examples:
  bloberry key revoke ak_5f2b
  bloberry key revoke ak_5f2b --yes
```

## Output states

**Success (interactive, after confirm)**

```
✓ Revoked key blob_live_••••4f2a (ak_5f2b)
  App:      acme-cms
  Last used: Mar 13 09:12 UTC from 203.0.113.8
  Effect:   this key now fails every request
```

**Last active key warning**

```
✓ Revoked key blob_live_••••4f2a (ak_5f2b)

  ⚠ This was acme-cms's only active key. Its next request will fail
  until a new key is issued:
    bloberry key create --app app_3d9f
```

**`--json`**

```json
{"id":"ak_5f2b","revoked":true,"revoked_at":"2026-03-12T14:10:00Z","was_last_active":false}
```

## Exit codes

| Code | When |
|---|---|
| `0` | Revoked. |
| `5` | No key with that ID. |
| `4` | Forbidden. |
| `2` | Missing ID. |

## Behavior notes

- **Immediate and irreversible**: invalidation is explicit (not TTL), and there is no un-revoke; the confirmation states both.
- **The confirmation shows last-used + IP** (PRD TA-E3) — "understand what you're breaking before you break it". The `--yes` path keeps the last-used line in the result so the information isn't lost by skipping the prompt.
- **Last-active-key detection**: if this was the app's only active key, the result names the consequence and the replacement command — the CLI does the "this pipeline will die" reasoning the web does (PRD TA-E3).
- **stdout**: the result. **stderr**: the confirmation prompt.
- **Idempotent**: revoking an already-revoked key is a success (exit 0), so cleanup scripts don't need state guarding.
