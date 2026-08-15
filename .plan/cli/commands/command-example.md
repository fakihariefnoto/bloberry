# Command — `<cli> <noun> <verb>`

## Purpose & context

- **User goal**: what the user is trying to accomplish with this command, one line.
- **When they reach for it**: the situation that leads here — interactively while exploring, or from a script/CI job. This changes the design: a command run in CI needs every prompt to have a flag equivalent and its output to survive being piped.
- **Needs**: which config values, auth state, and prior commands this depends on (e.g. "requires `auth login`", "reads `project` from config if `--project` is omitted").
- **Data**: the actual fields this command reads or writes — pull from `../../PRD.md` / `../../ERD.md`. Real field names, not placeholders; if they aren't decided yet, flag it as an open question rather than inventing them.

## Signature

```
<cli> <noun> <verb> [<positional>] [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `<positional>` | string | — | yes | What it identifies. |
| `--flag <value>` | string | from config | no | What it changes. |
| `--yes, -y` | bool | false | no | Skip the confirmation prompt (needed for scripts). |

Global flags (`--help`, `--version`, `--config`, `--json`, `--quiet`, `--verbose`, `--no-color`) apply here too — see `../README.md`; don't restate them per command.

## Help text

The real `--help` output, with at least one worked example — this is the first UI anyone sees:

```
Usage: <cli> <noun> <verb> [<positional>] [flags]

One-line summary of what this does.

A short paragraph of context when the summary isn't enough — what it affects,
what it doesn't, anything surprising.

Flags:
  --flag string   What it changes (default: from config)
  -y, --yes       Skip the confirmation prompt

Examples:
  # The common case
  <cli> <noun> <verb> my-thing

  # Machine-readable, for scripts
  <cli> <noun> <verb> my-thing --json
```

## Output states

One block per meaningfully different state, with realistic values — not `foo`/`bar`. Delete the ones that don't apply, and add domain-specific ones that do.

**Success**

```
✓ Created project "atlas" (prj_8f2a1c)
  Region:   eu-west-1
  Endpoint: https://atlas.example.com
```

**Empty result** (for anything that lists) — say what to do next, since an empty table tells the user nothing:

```
No projects yet. Create one with:
  <cli> project create <name>
```

**Error** — what failed and the next action, on stderr, per `templates/cli-defaults.md`:

```
Error: no project named "atals"

Did you mean "atlas"? Run `<cli> project list` to see all projects.
```

**`--json`** — only the payload on stdout, no decoration, so `| jq` works:

```json
{"id":"prj_8f2a1c","name":"atlas","region":"eu-west-1","created_at":"2026-01-14T09:31:00Z"}
```

## Exit codes

| Code | When |
|---|---|
| `0` | Succeeded. |
| `1` | The operation failed (not found, permission denied, API error). |
| `2` | Bad invocation (missing arg, unknown flag). |

Consistent with `../README.md`'s table — scripts branch on these, so a command inventing its own meaning for `1` is a bug.

## Behavior notes

- **Confirmation**: whether this prompts before acting, and that `--yes` skips it. In a non-TTY without `--yes`, fail with a clear message rather than hanging on an invisible prompt.
- **stdout vs stderr**: the result goes to stdout; progress, spinners and warnings go to stderr.
- **TTY-only**: color and progress indicators only when attached to a terminal.
- **Idempotency / destructiveness**: safe to re-run? What's irreversible? Anything destructive says exactly what will be lost before asking.
