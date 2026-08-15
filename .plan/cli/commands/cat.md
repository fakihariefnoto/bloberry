# Command — `bloberry cat`

## Purpose & context

- **User goal**: stream an object's bytes to stdout — the unix-y building block for `bloberry cat … | grep`, `| jq`, `| head`.
- **When they reach for it**: interactively to eyeball a small file; in scripts to pipe an object into another tool without downloading it first.
- **Needs**: auth, tenant context, read permission.
- **Data**: `objects` bytes (streamed, never buffered whole into memory — `architecture.md` streaming discipline).

## Signature

```
bloberry cat <path> [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `<path>` | path | — | yes | `bloberry://` object path. |

## Help text

```
Stream an object's contents to stdout.

Bytes go straight to stdout for piping; nothing else is written there.
Progress and errors go to stderr. Use --json with --stat for metadata
without streaming.

Examples:
  bloberry cat bloberry://logs/app.log | grep ERROR
  bloberry cat bloberry://data/latest.json | jq .version
```

## Output states

**Success (streams to stdout; nothing else)**

```
<object bytes on stdout, e.g.>
{"version":"2.1.0","build":412}
```

**Error — not found**

```
Error: no object "bloberry://logs/old.log" (run 'bloberry ls bloberry://logs' to see what exists)
```

## Exit codes

| Code | When |
|---|---|
| `0` | Streamed. |
| `5` | Object doesn't exist. |
| `4` | Forbidden to read. |
| `3` | Not authenticated. |

## Behavior notes

- **stdout is pure bytes** — the whole point of a `cat`-style command. Errors, warnings and progress go to stderr and never corrupt the pipe.
- **Binary safety**: no post-processing, no color, no truncation; `cat` passes bytes through untouched. A user who wants human metadata uses `stat`, not `cat`.
- **TTY caveat**: when stdout is a terminal, the shell's own behavior applies (a binary blob may garble the screen) — the CLI itself does not warn; this matches `cat`'s contract.
- **Streaming**: reads via the download endpoint in chunks (`io.Copy`); memory stays flat even for multi-GB objects (PRD G10 discipline).
