# Command — `bloberry init`

## Purpose & context

- **User goal**: one-command first-run setup — point the CLI at a server, pick a tenant, authenticate. The answer to "how does the CLI know the URL": `init` walks server → reachability → tenant → auth, so a fresh machine goes from `install` to a working `ls` in one command instead of `config set server` + `auth login` + `tenant use` by hand.
- **When they reach for it**: interactive first run; CI provisioning (fully non-interactive via flags).
- **Needs**: nothing pre-configured (it *writes* the config). Non-interactive path needs `--server` (and `--token` or an existing keychain session).
- **Data**: `server` + `tenant` config keys, the health endpoint for reachability, the auth flow.

## Signature

```
bloberry init [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `--server <url>` | string | — | if non-interactive | The Bloberry server URL to configure. |
| `--tenant <slug>` | string | — | no | Default tenant to set (prompted if multiple). |
| `--token <key>` | string | — | no | Authenticate with this access key (CI; stored in env, never the config file). |
| `--yes, -y` | bool | false | no | Accept defaults; no prompts. |

Global flags apply — `../README.md`.

## Help text

```
Set up the CLI for a Bloberry server in one command.

Walks: server URL → reachability check → tenant → authentication.
Interactive by default; fully non-interactive with --server and
--token (or an existing keychain session) — the CI path:

  bloberry init --server https://bloberry.example.com --tenant acme

Interactive mode:
  bloberry init

Flags:
      --server url   Bloberry server URL
      --tenant slug  Default tenant (prompted if you belong to several)
      --token key    Authenticate with an access key (CI)
  -y, --yes          Accept defaults, no prompts

Examples:
  bloberry init
  bloberry init --server https://bloberry.example.com
  BLOBERRY_TOKEN=blob_live_… bloberry init --server https://… --tenant acme --yes
```

## Output states

**Interactive success**

```
Server:   https://bloberry.example.com  (reachable ✓)
Tenant:   acme  (Acme Inc)   [1 of 3 — pick with --tenant]
Auth:     browser device flow — open the link, enter G H Q K-8 3 F 9
✓ Signed in as jane@acme.dev
✓ Config written: ~/.config/bloberry/config.yaml
Ready: run 'bloberry ls' to see your files.
```

**Non-interactive success (CI)**

```
✓ Config written: server=https://…, tenant=acme (via BLOBERRY_TOKEN)
Ready: run 'bloberry ls'.
```

**Server unreachable**

```
Error: could not reach https://wrong.example.com

Check the address or that the server is up. (backend_unreachable)
```

**`--json` (non-interactive success)**

```json
{"server":"https://bloberry.example.com","tenant":"acme","auth":"token","ready":true}
```

## Exit codes

| Code | When |
|---|---|
| `0` | Configured and authenticated; `ls` will work. |
| `1` | Setup failed (server unreachable, auth failed). |
| `2` | Bad invocation (missing `--server` in a non-TTY, bad URL). |
| `3` | Auth failed / no credentials in a non-TTY without `--token`. |
| `9` | Server reachable but the health check failed (backend unreachable). |

## Behavior notes

- **Non-interactive contract**: in a non-TTY, `--server` is required and either `--token` or an existing keychain session must authenticate — otherwise exit 3 with "pass `--server` and `--token` (or run interactively)". Never hangs on an invisible prompt.
- **Reachability first**: `init` checks the health endpoint before writing anything — a typo'd URL fails loudly here rather than as a confusing 401 on the first real command.
- **Writes the config file** via the same path `config set` uses (flag > env > file precedence holds); `--token` is **never written to the config file** (it belongs in `BLOBERRY_TOKEN` / the keychain — `cli/README.md` §Config secrets rule).
- **Idempotent**: re-running reconfigures (a new URL overwrites, a different tenant switches the default). Safe to run in a fresh CI checkout.
- **stdout**: the result. **stderr**: progress, the device-flow code.
- **TTY**: the device flow only opens a browser when a terminal is present; the interactive multi-tenant pick is a numbered prompt with a `--tenant` flag equivalent.
