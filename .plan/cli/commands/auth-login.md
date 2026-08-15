# Command — `bloberry auth login`

## Purpose & context

- **User goal**: authenticate a human at a terminal via the browser device flow — no token pasting, no password in the shell. Tokens land in the OS keychain (`cli/README.md` §Config).
- **When they reach for it**: first use of the CLI, or after the token expired (`exit 3` remediation: "run `bloberry auth login`").
- **Needs**: the `server` from config (`~/.config/bloberry/config.yaml` or `BLOBERRY_SERVER`).
- **Data**: device-flow endpoints (device code, user code, polling URL) from the auth API.

## Signature

```
bloberry auth login [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `--server <url>` | string | config | no | Override the server to authenticate against. |

## Help text

```
Sign in to Bloberry in the browser and store the session in the OS
keychain.

Prints a one-time code, opens the default browser, and polls until you
confirm. Nothing secret is typed into the terminal. For CI, don't use
this — set BLOBERRY_TOKEN to an access key instead.

Examples:
  bloberry auth login
  bloberry auth login --server https://bloberry.example.com
```

## Output states

**Success (interactive)**

```
Opening https://bloberry.example.com/device …
Enter this code in the browser:  G H Q K-8 3 F 9
Waiting for confirmation …
✓ Signed in as jane@acme.dev
  Tenant:   acme (default)
  Token:    stored in macOS keychain (expires in 90 days)
```

**Already signed in**

```
Already signed in as jane@acme.dev. Run 'bloberry auth status' to check
expiry, or 'bloberry auth logout' first to switch accounts.
```

**Failure — server unreachable**

```
Error: could not reach https://bloberry.example.com/device

Check that the server is up (run 'bloberry auth status' with --verbose
for details) or that the 'server' config value is right:
  bloberry config get server
```

## Exit codes

| Code | When |
|---|---|
| `0` | Signed in; token stored. |
| `1` | Device flow failed (server unreachable, flow timed out). |
| `2` | Bad invocation (missing `--server` and no config). |
| `3` | — (this is the command that *fixes* 3; it never exits 3 itself). |

## Behavior notes

- **Browser device flow** only — there is no `--password` flag and there must never be one (a password flag is a password in the process list and shell history; `cli-defaults.md` keychain rule).
- **Tokens to the keychain** via `zalando/go-keyring`; on a headless box with no keychain, login fails with "no keychain available — use `BLOBERRY_TOKEN`" rather than silently degrading to a plaintext file.
- **stdout**: the code and the result. **stderr**: the polling progress ("Waiting for confirmation …").
- **TTY**: the opening of the browser is the only TTY-coupled step; in a non-TTY it prints the code + URL and polls without opening anything (a user can complete the flow from another device).
- **Polling timeout** (~2 min) exits 1 with "timed out waiting for confirmation — run again".
- **Idempotent**: re-running when already signed in is a no-op that reports the existing session.
