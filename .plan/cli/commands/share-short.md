# Command — `bloberry share short`

## Purpose & context

- **User goal**: create a compact short URL (`/s/<slug>`) for an object (PRD M8) — shareable in chat, memorable enough to type.
- **When they reach for it**: interactive ("paste this in the group chat"), and scripts generating a stable share slug for a published asset.
- **Needs**: auth, tenant context, `share` permission. Slugs are unique per install (`ERD.md` share-links note) and served on the main domain (PRD D6).
- **Data**: `share_links` — kind `short`, `slug`, `expires_at` (nullable — short links can be permanent until revoked).

## Signature

```
bloberry share short <path> [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `<path>` | path | — | yes | `bloberry://` object path. |
| `--slug <slug>` | string | random | no | Request a specific slug (letters, digits, `-`; 4–32 chars). |
| `--ttl <duration>` | duration | none | no | Optional expiry; default short links live until revoked. |

## Help text

```
Create a short URL (/s/<slug>) for an object.

By default the slug is random and unguessable. With --slug you can
request a memorable one (it must still be available — slugs are unique
across the whole install, not per tenant). Short links live until
revoked unless --ttl is given.

Examples:
  bloberry share short bloberry://assets/hero.png
  bloberry share short bloberry://release/v2.1.0.zip --slug v2-1-0
```

## Output states

**Success**

```
✓ Short link created (sl_3d8b)
  URL:  https://bloberry.example.com/s/v2-1-0
```

**Slug taken (exit 8)**

```
Error: slug "v2-1-0" is already in use (name_conflict)

Run without --slug to get a random one.
```

**`--json`**

```json
{"id":"sl_3d8b","url":"https://bloberry.example.com/s/v2-1-0","slug":"v2-1-0","kind":"short","expires_at":null}
```

## Exit codes

| Code | When |
|---|---|
| `0` | Created. |
| `8` | Slug taken — retry without `--slug`. |
| `5` | Object doesn't exist. |
| `4` | Forbidden. |
| `2` | Bad slug (invalid chars, wrong length). |

## Behavior notes

- **Random slugs are unguessable by default** — a short URL is a capability (PRD D6), so the default is entropy, not a guessable sequence.
- **`--slug` collisions exit `8`** (conflict), distinct from a not-found `5`, so a script can loop on `--slug "$slug-$(date +%s)"` without misreading a 404.
- **Permanent by default**: no `--ttl` means "until revoked" — the help text says so explicitly, because "short" shouldn't be read as "ephemeral".
- **stdout**: the result / JSON. **stderr**: nothing on success.
- **Revocation** is `share revoke <id>`; a revoked or expired short link renders the `link-expired` 410 page (`web/mockup/link-expired.md`).
