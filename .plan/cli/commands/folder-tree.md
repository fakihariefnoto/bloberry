# Command — `bloberry folder tree`

## Purpose & context

- **User goal**: print the tenant's folder hierarchy (or a subtree) as an indented tree — the CLI's answer to the web's breadcrumb browser, for orienting in a script or a terminal.
- **When they reach for it**: interactive exploration; answering "what's actually under assets/" before a recursive operation.
- **Needs**: auth, tenant context, read on the traversed folders.
- **Data**: `folders` (name, id, path), optionally object counts per folder.

## Signature

```
bloberry folder tree [path] [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `[path]` | path | tenant root | no | Subtree root to print. |
| `--objects` | bool | false | no | Also list objects inside each folder (a full recursive listing). |

## Help text

```
Print the folder hierarchy as an indented tree.

Starts at the tenant root by default, or at [path]. With --objects,
objects are shown inside their folders (equivalent to a recursive
ls with structure).

Examples:
  bloberry folder tree
  bloberry folder tree bloberry://assets --objects
```

## Output states

**Success (tree)**

```
assets/
├── v2/
│   ├── 2026/
│   └── archive/
├── public/
└── scripts/
projects/
└── web/
    └── 2026/
```

**With `--objects`**

```
assets/
├── v2/
│   ├── 2026/
│   │   ├── hero.png       2.4 MB
│   │   └── index.html    18.2 KB
│   └── archive/
└── public/
    └── share-card.png      812 KB
```

**`--json`**

```json
{"root":"assets/","folders":[{"id":"fd_4b8a2c1e","name":"v2","children":[{"id":"fd_9c3f1d2b","name":"2026","objects":2}]}]}
```

## Exit codes

| Code | When |
|---|---|
| `0` | Printed. |
| `5` | The subtree root doesn't exist. |
| `4` | Forbidden to read the root. |
| `2` | Bad invocation. |

## Behavior notes

- **Tree characters** are the standard box-drawing set; the terminal table alignment is real (`check-wireframes.mjs` covers `cli/` too — an emoji or wide glyph breaks the columns the same way it breaks a wireframe).
- **stdout**: the tree. **stderr**: nothing on success.
- **Depth**: bounded (default 20 levels; deeper trees truncate with a `…` marker and the count of elided folders).
- **`--objects`** shows size per object; without it, folders only — a tree of 40,000 objects stays readable.
- **No color** in the tree (indentation carries the structure); `--json` emits the nested shape for scripts that need the hierarchy programmatically.
