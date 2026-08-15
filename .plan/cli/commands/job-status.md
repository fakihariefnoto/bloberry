# Command — `bloberry job status <id>`

## Purpose & context

- **User goal**: get one job's full status — progress, result, or the real failure reason (PRD AP-E2: the failure reason is retrievable by job ID).
- **When they reach for it**: interactive debugging; polling in scripts (`job status` in a loop, or `job watch` to block).
- **Needs**: auth, tenant context, read.
- **Data**: `jobs` — full record incl. `payload`, `result`, `failure_code`/`failure_message`, `attempts`, timestamps.

## Signature

```
bloberry job status <id> [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `<id>` | string | — | yes | Job ID (`job_…`). |

## Help text

```
Show one job's full status.

Includes progress, the result payload (bundle URLs, extraction target)
and, for failures, both the machine-readable code and the human
message.

Examples:
  bloberry job status job_8f2a
  bloberry job status job_8f2a --json
```

## Output states

**Success (running)**

```
Job:      job_8f2a · extract
State:    running (attempt 1)
Progress: 47 / 184 files  ████████░░░░░░░░░░░░
Started:  2026-03-12 14:00:08 UTC
```

**Success (failed)**

```
Job:      job_5b2e · extract
State:    failed (attempt 2 of 3)
Code:     archive_rejected
Message:  decompressed size exceeded the 8 GB ceiling
Note:     the target folder is unchanged — extraction commits atomically
```

**`--json` (failed)**

```json
{"id":"job_5b2e","kind":"extract","state":"failed","attempts":2,"failure_code":"archive_rejected","failure_message":"decompressed size exceeded the 8 GB ceiling","finished_at":"2026-03-12T14:00:11Z"}
```

## Exit codes

| Code | When |
|---|---|
| `0` | Shown. |
| `5` | No job with that ID. |
| `4` | Forbidden. |
| `2` | Missing ID. |

## Behavior notes

- **The failure view always states the atomic-commit guarantee** — "the target folder is unchanged" (PRD AP-E2) is the reassurance that turns a scary failure into a safe one; it's part of every failed-job render.
- **`failure_code` is machine-readable and stable** (`archive_rejected`, `backend_unreachable`, …) — scripts branch on it, matching the API's error-code contract (`backend/domains.md` §8).
- **stdout**: the report / JSON. **stderr**: nothing on success.
- **Not idempotent-sensitive**: repeated calls are safe (jobs are immutable after finishing).
