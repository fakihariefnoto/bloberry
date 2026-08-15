# Task group — 05 commands: folder + share

**Depends on:** `02-core-infra`, `03-auth`. One checkbox per command; acceptance = "output matches the sample in `cli/commands/<file>.md`".

## `folder create` (`cli/commands/folder-create.md`)

- [ ] Create a folder, creating missing parents; already-exists → exit 8 with "no change made"; `--json` reports `created_intermediates`.

## `folder tree` (`cli/commands/folder-tree.md`)

- [ ] Indented tree from the tenant root or `[path]`; `--objects` adds per-folder object lines with sizes; depth-bounded (20 levels) with a `…` elided marker; `--json` emits the nested shape.

## `share link` (`cli/commands/share-link.md`)

- [ ] Signed link with `--ttl` (default 24h, max from config); the URL is the prominent output; `--open` (TTY only) opens it. **Revocation honesty stated** (the R11 redirect caveat, same as the web UI). Each run creates a new link (not idempotent — documented).

## `share short` (`cli/commands/share-short.md`)

- [ ] Random unguessable slug by default (a short URL is a capability, PRD D6); `--slug` requests a memorable one; slug-taken → exit 8 with "run without `--slug`"; permanent until revoked unless `--ttl`.

## `share public` (`cli/commands/share-public.md`)

- [ ] Flip `objects.visibility` to public; **confirmation required** unless `--yes` (public has no expiry and can't be revoked, only un-published); the stable `/o/<file_id>` URL in the output. **No `share private` command in v1** — the output says un-publishing is a web-dashboard action rather than inventing one.

## `share list` (`cli/commands/share-list.md`)

- [ ] List links, scoped by `[path]`, filtered by `--status`; **hits-first default sort** (the revoke decision); URLs middle-truncated in human view, full in `--json`; empty state names the create command.

## `share revoke` (`cli/commands/share-revoke.md`)

- [ ] Revoke a link by id; the confirmation states the hit count ("12 hits · last used 2h ago · now dead"); **idempotent** (already-revoked = exit 0); immediate effect (explicit cache invalidation); no Undo.

**tests:** per command — happy path, error, `--json`; the share-link TTL boundary; short-slug collision exit 8; revoke confirm copy.
