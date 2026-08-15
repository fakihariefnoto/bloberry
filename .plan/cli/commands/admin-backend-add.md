# Command — `bloberry admin backend add`

## Purpose & context

- **User goal**: register a storage backend with its credentials (PRD M1/M20/PA1) — the CLI's way to provision infrastructure, for platform admins who prefer a terminal or need to script install bootstrap.
- **When they reach for it**: install setup; adding a provider without touching the dashboard.
- **Needs**: auth as `platform_admin`; the server reachable. Credentials are **envelope-encrypted at rest** and never echoed (PRD M20/R7).
- **Data**: `storage_backends` — driver, name, config, `credentials_encrypted`, rate_card, `tenant_id` null = install-level pool (PRD D4).

## Signature

```
bloberry admin backend add <driver> <name> [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `<driver>` | string | — | yes | `s3`, `r2`, `oss`, `gcs`, `azblob`, `disk`. |
| `<name>` | string | — | yes | Unique install-level name (`s3-eu-prod`). |
| `--config <key=value>` | string (repeatable) | — | no | Driver config (endpoint, bucket, prefix, region…). |
| `--credential <key=value>` | string (repeatable) | — | varies | Credentials (secret, key_id, service-account file…). |
| `--credential-file <key=path>` | string (repeatable) | — | varies | Read a credential value from a file (GCS service account JSON). |
| `--rate-card <s>$/<g>$/<r>$` | string | — | no | `storage$/GB-mo,egress$/GB,requests$per-1k`. |

## Help text

```
Register a storage backend (platform admin).

Credentials are encrypted at rest with the install's envelope key and
are write-only — they can be replaced but never read back. For GCS,
pass --credential-file service_account=/path/to/key.json; for Azure
(azblob), pass a SharedKey or SAS credential. Credentials that would
be visible in the process list are read from files instead of flags.

Examples:
  bloberry admin backend add s3 s3-eu-prod \
    --config bucket=app-uploads --config region=eu-west-1 \
    --credential-file access_key=~/.secrets/s3-key.json
  bloberry admin backend add azblob az-eu-prod \
    --config container=bloberry-prod \
    --credential-file shared_key=~/.secrets/azblob-key.json
  bloberry admin backend add disk vps-volume \
    --config root=/data/blob
```

## Output states

**Success**

```
✓ Registered storage backend s3-eu-prod (sb_3d9f)
  Driver:  s3 · bucket app-uploads · region eu-west-1
  Health:  unchecked (test with: bloberry admin backend test sb_3d9f)
```

**`--json`**

```json
{"id":"sb_3d9f","driver":"s3","name":"s3-eu-prod","health":"unchecked","credentials_set":true,"rate_card":null}
```

## Exit codes

| Code | When |
|---|---|
| `0` | Registered. |
| `8` | Name already in use (names unique per install, `ERD.md`). |
| `2` | Bad invocation (missing required config/credentials for the driver). |
| `4` | Forbidden — not `platform_admin`. |
| `9` | The driver conformance check failed on first contact (bad credentials). |

## Behavior notes

- **Credentials never on stdout**: the human output confirms `credentials_set: true`, never the values; `--json` likewise omits them (the API never returns them, `ERD.md`).
- **The process-list rule**: secrets passed as `--config key=secret` risk `/proc` exposure — the help pushes `--credential-file` for anything sensitive (GCS service-account JSON is the canonical case).
- **stdout**: the result. **stderr**: nothing on success.
- **Health starts `unchecked`** — registration doesn't contact the provider; `admin backend test` is the explicit next step, named in the output.
- **Rate card** optional at add time; settable later via `admin backend rate-card`.
