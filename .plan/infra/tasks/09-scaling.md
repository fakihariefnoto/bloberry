# Task group — 09 scaling (cloud-backend model)

**Depends on:** `02-process-model.md` (the single-node base), `04-environments-secrets.md`, `05-first-deploy.md` (staging proven before scaling). **Blocks:** nothing (post-v1 operational growth).

Per `infra/README.md` §Scaling out — the multi-node shape for installs outgrowing one box. **Local-disk driver is single-node by design and absent here**; every tenant is on a cloud backend.

- [ ] **Worker split** — a `cmd/bloberry-worker` binary consuming the same Redis `job:queue` (the R12 escape hatch — a deployment change, not a rewrite); API nodes run with `ROLE=api` (no worker, no tickers); worker concurrency scales independently of request handling.
- [ ] **Scheduler + leader election** — the four in-process tickers (reconciliation sweep, usage metering, backend health, audit retention) run on one `ROLE=scheduler` node holding a Redis lease (`SET key NX EX` heartbeat); the loser skips. Single-instance assumptions confined to the scheduler, never to API/worker nodes.
- [ ] **Mongo replica set** (3 nodes) — writes to the primary, automatic failover; the API's Mongo driver handles it; migrations run from the scheduler node only.
- [ ] **Redis sentinel + replication** — Redis gains replication + sentinel (a lost Redis in multi-node = a lost job queue AND everyone logged out); **AOF `appendonly yes` / `everysec` carries over unchanged** (the queue's durability is non-negotiable).
- [ ] **Caddy as load balancer** — `reverse_proxy` gains an upstream list with health checks (one entry per API node) instead of a single `localhost:8080`; the no-body-buffering + 5GB ceiling config carries over verbatim.
- [ ] **Caddy health-check eviction** — `health_uri /healthz`, `health_interval`/`timeout`, `fail_duration`/`max_fails` so a dead node drops out of the pool automatically and re-enters on recovery (attach + detach without human action at the edge).
- [ ] **DNS-based discovery** — API nodes named (`api-1.example.internal`, …); Caddy `upstreams dns api.example.internal` resolves the pool from A/SRV records; **attaching a new node = adding one A record** (no Caddyfile edit, no reload). Manual static-pool config is the fallback for ≤5 nodes.
- [ ] **Node-provisioning runbook** — provision → same binary + `ROLE=api` + shared `.env` → add A record (or Caddyfile line) → health check admits it. **No data to move, nothing to rebalance** — the node is stateless from boot.
- [ ] **Bottleneck diagnosis before adding nodes** — check whether Mongo primary / Redis is saturated vs. API concurrency; adding API nodes scales concurrency but not the stores (Mongo secondaries for reads; sharding is the later rung).
- [ ] **`deploy.yml` targets N nodes** — the same binary with `ROLE=api|worker|scheduler` per host; the deploy workflow iterates the node inventory rather than one fixed path.
- [ ] **No shared-filesystem disk** — stated as a constraint: a disk backend is refused (or pinned to a single node) on scaled installs, never offered as a multi-node path (ADR-15).

## 09.1 Docker stage (future, recorded — not the default)

The documented pivot toward containerization, for when the multi-node systemd model itself outgrows.

- [ ] **`Dockerfile` for the same binary** — containerize `bloberry-server`/`bloberry-worker` (and the scheduler via `ROLE`), NOT the stores.
- [ ] **Stores stay external** — Mongo RS + Redis sentinel remain on dedicated hosts/volumes; **never containerize the stateful stores with ephemeral volumes** (the single most expensive orchestrator mistake).
- [ ] **Scale-out via a service** — a Swarm service / compose scale set runs `ROLE=api` replicas behind the reverse proxy; the worker and scheduler stay separate services.
- [ ] **Caddy stays the edge** — the no-buffering + health-check config carries over; the orchestrator's ingress is a later concern.

## 09.2 Kubernetes stage (future, recorded — the identity change)

The full pivot. **This is where PRD G8's "one binary, 15 minutes" promise formally ends** — recorded so it's recognized as the identity change it is, not an accretion.

- [ ] **API as a Deployment + HPA** — `ROLE=api` replicas autoscaling on CPU/request latency.
- [ ] **Worker as its own Deployment** — queue consumers scale independently; the scheduler as a single-replica Deployment (the Redis lease still guards the tickers).
- [ ] **Ingress replaces Caddy's edge role** — the no-buffering constraint (TRD R5) must survive whatever ingress terminates TLS/streaming.
- [ ] **Mongo/Redis external or operator-managed** — never ephemeral volumes; the AOF + replica-set decisions carry over.

**verification:** a staging scale-out of 2 API + 1 worker + 1 scheduler serves traffic with any one API node killed (failover works); a worker crash reclaims its running job; a scheduler restart re-elects the lease without double-running a ticker; `bloberry ls`/upload round-trips through the LB; a node added by A record alone serves traffic within the health-check interval; (Docker/K8s stages) rolling scale-up/down of `ROLE=api` replicas with the stores untouched.
