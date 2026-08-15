# Task group — 04 commands: file verbs

**Depends on:** `02-core-infra`, `03-auth`. **Blocks:** `10-flows` (publish journey). One checkbox per command; acceptance = "output matches the sample in `cli/commands/<file>.md`".

## `cp` (`cli/commands/cp.md`)

- [ ] Copy in all three directions: local→remote, remote→local, remote→remote (server-side, bytes never cross the client). `bloberry://` path scheme.
- [ ] `-r/--recursive`, `--exclude` (repeatable glob), `--dry-run`, `--no-clobber`.
- [ ] **Partial failure → exit 7** with the per-file summary on stderr; successful files stay uploaded; rerun is idempotent (PRD CL-E1). This is the designed core of the command.
- [ ] `--json` emits only the payload (src/dst/copied/bytes/skipped/failed).
- [ ] Progress bar only on a TTY; one line per file when piped.

## `ls` (`cli/commands/ls.md`)

- [ ] List a folder (bare = tenant root); folders get a trailing `/`; `-l` adds size/visibility/modified; `--recursive` walks the subtree.
- [ ] Folders-first, alphabetical, stable ordering; cursor pagination internally, complete listing always printed.
- [ ] Empty folder prints "(empty)" + a next-step hint, exit 0 (an empty folder is a valid answer, not an error).
- [ ] `--json` emits the full array.

## `rm` (`cli/commands/rm.md`)

- [ ] Delete object or folder; `-r` for folders; `--dry-run`; soft-delete default, `--hard` permanent.
- [ ] **Typed-name confirmation for folders** (the web's `confirm-destructive` discipline, PRD TA-E1) with the real object count; `--yes` skips; non-TTY without `--yes` fails with "pass `--yes`".
- [ ] Large folder delete → `subtree_delete` job server-side; prints the job ID and (in scripts) the terminal state via `job watch`.
- [ ] Partial recursive failure → exit 7; already-deleted counts as succeeded (idempotent).

## `mv` (`cli/commands/mv.md`)

- [ ] Move/rename; the `file_id` never changes (PRD M4); `folder_cycle` refused with the actionable message (PRD TA-E2); destination-exists → exit 8, never silent overwrite.
- [ ] Confirmation only for folder moves above the descendant threshold; single objects move without a prompt.

## `cat` (`cli/commands/cat.md`)

- [ ] Stream bytes to stdout untouched (binary-safe, no post-processing); errors/progress to stderr so the pipe stays clean; memory stays flat via streaming (PRD G10 discipline).

## `stat` (`cli/commands/stat.md`)

- [ ] Object/folder metadata: `file_id`, size, type, visibility, backend, hash, uploader, timestamps; `--json` emits the single object (scripts read `.id`).

## `sync` (`cli/commands/sync.md`)

- [ ] One-way local→remote mirror; **compares by size + `content_hash`** (identical bytes skip even with a changed mtime); `--delete` removes remote-only objects; `--dry-run`; `--exclude`.
- [ ] `--delete` refuses the tenant root / top-level folder without explicit confirmation.
- [ ] Partial failure → exit 7, successful uploads persist.
- [ ] **Direction is one-way and stated** (PRD NG4 — publish, not a backup guarantee).

**tests:** per command — happy path, empty result, one error case, `--json`; for `cp`/`rm`/`sync`: the partial-failure case asserting exit 7 + the summary format.
