# Command — `bloberry archive extract`

## Purpose & context

- **User goal**: upload an archive and have the server extract it into a folder (PRD M11/AP4) — one request instead of N uploads for a bulk import.
- **When they reach for it**: interactive bulk import; CI ingestion of a release tarball.
- **Needs**: auth, tenant context, `write` on the target folder.
- **Data**: `jobs` (kind `extract`) + `folders`. The job is **queued, never inline** (PRD D3); extraction commits atomically (AP-E2). Safety ceilings are server-side (zip bomb, path traversal, symlinks — TRD R6).

## Signature

```
bloberry archive extract <archive> <target> [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `<archive>` | local path | — | yes | `.zip`/`.tar.gz` to upload. |
| `<target>` | `bloberry://` path | — | yes | Folder to extract into. |
| `--wait` | bool | false | no | Block on the job (`job watch`) until it finishes. |

## Help text

```
Upload an archive and extract it server-side into a folder.

The upload goes through the normal write path, then an extraction job
unpacks it. The target folder is unchanged if the job fails (atomic
commit). Pass --wait to block until extraction finishes — use it in
CI.

Examples:
  bloberry archive extract ./bundle.zip bloberry://imports
  bloberry archive extract ./site.tar.gz bloberry://site --wait
```

## Output states

**Success (queued, no `--wait`)**

```
✓ Uploaded bundle.zip → extraction queued
  Job: job_8f2a  (watch: bloberry job watch job_8f2a)
```

**Success (`--wait`)**

```
✓ Extraction complete into bloberry://imports/ (184 files)
  Job: job_8f2a · 2.1s
```

**Rejected (exit 8)**

```
Error: archive rejected (archive_rejected)

The archive contains entries that can't be safely extracted (paths
outside the target, or a decompressed size over the 8 GB ceiling).
No files were written.
```

**`--json` (`--wait`, success)**

```json
{"job_id":"job_8f2a","state":"succeeded","target":"imports","files":184,"bytes":212000000}
```

## Exit codes

| Code | When |
|---|---|
| `0` | Queued, or (with `--wait`) completed. |
| `8` | Archive rejected server-side — nothing written. |
| `6` | Quota blocked the upload. |
| `7` | — (extraction is one job, not a batch; partial state is "failed, unchanged"). |
| `2` | Bad invocation (missing archive/target, unsupported archive type). |

## Behavior notes

- **stdout**: the result. **stderr**: the upload progress, job progress.
- **`--wait` is the CI flag**: without it, the exit is "queued" (0) and the caller watches via `job watch`; with it, exit reflects the terminal state. The help shows the pairing.
- **Rejection reasons** map to the `archive_rejected` code with the specific cause (bomb / traversal / symlink / ratio) — the CLI says what was rejected and that nothing was written, matching AP-E2's atomicity promise.
- **Not idempotent**: each run uploads and queues a fresh extraction; the target folder's existing contents determine collision behavior (server-side name handling).
