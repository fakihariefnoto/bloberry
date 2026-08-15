# Command — `bloberry share link`

## Purpose & context

- **User goal**: create a temporary, revocable signed link to an object (PRD M7) — the CLI's scriptable sharing.
- **When they reach for it**: interactive ("share this file with the team"), and in CI/scripts to emit a download URL (e.g. a build artifact link in a PR comment).
- **Needs**: auth, tenant context, `share` permission on the object. `--ttl` sets the expiry; revocation is a `share revoke`.
- **Data**: `share_links` — kind `signed`, token, `expires_at`, target object. The signed URL is returned once; the underlying record is what revocation targets.

## Signature

```
bloberry share link <path> [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `<path>` | path | — | yes | `bloberry://` object path. |
| `--ttl <duration>` | duration | `24h` | no | How long the link lives (e.g. `1h`, `30m`, `7d`). |
| `--open` | bool | false | no | Open the created URL in the browser (TTY only). |

## Help text

```
Create a temporary signed link to an object.

The link works until it expires or is revoked (bloberry share revoke).
Downloads are redirected to a short-lived provider URL, so revoking
takes effect on the next request to the signed link — but an already-
issued provider URL can outlive the link briefly (see docs, R11).

Examples:
  bloberry share link bloberry://assets/hero.png --ttl 1h
  bloberry share link bloberry://release/build.zip --ttl 7d --json
```

## Output states

**Success (interactive)**

```
✓ Signed link created (sl_7f2a91)
  URL:  https://bloberry.example.com/sl/9fK2mQxL31A8
  TTL:  1 hour (expires 2026-03-12 15:03 UTC)
```

**`--json`**

```json
{"id":"sl_7f2a91","url":"https://bloberry.example.com/sl/9fK2mQxL31A8","kind":"signed","expires_at":"2026-03-12T15:03:00Z","revoked":false}
```

## Exit codes

| Code | When |
|---|---|
| `0` | Link created. |
| `5` | Object doesn't exist. |
| `4` | Forbidden — no `share` permission. |
| `2` | Bad `--ttl` (e.g. negative, over the max). |

## Behavior notes

- **The URL is the output** — it's the whole point; both human and `--json` forms put it prominently, never buried under other fields.
- **`--ttl` semantics**: human-friendly durations (`1h`, `30m`, `7d`); max TTL is config (default 30 days) and a longer request is refused with "max TTL is 30d".
- **Revocation honesty**: the help text states the R11 caveat (a redirect path can outlive the link briefly) rather than over-promising "revoked means revoked" — same honesty as the web UI.
- **stdout**: the result line / JSON. **stderr**: nothing on success. `--open` (TTY only) opens the URL and notes "opened in browser".
- **Not idempotent**: each run creates a new link (links are cheap, independent capabilities); scripts that want one URL per artifact should capture and reuse the output.
