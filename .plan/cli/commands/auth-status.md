# Command — `bloberry auth status`

## Purpose & context

- **User goal**: answer "who am I, which tenant, does my token still work, when does it expire" — the first command anyone runs when the CLI starts returning 3s.
- **When they reach for it**: debugging auth, checking which auth shape is active (keychain vs `BLOBERRY_TOKEN`), and in CI as a cheap preflight.
- **Needs**: current credentials; makes one authenticated call to confirm validity.

## Signature

```
bloberry auth status [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| — | — | — | — | No args. |

## Help text

```
Show who you're authenticated as, which tenant is current, and token
expiry.

Reports which credential shape is in use: the OS keychain session, or
BLOBERRY_TOKEN. Exits 0 when the session is valid, 3 when not.

Examples:
  bloberry auth status
  bloberry auth status --json
```

## Output states

**Success (keychain session)**

```
Signed in as  jane@acme.dev        (user_8f2a1c)
Credential:   keychain session
Tenant:       acme (default)
Expires:      2026-06-12 (in 89 days)
```

**Success (`BLOBERRY_TOKEN` in CI)**

```
Signed in as  acme-cms (application)   (app_3d9f)
Credential:   BLOBERRY_TOKEN (blob_live_••••4f2a)
Tenant:       acme
Scope:        2026/ · read/write
Key expires:  never
```

**Not authenticated**

```
Not signed in.

Run 'bloberry auth login' (interactive) or set BLOBERRY_TOKEN (CI).
```

**`--json` (success)**

```json
{"principal_type":"user","principal_id":"user_8f2a1c","email":"jane@acme.dev","credential":"keychain","tenant":"acme","expires_at":"2026-06-12T09:00:00Z"}
```

## Exit codes

| Code | When |
|---|---|
| `0` | Session valid. |
| `3` | Not authenticated or expired — the canonical way scripts check "do I need to re-auth". |
| `1` | Check itself failed (server unreachable, ambiguous). |

## Behavior notes

- **Exit 3 is the contract**: a script preflights with `bloberry auth status` and branches on 0 vs 3 — the exact code the rest of the CLI returns on auth failures, so the preflight and the real commands agree.
- **Names the credential shape** explicitly — a CI run that accidentally uses a stale keychain token is diagnosed by this line ("Credential: keychain session") instead of being mysterious.
- **stdout**: the status (data). **stderr**: nothing on success.
- **Does not refresh the token** — it reports validity as of now; a token a minute from expiry is reported as valid (the next real command will exercise the refresh path).
