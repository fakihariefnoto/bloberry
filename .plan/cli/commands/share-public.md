# Command — `bloberry share public`

## Purpose & context

- **User goal**: flip an object's `visibility` to public (PRD M6) — served at a stable public URL, no token needed.
- **When they reach for it**: interactive ("make this download-able by anyone"), and scripts publishing static assets.
- **Needs**: auth, tenant context, `write`/`share` on the object. This is a visibility change on the object, not a link record.
- **Data**: `objects.visibility` (private | public). The public URL is stable (the `file_id` never changes, PRD M4) — unlike a signed link, it can't "expire", only be re-privatized.

## Signature

```
bloberry share public <path> [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `<path>` | path | — | yes | `bloberry://` object path. |
| `--yes, -y` | bool | false | no | Skip the confirmation. |

## Help text

```
Make an object publicly readable at a stable URL.

Public is effectively irreversible once the URL is copied — it has no
expiry and can't be revoked like a signed link, only un-published.
That's why this confirms before acting.

Examples:
  bloberry share public bloberry://assets/hero.png
  bloberry share public bloberry://assets/badge.svg --yes
```

## Output states

**Success (interactive, after confirm)**

```
✓ hero.png is now public
  URL: https://bloberry.example.com/o/f_8Kd2pQxL31A
  (Un-publish any time from the web dashboard — the CLI only makes public)
```

**`--json`**

```json
{"id":"f_8Kd2pQxL31A","name":"hero.png","visibility":"public","url":"https://bloberry.example.com/o/f_8Kd2pQxL31A"}
```

## Exit codes

| Code | When |
|---|---|
| `0` | Made public. |
| `5` | Object doesn't exist. |
| `4` | Forbidden. |
| `2` | Bad invocation. |

## Behavior notes

- **Confirmation required** (unless `--yes`): the prompt states the consequence — "public has no expiry; anyone with the URL can view it, and it can't be revoked, only un-published." Non-TTY without `--yes` fails with "pass `--yes`".
- **Un-publishing is a web-dashboard action in v1** — the CLI makes public but has no `share private`; the success output says so rather than inventing a command. If un-publishing from a terminal becomes a real need, it's a one-command addition to the tree.
- **The public URL is stable** (`/o/<file_id>`) — it survives renames and moves (PRD M4); the help and output don't over-promise revocation because there is none, only un-publication (the redirect caveat applies, ADR-3/R11).
- **stdout**: the result / JSON. **stderr**: the confirmation.
- **Not idempotent in effect**: making an already-public object public again is a no-op that still succeeds (exit 0) — the state is already what was asked.
