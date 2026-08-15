# Command — `bloberry share list`

## Purpose & context

- **User goal**: see every share link on an object (or in the tenant) — active, expired, revoked — with hit counts (the "is this link still being used" answer).
- **When they reach for it**: interactive auditing before revoking; scripts that collect all URLs for a file.
- **Needs**: auth, tenant context, read permission.
- **Data**: `share_links` — kind, slug/token, target, `expires_at`, `hit_count`, `last_accessed_at`, `revoked_at`.

## Signature

```
bloberry share list [path] [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `[path]` | path | all tenant shares | no | Restrict to links on this object/folder. |
| `--status <s>` | string | all | no | Filter: `active`, `expired`, `revoked`. |

## Help text

```
List share links, optionally scoped to one object or folder.

Shows kind, URL, expiry and hit count — the information you need
before revoking something. Active links sort by hits descending so
the load-bearing ones are at the top.

Examples:
  bloberry share list
  bloberry share list bloberry://assets/hero.png
  bloberry share list --status active --json
```

## Output states

**Success (object-scoped)**

```
hero.png  (f_8Kd2pQxL31A)
  sl_7f2a91  short   https://…/s/v2-1-0         12 hits  last 2h   active
  sl_3d8b    signed  https://…/sl/9fK2…         3 hits   last 1h   active (exp 4h)
  sl_9c1f    signed  https://…/sl/…             45 hits  last 3d   expired
```

**No links**

```
No share links for hero.png. Create one:
  bloberry share link bloberry://assets/hero.png
```

**`--json`**

```json
[{"id":"sl_7f2a91","kind":"short","url":"https://…/s/v2-1-0","hits":12,"last_accessed_at":"2026-03-12T12:03:00Z","revoked":false}]
```

## Exit codes

| Code | When |
|---|---|
| `0` | Listed (including an empty result). |
| `5` | The scoping path doesn't exist. |
| `4` | Forbidden. |

## Behavior notes

- **Hits-first by default** — the revoke decision ("is this load-bearing?") is the first column in active mode.
- **stdout**: the listing. **stderr**: nothing on success.
- **Empty state** names the create command instead of printing an empty table — the same "an empty table tells the user nothing" rule as every list.
- **URLs are truncated middle** in the human view (full URL always in `--json`); a long link is never allowed to wrap the table.
- **No pagination in v1** — tenant share lists are small enough to print whole; if they ever grow past a page, this gets cursor pagination with the same output shape.
