# Command — `bloberry config set`

## Purpose & context

- **User goal**: write a config value without hand-editing YAML (`cli-defaults.md`; `cli/README.md` config section). The `tenant` key is also managed ergonomically by `tenant use`.
- **When they reach for it**: pointing the CLI at a different server, changing the output format.
- **Needs**: nothing (writes config; no auth).
- **Data**: config keys — `server`, `tenant`, `output` (`table`/`json`), `color` (`auto`/`always`/`never`). Secret-like names are rejected.

## Signature

```
bloberry config set <key> <value> [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `<key>` | string | — | yes | Config key. |
| `<value>` | string | — | yes | New value. |
| `--unset` | bool | false | no | Remove the key instead of setting it. |

## Help text

```
Set a config value in the config file.

Validated against known keys and value shapes. Secrets are never
stored here — use the keychain (auth login) or BLOBERRY_TOKEN.

Examples:
  bloberry config set server https://bloberry.example.com
  bloberry config set output json
  bloberry config set tenant folio
```

## Output states

**Success**

```
✓ config: server = https://bloberry.example.com
  (wrote ~/.config/bloberry/config.yaml)
```

**Unset**

```
✓ config: tenant removed (now uses default)
```

**Rejected — secret-like key**

```
Error: "token" is not a config key. Secrets live in the keychain or
BLOBERRY_TOKEN, never in the config file.
```

## Exit codes

| Code | When |
|---|---|
| `0` | Set or unset. |
| `2` | Unknown key, invalid value, or missing args. |

## Behavior notes

- **Validated on write** — `output bogus` is refused with the allowed values, not silently stored for a later command to choke on.
- **Secret-like keys are refused** (the message above) — the config file must stay clean by construction, not by convention.
- **stdout**: the result. **stderr**: nothing on success.
- **File writing** respects `--config <path>` override; the result states which file was written.
