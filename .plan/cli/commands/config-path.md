# Command — `bloberry config path`

## Purpose & context

- **User goal**: find the config file — the answer to "where is this tool storing its settings" (`cli/README.md` XDG location, `--config` override).
- **When they reach for it**: backing up config, or when an error message references the file.
- **Needs**: nothing.
- **Data**: the resolved config path.

## Signature

```
bloberry config path [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `--show-if-missing` | bool | false | no | Print the path even if the file doesn't exist yet. |

## Help text

```
Print the resolved config file path.

By default exits non-zero when no config file exists; pass
--show-if-missing to print the would-be path anyway.

Examples:
  bloberry config path
  bloberry config path --show-if-missing
```

## Output states

**Exists**

```
~/.config/bloberry/config.yaml
```

**Missing**

```
Error: no config file at ~/.config/bloberry/config.yaml

(Defaults are being used. Create one with:
  bloberry config set server https://bloberry.example.com)
```

## Exit codes

| Code | When |
|---|---|
| `0` | File exists (path printed). |
| `1` | No config file — with `--show-if-missing` still exit 0 and print the path. |
| `2` | Bad invocation. |

## Behavior notes

- **stdout**: exactly the path, nothing else — script-friendly (a script captures the path and opens/parses it).
- **Missing-file behavior is explicit**: exit 1 (the file isn't there) with a pointer to how a config gets created. `--show-if-missing` flips it to exit 0 + path for tools that want the would-be location.
- **The path respects `--config <path>`** — an overridden path prints that path, not the default.
