# Task group — 09 commands: config + misc + init

**Depends on:** `02-core-infra` (config resolver, output layer), `03-auth`. One checkbox per command; acceptance = "output matches the sample in `cli/commands/<file>.md`".

## `init` (`cli/commands/init.md`)

- [ ] **One-command first-run** — walks server URL → reachability (health endpoint) → tenant → auth (device flow interactive / `--token` or existing keychain session non-interactive). Writes the config via the same path `config set` uses.
- [ ] **Non-interactive contract** — in a non-TTY, `--server` required + `--token` or an existing session, else exit 3 with "pass `--server` and `--token`"; never hangs on an invisible prompt.
- [ ] **`--token` never written to the config file** (env/keychain only, per the secrets rule).
- [ ] **Reachability fails loudly before writing anything** — a typo'd URL errors here, not as a confusing 401 on the first real command.
- [ ] **Idempotent** — re-running reconfigures; safe in a fresh CI checkout.
- [ ] **`--json`** emits `{server, tenant, auth, ready}`.

## `config get` (`cli/commands/config-get.md`)

- [ ] Read a config key with **`--source` showing which precedence layer won** (flag > env > file > default — the debug answer); unset key → exit 0 with "(unset — default: …)"; unknown key → exit 2; secret-like keys rejected.

## `config set` (`cli/commands/config-set.md`)

- [ ] Write a config value, **validated on write** (`output bogus` refused with the allowed values); `--unset` removes; **secret-like keys refused** ("Secrets live in the keychain or BLOBERRY_TOKEN, never in the config file") — the file stays clean by construction.

## `config path` (`cli/commands/config-path.md`)

- [ ] Print the resolved path (respecting `--config`); missing file → exit 1 with a "create one with `config set`" hint; `--show-if-missing` flips to exit 0 + path for tooling.

## `completion` (`cli/commands/completion.md`)

- [ ] Generate bash/zsh/fish/powershell scripts to stdout (nothing else); unknown shell → exit 2 with the supported list. **Dynamic completions** where it matters: `bloberry://<TAB>` completes remote folders via the API, `key revoke <TAB>` completes key IDs (static fallback when unauthenticated, never an error).

## `version` (`cli/commands/version.md`)

- [ ] version/commit/date from `-ldflags` (never hardcoded); `--check` asks for a newer release and prints the **install-method-aware** update command (`brew upgrade` / `scoop update` / `apt install --only-upgrade` / a download link). **No `self-update`** — overwriting a package-managed binary breaks the package database (`cli/README.md` §Update). `--check` is offline-safe (version still prints, exit 0).

**tests:** per command — happy path, error, `--json` where it applies; `config get --source` precedence; `config set` value validation + secret-key refusal; `version --check` offline-safe.
