# Command — `bloberry auth logout`

## Purpose & context

- **User goal**: end the session **server-side**, not just delete a local file. A logout that leaves a live refresh token on the server isn't a logout (`cli/README.md` §Config; `backend/domains.md` §4.3 rotation makes the old refresh token dead after one use, but the logout must revoke it explicitly).
- **When they reach for it**: leaving a shared machine, or switching accounts.
- **Needs**: the current token (keychain or `BLOBERRY_TOKEN`).

## Signature

```
bloberry auth logout [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `--revoke-key` | bool | false | no | Also revoke the underlying access key (only when logged in via `BLOBERRY_TOKEN`). |

## Help text

```
Sign out and revoke the session on the server.

Deletes the local token and tells the server to invalidate the
refresh session. When authenticating via BLOBERRY_TOKEN, pass
--revoke-key to revoke the access key itself (irreversible — CI
pipeline will fail until a new key is issued).

Examples:
  bloberry auth logout
  bloberry auth logout --revoke-key
```

## Output states

**Success (keychain session)**

```
✓ Signed out. Server session revoked; local token removed.
```

**Success (`--revoke-key`)**

```
✓ Signed out and revoked access key blob_live_••••4f2a.

This key no longer works anywhere. Issue a new one if this was a CI
credential:
  bloberry key create
```

**No session**

```
Not signed in — nothing to log out. (exit 0)
```

## Exit codes

| Code | When |
|---|---|
| `0` | Signed out (including the no-session no-op). |
| `1` | Server revocation failed (e.g. unreachable) — local token still removed, message says so. |
| `2` | Bad invocation. |

## Behavior notes

- **Revokes server-side, always**: the refresh token is invalidated on the server before the local token is dropped. If the server is unreachable, the local token is still removed (so the user is at least logged out locally) and the message states the server session may live on.
- **`--revoke-key` is dangerous**: it's the one path that destroys a CI credential. It requires a confirmation unless `--yes`; the confirmation names the key's prefix + last-4.
- **stdout**: the result. **stderr**: the confirmation (if any).
- **No session is a success** (exit 0), not an error — logging out when already logged out is fine and scripts shouldn't treat it as failure.
- **In a non-TTY** without `--yes`, `--revoke-key` fails with "pass `--yes`" rather than hanging.
