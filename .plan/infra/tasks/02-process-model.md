# Task group — 02 process model (systemd + native Mongo/Redis)

**Depends on:** `01-vps-runner.md`. **Blocks:** `03`, `05`, `08`.

Per `infra/README.md`: **one** unit — the `bloberry` binary (API + embedded dashboard + `/s/` + public objects + the in-process worker) with MongoDB and Redis installed natively alongside. systemd is what `deploy.yml` targets. (A Docker Compose file is still published for self-hosters, but it is not the reference deployment.)

- [ ] **`bloberry.service` unit written** per the README's exact unit: `Type=notify`, `User/Group=bloberry`, `WorkingDirectory=/opt/bloberry`, `EnvironmentFile=/etc/bloberry/bloberry.env`, `Restart=on-failure`, **`Requires=mongod.service redis-server.service`** (failing at boot beats serving 500s), **`MemoryMax=2G`** (extraction runs in-process; an unbounded worker takes MongoDB down with it — TRD R12), `LimitNOFILE=65536`, and the hardening block (`NoNewPrivileges`, `ProtectSystem=strict`, `ProtectHome`, `ReadWritePaths=/var/lib/bloberry`, `RestrictSUIDSGID`, etc. — the process holds credentials to every connected bucket).
- [ ] **Service installed + enabled** — `/etc/systemd/system/bloberry.service`, `systemctl enable`, started.
- [ ] **Log rotation configured** — journald with a retention bound.
- [ ] **MongoDB installed natively**, bound to localhost, running as its own service. The migrations from `backend/tasks/02` run against it.
- [ ] **Redis installed natively**, bound to `127.0.0.1` only (a public Redis holding session tokens is full account takeover — README's explicit warning).
- [ ] **Redis `appendonly yes` with `appendfsync everysec`** — the AOF decision (README §2): the job queue forces it. A user whose 2 GB extraction vanished on a restart sees it stuck at `queued` forever with no way to know it's dead; AOF bounds the loss to a second of enqueues.
- [ ] **Job reclaim** — a startup sweep moves stale `running` jobs (no heartbeat past a threshold) back to `queued`, bounded by `attempts` so a poison job doesn't loop forever.
- [ ] **The four in-process tickers** (not systemd timers — they share the server's config/db/driver registry, keeping the one-binary story): reconciliation sweep (15 min), usage metering (hourly), backend health (5 min), audit retention (daily). **Single-instance assumption recorded** (running two processes double-runs every ticker — fine for v1, the first thing to break under scaling).
- [ ] **Local-disk objects volume** — `/var/lib/bloberry/objects` owned `bloberry:bloberry` mode `0700`, **mounted as its own volume** (it can grow without bound; filling root takes Mongo and Redis down). Alert at 80% (see `08`).
- [ ] **`/etc/bloberry/bloberry.env`** — mode `0600`, owner `root`; holds `CREDENTIAL_ENCRYPTION_KEY` among the secrets (README's special handling — never in the repo, never in MongoDB, excluded from backups by construction).
- [ ] **`deploy/` in the repo** — the systemd unit, `Caddyfile`, and `.env.example` per stage committed (the real values are not).
- [ ] **`--rotate-credential-key` runbook** — the offline rotation command (both keys supplied, service stopped, rehearsed on staging before production; README §`CREDENTIAL_ENCRYPTION_KEY`).

**verification:** after `systemctl restart bloberry`, `GET /readyz` passes (Mongo + Redis reachable); a killed worker's running job is reclaimed on restart; the disk volume is mounted separately.
