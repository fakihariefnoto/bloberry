# Task group — 33 cross-page flows

**Depends on:** the pages it chains — `07`–`33`. **Blocks:** nothing (integration file). 

The vertical-slice view a per-page split loses. Each flow is one checkbox group listing the pages it touches and its cancel/failure landing, referencing (not duplicating) the per-page files. Source: `web/navigation.md` §Flow chains.

## 33.1 Auth — login → dashboard

- [ ] `login` → `files` on success; `?next` round-trip (login returns to the original destination). Chain connects end-to-end.
- [ ] **Session expiry landing**: any 401 on any page → `login?next=<current>` → back to the same place after login. In-flight presigned uploads survive (the `complete` call fails until re-login; the queue shows them retryable, not failed).
- [ ] OTP and Google paths both land on `files` (or `?next`), not on a dead end.

## 33.1a Pair a mobile device (QR, M22)

- [ ] user menu → `pair-device` → QR renders (one-time token, ~2 min TTL) with the capability warning → phone scans in the mobile app → mobile session established → the QR card flips to "paired" and dies. Chain: `25-page-pair-device` + the auth `/pair/issue` + `/pair/verify` endpoints.
- [ ] **Expired** → auto-refresh with a notice (the old token is dead even unscanned). **Refresh** → a fresh token; the previous one is invalidated.
- [ ] **A used code must not still render scannable** — the paired state kills the QR (ADR-13 single-use).

## 33.1b Export a desktop login file (M23)

- [ ] user menu → `pair-device` → config-file card → passphrase (min strength) → download `.bloberry` file (signed, encrypted client-side, 24h import window) → desktop imports it. Chain: `25-page-pair-device` + the auth `/config/issue` endpoint.
- [ ] **Wrong passphrase at import** → "wrong passphrase" (never transits the server). **Window passed** → "download a fresh one".
- [ ] The imported session stays revocable via `auth logout` (DT-E2).

## 33.2 Upload

- [ ] `files` → drop or Upload → `upload-queue` docks (no navigation) → per-file progress → complete. Chain: `13-page-files` + `06` UploadQueue.
- [ ] **Name collision** → inline replace/keep-both/cancel per file, queue stays open (PRD MV-E2).
- [ ] **Quota exceeded** → that file fails with `quota_exceeded`, others continue (PRD MV-E1); queue shows a link to `usage` (`13` → `21`).
- [ ] **Navigate away mid-upload** → queue persists (store, not component state); uploads continue.
- [ ] **Close the browser tab** → presigned uploads die; `beforeunload` fires while anything is in flight.

## 33.3 Share a file

- [ ] `files` or `file-detail` → `share-dialog` → choose signed-link + TTL / short URL / make public → created → toast with copy button → link appears in `shares`. Chain: `13`/`14` + `share-dialog` + `15`.
- [ ] **Cancel** → back to the originating page, nothing created.
- [ ] **Make public** → `confirm-destructive` variant first (public is effectively irreversible once the URL is copied — ADR-3, R11). Consistent on both entry points.

## 33.4 Issue an access key

- [ ] `applications` → `application-detail` → Create key → scope form → `key-created` modal showing the secret **once** → close → key appears masked in the list forever. Chain: `17` → `18`.
- [ ] **Close without copying** → the secret is unrecoverable. The modal says so before it can be dismissed; dismissal requires acknowledgement, never a backdrop click.
- [ ] Revoke path (`18`) states last-used + the last-active-key consequence (PRD TA-E3).

## 33.5 Extract an archive

- [ ] `files` → row action Extract → target-folder confirm → 202 → toast "Extraction queued" linking `jobs` → `jobs` shows progress. Chain: `13` → `16`.
- [ ] **Failure** → job row shows the real reason; **the target folder is unchanged** (PRD AP-E2). The toast→jobs link must survive (running is the default filter on arrival).

## 33.6 Delete a large folder

- [ ] `files` → folder row Delete → `confirm-destructive` stating the **real object count** → above threshold becomes a job → toast links to `jobs`. Chain: `13` → `06 ConfirmDestructive` → `16`.
- [ ] **Cancel** → nothing happens. Typed-name confirmation required (PRD TA-E1).

## 33.7 Tenant switch

- [ ] Any route → tenant switcher → **always lands on `files` at the new tenant's root**, never the equivalent path in the new tenant (a `folderId` from tenant A is meaningless in tenant B). Verify across `03-routing`'s switcher + `13-page-files`.

**tests (per chain):** one end-to-end test per chain covering the happy path and its cancel/failure landing — 9 flows total.
