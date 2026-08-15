# Command — `bloberry version`

## Purpose & context

- **User goal**: see the installed version, commit and build date — and whether a newer release exists (the update story, `cli/README.md` §Update; injected via `-ldflags`, never hardcoded).
- **When they reach for it**: filing a bug, checking for updates, CI stamping.
- **Needs**: nothing for the version block itself; the update check hits the Releases endpoint (offline-safe — the check failing never fails the command).
- **Data**: version, commit, build date (build-time `-ldflags`); the install method (brew/scoop/apt/manual) for the update hint.

## Signature

```
bloberry version [flags]
```

| Arg / flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `--check` | bool | false | no | Also check for a newer release and print the update command. |

## Help text

```
Print the installed version, commit and build date.

With --check, also asks the server for the latest release and prints
the right update command for this install method (never overwrites a
package-managed binary — see docs).

Examples:
  bloberry version
  bloberry version --check
  bloberry version --json
```

## Output states

**Version**

```
bloberry v1.0.0
  commit:  8f2a1c9
  built:   2026-03-01T09:00:00Z
```

**`--check`, update available**

```
bloberry v1.0.0
  commit:  8f2a1c9
  built:   2026-03-01T09:00:00Z

Update available: v1.0.2
  Install via Homebrew:
    brew upgrade bloberry
```

**`--check`, up to date**

```
bloberry v1.0.2 (up to date)
```

**`--json`**

```json
{"version":"v1.0.0","commit":"8f2a1c9","built":"2026-03-01T09:00:00Z","update_available":"v1.0.2"}
```

## Exit codes

| Code | When |
|---|---|
| `0` | Version shown (including update-available — that's info, not an error). |
| `1` | Only with `--check` when the release check failed and `--json` isn't masking it — version still printed. |

## Behavior notes

- **Never hardcoded**: `version`, commit and build date come from build-time `-ldflags` (`cli-defaults.md` §Output) — a hardcoded version silently lies after the next release.
- **The update hint is install-method-aware**: `brew upgrade bloberry` vs `scoop update bloberry` vs `apt install --only-upgrade bloberry` vs a download link — because **there is no `self-update`** that overwrites a package-managed binary (`cli/README.md` §Update; overwriting a Homebrew/apt binary breaks the package database).
- **stdout**: the version block / JSON. **stderr**: nothing on success; the `--check` failure message.
- **`--check` is offline-safe**: a network failure prints the version plus "update check failed (offline?)" and still exits 0 — an offline machine must not exit non-zero on a version query.
