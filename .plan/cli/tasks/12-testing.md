# Task group — 12 testing

**Depends on:** the command groups, `11-completions-man` (for dynamic-completion tests). **Blocks:** `13-release-pipeline` (a release of untested code isn't a release).

**Golden-file tests** per `cli/README.md` §Testing: each command's sample-output block in `cli/commands/*.md` becomes a fixture in `testdata/golden/` — the design artifact and the test assertion describe the same thing.

- [ ] **Golden fixtures generated** — one fixture per command's designed sample block: happy path, empty result, one input-error case, and `--json` where it exists. The `commands/*.md` blocks are the source (they were written with real values for exactly this).
- [ ] **Partial-failure fixtures** — for `cp`, `rm`, `sync`: a fixture asserting exit 7 + the per-file summary format (stderr).
- [ ] **Exit-code contract tests** — each failure kind maps to its designed code (3 auth, 4 forbidden, 5 not-found, 6 quota, 8 conflict, 9 backend) and the test asserts the code, not just the message.
- [ ] **stdout/stderr separation tests** — `--json | jq` parses on every `--json`-capable command; no progress/decoration ever lands on stdout.
- [ ] **TTY gating tests** — color/progress only when attached to a terminal; `NO_COLOR` and `TERM=dumb` honored; piped output prints one line per file, not spinner redraws.
- [ ] **Dynamic-completion tests** — `cp bloberry://<TAB>` returns remote folders (authenticated) and static-completes without error (unauthenticated).
- [ ] **`--quiet` honored** — suppresses the non-data decoration on commands that support it.

**tests:** `make test` green with the golden fixtures; a fixture drift (output changed without updating the fixture) fails loudly — the check that keeps documented output honest.
