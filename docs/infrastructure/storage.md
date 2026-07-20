---
sidebar_position: 7
---

# File Storage

VEF ships a provider-neutral storage abstraction, three built-in providers, a multipart upload protocol, a typed CRUD lifecycle facade for keeping model file references in sync with the backend, and a transactional outbox for downstream cleanup.

> The typed lifecycle facade is `Files` / `FilesFor[T]` — there is no `Promoter[T]`. The upload protocol is chunked multipart with an explicit claim/queue lifecycle, and principal authorization is threaded through the whole lifecycle. This page describes the current public surface.

## Supported Providers

| Provider value | Backend |
| --- | --- |
| `memory` | in-process map; tests and ephemeral demos |
| `filesystem` | local filesystem |
| `minio` | MinIO / S3-compatible object storage |

`storage.provider` selects the backend. Without configuration the module
defaults to `memory` and logs a warning; objects are lost on restart.

Set `vef.storage.auto_migrate = true` when the storage tables should be created
by the module at startup. The migration is idempotent and checks
`sys_storage_upload_claim`, `sys_storage_upload_part`, and
`sys_storage_pending_delete`.

## `storage.Service` Interface

Application code depends on `storage.Service`, never on a provider-specific type:

```go
type Service interface {
    PutObject(ctx, opts PutObjectOptions) (*ObjectInfo, error)
    GetObject(ctx, opts GetObjectOptions) (io.ReadCloser, *ObjectInfo, error)
    DeleteObject(ctx, opts DeleteObjectOptions) error
    DeleteObjects(ctx, opts DeleteObjectsOptions) error
    CopyObject(ctx, opts CopyObjectOptions) (*ObjectInfo, error)
    StatObject(ctx, opts StatObjectOptions) (*ObjectInfo, error)
}
```

Option types: `PutObjectOptions`, `GetObjectOptions`, `DeleteObjectOptions`, `DeleteObjectsOptions`, `CopyObjectOptions`, `StatObjectOptions`. Use the option struct for every call — direct positional arguments are not supported on purpose so that adding fields stays additive.

`GetObject` returns the body reader together with best-effort `ObjectInfo`.
Callers must close the reader and nil-check the `ObjectInfo`.

## Multipart Upload

The framework's upload protocol is **chunked multipart only** — there is no single-PUT upload. Every backend implements `storage.Multipart`:

```go
type Multipart interface {
    PartSize() int64
    MaxPartCount() int
    InitMultipart(ctx, opts InitMultipartOptions) (*MultipartSession, error)
    PutPart(ctx, opts PutPartOptions) (*PartInfo, error)
    CompleteMultipart(ctx, opts CompleteMultipartOptions) (*ObjectInfo, error)
    AbortMultipart(ctx, opts AbortMultipartOptions) error
}
```

Obtain the typed handle with `storage.MultipartFor(svc)` (returns `nil` when the backend does not implement chunked uploads). The contract guarantees:

- Distinct part numbers may upload concurrently; same-part calls are last-writer-wins.
- Every non-final part must be at least `PartSize()` bytes.
- `CompleteMultipart` verifies every recorded `(PartNumber, ETag)` and that parts cover `1..N` contiguously.
- Sessions close after `CompleteMultipart` or `AbortMultipart`; further calls return `ErrUploadSessionNotFound`. `AbortMultipart` is idempotent.

> The `sys/storage.list_parts` RPC action exists to let clients resume an in-flight upload, but it is served from the framework's part-store table, not from a `ListParts` method on `storage.Multipart` — the backend interface itself only exposes the six methods above.

## Built-In Resource: `sys/storage`

The storage module registers an RPC resource with the multipart upload actions:

| Action | Access | Input | Output | Purpose |
| --- | --- | --- | --- | --- |
| `init_upload` | Bearer auth (engine default) | `InitUploadParams` | `InitUploadResult` | create a pending claim, open a multipart session, and return opaque `claimId` plus the negotiated `partSize` |
| `upload_part` | Bearer auth (engine default) | `UploadPartParams` (multipart form) | `UploadPartResult` | upload one part of an open session |
| `list_parts` | Bearer auth (engine default) | `ListPartsParams` | `ListPartsResult` | inspect parts already uploaded for a session |
| `complete_upload` | Bearer auth (engine default) | `CompleteUploadParams` | `CompleteUploadResult` | seal a session; the server assembles the final manifest from recorded parts |
| `abort_upload` | Bearer auth (engine default) | `AbortUploadParams` | success, no `data` payload | abort and release a session |

Download is served via the proxy middleware described below.

None of the five actions declares a per-action permission, a public flag, a
custom rate limit, or audit logging, so all of them inherit the API engine
defaults: Bearer authentication plus the default rate limit (100 requests per
5-minute sliding window unless `vef.api.rate_limit` overrides it). On top of
authentication,
every action that takes a `claimId` enforces per-claim ownership — only the
principal that created the claim may operate on it. A claim ID that does not
exist is answered exactly like a claim owned by someone else
(`result.ErrAccessDenied`, code `1100`, HTTP 403): the API deliberately avoids
revealing whether a given claim ID exists. The one exception is
`abort_upload`, which treats an unknown claim as already aborted and returns
success.

All HTTP uploads use this same protocol: `init_upload -> upload_part ->
complete_upload`. Small files still return `partCount = 1`; there is no
single-PUT HTTP action. The `public` flag defaults to private behavior, and
`vef.storage.allow_public_uploads` must be true before clients can request
`pub/` keys.

Client `contentType` values are sanitized. Safe binary, image, audio, video,
font, archive, and PDF types are accepted; unsafe same-origin types such as
`text/html` and `application/javascript` are replaced by extension detection or
`application/octet-stream`.

### `init_upload`

`InitUploadParams`:

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `filename` | `string` | Yes | Original filename, at most 255 characters (`validate:"required,max=255"`). Its extension is reused for the object key when it is purely alphanumeric (`.pdf`, `.tar` — pattern `^\.[a-zA-Z0-9]+$`); anything else falls back to `.bin`. Persisted on the claim row and echoed back as `originalFilename`. |
| `size` | `int64` | Yes | Exact total byte count of the object, at least 1 (`validate:"required,min=1"`). Validated against `vef.storage.max_upload_size` (default 1 GiB) — `ErrCodeUploadSizeExceedsLimit` when over — and used to compute the part plan. `complete_upload` later verifies the uploaded total matches this declaration. |
| `contentType` | `string` | No | Client-suggested MIME type, at most 127 characters (`validate:"max=127"`). Sanitized server-side: `image/*`, `audio/*`, `video/*`, `font/*` prefixes and the exact types `application/pdf`, `application/zip`, `application/gzip`, `application/x-tar`, `application/octet-stream` are accepted as-is; anything else is replaced by extension-based detection, falling back to `application/octet-stream`. |
| `public` | `bool` | No | `true` places the key under `pub/` instead of `priv/`. Rejected with `ErrCodePublicUploadsNotAllowed` unless `vef.storage.allow_public_uploads = true`. |

Behavior:

- Rejected with `ErrCodeMultipartNotSupported` when the configured backend does
  not implement `storage.Multipart`; with `ErrCodeUploadTooManyParts` when the
  part plan (`ceil(size / partSize)`) exceeds the backend's `MaxPartCount()`
  (10000 on MinIO; filesystem and memory are unbounded); and with
  `ErrCodeTooManyPendingUploads` when the caller already holds
  `vef.storage.max_pending_claims` pending claims (default 100, best-effort
  count).
- The generated key is date-partitioned:
  `<pub/|priv/>YYYY/MM/DD/<uuid><ext>`.
- The pending claim expires after `vef.storage.claim_ttl` (default 24h); parts
  of an unfinished upload survive until then.

`InitUploadResult`:

| Field | Type | Description |
| --- | --- | --- |
| `key` | `string` | Final object key under `priv/` or `pub/`, fixed at init time. |
| `claimId` | `string` | Opaque session handle for all follow-up actions. This is the **only** client-visible identifier — the backend's multipart `UploadID` stays on the server. |
| `originalFilename` | `string` | The client-supplied `filename`, persisted on the claim row (not in backend metadata), so it survives independent of the storage backend. |
| `partSize` | `int64` | Backend-authoritative slice size in bytes (16 MiB on MinIO, 4 MiB on filesystem, 64 KiB on memory). Every part except the last must be exactly this size. Clients must not assume a value — always use the returned figure. |
| `partCount` | `int` | Number of parts to upload; `ceil(size / partSize)`. Small files still get `partCount = 1`. |
| `expiresAt` | timestamp (RFC 3339) | When the claim (and with it the upload session) lapses. |

### `upload_part`

`upload_part` rejects JSON bodies (`ErrCodeUploadRequiresMultipart`). Send
`multipart/form-data` with `resource`, `action`, and `version` as plain form
fields, `params` as a JSON string, and the raw part bytes in a form part named
`file` (`ErrCodeUploadRequiresFile` when missing).

`UploadPartParams`:

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `file` | form file part | Yes | Raw bytes of this part. Larger than `partSize` fails with `ErrCodeUploadPartTooLarge`; a non-final part smaller than `partSize` fails with `ErrCodeUploadPartTooSmall` (only the last part may carry the remainder). |
| `claimId` | `string` | Yes | The handle returned by `init_upload` (`validate:"required"`). The claim must be owned by the caller, `pending`, and unexpired. |
| `partNumber` | `int` | Yes | 1-based part position (`validate:"required,min=1"`). Values above `partCount` fail with `ErrCodeUploadPartNumberOutOfRange`. |

Behavior:

- Distinct part numbers may upload concurrently. Re-sending a part number
  overwrites the earlier bytes — the part row is upserted and the latest
  backend ETag wins, matching the backend's last-writer-wins semantics.
- The claim's declared `size` is re-validated against the **current**
  `vef.storage.max_upload_size` on every part, so tightening the cap at
  runtime also stops in-flight uploads.
- The backend ETag is recorded in the framework's part table and intentionally
  **not** returned: `complete_upload` assembles the manifest server-side, so
  clients never round-trip ETags.

`UploadPartResult`:

| Field | Type | Description |
| --- | --- | --- |
| `partNumber` | `int` | The accepted part position, echoed back. |
| `size` | `int64` | Byte count the backend recorded for this part. |

### `list_parts`

`ListPartsParams`:

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `claimId` | `string` | Yes | The in-flight session to inspect (`validate:"required"`). The claim must be owned by the caller, `pending`, and unexpired — completed claims answer `ErrCodeClaimNotPending`. |

`ListPartsResult`:

| Field | Type | Description |
| --- | --- | --- |
| `parts` | `ListedPart[]` | Parts already accepted, ordered by `partNumber` ascending. |

`ListedPart`:

| Field | Type | Description |
| --- | --- | --- |
| `partNumber` | `int` | 1-based part position. |
| `size` | `int64` | Recorded byte count. |

The list is served from the framework's part-store table, not the backend's
native listing, and it is the same table `complete_upload` assembles from —
every part listed here is honored as-is. Part ETags are intentionally omitted.

### `complete_upload`

`CompleteUploadParams`:

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `claimId` | `string` | Yes | The session to seal (`validate:"required"`). The server reconstructs the parts manifest from its own part-store rows — client-supplied ETags are never accepted. |

Behavior:

- Fails with `ErrCodeUploadPartsIncomplete` when fewer than `partCount` parts
  are recorded, and with `ErrCodeClaimExpired` when the claim TTL has elapsed.
  The declared size is re-validated against the current
  `vef.storage.max_upload_size` cap.
- After assembly the server compares the object size against the declared
  `size`; a mismatch deletes the assembled object immediately and returns
  `ErrCodeUploadSizeMismatch`.
- On success, one transaction marks the claim `uploaded` and clears its part
  rows. The object now waits for business adoption (see `Files` below).
- Idempotent: calling `complete_upload` again on an `uploaded` claim re-stats
  the object and returns the same shape. A retry that arrives after the
  backend session closed but before the bookkeeping committed re-stats the
  object and commits the same transaction; if neither session nor object
  exists the call fails with `ErrCodeUploadObjectNotFound`.

`CompleteUploadResult` — `storage.ObjectInfo` plus the framework-tracked
original filename:

| Field | Type | Description |
| --- | --- | --- |
| `bucket` | `string` | Backend bucket. MinIO reports the real bucket; the bucket-less backends use the sentinels `filesystem` and `memory`. |
| `key` | `string` | Final object key (same as `init_upload` returned). |
| `eTag` | `string` | ETag of the assembled object — not a part ETag, and never an input to this action. |
| `size` | `int64` | Final object size in bytes. |
| `contentType` | `string` | The sanitized content type stored with the object. |
| `lastModified` | timestamp (RFC 3339) | Backend last-modified time. |
| `metadata` | `object` (string→string) | Backend user metadata with canonicalized keys; omitted when empty. The upload RPC never accepts client metadata, so uploads through this resource carry none. |
| `originalFilename` | `string` | The filename captured at `init_upload`. |

### `abort_upload`

`AbortUploadParams`:

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `claimId` | `string` | Yes | The session to cancel (`validate:"required"`). |

Behavior — abort is the retry-safe cleanup path:

- An unknown `claimId` returns success (the only action that does not answer
  a missing claim with access-denied); a claim owned by another principal is
  still rejected with `result.ErrAccessDenied`.
- Only `pending` claims are aborted. Calling it on an `uploaded` claim is a
  no-op success — abort never deletes a finalized object.
- A pending claim's backend session is aborted (failure surfaces
  `ErrCodeAbortFailed`), any published object bytes are deleted, and one
  transaction drops the part rows and the claim row.
- The success response carries no payload: `data` is `null`.

## Client Walkthrough: Multipart Upload

The exact wire sequence a client implements against `POST /api`, shown for a 40 MiB `report.pdf` on a MinIO-backed server (`partSize` 16 MiB, so 3 parts). All five actions require authentication (Bearer by default), and every follow-up call routes by the `claimId` returned from `init_upload` — the backend's multipart `UploadID` never leaves the server. Success responses use the standard envelope (`message` text is language-dependent); failures reuse it with a non-zero `code` from the storage range (2200–2299, see the error table below).

### 1. `init_upload`

```bash
curl http://localhost:8080/api \
  -H 'Authorization: Bearer <token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "resource": "sys/storage",
    "action": "init_upload",
    "version": "v1",
    "params": {
      "filename": "report.pdf",
      "size": 41943040,
      "contentType": "application/pdf",
      "public": false
    }
  }'
```

```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "key": "priv/2026/07/09/6c9e6f0e-8d5a-4d5e-9a3b-2f4a1c7e9b21.pdf",
    "claimId": "b3a2c1d0-4e5f-47a9-8bcd-ef0123456789",
    "originalFilename": "report.pdf",
    "partSize": 16777216,
    "partCount": 3,
    "expiresAt": "2026-07-10T12:04:05Z"
  }
}
```

`size` must be the exact byte count: `complete_upload` deletes the assembled object and fails with a size mismatch when the uploaded total differs. `partSize` is the backend's authoritative slice size — every part except the last must be exactly `partSize` bytes; the last part carries the remainder.

### 2. `upload_part` (× `partCount`)

`upload_part` rejects JSON bodies. Send `multipart/form-data`: `resource`, `action`, and `version` as plain form fields, `params` as a JSON string, and the raw bytes in a form part named `file`:

```bash
split -b 16777216 report.pdf part-   # part-aa, part-ab, part-ac

curl http://localhost:8080/api \
  -H 'Authorization: Bearer <token>' \
  -F 'resource=sys/storage' \
  -F 'action=upload_part' \
  -F 'version=v1' \
  -F 'params={"claimId":"b3a2c1d0-4e5f-47a9-8bcd-ef0123456789","partNumber":1}' \
  -F 'file=@part-aa'
```

```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "partNumber": 1,
    "size": 16777216
  }
}
```

Repeat with `partNumber` 2 and 3 (`part-ac` is the 8388608-byte remainder). Distinct part numbers may upload concurrently; re-sending a part number overwrites the earlier bytes (last-writer-wins). The backend ETag is recorded server-side and intentionally not returned — clients never round-trip ETags.

### 3. `complete_upload`

```bash
curl http://localhost:8080/api \
  -H 'Authorization: Bearer <token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "resource": "sys/storage",
    "action": "complete_upload",
    "version": "v1",
    "params": { "claimId": "b3a2c1d0-4e5f-47a9-8bcd-ef0123456789" }
  }'
```

```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "bucket": "app-files",
    "key": "priv/2026/07/09/6c9e6f0e-8d5a-4d5e-9a3b-2f4a1c7e9b21.pdf",
    "eTag": "9b2cf535f27731c974343645a3985328-3",
    "size": 41943040,
    "contentType": "application/pdf",
    "lastModified": "2026-07-09T12:08:15Z",
    "originalFilename": "report.pdf"
  }
}
```

The server assembles the parts manifest from its own table; with fewer than `partCount` parts recorded, the call fails with `ErrCodeUploadPartsIncomplete`. Retries are idempotent — a retry arriving after the backend session closed re-stats the object and returns the same shape. The claim is now `uploaded` and waits for business adoption (see `Files` below).

### Resume an interrupted upload: `list_parts`

Accepted parts survive client restarts while the claim is still pending and unexpired (`expiresAt`). Ask the server which parts it holds, skip those, and upload only the rest:

```bash
curl http://localhost:8080/api \
  -H 'Authorization: Bearer <token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "resource": "sys/storage",
    "action": "list_parts",
    "version": "v1",
    "params": { "claimId": "b3a2c1d0-4e5f-47a9-8bcd-ef0123456789" }
  }'
```

```json
{
  "code": 0,
  "message": "Success",
  "data": {
    "parts": [
      { "partNumber": 1, "size": 16777216 },
      { "partNumber": 2, "size": 16777216 }
    ]
  }
}
```

Here part 3 is missing: upload it, then call `complete_upload`. The list is ordered by `partNumber` ascending, and every listed part is recorded with its ETag server-side and honored by `complete_upload` as-is.

### Cancel an upload: `abort_upload`

Same JSON envelope with `"action": "abort_upload"` and the `claimId` in `params`:

```json
{ "code": 0, "message": "Success", "data": null }
```

Abort is idempotent: an unknown or already-aborted `claimId` still returns `code: 0`. Only pending claims are aborted — calling it on an `uploaded` claim is a no-op and never deletes a finalized object.

### Downloading through the proxy

Downloads are plain HTTP `GET`s against the proxy route, not RPC actions:

```bash
curl -O http://localhost:8080/storage/files/pub/2026/07/09/6c9e6f0e-8d5a-4d5e-9a3b-2f4a1c7e9b21.pdf
```

`pub/*` keys are served anonymously. For any other key the proxy calls `FileACL.CanRead` with the request principal; the framework does not resolve a token on this route itself, so private downloads need a registered `FileACL` plus app-level middleware that sets the principal (see the proxy middleware section below).

## Visibility Prefixes

Object keys carry their intended visibility as a prefix:

| Constant | Value | Meaning |
| --- | --- | --- |
| `storage.PublicPrefix` | `pub/` | world-readable; default ACL grants read |
| `storage.PrivatePrefix` | `priv/` | controlled by business state via `FileACL` |

The storage resource emits keys under `pub/` or `priv/` depending on the upload's `public` flag. Proxy downloads serve `pub/*` anonymously and call `FileACL` for non-public keys; the storage backend itself does not enforce visibility.

## FileACL

`storage.FileACL` decides whether a principal may read a private key.

```go
type FileACL interface {
    CanRead(ctx context.Context, principal *security.Principal, key string) (bool, error)
}
```

Default behavior (`storage.DefaultFileACL`): grant read access only to keys under `pub/`. The proxy short-circuits `pub/*` before calling the ACL so public files work without an auth token; business code overrides `FileACL` via `vef.SupplyFileACL(...)` for private keys and ownership-aware reads.

## Storage Proxy Middleware

The module mounts an app-level download route:

```
GET /storage/files/<key>
```

Behavior:

| Surface | Behavior |
| --- | --- |
| Routing | app middleware named `storage_proxy` at order `900`; not an RPC action and not dispatched by the API engine |
| Key validation | URL-decodes `<key>` once; rejects empty keys, absolute paths, `..` segments, backslashes, NUL bytes, redundant slashes, and trailing slashes |
| Access | serves `pub/*` anonymously; every other key calls `FileACL.CanRead` with the request principal |
| Content type | uses backend metadata or extension detection, then sanitizes unsafe types to `application/octet-stream`; always sends `X-Content-Type-Options: nosniff` |
| Cache headers | `pub/*` gets `Cache-Control: public, max-age=3600, immutable` and an `ETag` when stat data has one; non-public keys get `Cache-Control: private, no-store` and no `ETag` |

## Upload Claims and Pending Delete (Lifecycle)

`init_upload` persists an `upload_claim` row owned by the calling principal with
status `pending`. `complete_upload` marks that same claim as `uploaded`. Until
the business model adopts the key (via `Files.OnCreate` / `OnUpdate`), the
object lives in a quarantined state — a periodic sweeper either recovers an
expired-but-completed multipart object by marking it uploaded, or enqueues the
abandoned object for asynchronous deletion (`DeleteReasonClaimExpired`).

Business writes therefore split into two transactional surfaces:

- **Claim consumer**: deletes the `upload_claim` row in the same transaction as the business insert.
- **Delete enqueuer**: inserts a `pending_delete` row for objects that should be reclaimed asynchronously (replaced field values, deleted business rows).

A background `DeleteWorker` then drains `pending_delete` rows against the backend and applies retry/backoff. Successfully drained rows emit `vef.storage.file.deleted`; rows that exhaust the retry budget are removed from the queue after emitting `vef.storage.delete.dead_letter`, which is the durable signal for manual investigation.

Storage fails fast at startup unless `vef.storage.file.claimed`,
`vef.storage.file.deleted`, and `vef.storage.delete.dead_letter` route through a
transactional event transport. In practice, enable the outbox transport and add
a route for `vef.storage.*` to `outbox`, or set the default event transport to
`outbox`.

## `Files` and `FilesFor[T]`

The high-level CRUD lifecycle facade — this replaced the older `Promoter[T]`:

```go
type Files interface {
    OnCreate(ctx, tx orm.DB, principal *security.Principal, model any) error
    OnUpdate(ctx, tx orm.DB, principal *security.Principal, oldModel, newModel any) error
    OnDelete(ctx, tx orm.DB, model any) error
}
```

Key semantics:

- All three methods **must run inside a business transaction** (`orm.DB.RunInTx`). The supplied `tx` is the business-DB instance, so claim consumption and pending-delete bookkeeping commit or roll back atomically with the business write.
- `OnCreate` / `OnUpdate` take a `*security.Principal` — only claims owned by that principal can be adopted. Nil / anonymous principals fail with `ErrAccessDenied`. Background jobs that legitimately operate on behalf of the system pass a synthetic system principal explicitly.
- `OnDelete` does not consume claims and therefore takes no principal; row ownership must be verified at the CRUD layer first.
- `FileClaimedEvent` is published through the **outbox transport inside the caller's transaction** (`event.WithTx`) — subscribers see the event only if the business transaction commits.

### Typed counterpart

`storage.FilesFor[T]` resolves the meta spec once at construction so the per-call reflect lookup disappears:

```go
files := storage.NewFilesFor[User](filesFacade)
err := files.OnCreate(ctx, tx, principal, &user)
```

CRUD lifecycle hooks are built on `FilesFor[T]`; custom hooks should follow the same pattern.

## Two Ways To Claim A File

When the user uploads a file, the framework keeps it in a "pending" state until your business code claims it. There are two ways to do that — pick the one that matches what your code already has:

- **Have a model struct?** Pass it to `Files` / `FilesFor[T]` and the framework will figure out the file fields by itself.
- **Just have a file key (or a list of keys)?** Call `ClaimConsumer.Consume(...)` directly.

Both end up doing the same thing — the second one is just the manual version of the first. Use whichever fits the call site better.

### Way 1: Pass in the struct (the easy way)

Tag the file fields with `meta:"uploaded_file"`, then hand the struct to `FilesFor[T]`. That's it.

```go
type Article struct {
    orm.FullAuditedModel
    CoverImage string   `json:"coverImage" bun:"cover_image" meta:"uploaded_file"`
    Gallery    []string `json:"gallery"    bun:"gallery,array" meta:"uploaded_file"`
    Body       string   `json:"body"       bun:"body"        meta:"rich_text"`
}

files := storage.NewFilesFor[Article](filesFacade)

err := db.RunInTx(ctx, func(ctx context.Context, tx orm.DB) error {
    if _, err := tx.NewInsert().Model(article).Exec(ctx); err != nil {
        return err
    }
    // Claim every file referenced by `article` in one call.
    return files.OnCreate(ctx, tx, principal, article)
})
```

On update, pass both the old and the new model — the framework claims the new files and queues the replaced ones for deletion:

```go
err := files.OnUpdate(ctx, tx, principal, oldArticle, newArticle)
```

On delete, pass the model — every referenced file gets queued for deletion:

```go
err := files.OnDelete(ctx, tx, article)
```

This is what regular CRUD already uses under the hood. If a struct fits, this is what you want.

### Way 2: Pass in the file key (when there's no struct)

Sometimes you don't have a model — maybe it's a background job, a custom upload flow, or you just want to claim one specific key. Inject `storage.ClaimConsumer` and call `Consume` with a slice of keys:

```go
err := db.RunInTx(ctx, func(ctx context.Context, tx orm.DB) error {
    if _, err := tx.NewInsert().Model(report).Exec(ctx); err != nil {
        return err
    }
    // Claim the file directly by its key.
    return claims.Consume(ctx, tx, principal, []string{report.FileKey})
})
```

If you also need to delete a file (e.g. the previous version), use `storage.DeleteEnqueuer`:

```go
err := deletes.Enqueue(ctx, tx,
    []string{oldKey},
    storage.DeleteReasonReplaced, // or DeleteReasonDeleted
)
```

A few things to keep in mind:

- Always call these inside `RunInTx` and pass the same `tx` — that's how the claim and your business write commit together.
- `Consume` only succeeds for files uploaded by **the same principal**. Trying to claim someone else's file returns `storage.ErrClaimNotFound`.
- A nil or anonymous principal returns `storage.ErrAccessDenied`. Background jobs need to construct a real system principal first.
- Empty / nil key slices are fine — they do nothing.
- Use `DeleteReasonReplaced` when overwriting a field, `DeleteReasonDeleted` when removing the owning record. `DeleteReasonClaimExpired` is for the framework only, don't pass it.

> If you ever catch yourself writing reflection to scan a struct's file fields, stop — that's exactly what `FilesFor[T]` does. Switch back to Way 1.

## Meta-Tagged Model Fields

Fields participate in the lifecycle by carrying a `meta` tag:

```go
type User struct {
    orm.FullAuditedModel

    Avatar   string            `json:"avatar" bun:"avatar" meta:"uploaded_file"`
    Gallery  []string          `json:"gallery" bun:"gallery,array" meta:"uploaded_file,category:gallery"`
    Profiles map[string]string `json:"profiles" bun:"profiles" meta:"uploaded_file"`
    Bio      string            `json:"bio" bun:"bio" meta:"rich_text"`
    Notes    string            `json:"notes" bun:"notes" meta:"markdown"`
}
```

| `meta` value | Field shape | Extraction strategy |
| --- | --- | --- |
| `uploaded_file` | `string` / `*string` / `[]string` / `map[string]string` | the value(s) are treated as file keys — for maps the **values** are the keys; the map's own keys are arbitrary labels |
| `rich_text` | `string` | scan HTML for embedded resource URLs and translate via `URLKeyMapper` |
| `markdown` | `string` | scan Markdown for embedded resource URLs and translate via `URLKeyMapper` |

Use `meta:"dive"` on a nested struct field when the file references live inside that nested struct; the scanner will recurse into the nested value and pick up its own `meta:"uploaded_file"`, `meta:"rich_text"`, and `meta:"markdown"` fields. Unsupported field shapes are ignored instead of producing refs.

`URLKeyMapper` translates rich-text/markdown URLs to storage keys during reconciliation. The framework DI graph supplies `storage.ProxyURLKeyMapper` by default, so content that embeds `/storage/files/<key>` is reconciled without extra wiring. If you call `storage.NewFiles(...)` directly, a nil mapper is normalized to `IdentityURLKeyMapper`; pass `&storage.IdentityURLKeyMapper{}` only when business content embeds bare keys directly.

The mapper surface is explicit in both directions: `URLToKey` consumes content
URLs during reconciliation, and `KeyToURL` is used when code needs to render
stored keys back into URLs.

Use `storage.ProxyURLKeyMapper{Prefix: storage.DefaultProxyPrefix}` when content
embeds the framework proxy URL form (`/storage/files/<key>`). The public helpers
`ReplaceHtmlURLs(content, replacements)` and `ReplaceMarkdownURLs(content,
replacements)` rewrite embedded URLs in rendered content, typically after mapping
storage keys through `URLKeyMapper.KeyToURL`.

## Storage Events

| Type constant / topic | Payload / constructor | JSON payload | Trigger |
| --- | --- | --- | --- |
| `EventTypeFileClaimed` / `vef.storage.file.claimed` | `FileClaimedEvent`; `NewFileClaimedEvent(key)` | `fileKey` | a previously pending claim was adopted by a business transaction (`Files.OnCreate` or update new-side) |
| `EventTypeFileDeleted` / `vef.storage.file.deleted` | `FileDeletedEvent`; `NewFileDeletedEvent(key, reason)` | `fileKey`, `reason` | the delete worker successfully removed an object from the backend |
| `EventTypeDeleteDeadLetter` / `vef.storage.delete.dead_letter` | `DeleteDeadLetterEvent`; `NewDeleteDeadLetterEvent(id, key, reason, attempts, lastErr)` | `pendingDeleteId`, `fileKey`, `reason`, `attempts`, optional `lastError` | the delete worker exhausted retries for a row; the queue row is removed after this event is published |

All three are published through the outbox transport with `event.WithTx(...)`. `FileClaimedEvent` shares the caller's business transaction; `FileDeletedEvent` and `DeleteDeadLetterEvent` share the delete worker's bookkeeping transaction. Subscribers attach with `event.WithGroup("...")` on the downstream sink transport and rely on the Inbox middleware for dedupe.

`DeleteReason` values forwarded onto the events:

| Reason | Wire value | Meaning |
| --- | --- | --- |
| `DeleteReasonReplaced` | `replaced` | an `uploaded_file` field was overwritten with a new key |
| `DeleteReasonDeleted` | `deleted` | the owning business row was deleted |
| `DeleteReasonClaimExpired` | `claim_expired` | a pending claim expired (framework-internal sweeper only) |

Dead-letter events carry a sanitized `lastError` classification rather than raw
backend errors. Current values are `access_denied`, `bucket_not_found`,
`session_not_found`, and `transient`.

Public supporting APIs:

| API group | Public surface |
| --- | --- |
| event constructors | `EventTypeFileClaimed`, `EventTypeFileDeleted`, `EventTypeDeleteDeadLetter`, `NewFileClaimedEvent`, `NewFileDeletedEvent`, `NewDeleteDeadLetterEvent` |
| facade constructors | `NewFiles`, `NewFilesFor`, `MultipartFor` |
| lifecycle services | `ClaimConsumer`, `DeleteEnqueuer`, `Files`, `FilesFor[T]` |
| storage interfaces | `Service`, `Multipart`, `FileACL`, `URLKeyMapper` |
| URL mappers | `DefaultFileACL`, `IdentityURLKeyMapper`, `ProxyURLKeyMapper`, `DefaultProxyPrefix` |
| metadata helpers | `CanonicalizeMetadataKeys` |
| option structs | `PutObjectOptions`, `GetObjectOptions`, `DeleteObjectOptions`, `DeleteObjectsOptions`, `CopyObjectOptions`, `StatObjectOptions`, `InitMultipartOptions`, `PutPartOptions`, `CompleteMultipartOptions`, `AbortMultipartOptions` |
| result structs | `ObjectInfo`, `MultipartSession`, `PartInfo`, `CompletedPart`, `FileRef` |
| meta constants | `MetaType`, `MetaTypeUploadedFile`, `MetaTypeRichText`, `MetaTypeMarkdown` |

`storage.CanonicalizeMetadataKeys(m)` returns a new metadata map whose keys use
the S3/HTTP-header canonical form, such as `author` to `Author`; nil or empty
input returns nil. Every backend applies this helper at the store boundary so
metadata round-trips in one provider-neutral shape.

## Errors

The storage package exposes two kinds of error values; match both with
`errors.Is`, but note only the first kind are plain Go sentinels.

Plain Go sentinels (`errors.New`, no API code or HTTP status of their own):

| Error | Cause |
| --- | --- |
| `storage.ErrUploadSessionNotFound` | multipart session already closed or never opened |
| `storage.ErrPartTooSmall` | non-final part smaller than `PartSize()` |
| `storage.ErrPartETagMismatch` | recorded part ETag disagrees with backend state during completion |
| `storage.ErrPartNumberOutOfRange` | parts don't cover `1..N` contiguously |
| `storage.ErrClaimNotFound` | a claim referenced by `Consume` doesn't exist or belongs to another principal |
| `storage.ErrAccessDenied` | anonymous / nil principal passed to a lifecycle method |
| `storage.ErrBucketNotFound` / `ErrObjectNotFound` / `ErrInvalidBucketName` | provider-level lookup failures |

`result.Err` business errors are carried through the API envelope with
`ErrCode*` constants in the `2200-2299` range; the response stays HTTP 200 and
the failure rides in the body `code`. Each `ErrCode*` constant pairs with a
`result.Err` value of the same name (`storage.ErrCodeUploadSizeMismatch` ↔
`storage.ErrUploadSizeMismatch`):

| Code | Error | i18n key | Trigger |
| --- | --- | --- | --- |
| `2200` | `storage.ErrInvalidFileKey` | `storage_invalid_file_key` | malformed object key on the download proxy (`/storage/files/<key>`) |
| `2201` | `storage.ErrFileNotFound` | `storage_file_not_found` | proxy download: object missing from the backend |
| `2202` | `storage.ErrFailedToGetFile` | `storage_failed_to_get_file` | proxy download: backend read or `FileACL` evaluation failed |
| `2203` | `storage.ErrClaimNotPending` | `storage_claim_not_pending` | `upload_part` / `list_parts` against a claim that is no longer `pending` (e.g. already completed) |
| `2204` | `storage.ErrClaimExpired` | `storage_claim_expired` | the claim's `expiresAt` (TTL `vef.storage.claim_ttl`, default 24h) has elapsed |
| `2205` | `storage.ErrUploadSizeExceedsLimit` | `storage_upload_size_exceeds_limit` | declared `size` exceeds `vef.storage.max_upload_size`; checked at `init_upload` and re-checked at `upload_part` / `complete_upload` |
| `2206` | `storage.ErrMultipartNotSupported` | `storage_multipart_not_supported` | the configured backend does not implement `storage.Multipart` |
| `2207` | `storage.ErrPublicUploadsNotAllowed` | `storage_public_uploads_not_allowed` | `public = true` while `vef.storage.allow_public_uploads` is false |
| `2208` | `storage.ErrUploadTooManyParts` | `storage_upload_too_many_parts` | the part plan exceeds the backend's `MaxPartCount()` |
| `2209` | `storage.ErrTooManyPendingUploads` | `storage_too_many_pending_uploads` | the principal already holds `vef.storage.max_pending_claims` pending claims |
| `2210` | `storage.ErrUploadRequiresMultipart` | `storage_upload_requires_multipart` | `upload_part` called with a JSON body instead of `multipart/form-data` |
| `2211` | `storage.ErrUploadRequiresFile` | `storage_upload_requires_file` | `upload_part` without a form part named `file` |
| `2212` | `storage.ErrClaimNotMultipart` | `storage_claim_not_multipart` | claim row without a bound backend session (defense in depth; arises only from an interrupted `init_upload`) |
| `2213` | `storage.ErrUploadPartNumberOutOfRange` | `storage_part_number_out_of_range` | `partNumber` above `partCount` (values below 1 already fail parameter validation) |
| `2214` | `storage.ErrUploadPartTooLarge` | `storage_upload_part_too_large` | a part larger than `partSize` |
| `2215` | `storage.ErrUploadPartTooSmall` | `storage_upload_part_too_small` | a non-final part smaller than `partSize` |
| `2216` | `storage.ErrUploadPartsIncomplete` | `storage_upload_parts_incomplete` | `complete_upload` with fewer recorded parts than `partCount` |
| `2217` | `storage.ErrUploadObjectNotFound` | `storage_object_not_found` | idempotent `complete_upload` retry: backend session closed and no object exists |
| `2218` | `storage.ErrUploadSizeMismatch` | `storage_upload_size_mismatch` | assembled object size differs from the declared `size`; the object is deleted before the error returns |
| `2219` | `storage.ErrAbortFailed` | `storage_abort_failed` | the backend refused to abort the multipart session |

Ownership violations and unknown claim IDs are not in this range: they answer
with the framework-generic `result.ErrAccessDenied` (code `1100`, HTTP 403) so
the API does not reveal whether a claim ID exists (`abort_upload` excepted —
see above).

## Minimal Service Example

```go
package avatars

import (
    "context"
    "strings"

    "github.com/coldsmirk/vef-framework-go/storage"
)

func SaveAvatar(ctx context.Context, svc storage.Service) error {
    _, err := svc.PutObject(ctx, storage.PutObjectOptions{
        Key:         "pub/avatars/user-1001.txt",
        Reader:      strings.NewReader("demo"),
        Size:        int64(len("demo")),
        ContentType: "text/plain",
    })

    return err
}
```

## CRUD Integration Pattern

For models with `meta`-tagged file fields, integrate via `FilesFor[T]` from a typed hook:

```go
filesUser := storage.NewFilesFor[User](filesFacade)

create := crud.NewCreate[User, UserParams]().
    AfterTx(func(ctx context.Context, tx orm.DB, principal *security.Principal, model *User) error {
        return filesUser.OnCreate(ctx, tx, principal, model)
    })
```

Generic CRUD already wires `FilesFor[T]` for the standard write builders (see [Hooks](../data-access/hooks)); custom write paths should follow the same pattern.

## Practical Advice

- Depend on `storage.Service` and `storage.Multipart`, not provider types.
- Keep all `Files` / `FilesFor[T]` calls inside the business transaction — that is the whole point of the facade.
- Treat unconfirmed objects as quarantined: the claim sweeper will eventually evict them; relying on raw `PutObject` keys without a claim bypasses lifecycle tracking.
- Register a real `FileACL` once you store private files; the default denies every `priv/*` read.
- Subscribe to `vef.storage.delete.dead_letter` for ops dashboards — the queue row is already retired, and the event carries the details operators need.
- Extension group names used by the module are `vef:api:resources` and `vef:app:middlewares`; use `vef.SupplyURLKeyMapper(...)` when replacing URL mapping.

## Next Step

Read [Custom Handlers](../building-apis/custom-handlers) to combine direct `storage.Service` use with business workflows, or [Event Bus](./event-bus) for the outbox transport that backs the lifecycle events.
