# Task group — 02 core infra

**Depends on:** `01-project-setup.md`, `backend/tasks/19-openapi.md` (the SDK). **Blocks:** every command file. The pieces every command depends on — retrofitting after ten commands means editing all ten.

- [ ] **Config resolution** — the precedence from `cli/README.md` §Config: **flag → env (`BLOBERRY_`) → config file → default**. `~/.config/bloberry/config.yaml` (XDG) / `%APPDATA%\bloberry\config.yaml`, overridable with `--config`. One resolver, consumed everywhere.
- [ ] **`config get/set/path`** implemented (details in `09-commands-config-misc.md`) — the ergonomic way users manage config instead of hand-editing YAML blind.
- [ ] **Output layer** — one place that renders results for humans and `--json`, honoring `--quiet` and TTY detection (color/spinners/progress only on a terminal; `NO_COLOR` and `TERM=dumb` honored). **stdout is data, stderr is everything else** (`cli/README.md` §Output). Centralizing this is what keeps the formats consistent; per-command `printf` is how they diverge.
- [ ] **Error wrapping** — "what failed + what to do next" (`no folder 'assets/v3' (run 'bloberry ls bloberry://assets'…)`), never a bare 404 or stack trace. Panics/traces only under `--verbose`.
- [ ] **Error → exit-code mapping** — the table from `cli/README.md` §Exit codes: 0 success · 1 generic failure · 2 invocation · 3 not-authenticated · 4 forbidden · 5 not-found · 6 quota · 7 partial · 8 conflict · 9 backend unreachable. **Code 7 is the one that matters** — recursive ops that partially fail must exit 7 (never 0 = CI green, never 1 = "nothing happened").
- [ ] **API client via the SDK** — the shared `sdk/go` client + envelope parsing (`{data?, messages?}`), one layer, used by every command. Never per-command HTTP.

**tests:** precedence resolution (flag beats env beats file beats default); stdout/stderr separation (`--json | jq` works); each exit-code mapping fires on its failure kind; a partial-failure path returns 7.
