# Task group — 18 flows

**Depends on:** the domains it chains (`04`–`15`). The vertical-slice view across domains — the backend-internal flows from `backend/domains.md` §4 plus the genuinely multi-step app-specific flows. This file references (doesn't duplicate) the per-domain tasks and adds only the end-to-end wiring checks.

## 18.1 Auth flows (from `domains.md` §4)

- [ ] **Signup by invitation** (§4.1) — invite token validated → existing-email adds a membership / new user created (email_verified:false + verification email) → audit `member.join` → tokens. End-to-end across `auth` + `tenant` + `audit`.
- [ ] **Login** (§4.2) — platform-aware token issue; identical response for unknown-email vs wrong-password; constant-time. Across `auth`.
- [ ] **Refresh with rotation** (§4.3) — old token deleted, new pair issued; a rotated token presented again fails. Across `auth`.
- [ ] **Forgot → reset** (§4.4) — identical response whether or not the email exists; successful reset **invalidates every session**; reset link 30-min TTL. Across `auth`.
- [ ] **OTP login** (§4.5) — request (rate-limited, hashed code in Redis) → verify (5-attempt cap). Across `auth`.
- [ ] **Google login** (§4.6) — id_token verification; auto-link only when the email is verified; `no_invitation` otherwise. Across `auth`.
- [ ] **Two-schemes-one-principal seam** (§4.7) — a request authenticated by JWT and one by access key produce the same `Principal` shape, and no handler branches on which was used. End-to-end through `platform/httpx` middleware → `authz`.
- [ ] **QR pairing (M22, §4.8)** — web (authed) mints a one-time 2m Redis token → mobile scans → verify exchanges it for a mobile-TTL session, DELs on use. End-to-end across `auth` + Redis + the audit append.
- [ ] **Desktop config import (M23, §4.9)** — web issues a signed payload (import-window + refresh token) → client encrypts with passphrase-derived key → desktop decrypts locally, validates signature + window, stores in keychain → refresh. End-to-end across `auth` + Redis; the passphrase never touches the server.
- [ ] **TOTP 2FA (M24, §4.10)** — provision (secret once) → enable (confirm code, mint backup codes) → every human login ends with `totp_required` → verify-totp (code or backup) → session issued only after. End-to-end across `auth` + `user` + Redis (`totp:pending:`) + the admin reset path.

## 18.2 App-specific flows (from `architecture.md` §3)

- [ ] **Presigned upload** (§3.1) — quota checked → presign-PUT issued → browser writes bytes to storage → complete promotes `pending`→`active`. Across `object` + `tenant` (quota) + `storage`. The direct and multipart variants are the same seam with different byte paths.
- [ ] **Redirect download** (§3.2) — object → signed redirect (cloud) or raw proxy (disk); revoked/expired → 410. Across `share` + `object` + `storage`. The R11 caveat (presigned URL outlives revocation) recorded in the flow.
- [ ] **Short URL** (§3.3) — slug → object; expired/revoked → HTML 410 page. Across `share` + the public router.
- [ ] **Access-key authz** (§3.4) — access key → principal → resolver → allowed/denied. Across `apikey` + `authz`.
- [ ] **Queued extraction** (§3.5) — upload archive → enqueue → worker extracts to staging → atomic commit → job status. Across `job` + `object` + `folder` + `storage`. Failure leaves the target folder unchanged (PRD AP-E2).
- [ ] **Desktop sync / mobile resume** (§3.7) — the server side: multipart resume reconciliation (`parts_received` vs provider), which is what a resuming client queries. Across `object` + `storage`.
- [ ] **Usage metering** (§3.8) — hourly snapshot per tenant + counter reconciliation + cost from rate card. Across `usage` + `tenant` + `admin`.

**tests:** one end-to-end integration test per flow (10 auth + 7 app-specific = 17), exercising the real seams across domains with Mongo + Redis + MinIO — not mocked units.
