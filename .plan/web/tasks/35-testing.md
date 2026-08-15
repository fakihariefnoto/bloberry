# Task group — 34 testing

**Depends on:** the pages it tests (`07`–`33`), `34-flows.md`. 

Named high-risk tests, not blanket "add tests." Tooling: the framework's test runner (Vitest for Vue + Vite, under Bun), Vue Test Utils, and the golden-output discipline shared with the CLI (`cli/README.md` §Testing — sample output blocks become fixtures).

- [ ] **Auth suite** — login success/`?next`/invalid-credentials-identical-message/network-failure; OTP rate-limit + attempt cap; reset-token expired state; invite new-user vs existing-email. (Covers `07`–`11`.)
- [ ] **Session expiry end-to-end** — an authenticated request 401s → redirect `login?next` → re-login returns to the original page (`33.1`).
- [ ] **Permission-aware rendering** — a viewer on `files`/`file-detail` sees write actions disabled with a reason, not hidden and not an error wall (PRD MV4). Test at the component level (`PermissionDenied`) and one page level.
- [ ] **Envelope parsing** — success-data, success-with-messages (toast fired), error-messages (mapped, never raw provider text), unknown-code fallback (`05-core-infra`).
- [ ] **Token handling** — 401 → refresh → retry-once; `key_revoked` terminal no-retry branch (`05-core-infra`).
- [ ] **DataTable** — four states on one table page; cursor pagination; selection + bulk bar; permission-disabled rows (`06`, `13`).
- [ ] **Upload queue** — persistence across navigation; determinate progress; per-file retry; `beforeunload` while in flight; quota/collision per-file outcomes (`06`, `33.2`).
- [ ] **Secret handling** — `key-created` modal: once-only, no-backdrop-dismiss, ack required (`18`); masked keys everywhere else (`18`, `06`).
- [ ] **Destructive confirms** — typed-name required for folder delete / key revoke / member removal / backend change / make-public; wrong name refused; real object count stated (PRD TA-E1) (`06`, `13`, `18`, `19`, `22`).
- [ ] **Jobs** — state-flip without reload; hidden-tab polling stops; terminal codes have no Retry; "target folder unchanged" on extraction failure (`16`).
- [ ] **Rate-card rule** — `unknown`, never $0, in both `usage` and `admin-usage`; backend rate-card edit re-renders estimates (`21`, `30`, `29`).
- [ ] **Pair-device auth** — QR token single-use (second verify fails); expiry auto-refresh; paired state no longer renders scannable; config download enforces passphrase strength; the passphrase is never sent to the server (`25`, `33.1a`/`33.1b`).
- [ ] **A11y checks** — focus management in the modals (`06`), announced errors on the auth screens, `role="progressbar"` values (`16`).
- [ ] **Golden-output** — where a component renders fixed output (CopyableCode confirmation, empty states, StatusPill text), assert exact strings so copy changes are reviewable.

**tests:** this group IS the test plan — each checkbox above maps to a runnable suite under `bun test`.
