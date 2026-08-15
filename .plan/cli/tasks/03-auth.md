# Task group — 03 auth

**Depends on:** `01-project-setup.md`, `02-core-infra.md`, `backend/tasks/04-domain-auth.md` (the endpoints). **Blocks:** every command that needs a session (everything except `config`/`completion`/`version`).

Two auth shapes per `cli/README.md` §Config: browser device flow for humans (keychain), `BLOBERRY_TOKEN` access key for CI.

- [ ] **`auth login`** — browser device flow per its designed file: prints a code, opens the browser, polls until confirmed, tokens to the **OS keychain** (`zalando/go-keyring`), never the config file. On a keychain-less headless box → "no keychain available — use `BLOBERRY_TOKEN`", never a silent plaintext fallback. Polling timeout → exit 1 with a re-run message.
- [ ] **`auth status`** — reports who/tenant/expiry **and which credential shape** (keychain vs `BLOBERRY_TOKEN`); exits 0 valid / 3 not-authenticated (the canonical preflight contract).
- [ ] **`auth logout`** — **revokes server-side** (deletes the Redis refresh session), then drops the local token. If the server is unreachable, drops the local token anyway and says the server session may live on. `--revoke-key` (BLOBERRY_TOKEN case) destroys the CI credential with a confirmation.
- [ ] **Refresh handling** — platform-aware TTLs per `backend/domains.md` §7 (CLI presents as a non-web platform); a stale keychain session refreshes transparently; a dead one returns exit 3 → "run `bloberry auth login`".
- [ ] **`BLOBERRY_TOKEN` precedence** — the env var wins over the keychain when present (CI never accidentally uses a stale desk session); `auth status` names which one it used.

**tests:** device-flow success/timeout; keychain store/load; logout revokes server-side (assert the server session is gone); `--revoke-key` confirmation; token-vs-keychain precedence; status exit 0/3 contract.
