# Build index — Infra

Build-order table, dependency edges, and status for the infra/CI-CD task list.

## Build order

| File | Covers | Status |
|---|---|---|
| `01-vps-runner.md` | VPS provisioned, self-hosted runner installed as a service, labels, sudoers rule | ☐ |
| `02-process-model.md` | systemd unit, native Mongo/Redis (with AOF + job reclaim), local-disk volume, `.env.example`, runbooks | ☐ |
| `03-reverse-proxy-dns.md` | Caddy (no body buffering, 5GB ceiling), one domain, TLS, DNS | ☐ |
| `04-environments-secrets.md` | staging/prod GitHub Environments, every secret set, `ci.yml`/`deploy.yml` committed | ☐ |
| `05-first-deploy.md` | Manual dry-run of `ci.yml` → `deploy.yml` against staging, confirmed end-to-end | ☐ |
| `06-cli-release.md` | `release-cli.yml` + the PAT for the tap/bucket repos (Linux runner is fine) | ☐ |
| `07-desktop-ci.md` | `release-desktop.yml` — the OS-native exception (macOS/Windows billable), per-OS toolchains, signing secrets | ☐ |
| `08-monitoring-backups.md` | health endpoints, alerts, `mongodump` nightly, restore runbook | ☐ |
| `09-scaling.md` | Cloud-backend multi-node: worker split, scheduler lease, Mongo RS + Redis sentinel, Caddy LB | ☐ |

## Dependency edges

- **`01-vps-runner` blocks everything** — deploy runs on the self-hosted runner; no runner, no deploy.
- **`02-process-model` blocks `03`, `05`, `08`** — the unit must exist before Caddy proxies to it or deploy restarts it.
- **`03-reverse-proxy-dns` blocks `05`** — a first deploy needs the domain routing to reach the box.
- **`04-environments-secrets` blocks `05`** — deploy can't run without secrets set.
- **`06-cli-release` depends on `cli/tasks/13`** (GoReleaser config) — infra owns the runner and the PAT; cli owns the artifacts.
- **`07-desktop-ci` depends on `desktop/tasks/08`** (the packaging files) — infra owns the runners + secrets; desktop owns the artifacts.

## External edges

- **Monorepo, one runner serves everything** — backend, web, CLI-Linux, desktop-Linux all deploy from the same registered runner (`infra/README.md` §Self-hosted runner).
- **`release-cli.yml` runs on GitHub-hosted Linux** — a pure-Go CLI cross-compiles all five platforms from one Linux runner; the OS-native exception does NOT apply to it (don't copy it onto a CLI that doesn't need it).
- **`release-desktop.yml` is the exception** — macOS/Windows are GitHub-hosted and **billable** (the one place this project pays GitHub for minutes); only the Linux job uses the self-hosted runner.
- **CLI channel publishing tasks live in `cli/tasks/14-distribution.md`** — infra owns the runner and its credentials; this file only covers the workflow + token.

## Gaps flagged

None. Process model (systemd), reverse proxy (Caddy, one domain), secrets list per `.env.example`, 4 workflows, the 3 non-default constraints (no-buffering, Redis AOF + job reclaim, in-process tickers), backups and monitoring are all concrete.
