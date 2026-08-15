# Command — `bloberry job watch <id>`

## Purpose & context

- **User goal**: block until a job reaches a terminal state (succeeded/failed), streaming progress — the CI-friendly command that turns a 202 into a real result.
- **When they reach for it**: scripts that queue an extraction/bundle then must act on the result; interactive watching of a long delete.
- **Needs**: auth, tenant context, read.
- **Data**: `jobs` — polled until `state` is `succeeded`/`failed`; `--timeout` bounds the wait.

## Signature

```
bloberry job watch <id> [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `<id>` | string | — | yes | Job ID (`job_…`). |
| `--timeout <duration>` | duration | none | no | Fail if not terminal within this long. |
| `--interval <duration>` | duration | `2s` | no | Poll interval. |

## Help text

```
Wait for a job to finish, printing progress as it goes.

Exits 0 when the job succeeds, 1 when it fails (with the reason),
and the documented code on timeout. This is what CI uses after
queueing an extraction:

  bloberry archive extract … 
  bloberry job watch "$(…)" --timeout 10m

Examples:
  bloberry job watch job_8f2a
  bloberry job watch job_8f2a --timeout 5m --json
```

## Output states

**Success (interactive)**

```
job_8f2a  extract  47/184 ████████░░░░░░░░░░░░  (running)
job_8f2a  extract  184/184 ████████████████████  ✓ succeeded
  Result: extracted into bloberry://projects/2026/ (184 files)
```

**Failure**

```
job_5b2e  extract  ✗ failed
  archive_rejected: decompressed size exceeded the 8 GB ceiling
  Target folder unchanged.
```

**Timeout (exit 9)**

```
Timed out after 5m — job_8f2a still running (state: running).

Re-run 'bloberry job watch job_8f2a' to keep watching; the job keeps
running server-side.
```

**`--json` (success)**

```json
{"id":"job_8f2a","state":"succeeded","progress_done":184,"progress_total":184,"result":{"target":"projects/2026","files":184}}
```

## Exit codes

| Code | When |
|---|---|
| `0` | Job succeeded. |
| `1` | Job failed — reason on stderr. |
| `9` | `--timeout` elapsed while still running (storage-ish infra concern; the job continues server-side). |
| `5` | No job with that ID. |
| `2` | Missing ID / bad timeout. |

## Behavior notes

- **CI contract**: exit `0` only on terminal success, `1` on failure, `9` on timeout — a script that queues-and-watches branches on exactly these three and gets a real result instead of a 202 it must hand-poll.
- **stdout**: the final result line / JSON. **stderr**: the progress lines (they'd corrupt a piped `--json`).
- **TTY**: progress redraws in place only when stderr is a terminal; piped output prints one line per poll.
- **Idempotent-safe**: watching a finished job returns its terminal state immediately (no wait).
- **The result payload** (extraction target, bundle URL) is printed on success — watch is the command that hands the next step its input.
