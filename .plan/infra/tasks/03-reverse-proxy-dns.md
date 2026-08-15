# Task group — 03 reverse proxy / TLS / DNS

**Depends on:** `02-process-model.md` (a service to proxy to). **Blocks:** `05-first-deploy.md`.

Per `infra/README.md`: **Caddy**, one domain (short URLs live at `/s/<slug>` on the main domain, PRD D6 — one site, one cert).

- [ ] **`Caddyfile` written** per the README's exact config for `bloberry.example.com`:
  - `reverse_proxy localhost:8080` with **`flush_interval -1`** — Caddy's default buffers request bodies; multi-GB uploads would spool to the proxy's disk (TRD R5). Configured even though presigned PUTs bypass it — the direct path and the local-disk proxy always use it.
  - **`request_body { max_size 5GB }`** — matches the API's own ceiling (`backend/domains.md` §9).
  - `encode gzip`, log output to `/var/log/caddy/bloberry.log`.
- [ ] **Caddy installed and enabled** on the VPS, running on :80/:443.
- [ ] **TLS via Caddy's automatic HTTPS** — one certificate for `bloberry.example.com`, no manual certbot.
- [ ] **DNS records pointed at the VPS** — `A` record for the main domain; no second domain (one site, one cert).
- [ ] **The SPA catch-all works through the proxy** — a deep link like `/files/abc123` reaches the embedded frontend's router, not a Caddy 404 (`web/tasks/01-setup.md`'s catch-all through the real proxy).

**verification:** `https://bloberry.example.com` serves the dashboard with a valid cert; `curl -H "Expect:" --data-binary @5GB-file` through the proxy doesn't spool to `/var/lib/caddy`; a deep link loads the SPA.
