# Task group — 17 storage (the driver abstraction)

**Depends on:** `01-setup`, `03-platform` (crypto for credential decryption). **Blocks:** `08-domain-object`, `09-domain-share`, `12-domain-job`, `15-domain-admin`. **Not a domain** — infrastructure the object/share/job/admin domains call. The thing the product exists to provide (TRD R1/R2, PRD G1/G2, ADR-2).

- [ ] **`internal/storage/driver.go`** — the `Driver` interface + `Capabilities` descriptor (`domains.md` §6.1): `Put`/`Get`/`Delete`/`Stat`; `PresignGet`/`PresignPut`; `MultipartInit`/`PresignPart`/`Complete`/`Abort`; `HealthCheck`. Capabilities: Presign, Multipart, MinPartSize, MaxPartCount, StorageClasses, ServerSideCopy, RangeRequests, ObjectAttributes.
- [ ] **Design against the hardest driver** (TRD R1) — **local disk** for presigning (no external signer at all) and **GCS** for credentials (service-account signer). The interface is settled by these two before the easy ones; designing against S3 first is how the abstraction leaks.
- [ ] **`s3/driver.go`** — SigV4 static key pair; serves S3, R2 (via account-scoped endpoint), MinIO, B2, Spaces, Wasabi.
- [ ] **`r2` capability flags** — R2 is *not* fully S3-compatible (TRD R2): no storage classes, no `GetObjectAttributes`, different multipart part-size behavior. Declared via `Capabilities`, asserted by the conformance suite — never assumed.
- [ ] **`oss/driver.go`** — Alibaba OSS's own signature version, separate SDK.
- [ ] **`gcs/driver.go`** — service-account signer (or IAM `signBlob` without a key file); fundamentally different credential shape from the other three.
- [ ] **`azblob/driver.go`** — Azure Blob Storage via `azure-sdk-for-go/sdk/storage/azblob`; SharedKey/SAS (or AAD) auth, container-scoped, block-blob staging for multipart. A separate SDK with its own primitive model — deliberately NOT an S3 endpoint override (TRD tech stack). Presign via SAS tokens.
- [ ] **`disk/driver.go`** — local VPS volume; **Bloberry-issued HMAC tokens** against its own `/v1/objects/:id/raw` endpoint (`domains.md` §6.3). `Capabilities.Presign` true only because Bloberry implements the signer itself. A proxy wearing a presign interface — deliberate, so `object` calls `PresignGet` unconditionally (PRD G2).
- [ ] **`registry.go`** — `backend_id` → constructed Driver; decrypts credentials in memory at construction (`platform/crypto`), never logs or returns them.
- [ ] **`conformance/suite.go`** — the **one suite every driver must pass** (PRD G2): round-trip put/get/stat/delete; range requests; presigned PUT/GET from an external HTTP client; multipart across the declared minimum part size; multipart abort leaves nothing; overwrite; delete-nonexistent (must not error); listing correctness; **capability honesty** — every `true` capability exercised, every `false` asserted to fail in the documented way.
- [ ] **Conformance CI wiring** — run against MinIO locally (fast, no credentials) and against **real S3/R2/OSS/GCS/Azure Blob on a schedule**, not every push (costs money, needs real credentials) — cross-ref `20-guards` / `infra/tasks/`.

**tests:** the conformance suite itself is the test — green for MinIO locally; a capability-honesty failure (a driver declaring `true` it can't do) fails the suite by design.
