# Infra — Bloberry

Hosting: **self-hosted VPS**. Process model: **systemd**. Reverse proxy: **Caddy**. CI/CD: **GitHub Actions, everything manually triggered** via `workflow_dispatch` — including CI checks, not just deploy — with deploy running on a self-hosted runner installed on the VPS (no SSH secrets). Those conventions are fixed; see [`templates/infra-defaults.md`](../../../templates/infra-defaults.md).

---

## What actually runs

**One unit**: the `bloberry` binary (API + embedded dashboard + `/s/` short URLs + public objects + the in-process job worker), with MongoDB and Redis installed natively alongside. That's the whole reference deployment, and it's also the story handed to self-hosters (PRD G8: a fresh VPS to running in under 15 minutes), so it has to stay genuinely simple.

```
Internet ──443──▶ Caddy ──127.0.0.1:8080──▶ bloberry.service (systemd)
                                                  │
                                            ┌─────┴─────┐
                                       MongoDB       Redis
                                    (localhost)   (localhost)
                                                  │
                                          /var/lib/bloberry/objects
                                            (local disk driver)
```

A Docker Compose file is still published for self-hosters who prefer it, but **systemd is what `deploy.yml` targets** and what the runbooks describe.

---

## Scaling out — the cloud-backend model (no local-disk driver)

The single-node shape above is the reference deployment and stays the 15-minute story (PRD G8). When an install outgrows it, the documented path is **multi-node with every tenant on a cloud backend** (S3/R2/OSS/GCS/Azure Blob) — **the local-disk driver is single-node by design** and not part of the scaled shape (its HMAC presigner proxies to whichever node holds the bytes; there is no shared-filesystem story that isn't a new SPOF).

The API itself is **stateless** — sessions, authz cache and the job queue live in Redis, metadata in Mongo, bytes in providers — so scale-out is "more API nodes + HA for the two stores + the worker/ticker split":

```
Internet ──443── Caddy (load balancer, health-checked)
                     │
        ┌────────────┼──────────────┬───────────────┐
        │            │              │               │
   api-node-1    api-node-2    worker-1        scheduler
 (bloberry-server)  (server)  (bloberry-worker) (server + leader lease
   role=api, no worker  role=api   jobs only         owns the 4 tickers)
        │            │              │               │
        └────────────┴──────┬───────┴───────────────┘
                            │
                    ┌───────┴───────┐
              Mongo replica set   Redis + sentinel
               (3 nodes)          (AOF, replicated)
                            │
              storage backends (S3/R2/OSS/GCS/Azure Blob)
```

Three things make this not-free, each already an escape hatch in the plan:

1. **The job worker splits out** (`TRD.md` R12): a `cmd/bloberry-worker` binary consuming the same `job:queue`, so API nodes stay lean and worker concurrency scales independently. A deployment change, not a rewrite — the queue boundary already exists.
2. **The four in-process tickers need a leader** (reconciliation sweep, usage metering, backend health, audit retention): a Redis lease (`SET key NX EX` heartbeat) so exactly **one** node runs them; the rest skip. Same binary, a `role=scheduler` config.
3. **Mongo → replica set, Redis → sentinel**: writes to the Mongo primary (automatic failover); Redis gains replication + sentinel because a lost Redis in multi-node = a lost job queue *and* everyone logged out. The AOF decision (`appendonly yes`, `everysec`) carries over unchanged — it's what keeps the queue durable.

Each node runs the **same binary** with a `ROLE=api|worker|scheduler` env var — no separate code paths, no third project. `deploy.yml` targets N nodes; Caddy's `reverse_proxy` gains an upstream list with health checks instead of one `localhost:8080`. Single-instance assumptions are confined to the scheduler (leader-elected), never to API or worker nodes. This is the model to reach for when a tenant's cloud-backend traffic outgrows one box — the disk driver simply isn't offered on scaled installs.

### 3.1 The honest bottleneck: it's the stores, not the API

Before building the API tier, know what actually saturates. Bloberry is designed so API call volume is smaller than it looks:

- **Uploads/downloads mostly bypass the API entirely** — presigned URLs go straight to the storage provider (`architecture.md` §3.1). A 5 GB transfer is ~0 bytes through an API node.
- API load = **metadata operations** (list, stat, authz, presign). Each is: Redis read (authz principal cache) → satisfied; Mongo only on cache miss.

So the scaling ceiling is the **stores**, not the node count:

| Resource | Ceiling | How it scales |
|---|---|---|
| **MongoDB** | Single primary writer | Reads via secondaries + read preference; **writes are the real ceiling** — sharding is the later rung |
| **Redis** | Single-instance by design (sessions, authz cache, job queue) | Read-mostly is fine; the queue's list ops and the authz cache are inherently single-writer |

Adding API nodes scales **concurrency** (more parallel requests) but not the stores. If Mongo's primary or Redis is saturated, more nodes just queue harder. Diagnose which one is hot before adding nodes — `GET /readyz` + the metering snapshots show it.

### 3.2 Adding an API node — the "how does the new IP attach" mechanics

**1. Provision.** Same binary, `ROLE=api`, the same `.env` pointing at the shared Mongo/Redis. **No data to move, nothing to rebalance** — a Bloberry node is fully stateless (state lives in the stores/providers), so it serves any request from the instant it boots. This is the fundamental difference from MinIO, where adding a node means `mc admin decommission` + shard rebalance.

**2. Attach it to the load balancer — two tiers:**

- **Manual (the reference, up to ~5 API nodes):** add `192.0.2.10:8080` to the Caddyfile's `reverse_proxy` upstream list, `systemctl reload caddy`. The new IP attaches by editing one list. Boring, reliable, fits the 15-minute self-host story.

```caddyfile
bloberry.example.com {
	reverse_proxy {
		# Health-checked static pool — add a line per API node.
		to 192.0.2.11:8080 192.0.2.12:8080 192.0.2.13:8080
		health_uri /healthz
		health_interval 10s
		health_timeout 2s
		fail_duration 30s
		max_fails 3
		# Body streaming — unchanged from the single-node config (TRD R5).
		flush_interval -1
	}
	request_body {
		max_size 5GB
	}
}
```

- **DNS-based discovery (the scalable answer, still Caddy):** give API nodes names (`api-1.example.internal`, `api-2…`) and resolve upstreams from **DNS** instead of a static list:

```caddyfile
bloberry.example.com {
	reverse_proxy {
		upstreams dns api.example.internal {
			# A/SRV records under api.example.internal = the API node pool.
			# Caddy re-resolves on its refresh interval — a new node's
			# A record is picked up with no Caddyfile edit, no reload.
		}
		health_uri /healthz
		health_interval 10s
		fail_duration 30s
		max_fails 3
		flush_interval -1
	}
	request_body {
		max_size 5GB
	}
}
```

Now **attaching a new node = adding one A record** (e.g. `api-3 A 192.0.2.14`). Caddy picks it up on the next upstream refresh; the health check drops a dead node out of the pool automatically and re-admits it on recovery. New IPs attach by DNS, which is the standard answer to "how does a new host get into the cluster."

**3. Health checks do the eviction.** `health_uri /healthz` + a fail threshold means attach and detach both happen without human action at the edge. A node under maintenance is removed from DNS or the pool; Caddy notices within the interval; no traffic is sent to a dead host.

**4. `deploy.yml` targets the inventory** — a node list (IPs or DNS names) drives the same `pull → generate → build → install → restart → health-wait` sequence per node; the scheduler node is the only one that also runs migrations.

### 3.3 Manual scale-out vs. autoscaling — the future path

Two distinct futures, one of which abandons the single-binary story:

| Rung | What it is | Fits G8's "one binary, 15 min"? |
|---|---|---|
| **Manual scale-out** (above) | Add a node when load demands; DNS + health checks handle discovery | ✅ Yes — same binary, same systemd, just more of them |
| **Orchestrator (Docker/Kubernetes)** | Containers, a scheduler (K8s replicas / Swarm services), a service registry | ❌ No — a full platform replaces the "one binary, one Caddy site" story |

**The orchestrator is the documented future, not the default.** The plan for *when* it happens (recorded here so it's a decision, not an accretion):

- **Docker stage (Swarm-lite or compose-scale):** containerize the same binary (`Dockerfile`), keep the stores as external stateful services (Mongo RS + Redis sentinel stay *outside* the orchestrator — never containerize the stateful stores with ephemeral volumes), and let a Swarm service / compose scale out `ROLE=api` replicas. The reverse proxy stays Caddy (or the orchestrator's ingress). This buys scheduling and rolling restarts without abandoning the binary.
- **Kubernetes stage:** the full pivot. `ROLE=api` as a Deployment with HPA (autoscaling on CPU/request latency), the worker as its own Deployment (queue consumers scale independently), the scheduler as a single-replica Deployment (the Redis lease still guards the tickers), Mongo as an external RS or an operator-managed statefulset, Redis external (sentinel) or via an operator. Ingress replaces Caddy's edge role. **This is where the "one binary" promise formally ends** — recorded so the pivot is recognized as the identity change it is.
- **Both stages keep the stores external.** The single most expensive mistake in orchestrating a stateful service is putting Mongo/Redis behind ephemeral container volumes. Whether Docker or K8s, the durable state stays on dedicated hosts/volumes; the orchestrator only runs the stateless `ROLE=api|worker|scheduler` tiers.

The 15-minute self-host story (PRD G8) is untouched by any of this — it's the single-node reference deployment, which stays the default and the install docs' subject. Orchestration is a later decision made by an install that outgrew even the multi-node model, and this section is the record of how it's meant to go.

---

## The three constraints that aren't defaults

### 1. Caddy must not buffer request bodies

Caddy's defaults spool request bodies before forwarding. With multi-GB uploads that fills the proxy's disk, and it only manifests at the sizes that matter (`TRD.md` R5).

```caddyfile
bloberry.example.com {
	reverse_proxy localhost:8080 {
		# Stream request and response bodies rather than buffering to disk.
		flush_interval -1
	}

	request_body {
		max_size 5GB          # matches the API's own ceiling (backend/domains.md §9)
	}

	encode gzip
	log {
		output file /var/log/caddy/bloberry.log
	}
}
```

Note this matters **less** than it first appears and should still be configured: the default browser upload path is a presigned PUT **straight to the storage provider**, bypassing Caddy and Bloberry entirely (`architecture.md` §3.1). Caddy only sees direct uploads and proxied downloads. Configure it anyway — the direct path exists, and the local-disk driver always proxies.

**One domain, not two.** Short URLs live at `/s/<slug>` on the main domain (PRD D6), so there's one Caddy site and one certificate.

### 2. Redis is load-bearing for three things, not one

The default assumes Redis holds sessions. Here it also holds the **authz cache** and the **job queue** (`architecture.md` §4, flagged there for this file to resolve).

**Decision: `appendonly yes` with `appendfsync everysec`.** Rationale by what's stored:

| Data | Loss impact | Needs persistence? |
|---|---|---|
| Sessions / refresh tokens | Everyone logged out | Nice to have |
| Authz cache | Cold start, rebuilt from Mongo on demand | No |
| **Job queue** | **Queued extractions and bundles vanish silently** | **Yes** |

The job queue is what forces it. A user who submitted a 2 GB extraction and sees the job stuck at `queued` forever after a Redis restart has no way to know it's dead. AOF with `everysec` bounds the loss to a second of enqueues at negligible cost.

**Also required:** a `job` in `running` state whose worker died must be reclaimed. A startup sweep moves stale `running` jobs (no heartbeat past a threshold) back to `queued`, bounded by `attempts` so a poison job doesn't loop forever.

Redis binds to `127.0.0.1` only. A publicly-reachable Redis holding session tokens is full account takeover.

### 3. Two scheduled tasks the default has no place for

Both flagged in `architecture.md` §4. **Decision: in-process tickers, not systemd timers** — they need the same config, database handles and driver registry the server already holds, and a separate binary would duplicate all of it for two jobs. This keeps the single-binary story (PRD G8) intact.

| Task | Interval | Purpose |
|---|---|---|
| **Reconciliation sweep** | 15 min | Hard-delete `pending` objects past TTL and their orphaned blobs; abort dangling provider multipart uploads (`architecture.md` ADR-5). Without it, failed uploads accrue storage cost invisibly. |
| **Usage metering** | hourly | Aggregate per-tenant bytes/objects/egress into `usage_snapshots`, compute estimated cost, and reconcile the denormalized `tenants.used_bytes` counters (`architecture.md` §3.8, `ERD.md`). |
| **Backend health check** | 5 min | Probe each registered storage backend; update `health_status`/`health_error` (PRD M19, PA-E1). |
| **Audit retention** | daily | Hard-delete `audit_events` past the retention window (`backend/domains.md` §9). |

**A single-instance assumption is baked in here.** Running two `bloberry` processes would double-run every ticker. Fine for v1 — one VPS, one unit — but it's the first thing that breaks under horizontal scaling, so it's recorded rather than discovered.

---

## systemd

`/etc/systemd/system/bloberry.service`:

```ini
[Unit]
Description=Bloberry object storage service
After=network-online.target mongod.service redis-server.service
Wants=network-online.target
Requires=mongod.service redis-server.service

[Service]
Type=notify
User=bloberry
Group=bloberry
WorkingDirectory=/opt/bloberry
ExecStart=/opt/bloberry/bloberry-server
EnvironmentFile=/etc/bloberry/bloberry.env
Restart=on-failure
RestartSec=5s

# Uploads and extraction are memory-hungry; bound them rather than OOM-killing
# MongoDB, which shares this box.
MemoryMax=2G
LimitNOFILE=65536

# Hardening — the process holds credentials to every connected bucket.
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/lib/bloberry
ProtectKernelTunables=true
ProtectControlGroups=true
RestrictSUIDSGID=true

[Install]
WantedBy=multi-user.target
```

`Requires=` on Mongo and Redis, not merely `After=` — Bloberry cannot serve a single request without either, and failing loudly at boot beats serving 500s.

`MemoryMax=2G` matters specifically because extraction runs **in-process** (`architecture.md` ADR-8, `TRD.md` R12); an unbounded worker on a shared box takes MongoDB down with it.

---

## Local-disk driver storage

`/var/lib/bloberry/objects`, owned by `bloberry:bloberry`, mode `0700`.

**This wants its own volume.** It's the one component that can grow without bound, and filling the root filesystem takes MongoDB and Redis down with it — an outage caused by a tenant uploading files, which is the product working as intended. Mount it separately, and alert at 80%.

Per-tenant quotas (PRD M17) bound *logical* usage but not physical: soft-deleted objects and orphaned blobs awaiting the sweep both consume disk without counting against quota.

---

## Environments and secrets

Two GitHub Environments, `staging` and `production`, each with its own secrets. Dev is local only.

| Secret | Purpose |
|---|---|
| `MONGO_URI` | `mongodb://localhost:27017/bloberry` |
| `REDIS_ADDR` | `localhost:6379` |
| `JWT_SECRET` | Access-token signing (HS256) |
| **`CREDENTIAL_ENCRYPTION_KEY`** | **Envelope-encrypts every storage-backend credential.** See below. |
| `GOOGLE_CLIENT_ID` / `GOOGLE_CLIENT_SECRET` | Google login |
| `SMTP_HOST` / `_PORT` / `_USER` / `_PASSWORD` / `_FROM` | Verification, reset, OTP, invitations |
| `PUBLIC_BASE_URL` | Short URLs, share links, email links |
| `DISK_STORAGE_PATH` | `/var/lib/bloberry/objects` |
| `DISK_SIGNING_SECRET` | HMAC for the disk driver's presigned URLs (`backend/domains.md` §6.3) |
| `MAX_OBJECT_SIZE` | `5368709120` (5 GiB) |

Must match `backend/.env.example`, which is committed; the real values are not.

### `CREDENTIAL_ENCRYPTION_KEY` is not an ordinary secret

It protects credentials to **every connected storage bucket across every tenant** (`TRD.md` R7). A MongoDB dump plus this key is a full compromise of all customer storage; a dump without it is inert. Therefore:

- It lives in `/etc/bloberry/bloberry.env` (mode `0600`, owner `root`) — **never in MongoDB**, and never in the repo.
- It is **excluded from database backups by construction**, since it was never in the database.
- **Rotation** (PRD Q3, previously open — resolved here): a `bloberry-server --rotate-credential-key` maintenance command reads every `storage_backends` document with the old key, re-encrypts with the new, and writes back in one pass. Both keys are supplied during rotation (`CREDENTIAL_ENCRYPTION_KEY` and `CREDENTIAL_ENCRYPTION_KEY_PREVIOUS`); the previous key is removed afterwards. Rotation is offline — the service is stopped, because a half-rotated collection with a running server produces decryption failures on arbitrary requests. Documented as a runbook, and rehearsed on staging before it's ever needed on production.

---

## Workflows

From `templates/github-actions/`, adjusted for this app.

### `ci.yml` — `workflow_dispatch` only

Build order is not optional here (`architecture.md` §7):

```
1. make generate        # oapi-codegen from api/openapi.yaml
                        # FAILS if the working tree changes — keeps spec-first honest
2. cd web && bun install && bun run build    # produces web/dist
3. make lint            # golangci-lint
4. make security        # gosec, standalone
5. make test            # incl. the 100%-branch gate on internal/authz
6. make build           # go build — embeds web/dist
```

Steps 1 and 2 come **before** any Go build. A backend build that runs first embeds the *previous* frontend, or fails outright in CI where nothing has been built yet — this is exactly the trap `templates/repo-layouts.md` warns about for embedded topologies.

The storage **conformance suite** (`backend/domains.md` §6.4) runs against MinIO in a service container on every CI run, and against real S3/R2/OSS/GCS/Azure Blob on a **weekly schedule** — it costs money and needs live credentials, so it doesn't belong on every push.

### `deploy.yml` — `workflow_dispatch`, `runs-on: self-hosted`

`environment` input (`staging` | `production`), then:

```
git pull → make generate → bun run build → make build
  → run pending Mongo index migrations
  → install binary to /opt/bloberry
  → systemctl restart bloberry
  → wait for the health endpoint
```

**Index migrations before the restart**, not after: the new binary may issue queries that need an index the old schema lacks, and MongoDB will happily collection-scan instead of erroring — turning a missing migration into a silent latency cliff rather than a loud failure.

### `release-cli.yml` — `workflow_dispatch`, GitHub-hosted Linux

GoReleaser cross-compiles all five CLI platforms from **one Linux runner** (`cli/README.md`), publishes GitHub Releases with checksums, and pushes the Homebrew tap and Scoop bucket manifests. Needs `GITHUB_TOKEN` plus a PAT with write access to the `homebrew-tap` and `scoop-bucket` repos.

### `release-desktop.yml` — the exception to everything above

**No desktop framework cross-compiles.** Three runners:

| Job | Runner | Cost |
|---|---|---|
| macOS `.dmg` (universal, signed + notarized) | `macos-14` (GitHub-hosted) | **10× minute multiplier** |
| Windows NSIS installer | `windows-latest` (GitHub-hosted) | **2× multiplier** |
| Linux `.deb` / `.rpm` | **self-hosted VPS runner** | free |

This is the one place this project pays GitHub for CI. Manual-trigger-only keeps it bounded — which is a real benefit of the repo-wide no-automatic-triggers convention, not just a philosophical stance. Secrets: `APPLE_DEVELOPER_ID`, `APPLE_ID`, `APPLE_APP_PASSWORD`, `APPLE_TEAM_ID`. Windows signing is deferred (`desktop/README.md`), so no cert secret exists yet.

---

## Self-hosted runner

1. Settings → Actions → Runners → New self-hosted runner; run the generated `./config.sh` on the VPS.
2. `sudo ./svc.sh install && sudo ./svc.sh start` so it survives reboots.
3. One runner serves everything — this is a **monorepo** (`architecture.md` §7), so backend, web, CLI-Linux and desktop-Linux all deploy from the same registered runner.
4. Label it `self-hosted,linux,x64,bloberry` so `release-desktop.yml`'s Linux job can target it specifically while the macOS and Windows jobs go to GitHub-hosted runners.
5. The runner user needs `sudo systemctl restart bloberry` — grant that one command via a sudoers rule, not blanket sudo.

---

## Backups

Not in the defaults, and this app genuinely needs a position stated.

| What | Backed up? |
|---|---|
| **MongoDB** | **Yes** — nightly `mongodump` to a *different* storage backend than any tenant uses, retained 30 days. This is the only irreplaceable state: lose it and every stored object becomes an unaddressable blob. |
| **Redis** | No. Ephemeral by design (`architecture.md` ADR-12). |
| **Local-disk objects** | **Tenant's responsibility, stated explicitly.** Bloberry is not a backup product (PRD NG4) and doesn't silently become one. Self-hosters using the disk driver must be told plainly in the install docs that those bytes have exactly the durability of that one volume. |
| `CREDENTIAL_ENCRYPTION_KEY` | Stored separately from database backups, by hand, in a password manager. Backing it up *alongside* the dump would defeat the whole encryption design. |

**A Mongo restore must be paired with a reconciliation run**, since metadata and bytes are separately stored (`architecture.md` ADR-5): restoring yesterday's dump leaves today's uploaded blobs orphaned and today's deleted objects resurrected as broken references. The sweep resolves the first; the second needs a documented manual pass. Part of the restore runbook, not an afterthought.

---

## Monitoring

Minimal and honest — the point is knowing when the box is in trouble, not building an observability platform.

- **`GET /healthz`** — liveness. **`GET /readyz`** — checks Mongo and Redis; Caddy and the deploy workflow both wait on it.
- Structured JSON logs to stdout → journald. **Never logged:** access-key secrets, presigned URLs (they *are* credentials for their TTL), storage credentials.
- **Alert on:** disk > 80% on the objects volume, disk > 80% on root, any storage backend `unreachable` for > 15 min, job queue depth > 100, and `bloberry.service` restart-looping.

---

## Resolved

**PRD Q3** (envelope-encryption key rotation) — resolved above: an offline `--rotate-credential-key` maintenance command with both keys supplied, rehearsed on staging.

## Files

- `deploy/` in the implementation repo — the systemd unit, `Caddyfile`, `.env.example` per stage
- `.github/workflows/` — `ci.yml`, `deploy.yml`, `release-cli.yml`, `release-desktop.yml`
- `tasks/` — the implementation task list, once `build-infra` has run
