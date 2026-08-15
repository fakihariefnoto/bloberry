# Command — `bloberry rm`

## Purpose & context

- **User goal**: delete objects or folders. Destructive — the confirmation and `--yes` path matter (PRD TA-E1's "state the real count" applies here as much as the web).
- **When they reach for it**: interactive cleanup (needs confirmation), and scripted cleanup in CI (needs `--yes` and code 7 on partial).
- **Needs**: auth, tenant context. Delete permissions on the target (`delete` via role or grant).
- **Data**: `objects` (`state`, `deleted_at`), `folders` (descendant count for the confirmation).

## Signature

```
bloberry rm <path> [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `<path>` | path | — | yes | `bloberry://` path to delete. |
| `-r, --recursive` | bool | false | no | Delete a folder and everything in it. |
| `--dry-run` | bool | false | no | Show what would be deleted, delete nothing. |
| `--yes, -y` | bool | false | no | Skip the confirmation prompt. |

## Help text

```
Delete an object or folder from Bloberry.

Deleting a folder requires -r and a confirmation that states the real
object count. Deletes are soft by default (restorable within the
retention window); pass --hard to delete permanently.

Flags:
  -r, --recursive  Delete folders recursively
      --dry-run    Show what would be deleted, delete nothing
      --hard       Bypass the trash; permanent
  -y, --yes        Skip the confirmation

Examples:
  bloberry rm bloberry://assets/hero.png
  bloberry rm bloberry://assets/old -r
  bloberry rm bloberry://assets/old -r --yes --json
```

## Output states

**Success (single object)**

```
✓ Deleted hero.png (soft delete — restorable for 30 days)
```

**Recursive, confirming (interactive, TTY)**

```
Delete bloberry://assets/old?
  14 objects · 3 subfolders · 212 MB
Type the folder name to confirm: old
✓ Deleted assets/old (14 objects, 212 MB)
```

**`--dry-run`**

```
would delete 14 objects · 3 subfolders · 212 MB
```

**`--json` (success)**

```json
{"deleted":14,"bytes":212000000,"permanent":false}
```

## Exit codes

| Code | When |
|---|---|
| `0` | Deleted. |
| `5` | The path doesn't exist. |
| `4` | Forbidden — no `delete` on the target. |
| `7` | Recursive delete partially failed (summary on stderr). |
| `2` | Folder without `-r`, or missing name confirmation in a TTY. |

## Behavior notes

- **Confirmation**: single object → plain confirm (overridable `--yes`). Folder → **typed-name confirmation** (type the folder's name), matching the web's `confirm-destructive` discipline (PRD TA-E1). In a non-TTY without `--yes`, fail with "pass `--yes`" rather than hang on an invisible prompt.
- **Soft delete** is the default (S5); `--hard` is permanent and irreversible — the confirmation wording changes ("permanent — not restorable").
- **Large folder delete** (above the threshold) becomes a `subtree_delete` job server-side; `rm` prints the job ID and, in scripts, the terminal state (`job watch` semantics) so CI waits for the real result rather than a 202.
- **Partial failure** (exit 7): per-file summary on stderr; already-deleted items count as succeeded (idempotent rerun).
- **stdout**: the summary. **stderr**: progress, confirmations, errors.
