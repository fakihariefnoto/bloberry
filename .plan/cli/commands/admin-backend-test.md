# Command — `bloberry admin backend test`

## Purpose & context

- **User goal**: run an immediate reachability check against a backend — the fix loop ("replace credentials → test → healthy") that shouldn't wait for the server's periodic health ticker (PRD M19/PA-E1).
- **When they reach for it**: after registering or editing a backend; when a tenant reports upload failures.
- **Needs**: auth as `platform_admin`; the backend's credentials (stored).
- **Data**: `storage_backends` health check — runs the driver's conformance-level reachability probe, returns health + the raw error on failure.

## Signature

```
bloberry admin backend test <id> [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `<id>` | string | — | yes | Backend ID (`sb_…`). |

## Help text

```
Run an immediate health check against a storage backend.

Contacts the provider now, instead of waiting for the periodic check.
On failure, prints the real provider error (platform admins only).

Examples:
  bloberry admin backend test sb_3d9f
  bloberry admin backend test sb_3d9f --json
```

## Output states

**Success**

```
✓ s3-eu-prod (sb_3d9f) is healthy
  Bucket app-uploads reachable · 0.31s
```

**Failure**

```
✗ gcs-foundry (sb_5f8a) is unreachable

  SigV4/oidc auth failed: service account key expired 2026-07-31.
  Fix: replace the credential (admin backend add, or the dashboard)
  then test again.
```

**`--json`**

```json
{"id":"sb_5f8a","name":"gcs-foundry","healthy":false,"error":"service account key expired 2026-07-31","latency_ms":null}
```

## Exit codes

| Code | When |
|---|---|
| `0` | Healthy. |
| `9` | Unreachable — the backend-health code (scripts branch on this to page someone). |
| `5` | No backend with that ID. |
| `4` | Forbidden — not `platform_admin`. |

## Behavior notes

- **Exit `9` is the health signal**: an ops script runs `admin backend test sb_x || bloberry admin backend test sb_x` style checks and alerts on `9` — the same code the data path returns when a backend fails mid-request (`backend/domains.md` error table).
- **The raw provider error is the value** — PRD PA-E1 exists so the person who can fix it sees it; the fix hint ("replace the credential") is appended so the message isn't just a diagnostic dump.
- **stdout**: the result / JSON. **stderr**: nothing on success.
- **Not a full conformance run** — it's a reachability probe (list/stat on a probe key), not the multi-case suite (`internal/storage/conformance`, `backend/domains.md` §6.4); fast enough for a fix loop.
