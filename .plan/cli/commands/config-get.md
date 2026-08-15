# Command — `bloberry config get`

## Purpose & context

- **User goal**: read a config value — the debugging half of config management (`cli-defaults.md` "don't make users hand-edit YAML blind"; `cli/README.md` precedence flag → env → file → default).
- **When they reach for it**: "why is the CLI hitting the wrong server" — the answer is in precedence, and this command shows the effective value and where it came from.
- **Needs**: nothing (reads config; no auth).
- **Data**: config keys — `server`, `tenant`, `output`, `color`, and any future key.

## Signature

```
bloberry config get <key> [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `<key>` | string | — | yes | Config key (`server`, `tenant`, `output`, `color`). |
| `--source` | bool | false | no | Also show which layer supplied the value. |

## Help text

```
Show the effective value of a config key.

Respects precedence: flag > environment (BLOBERRY_*) > config file >
default. Use --source to see which layer won — essential when
debugging "why is it using that server".

Examples:
  bloberry config get server
  bloberry config get server --source
  bloberry config get tenant --json
```

## Output states

**Show (with source)**

```
server  https://bloberry.example.com   (from config file: ~/.config/bloberry/config.yaml)
```

**From environment**

```
server  https://staging.bloberry.example.com   (from environment: BLOBERRY_SERVER)
```

**Unset key**

```
tenant  (unset — default: none)
```

**`--json`**

```json
{"key":"server","value":"https://bloberry.example.com","source":"config-file","path":"~/.config/bloberry/config.yaml"}
```

## Exit codes

| Code | When |
|---|---|
| `0` | Read (including an unset key — that's a valid answer). |
| `2` | Unknown key (not a valid config key). |

## Behavior notes

- **`--source` is the debug flag** — the precedence chain is the thing that confuses people, and naming the winning layer ends the confusion in one line.
- **stdout**: the value / JSON. **stderr**: nothing on success.
- **An unset key exits 0 with "(unset — default: …)"** — not an error; the effective value is the default.
- **Secrets never readable here** — config never holds secrets (keychain/env do), so there's no `config get token` and the key set rejects secret-like names.
