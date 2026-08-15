# Command — `bloberry ls`

## Purpose & context

- **User goal**: list the contents of a folder (objects + subfolders) — the most-run command, the mental model of the tool.
- **When they reach for it**: interactively to see what's where, and in scripts to check whether something exists before acting on it.
- **Needs**: auth, tenant context. Path is `bloberry://` or bare (bare resolves against the current tenant root).
- **Data**: `folders`, `objects` — name, `size_bytes`, `visibility`, `modified_at`, type (folder vs object).

## Signature

```
bloberry ls [path] [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `[path]` | path | tenant root | no | `bloberry://` path; bare = current tenant root. |
| `-l` | bool | false | no | Long format — size, visibility, modified time. |
| `--recursive` | bool | false | no | Recurse into subfolders. |

## Help text

```
List the contents of a Bloberry folder.

With no path, lists the tenant root. Directories are marked with a
trailing '/'. Use -l for size, visibility and timestamps; --recursive
walks the whole subtree.

Flags:
  -l              Long format: size, visibility, modified
      --recursive Recurse into subfolders

Examples:
  bloberry ls
  bloberry ls bloberry://assets/v2 -l
  bloberry ls bloberry://assets -r --json
```

## Output states

**Success (bare)**

```
assets/  hero.png  index.html  main.ts  public/
```

**Success (`-l`)**

```
drwx  assets/                    Mar 12 09:41
-rw-  2.4 MB  hero.png  public  Mar 12 14:03
-rw-  18.2 KB  index.html        Mar 12 09:42
-rw-  9.8 KB  main.ts            Mar 11 17:20
drwx  public/  public            Mar 12 09:41
```

**Empty folder**

```
(empty)  — nothing in bloberry://assets/v2 yet
```

**`--json` (success)**

```json
[{"name":"assets/","type":"folder","modified_at":"2026-03-12T09:41:00Z"},
 {"name":"hero.png","type":"object","size":2400000,"visibility":"public","modified_at":"2026-03-12T14:03:00Z"}]
```

## Exit codes

| Code | When |
|---|---|
| `0` | Listed successfully (including an empty folder). |
| `5` | The path doesn't exist. |
| `3` | Not authenticated. |
| `4` | Forbidden to read that folder. |

## Behavior notes

- **stdout** is the listing (data); nothing else pollutes it.
- **Directory marker**: trailing `/` in bare mode; `d` type column in `-l` and `type:"folder"` in `--json` — one decision, three renderings.
- **Ordering**: folders first, then objects, each alphabetically; stable across runs (a script diffing two listings shouldn't see noise).
- **Empty state** says what's empty and implies the next step; it does not error (an empty folder is a valid answer).
- **TTY**: no color in bare mode (the `/` marker carries the distinction); `-l` may tint the visibility field only when TTY.
- **Pagination**: large folders use cursor pagination internally; `ls` always prints the complete listing (it's the tool's model of "what exists"), and `--json` emits the full array.
