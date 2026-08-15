# Command — `bloberry mv`

## Purpose & context

- **User goal**: move or rename an object/folder. The `file_id` survives the move (PRD M4/G4); a folder move rewrites descendants' ancestry server-side with zero storage-backend copies (ADR-7).
- **When they reach for it**: interactive reorganization; occasionally scripted (release-renames).
- **Needs**: auth, tenant context. Move permission on source (`delete` + `write` on the destination per the resolver).
- **Data**: `folders`, `objects` — the path rewrite (ancestors, path), and whether the target exists (name-conflict semantics).

## Signature

```
bloberry mv <src> <dst> [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `<src>` | path | — | yes | `bloberry://` source. |
| `<dst>` | path | — | yes | `bloberry://` destination (new name or target folder). |
| `--yes, -y` | bool | false | no | Skip the confirmation. |

## Help text

```
Move or rename an object or folder.

The file_id never changes, so references and share links survive the
move. Moving a folder into its own descendant is refused. When <dst>
is an existing folder, <src> is moved into it.

Flags:
  -y, --yes  Skip the confirmation

Examples:
  bloberry mv bloberry://assets/hero.png bloberry://assets/brand/hero.png
  bloberry mv bloberry://assets bloberry://archive/assets-2025
```

## Output states

**Success (rename)**

```
✓ assets/hero.png → assets/brand/hero.png (f_8Kd2pQxL31A unchanged)
```

**Refused — folder cycle**

```
Error: cannot move "assets/old" into its own descendant "assets/old/nested"

This would create a loop that can't be undone. (folder_cycle)
```

**`--json` (success)**

```json
{"id":"f_8Kd2pQxL31A","from":"assets/hero.png","to":"assets/brand/hero.png"}
```

## Exit codes

| Code | When |
|---|---|
| `0` | Moved/renamed. |
| `5` | Source or destination doesn't exist. |
| `8` | Conflict — destination exists, or `folder_cycle`. |
| `4` | Forbidden — no move permission. |
| `2` | Bad invocation. |

## Behavior notes

- **Cycle prevention is server-side** (PRD TA-E2): the CLI surfaces the `folder_cycle` code as the actionable message above, not a raw error.
- **Name conflict**: destination exists → exit `8` with "destination exists (run `bloberry ls <parent>`)". `mv` never silently overwrites; no `--force` in v1 (rename-then-merge is the safe path).
- **Confirmation**: only for folder moves with more than a threshold of descendants (server reports the count); `--yes` skips. Single objects move without a prompt.
- **stdout**: the result line. **stderr**: confirmation prompt, errors.
- **Idempotent**: re-running a completed move errors with "not found" (the source no longer exists) — not a silent no-op and not a fake success.
