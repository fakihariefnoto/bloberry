# Command — `bloberry completion <shell>`

## Purpose & context

- **User goal**: generate shell completions (bash/zsh/fish/powershell) — the dynamic ones included (`cp bloberry://<TAB>` completes remote folders via the API, `key revoke <TAB>` completes key IDs, `cli/README.md` §Completions).
- **When they reach for it**: first setup, after an install that didn't drop the completion file (packages do install it; this is for manual installs and regenerating).
- **Needs**: auth only for the dynamic completions (remote folder/key lookups); static subcommand completion works unauthenticated.
- **Data**: the command tree + (for dynamic completion) the API.

## Signature

```
bloberry completion <shell> [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `<shell>` | string | — | yes | `bash`, `zsh`, `fish`, `powershell`. |

## Help text

```
Generate shell completion script for the given shell.

Prints to stdout — pipe it into the shell's completion dir, e.g.:
  bloberry completion bash > /etc/bash_completion.d/bloberry
  bloberry completion zsh  > "${fpath[1]}/_bloberry"

Completions are dynamic where it matters: remote paths complete
folders from the server, and key/app/job IDs complete from the API.

Examples:
  bloberry completion bash
  bloberry completion fish
```

## Output states

**Success (stdout is the script; nothing else)**

```
# bash completion for bloberry                            _bloberry() { … }
```

**Bad shell (exit 2)**

```
Error: unknown shell "tcsh" — supported: bash, zsh, fish, powershell
```

## Exit codes

| Code | When |
|---|---|
| `0` | Script emitted to stdout. |
| `2` | Unknown shell. |

## Behavior notes

- **stdout is the script, nothing else** — any progress/warning would corrupt the pipe target; stderr stays empty on success.
- **TTY**: no color (a script, not a display).
- **The dynamic completions need auth**: an unauthenticated `bloberry://<TAB>` falls back to static completion rather than erroring — completing folders is a convenience, not a gate.
- **Packages install this automatically** (Homebrew/.deb/.rpm drop the file, `cli/README.md` §Completions) — the command exists for manual installs; the help's examples are the manual paths.
