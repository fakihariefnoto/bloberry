# Command — `bloberry job list`

## Purpose & context

- **User goal**: see the tenant's queued/running/finished jobs — extractions, bundles, subtree deletes.
- **When they reach for it**: checking whether a queued operation finished; CI watching a batch of jobs.
- **Needs**: auth, tenant context, read.
- **Data**: `jobs` — kind, state, progress, timestamps, failure fields (`ERD.md` jobs).

## Signature

```
bloberry job list [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `--state <s>` | string | all | no | Filter: `queued`, `running`, `succeeded`, `failed`. |

## Help text

```
List jobs for the current tenant.

Shows progress for running jobs and the failure reason for failed
ones. Poll with 'bloberry job status <id>' or block with
'bloberry job watch <id>'.

Examples:
  bloberry job list
  bloberry job list --state failed --json
```

## Output states

**Success**

```
ID          KIND            STATE       PROGRESS      STARTED
job_8f2a    extract         running     47/184 files  now
job_3d1c    bundle          succeeded   done          03:12
job_9c4f    subtree_delete  succeeded   done          03:05
job_5b2e    extract         failed      0/0           ⚠ zip bomb: size ceiling
```

**No jobs**

```
No jobs yet. Extraction and large deletes will appear here.
```

**`--json`**

```json
[{"id":"job_8f2a","kind":"extract","state":"running","progress_done":47,"progress_total":184,"started_at":"2026-03-12T14:00:00Z"}]
```

## Exit codes

| Code | When |
|---|---|
| `0` | Listed (including empty). |
| `4` | Forbidden. |

## Behavior notes

- **stdout**: the table / JSON. **stderr**: nothing on success.
- **The failure column carries the reason** (short, human-readable) — `job list` is the first place someone learns a batch failed; it shouldn't force opening `job status`.
- **Progress is honest**: `47/184 files`, never a fake percentage — the same real-progress rule as the web (PRD M21, `design/style-guide.md` Motion).
- **TTY**: `running` rows may show a live-updating progress line only when stderr is a terminal; piped output prints once.
