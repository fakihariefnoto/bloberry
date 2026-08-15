# Task group — 10 cross-command flows

**Depends on:** the command groups it chains (`04`–`09`). The journeys a per-command split loses.

## 10.1 First run (install → first command)

- [ ] `bloberry init` → server URL + reachability → `tenant use` equivalent (or pick) → device flow / keychain → `bloberry ls` lists the tenant root. **One command instead of `config set` + `auth login` + `tenant use`** — the whole point of `init`. Chain connects end-to-end; a fresh machine goes from install to a working `ls`.
- [ ] **CI first run** — `BLOBERRY_TOKEN=… bloberry init --server https://… --tenant acme --yes` non-interactively; `init` exits 0 and `ls` works; the token is never written to the config file.
- [ ] **Reachability fails loudly** — a typo'd URL errors at `init`, not as a confusing 401 on the first real command.
- [ ] **No-auth first run** — `bloberry ls` without a session returns exit 3 with "run `bloberry auth login` or `bloberry init`" (the designed error text, not a raw 401).

## 10.2 CI / non-interactive

- [ ] **Every prompt has a flag equivalent** — `rm -r`, `share public`, `key revoke`, `app delete`, `grant revoke`, `admin backend test`-adjacent destructives all run under `BLOBERRY_TOKEN` + `--yes` with zero stdin. In a non-TTY without `--yes`, the command fails with "pass `--yes`", never hangs on an invisible prompt.
- [ ] **`BLOBERRY_TOKEN` keychain never used in CI** — the env var wins; `auth status` confirms which shape.
- [ ] **Publish journey** — `cp ./dist bloberry://assets/v2 -r` → exit 7 on partial (summary + idempotent rerun) → `share link` the artifact → the URL lands in a PR comment. The designed exit-code contract holds end to end.

## 10.3 Key rotation (the incident path)

- [ ] `key list --status active` → identify the leak → `key revoke ak_…` (last-used + IP shown, no Undo) → `key create --app app_…` (new secret captured once) → new `BLOBERRY_TOKEN` in CI. Across `06`; the last-active-key warning fires if applicable.

## 10.4 Backend incident (platform admin)

- [ ] `admin backend test sb_…` → exit 9 (page) → `admin backend rate-card sb_… --set …` or credential fix → retest → healthy. The exit-9 signal + the fix-command hints chain end-to-end.

**tests:** one end-to-end golden test per flow; the CI path asserts zero stdin needed and correct exit codes throughout.
