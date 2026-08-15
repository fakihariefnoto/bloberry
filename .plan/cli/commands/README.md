# Commands — Bloberry

Each command is a `.md` file. This is a design artifact, not documentation written after the fact — for a CLI the terminal output *is* the interface, and a command whose output nobody designed gets whatever formatting the implementer improvised that day. Every command file has:

1. **Purpose & context** — what the user is trying to accomplish, when they reach for this command, and what it needs (args, config, auth).
2. **Signature** — usage line, positional args, and flags as a table (name, type, default, required, what it does).
3. **Help text** — the real `--help` output, including at least one worked example. This is the first UI anyone sees.
4. **Output states** — sample terminal output for each meaningful state (success, empty result, error, and `--json` where it applies), in a fenced code block with realistic values. Not just the happy path.
5. **Exit codes** — which code each outcome returns, consistent with `../README.md`'s table.
6. **Behavior notes** — confirmation prompts and their `--yes` equivalent, what goes to stdout vs stderr, TTY-only behavior (color/spinners), idempotency, and anything destructive.

See `command-example.md` for the format, and use the `generate-commands` skill to design the whole set.
