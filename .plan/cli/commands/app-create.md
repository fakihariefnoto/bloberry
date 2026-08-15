# Command — `bloberry app create`

## Purpose & context

- **User goal**: register an application (a non-human principal) so it can hold access keys (PRD M10/TA4). Usually the first step of a CI setup: `app create` → `key create`.
- **When they reach for it**: provisioning a new pipeline, scripted setup.
- **Needs**: auth, `tenant_admin`+ role.
- **Data**: `applications` — name, description, `created_at`.

## Signature

```
bloberry app create <name> [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `<name>` | string | — | yes | Application name. |
| `--description <text>` | string | — | no | What this app is for. |

## Help text

```
Register an application (a machine principal) in the current tenant.

An application is the identity that holds access keys. Creating it
here is the first half of provisioning CI; the second is issuing a
key:
  bloberry app create my-cms
  bloberry key create --app <id>

Examples:
  bloberry app create acme-cms
  bloberry app create ci-deploy --description "deploys the site"
```

## Output states

**Success**

```
✓ Created application acme-cms (app_3d9f)
  Next: issue it a scoped key
    bloberry key create --app app_3d9f
```

**Name conflict (exit 8)**

```
Error: an application named "acme-cms" already exists (app_3d9f)

No change made. (name_conflict)
```

**`--json`**

```json
{"id":"app_3d9f","name":"acme-cms","description":"","created_at":"2026-03-12T09:31:00Z"}
```

## Exit codes

| Code | When |
|---|---|
| `0` | Created. |
| `8` | Name already exists — no change. |
| `4` | Forbidden — not `tenant_admin`+. |
| `2` | Missing name. |

## Behavior notes

- **The output chains to the next step** — "Next: issue it a scoped key" with the exact command. Provisioning is a two-step flow; the CLI makes step two discoverable.
- **stdout**: the result. **stderr**: nothing on success.
- **No confirmation** (creating an app is harmless and reversible by `app delete`).
- **Idempotency**: a name-conflict exits `8` with the existing app's ID — a script can catch `8`, read the ID from the message, and continue rather than treating it as a hard failure.
