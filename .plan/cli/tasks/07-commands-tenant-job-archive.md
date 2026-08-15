# Task group — 07 commands: tenant + job + archive

**Depends on:** `02-core-infra`, `03-auth`. One checkbox per command; acceptance = "output matches the sample in `cli/commands/<file>.md`".

## `tenant list` (`cli/commands/tenant-list.md`)

- [ ] List the tenants the user belongs to with roles; `* (default)` marks the current tenant; only the user's own memberships (never cross-tenant, PRD M2).

## `tenant use` (`cli/commands/tenant-use.md`)

- [ ] Set the default tenant, **persisted to the config file**; membership-checked (not a member → exit 4 with "run `tenant list`"); idempotent; help nudges `--tenant`/fully-qualified paths for scripts (ambient state is what CI shouldn't rely on).

## `tenant usage` (`cli/commands/tenant-usage.md`)

- [ ] Stored/objects/egress/est. cost from the metering snapshots + rate card; **the "unknown, never $0" rule** (PRD M18) with the exact `admin backend rate-card` fix command in the output; egress labeled estimated (±10%).

## `job list` (`cli/commands/job-list.md`)

- [ ] List jobs with kind/state/progress/failure reason; **honest progress** (`47/184 files`, never a fake bar); live-updating progress only on a TTY.

## `job status` (`cli/commands/job-status.md`)

- [ ] Full job status incl. `failure_code` (machine-readable, stable) + `failure_message`; the failure view **always states "the target folder is unchanged"** (PRD AP-E2).

## `job watch` (`cli/commands/job-watch.md`)

- [ ] Block until terminal state, streaming progress; **exit 0 success / 1 failed / 9 timeout** (the CI contract); `--timeout`, `--interval`; progress to stderr so `--json` survives a pipe; watching a finished job returns immediately.

## `archive extract` (`cli/commands/archive-extract.md`)

- [ ] Upload + server-side extraction into a folder; `--wait` blocks (the CI flag) via `job watch`; `archive_rejected` (bomb/traversal/symlink) → exit 8 with "no files were written" (TRD R6, AP-E2); not idempotent (fresh job per run).

## `archive bundle` (`cli/commands/archive-bundle.md`)

- [ ] N objects → one archive; empty input refused up front (exit 2); **an unreadable path fails the whole request** (exit 4 — silent omissions in a download are data loss); `--wait` blocks; the result URL is signed + expires (stated).

**tests:** per command — happy path, error, `--json`; `job watch` exit 0/1/9 contract; `tenant usage` unknown-not-$0; `archive extract` rejection leaves nothing written.
